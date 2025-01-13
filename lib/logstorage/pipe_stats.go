package logstorage

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/cespare/xxhash/v2"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/bytesutil"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/cgroup"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/encoding"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/memory"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/slicesutil"
)

// pipeStats processes '| stats ...' queries.
//
// See https://docs.victoriametrics.com/victorialogs/logsql/#stats-pipe
type pipeStats struct {
	// byFields contains field names with optional buckets from 'by(...)' clause.
	byFields []*byStatsField

	// funcs contains stats functions to execute.
	funcs []pipeStatsFunc
}

type pipeStatsFunc struct {
	// f is stats function to execute
	f statsFunc

	// iff is an additional filter, which is applied to results before executing f on them
	iff *ifFilter

	// resultName is the name of the output generated by f
	resultName string
}

type statsFunc interface {
	// String returns string representation of statsFunc
	String() string

	// updateNeededFields update neededFields with the fields needed for calculating the given stats
	updateNeededFields(neededFields fieldsSet)

	// newStatsProcessor must create new statsProcessor for calculating stats for the given statsFunc
	newStatsProcessor(a *chunkedAllocator) statsProcessor
}

// statsProcessor must process stats for some statsFunc.
//
// All the statsProcessor methods are called from a single goroutine at a time,
// so there is no need in the internal synchronization.
type statsProcessor interface {
	// updateStatsForAllRows must update statsProcessor stats for all the rows in br.
	//
	// It must return the change of internal state size in bytes for the statsProcessor.
	updateStatsForAllRows(br *blockResult) int

	// updateStatsForRow must update statsProcessor stats for the row at rowIndex in br.
	//
	// It must return the change of internal state size in bytes for the statsProcessor.
	updateStatsForRow(br *blockResult, rowIndex int) int

	// mergeState must merge sfp state into statsProcessor state.
	mergeState(sfp statsProcessor)

	// finalizeStats must append string represetnation of the collected stats result to dst and return it.
	//
	// finalizeStats must immediately return if stopCh is closed.
	finalizeStats(dst []byte, stopCh <-chan struct{}) []byte
}

func (ps *pipeStats) String() string {
	s := "stats "
	if len(ps.byFields) > 0 {
		a := make([]string, len(ps.byFields))
		for i := range ps.byFields {
			a[i] = ps.byFields[i].String()
		}
		s += "by (" + strings.Join(a, ", ") + ") "
	}

	if len(ps.funcs) == 0 {
		logger.Panicf("BUG: pipeStats must contain at least a single statsFunc")
	}
	a := make([]string, len(ps.funcs))
	for i, f := range ps.funcs {
		line := f.f.String()
		if f.iff != nil {
			line += " " + f.iff.String()
		}
		line += " as " + quoteTokenIfNeeded(f.resultName)
		a[i] = line
	}
	s += strings.Join(a, ", ")
	return s
}

func (ps *pipeStats) canLiveTail() bool {
	return false
}

func (ps *pipeStats) updateNeededFields(neededFields, unneededFields fieldsSet) {
	neededFieldsOrig := neededFields.clone()
	neededFields.reset()

	// byFields are needed unconditionally, since the output number of rows depends on them.
	for _, bf := range ps.byFields {
		neededFields.add(bf.name)
	}

	for _, f := range ps.funcs {
		if neededFieldsOrig.contains(f.resultName) && !unneededFields.contains(f.resultName) {
			f.f.updateNeededFields(neededFields)
			if f.iff != nil {
				neededFields.addFields(f.iff.neededFields)
			}
		}
	}

	unneededFields.reset()
}

func (ps *pipeStats) hasFilterInWithQuery() bool {
	for _, f := range ps.funcs {
		if f.iff.hasFilterInWithQuery() {
			return true
		}
	}
	return false
}

func (ps *pipeStats) initFilterInValues(cache *inValuesCache, getFieldValuesFunc getFieldValuesFunc) (pipe, error) {
	funcsNew := make([]pipeStatsFunc, len(ps.funcs))
	for i := range ps.funcs {
		f := &ps.funcs[i]
		iffNew, err := f.iff.initFilterInValues(cache, getFieldValuesFunc)
		if err != nil {
			return nil, err
		}
		fNew := *f
		fNew.iff = iffNew
		funcsNew[i] = fNew
	}
	psNew := *ps
	psNew.funcs = funcsNew
	return &psNew, nil
}

func (ps *pipeStats) addByTimeField(step int64) {
	if step <= 0 {
		return
	}

	// add step to byFields
	stepStr := fmt.Sprintf("%d", step)
	dstFields := make([]*byStatsField, 0, len(ps.byFields)+1)
	hasByTime := false
	for _, f := range ps.byFields {
		if f.name == "_time" {
			f = &byStatsField{
				name:          "_time",
				bucketSizeStr: stepStr,
				bucketSize:    float64(step),
			}
			hasByTime = true
		}
		dstFields = append(dstFields, f)
	}
	if !hasByTime {
		dstFields = append(dstFields, &byStatsField{
			name:          "_time",
			bucketSizeStr: stepStr,
			bucketSize:    float64(step),
		})
	}
	ps.byFields = dstFields
}

func (ps *pipeStats) initRateFuncs(step int64) {
	if step <= 0 {
		return
	}

	stepSeconds := float64(step) / 1e9
	for _, f := range ps.funcs {
		switch t := f.f.(type) {
		case *statsRate:
			t.stepSeconds = stepSeconds
		case *statsRateSum:
			t.stepSeconds = stepSeconds
		}
	}
}

const stateSizeBudgetChunk = 1 << 20

func (ps *pipeStats) newPipeProcessor(workersCount int, stopCh <-chan struct{}, cancel func(), ppNext pipeProcessor) pipeProcessor {
	maxStateSize := int64(float64(memory.Allowed()) * 0.3)

	shards := make([]pipeStatsProcessorShard, workersCount)
	for i := range shards {
		shards[i] = pipeStatsProcessorShard{
			pipeStatsProcessorShardNopad: pipeStatsProcessorShardNopad{
				ps: ps,
			},
		}
		shards[i].init()
	}

	psp := &pipeStatsProcessor{
		ps:     ps,
		stopCh: stopCh,
		cancel: cancel,
		ppNext: ppNext,

		shards: shards,

		maxStateSize: maxStateSize,
	}
	psp.stateSizeBudget.Store(maxStateSize)

	return psp
}

type pipeStatsProcessor struct {
	ps     *pipeStats
	stopCh <-chan struct{}
	cancel func()
	ppNext pipeProcessor

	shards []pipeStatsProcessorShard

	maxStateSize    int64
	stateSizeBudget atomic.Int64
}

type pipeStatsProcessorShard struct {
	pipeStatsProcessorShardNopad

	// The padding prevents false sharing on widespread platforms with 128 mod (cache line size) = 0 .
	_ [128 - unsafe.Sizeof(pipeStatsProcessorShardNopad{})%128]byte
}

type pipeStatsProcessorShardNopad struct {
	ps *pipeStats

	m pipeStatsGroupMap

	// bms and brTmp are used for applying per-func filters.
	bms   []bitmap
	brTmp blockResult

	columnValues [][]string
	keyBuf       []byte

	stateSizeBudget int
}

type pipeStatsGroupMap struct {
	shard *pipeStatsProcessorShard

	u64        map[uint64]*pipeStatsGroup
	negative64 map[uint64]*pipeStatsGroup
	strings    map[string]*pipeStatsGroup

	// a and sfpsBuf are used for reducing memory allocations when calculating stats among big number of different groups.
	a       chunkedAllocator
	sfpsBuf []statsProcessor
}

func (psm *pipeStatsGroupMap) reset() {
	psm.shard = nil

	psm.u64 = nil
	psm.negative64 = nil
	psm.strings = nil

	psm.sfpsBuf = nil
}

func (psm *pipeStatsGroupMap) init(shard *pipeStatsProcessorShard) {
	psm.shard = shard

	psm.u64 = make(map[uint64]*pipeStatsGroup)
	psm.negative64 = make(map[uint64]*pipeStatsGroup)
	psm.strings = make(map[string]*pipeStatsGroup)
}

func (psm *pipeStatsGroupMap) entriesCount() uint64 {
	n := len(psm.u64) + len(psm.negative64) + len(psm.strings)
	return uint64(n)
}

func (psm *pipeStatsGroupMap) getPipeStatsGroupGeneric(key string) *pipeStatsGroup {
	if n, ok := tryParseUint64(key); ok {
		return psm.getPipeStatsGroupUint64(n)
	}
	if len(key) > 0 && key[0] == '-' {
		if n, ok := tryParseInt64(key); ok {
			return psm.getPipeStatsGroupNegativeInt64(n)
		}
	}
	return psm.getPipeStatsGroupString(bytesutil.ToUnsafeBytes(key))
}

func (psm *pipeStatsGroupMap) getPipeStatsGroupInt64(n int64) *pipeStatsGroup {
	if n >= 0 {
		return psm.getPipeStatsGroupUint64(uint64(n))
	}
	return psm.getPipeStatsGroupNegativeInt64(n)
}

func (psm *pipeStatsGroupMap) getPipeStatsGroupUint64(n uint64) *pipeStatsGroup {
	psg := psm.u64[n]
	if psg != nil {
		return psg
	}

	psg = psm.newPipeStatsGroup()
	psm.u64[n] = psg
	psm.shard.stateSizeBudget -= 8

	return psg
}

func (psm *pipeStatsGroupMap) getPipeStatsGroupNegativeInt64(n int64) *pipeStatsGroup {
	psg := psm.negative64[uint64(n)]
	if psg != nil {
		return psg
	}

	psg = psm.newPipeStatsGroup()
	psm.negative64[uint64(n)] = psg
	psm.shard.stateSizeBudget -= 8

	return psg
}

func (psm *pipeStatsGroupMap) getPipeStatsGroupString(key []byte) *pipeStatsGroup {
	psg := psm.strings[string(key)]
	if psg != nil {
		return psg
	}

	psg = psm.newPipeStatsGroup()
	keyCopy := psm.a.cloneBytesToString(key)
	psm.strings[keyCopy] = psg
	psm.shard.stateSizeBudget -= len(keyCopy) + int(unsafe.Sizeof(keyCopy))

	return psg
}

func (psm *pipeStatsGroupMap) newPipeStatsGroup() *pipeStatsGroup {
	sfps := psm.newStatsProcessors()

	for i, f := range psm.shard.ps.funcs {
		bytesAllocated := psm.a.bytesAllocated
		sfps[i] = f.f.newStatsProcessor(&psm.a)
		psm.shard.stateSizeBudget -= psm.a.bytesAllocated - bytesAllocated
	}

	psg := psm.a.newPipeStatsGroup()
	psg.funcs = psm.shard.ps.funcs
	psg.sfps = sfps
	psm.shard.stateSizeBudget -= int(unsafe.Sizeof(psg) + unsafe.Sizeof(sfps[0])*uintptr(len(sfps)))

	return psg
}

func (psm *pipeStatsGroupMap) newStatsProcessors() []statsProcessor {
	funcsLen := len(psm.shard.ps.funcs)
	if len(psm.sfpsBuf)+funcsLen > cap(psm.sfpsBuf) {
		psm.sfpsBuf = nil
	}
	if psm.sfpsBuf == nil {
		psm.sfpsBuf = make([]statsProcessor, 0, pipeStatsProcessorChunkLen)
	}

	sfpsBufLen := len(psm.sfpsBuf)
	psm.sfpsBuf = slicesutil.SetLength(psm.sfpsBuf, sfpsBufLen+funcsLen)
	return psm.sfpsBuf[sfpsBufLen:]
}

const pipeStatsProcessorChunkLen = 64 * 1024 / int(unsafe.Sizeof((statsProcessor)(nil)))

func (psm *pipeStatsGroupMap) mergeState(src *pipeStatsGroupMap, stopCh <-chan struct{}) {
	for n, psgSrc := range src.u64 {
		if needStop(stopCh) {
			return
		}
		psgDst := psm.u64[n]
		if psgDst == nil {
			psm.u64[n] = psgSrc
		} else {
			psgDst.mergeState(psgSrc)
		}
	}
	for n, psgSrc := range src.negative64 {
		if needStop(stopCh) {
			return
		}
		psgDst := psm.negative64[n]
		if psgDst == nil {
			psm.negative64[n] = psgSrc
		} else {
			psgDst.mergeState(psgSrc)
		}
	}
	for k, psgSrc := range src.strings {
		if needStop(stopCh) {
			return
		}
		psgDst := psm.strings[k]
		if psgDst == nil {
			psm.strings[k] = psgSrc
		} else {
			psgDst.mergeState(psgSrc)
		}
	}
}

func (shard *pipeStatsProcessorShard) init() {
	shard.m.init(shard)

	funcsLen := len(shard.ps.funcs)
	shard.bms = make([]bitmap, funcsLen)
}

func (shard *pipeStatsProcessorShard) writeBlock(br *blockResult) {
	byFields := shard.ps.byFields

	// Update shard.bms by applying per-function filters
	shard.applyPerFunctionFilters(br)

	// Process stats for the defined functions
	if len(byFields) == 0 {
		// Fast path - pass all the rows to a single group with empty key.
		psg := shard.m.getPipeStatsGroupString(nil)
		shard.stateSizeBudget -= psg.updateStatsForAllRows(shard.bms, br, &shard.brTmp)
		return
	}
	if len(byFields) == 1 {
		// Special case for grouping by a single column.
		shard.updateStatsSingleColumn(br, byFields[0])
		return
	}

	// Obtain columns for byFields
	columnValues := shard.columnValues[:0]
	for _, bf := range byFields {
		c := br.getColumnByName(bf.name)
		values := c.getValuesBucketed(br, bf)
		columnValues = append(columnValues, values)
	}
	shard.columnValues = columnValues

	// Verify whether all the 'by (...)' columns are constant.
	areAllConstColumns := true
	for _, values := range columnValues {
		if !areConstValues(values) {
			areAllConstColumns = false
			break
		}
	}
	if areAllConstColumns {
		// Fast path for constant 'by (...)' columns.
		keyBuf := shard.keyBuf[:0]
		for _, values := range columnValues {
			keyBuf = encoding.MarshalBytes(keyBuf, bytesutil.ToUnsafeBytes(values[0]))
		}
		psg := shard.m.getPipeStatsGroupString(keyBuf)
		shard.stateSizeBudget -= psg.updateStatsForAllRows(shard.bms, br, &shard.brTmp)
		shard.keyBuf = keyBuf
		return
	}

	// The slowest path - group by multiple columns with different values across rows.
	var psg *pipeStatsGroup
	keyBuf := shard.keyBuf[:0]
	for i := 0; i < br.rowsLen; i++ {
		// Verify whether the key for 'by (...)' fields equals the previous key
		sameValue := i > 0
		for _, values := range columnValues {
			if i <= 0 || values[i-1] != values[i] {
				sameValue = false
				break
			}
		}
		if !sameValue {
			// Construct new key for the 'by (...)' fields
			keyBuf = keyBuf[:0]
			for _, values := range columnValues {
				keyBuf = encoding.MarshalBytes(keyBuf, bytesutil.ToUnsafeBytes(values[i]))
			}
			psg = shard.m.getPipeStatsGroupString(keyBuf)
		}
		shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
	}
	shard.keyBuf = keyBuf
}

func (shard *pipeStatsProcessorShard) updateStatsSingleColumn(br *blockResult, bf *byStatsField) {
	c := br.getColumnByName(bf.name)
	if c.isConst {
		// Fast path for column with a constant value.
		v := br.getBucketedValue(c.valuesEncoded[0], bf)
		psg := shard.m.getPipeStatsGroupGeneric(v)
		shard.stateSizeBudget -= psg.updateStatsForAllRows(shard.bms, br, &shard.brTmp)
		return
	}

	if !bf.hasBucketConfig() {
		switch c.valueType {
		case valueTypeUint8:
			var psg *pipeStatsGroup
			values := c.getValuesEncoded(br)
			for i, v := range values {
				if i <= 0 || values[i-1] != v {
					n := unmarshalUint8(v)
					psg = shard.m.getPipeStatsGroupUint64(uint64(n))
				}
				shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
			}
			return
		case valueTypeUint16:
			var psg *pipeStatsGroup
			values := c.getValuesEncoded(br)
			for i, v := range values {
				if i <= 0 || values[i-1] != v {
					n := unmarshalUint16(v)
					psg = shard.m.getPipeStatsGroupUint64(uint64(n))
				}
				shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
			}
			return
		case valueTypeUint32:
			var psg *pipeStatsGroup
			values := c.getValuesEncoded(br)
			for i, v := range values {
				if i <= 0 || values[i-1] != v {
					n := unmarshalUint32(v)
					psg = shard.m.getPipeStatsGroupUint64(uint64(n))
				}
				shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
			}
			return
		case valueTypeUint64:
			var psg *pipeStatsGroup
			values := c.getValuesEncoded(br)
			for i, v := range values {
				if i <= 0 || values[i-1] != v {
					n := unmarshalUint64(v)
					psg = shard.m.getPipeStatsGroupUint64(n)
				}
				shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
			}
			return
		case valueTypeInt64:
			var psg *pipeStatsGroup
			values := c.getValuesEncoded(br)
			for i, v := range values {
				if i <= 0 || values[i-1] != v {
					n := unmarshalInt64(v)
					psg = shard.m.getPipeStatsGroupInt64(n)
				}
				shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
			}
			return
		}
	}

	// Slower generic path for a column with different values.
	var psg *pipeStatsGroup
	values := c.getValuesBucketed(br, bf)
	for i := 0; i < br.rowsLen; i++ {
		if i <= 0 || values[i-1] != values[i] {
			psg = shard.m.getPipeStatsGroupGeneric(values[i])
		}
		shard.stateSizeBudget -= psg.updateStatsForRow(shard.bms, br, i)
	}
}

func (shard *pipeStatsProcessorShard) applyPerFunctionFilters(br *blockResult) {
	funcs := shard.ps.funcs
	for i := range funcs {
		iff := funcs[i].iff
		if iff == nil {
			continue
		}

		bm := &shard.bms[i]
		bm.init(br.rowsLen)
		bm.setBits()

		iff.f.applyToBlockResult(br, bm)
	}
}

type pipeStatsGroup struct {
	funcs []pipeStatsFunc
	sfps  []statsProcessor
}

func (psg *pipeStatsGroup) mergeState(src *pipeStatsGroup) {
	for i, sfp := range psg.sfps {
		sfp.mergeState(src.sfps[i])
	}
}

func (psg *pipeStatsGroup) updateStatsForAllRows(bms []bitmap, br, brTmp *blockResult) int {
	n := 0
	for i, sfp := range psg.sfps {
		iff := psg.funcs[i].iff
		if iff == nil {
			n += sfp.updateStatsForAllRows(br)
		} else {
			brTmp.initFromFilterAllColumns(br, &bms[i])
			n += sfp.updateStatsForAllRows(brTmp)
		}
	}
	return n
}

func (psg *pipeStatsGroup) updateStatsForRow(bms []bitmap, br *blockResult, rowIdx int) int {
	n := 0
	for i, sfp := range psg.sfps {
		iff := psg.funcs[i].iff
		if iff == nil || bms[i].isSetBit(rowIdx) {
			n += sfp.updateStatsForRow(br, rowIdx)
		}
	}
	return n
}

func (psp *pipeStatsProcessor) writeBlock(workerID uint, br *blockResult) {
	if br.rowsLen == 0 {
		return
	}

	shard := &psp.shards[workerID]

	for shard.stateSizeBudget < 0 {
		// steal some budget for the state size from the global budget.
		remaining := psp.stateSizeBudget.Add(-stateSizeBudgetChunk)
		if remaining < 0 {
			// The state size is too big. Stop processing data in order to avoid OOM crash.
			if remaining+stateSizeBudgetChunk >= 0 {
				// Notify worker goroutines to stop calling writeBlock() in order to save CPU time.
				psp.cancel()
			}
			return
		}
		shard.stateSizeBudget += stateSizeBudgetChunk
	}

	shard.writeBlock(br)
}

func (psp *pipeStatsProcessor) flush() error {
	if n := psp.stateSizeBudget.Load(); n <= 0 {
		return fmt.Errorf("cannot calculate [%s], since it requires more than %dMB of memory", psp.ps.String(), psp.maxStateSize/(1<<20))
	}

	// Merge states across shards in parallel
	psms, err := psp.mergeShardsParallel()
	if err != nil {
		return err
	}
	if needStop(psp.stopCh) {
		return nil
	}

	if len(psp.ps.byFields) == 0 && len(psms) == 0 {
		// Special case - zero matching rows.
		shard := &psp.shards[0]
		shard.init()
		_ = shard.m.getPipeStatsGroupString(nil)
		psms = append(psms, &shard.m)
	}

	// Write the calculated stats in parallel to the next pipe.
	var wg sync.WaitGroup
	for i := range psms {
		wg.Add(1)
		go func(workerID uint) {
			defer wg.Done()

			psw := newPipeStatsWriter(psp, workerID)
			psw.writeShardData(psms[workerID])
			psw.flush()
		}(uint(i))
	}
	wg.Wait()

	return nil
}

type pipeStatsWriter struct {
	psp      *pipeStatsProcessor
	workerID uint

	rcs []resultColumn
	br  blockResult

	resultLen int
	rowsCount int

	values    []string
	valuesBuf []byte
}

func newPipeStatsWriter(psp *pipeStatsProcessor, workerID uint) *pipeStatsWriter {
	byFields := psp.ps.byFields
	rcs := make([]resultColumn, 0, len(byFields)+len(psp.ps.funcs))
	for _, bf := range byFields {
		rcs = appendResultColumnWithName(rcs, bf.name)
	}
	for _, f := range psp.ps.funcs {
		rcs = appendResultColumnWithName(rcs, f.resultName)
	}

	psw := &pipeStatsWriter{
		psp:      psp,
		workerID: workerID,
		rcs:      rcs,
	}
	return psw
}

func (psw *pipeStatsWriter) writePipeStatsGroup(psg *pipeStatsGroup) {
	for _, sfp := range psg.sfps {
		bufLen := len(psw.valuesBuf)
		psw.valuesBuf = sfp.finalizeStats(psw.valuesBuf, psw.psp.stopCh)
		value := bytesutil.ToUnsafeString(psw.valuesBuf[bufLen:])
		psw.values = append(psw.values, value)
	}
	if len(psw.values) != len(psw.rcs) {
		logger.Panicf("BUG: len(values)=%d must be equal to len(rcs)=%d", len(psw.values), len(psw.rcs))
	}

	n := 0
	for i, v := range psw.values {
		psw.rcs[i].addValue(v)
		n += len(v)
	}
	psw.resultLen += n
	psw.rowsCount++

	// The 64_000 limit provides the best performance results when generating stats
	// over big number of distinct groups.
	if psw.resultLen >= 64_000 {
		psw.flush()
	}
}

func (psw *pipeStatsWriter) flush() {
	psw.br.setResultColumns(psw.rcs, psw.rowsCount)
	psw.resultLen = 0
	psw.rowsCount = 0
	psw.psp.ppNext.writeBlock(psw.workerID, &psw.br)
	psw.br.reset()
	for i := range psw.rcs {
		psw.rcs[i].resetValues()
	}
	psw.values = psw.values[:0]
	psw.valuesBuf = psw.valuesBuf[:0]
}

func (psw *pipeStatsWriter) writeShardData(psm *pipeStatsGroupMap) {
	byFields := psw.psp.ps.byFields
	if len(byFields) == 1 {
		for n, psg := range psm.u64 {
			if needStop(psw.psp.stopCh) {
				return
			}
			psw.values = psw.values[:0]

			// Reconstruct value for byFields[0]
			valuesBufLen := len(psw.valuesBuf)
			psw.valuesBuf = marshalUint64String(psw.valuesBuf, n)
			psw.values = append(psw.values, bytesutil.ToUnsafeString(psw.valuesBuf[valuesBufLen:]))

			psw.writePipeStatsGroup(psg)
		}
		for n, psg := range psm.negative64 {
			if needStop(psw.psp.stopCh) {
				return
			}
			psw.values = psw.values[:0]

			// Reconstruct value for byFields[0]
			valuesBufLen := len(psw.valuesBuf)
			psw.valuesBuf = marshalInt64String(psw.valuesBuf, int64(n))
			psw.values = append(psw.values, bytesutil.ToUnsafeString(psw.valuesBuf[valuesBufLen:]))

			psw.writePipeStatsGroup(psg)
		}
		for key, psg := range psm.strings {
			if needStop(psw.psp.stopCh) {
				return
			}
			psw.values = psw.values[:0]

			psw.values = append(psw.values, key)
			psw.writePipeStatsGroup(psg)
		}
	} else {
		for key, psg := range psm.strings {
			if needStop(psw.psp.stopCh) {
				return
			}
			psw.values = psw.values[:0]

			// Unmarshal values for byFields from key.
			keyBuf := bytesutil.ToUnsafeBytes(key)
			for len(keyBuf) > 0 {
				v, nSize := encoding.UnmarshalBytes(keyBuf)
				if nSize <= 0 {
					logger.Panicf("BUG: cannot unmarshal value from keyBuf=%q", keyBuf)
				}
				keyBuf = keyBuf[nSize:]
				psw.values = append(psw.values, bytesutil.ToUnsafeString(v))
			}
			if len(psw.values) != len(byFields) {
				logger.Panicf("BUG: unexpected number of values decoded from keyBuf; got %d; want %d", len(psw.values), len(byFields))
			}

			psw.writePipeStatsGroup(psg)
		}
	}
}

func (psp *pipeStatsProcessor) mergeShardsParallel() ([]*pipeStatsGroupMap, error) {
	shards := psp.shards
	shardsLen := len(shards)
	cpusCount := cgroup.AvailableCPUs()

	if shardsLen == 1 {
		var psms []*pipeStatsGroupMap
		shard := &shards[0]
		if shard.m.entriesCount() > 0 {
			psms = append(psms, &shard.m)
		}
		return psms, nil
	}

	var wg sync.WaitGroup
	perShardMaps := make([][]pipeStatsGroupMap, shardsLen)
	for i := range shards {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			perCPU := make([]pipeStatsGroupMap, cpusCount)
			for i := range perCPU {
				perCPU[i].init(&shards[idx])
			}

			psm := &shards[idx].m

			for n, psg := range psm.u64 {
				if needStop(psp.stopCh) {
					return
				}
				k := unsafe.Slice((*byte)(unsafe.Pointer(&n)), 8)
				h := xxhash.Sum64(k)
				shardIdx := h % uint64(len(perCPU))
				perCPU[shardIdx].u64[n] = psg
			}
			for n, psg := range psm.negative64 {
				if needStop(psp.stopCh) {
					return
				}
				k := unsafe.Slice((*byte)(unsafe.Pointer(&n)), 8)
				h := xxhash.Sum64(k)
				shardIdx := h % uint64(len(perCPU))
				perCPU[shardIdx].negative64[n] = psg
			}
			for k, psg := range psm.strings {
				if needStop(psp.stopCh) {
					return
				}
				h := xxhash.Sum64(bytesutil.ToUnsafeBytes(k))
				shardIdx := h % uint64(len(perCPU))
				perCPU[shardIdx].strings[k] = psg
			}

			perShardMaps[idx] = perCPU
			psm.reset()
		}(i)
	}
	wg.Wait()
	if needStop(psp.stopCh) {
		return nil, nil
	}
	if n := psp.stateSizeBudget.Load(); n < 0 {
		return nil, fmt.Errorf("cannot calculate [%s], since it requires more than %dMB of memory", psp.ps.String(), psp.maxStateSize/(1<<20))
	}

	// Merge per-shard entries into perShardMaps[0]
	for i := 0; i < cpusCount; i++ {
		wg.Add(1)
		go func(cpuIdx int) {
			defer wg.Done()

			psm := &perShardMaps[0][cpuIdx]
			for _, perCPU := range perShardMaps[1:] {
				psm.mergeState(&perCPU[cpuIdx], psp.stopCh)
				perCPU[cpuIdx].reset()
			}
		}(i)
	}
	wg.Wait()
	if needStop(psp.stopCh) {
		return nil, nil
	}
	if n := psp.stateSizeBudget.Load(); n < 0 {
		return nil, fmt.Errorf("cannot calculate [%s], since it requires more than %dMB of memory", psp.ps.String(), psp.maxStateSize/(1<<20))
	}

	// Filter out maps without entries
	psms := perShardMaps[0]
	result := make([]*pipeStatsGroupMap, 0, len(psms))
	for i := range psms {
		if psms[i].entriesCount() > 0 {
			result = append(result, &psms[i])
		}
	}

	return result, nil
}

func parsePipeStats(lex *lexer, needStatsKeyword bool) (pipe, error) {
	if needStatsKeyword {
		if !lex.isKeyword("stats") {
			return nil, fmt.Errorf("expecting 'stats'; got %q", lex.token)
		}
		lex.nextToken()
	}

	var ps pipeStats
	if lex.isKeyword("by", "(") {
		if lex.isKeyword("by") {
			lex.nextToken()
		}
		bfs, err := parseByStatsFields(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'by' clause: %w", err)
		}
		ps.byFields = bfs
	}

	seenByFields := make(map[string]*byStatsField, len(ps.byFields))
	for _, bf := range ps.byFields {
		seenByFields[bf.name] = bf
	}

	seenResultNames := make(map[string]statsFunc)

	var funcs []pipeStatsFunc
	for {
		var f pipeStatsFunc

		sf, err := parseStatsFunc(lex)
		if err != nil {
			return nil, err
		}
		f.f = sf

		if lex.isKeyword("if") {
			iff, err := parseIfFilter(lex)
			if err != nil {
				return nil, fmt.Errorf("cannot parse 'if' filter for [%s]: %w", sf, err)
			}
			f.iff = iff
		}

		resultName := ""
		if lex.isKeyword(",", "|", ")", "") {
			resultName = sf.String()
			if f.iff != nil {
				resultName += " " + f.iff.String()
			}
		} else {
			if lex.isKeyword("as") {
				lex.nextToken()
			}
			fieldName, err := parseFieldName(lex)
			if err != nil {
				return nil, fmt.Errorf("cannot parse result name for [%s]: %w", sf, err)
			}
			resultName = fieldName
		}
		if bf := seenByFields[resultName]; bf != nil {
			return nil, fmt.Errorf("the %q is used as 'by' field [%s], so it cannot be used as result name for [%s]", resultName, bf, sf)
		}
		if sfPrev := seenResultNames[resultName]; sfPrev != nil {
			return nil, fmt.Errorf("cannot use identical result name %q for [%s] and [%s]", resultName, sfPrev, sf)
		}
		seenResultNames[resultName] = sf
		f.resultName = resultName

		funcs = append(funcs, f)

		if lex.isKeyword("|", ")", "") {
			ps.funcs = funcs
			return &ps, nil
		}
		if !lex.isKeyword(",") {
			return nil, fmt.Errorf("unexpected token %q after [%s]; want ',', '|' or ')'", lex.token, sf)
		}
		lex.nextToken()
	}
}

func parseStatsFunc(lex *lexer) (statsFunc, error) {
	switch {
	case lex.isKeyword("avg"):
		sas, err := parseStatsAvg(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'avg' func: %w", err)
		}
		return sas, nil
	case lex.isKeyword("count"):
		scs, err := parseStatsCount(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'count' func: %w", err)
		}
		return scs, nil
	case lex.isKeyword("count_empty"):
		scs, err := parseStatsCountEmpty(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'count_empty' func: %w", err)
		}
		return scs, nil
	case lex.isKeyword("count_uniq"):
		sus, err := parseStatsCountUniq(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'count_uniq' func: %w", err)
		}
		return sus, nil
	case lex.isKeyword("count_uniq_hash"):
		sus, err := parseStatsCountUniqHash(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'count_uniq_hash' func: %w", err)
		}
		return sus, nil
	case lex.isKeyword("max"):
		sms, err := parseStatsMax(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'max' func: %w", err)
		}
		return sms, nil
	case lex.isKeyword("median"):
		sms, err := parseStatsMedian(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'median' func: %w", err)
		}
		return sms, nil
	case lex.isKeyword("min"):
		sms, err := parseStatsMin(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'min' func: %w", err)
		}
		return sms, nil
	case lex.isKeyword("quantile"):
		sqs, err := parseStatsQuantile(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'quantile' func: %w", err)
		}
		return sqs, nil
	case lex.isKeyword("rate"):
		srs, err := parseStatsRate(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'rate' func: %w", err)
		}
		return srs, nil
	case lex.isKeyword("rate_sum"):
		srs, err := parseStatsRateSum(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'rate_sum' func: %w", err)
		}
		return srs, nil
	case lex.isKeyword("row_any"):
		sas, err := parseStatsRowAny(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'row_any' func: %w", err)
		}
		return sas, nil
	case lex.isKeyword("row_max"):
		sms, err := parseStatsRowMax(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'row_max' func: %w", err)
		}
		return sms, nil
	case lex.isKeyword("row_min"):
		sms, err := parseStatsRowMin(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'row_min' func: %w", err)
		}
		return sms, nil
	case lex.isKeyword("sum"):
		sss, err := parseStatsSum(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'sum' func: %w", err)
		}
		return sss, nil
	case lex.isKeyword("sum_len"):
		sss, err := parseStatsSumLen(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'sum_len' func: %w", err)
		}
		return sss, nil
	case lex.isKeyword("uniq_values"):
		sus, err := parseStatsUniqValues(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'uniq_values' func: %w", err)
		}
		return sus, nil
	case lex.isKeyword("values"):
		svs, err := parseStatsValues(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'values' func: %w", err)
		}
		return svs, nil
	default:
		return nil, fmt.Errorf("unknown stats func %q", lex.token)
	}
}

var statsNames = []string{
	"avg",
	"count",
	"count_empty",
	"count_uniq",
	"count_uniq_hash",
	"max",
	"median",
	"min",
	"quantile",
	"rate",
	"rate_sum",
	"row_any",
	"row_max",
	"row_min",
	"sum",
	"sum_len",
	"uniq_values",
	"values",
}

var zeroByStatsField = &byStatsField{}

// byStatsField represents 'by (...)' part of the pipeStats.
//
// It can have either 'name' representation or 'name:bucket' or 'name:bucket offset off' representation,
// where `bucket` and `off` can contain duration, size or numeric value for creating different buckets
// for 'value/bucket'.
type byStatsField struct {
	name string

	// bucketSizeStr is string representation of the bucket size
	bucketSizeStr string

	// bucketSize is the bucket for grouping the given field values with value/bucketSize calculations
	bucketSize float64

	// bucketOffsetStr is string representation of the offset for bucketSize
	bucketOffsetStr string

	// bucketOffset is the offset for bucketSize
	bucketOffset float64
}

func (bf *byStatsField) String() string {
	s := quoteTokenIfNeeded(bf.name)
	if bf.bucketSizeStr != "" {
		s += ":" + bf.bucketSizeStr
		if bf.bucketOffsetStr != "" {
			s += " offset " + bf.bucketOffsetStr
		}
	}
	return s
}

func (bf *byStatsField) hasBucketConfig() bool {
	return len(bf.bucketSizeStr) > 0 || len(bf.bucketOffsetStr) > 0
}

func parseByStatsFields(lex *lexer) ([]*byStatsField, error) {
	if !lex.isKeyword("(") {
		return nil, fmt.Errorf("missing `(`")
	}
	var bfs []*byStatsField
	for {
		lex.nextToken()
		if lex.isKeyword(")") {
			lex.nextToken()
			return bfs, nil
		}
		fieldName, err := getCompoundPhrase(lex, false)
		if err != nil {
			return nil, fmt.Errorf("cannot parse field name: %w", err)
		}
		fieldName = getCanonicalColumnName(fieldName)
		bf := &byStatsField{
			name: fieldName,
		}
		if lex.isKeyword(":") {
			// Parse bucket size
			lex.nextToken()
			bucketSizeStr := lex.token
			lex.nextToken()
			if bucketSizeStr == "/" {
				bucketSizeStr += lex.token
				lex.nextToken()
			}
			if bucketSizeStr != "year" && bucketSizeStr != "month" {
				bucketSize, ok := tryParseBucketSize(bucketSizeStr)
				if !ok {
					return nil, fmt.Errorf("cannot parse bucket size for field %q: %q", fieldName, bucketSizeStr)
				}
				bf.bucketSize = bucketSize
			}
			bf.bucketSizeStr = bucketSizeStr

			// Parse bucket offset
			if lex.isKeyword("offset") {
				lex.nextToken()
				bucketOffsetStr := lex.token
				lex.nextToken()
				if bucketOffsetStr == "-" {
					bucketOffsetStr += lex.token
					lex.nextToken()
				}
				bucketOffset, ok := tryParseBucketOffset(bucketOffsetStr)
				if !ok {
					return nil, fmt.Errorf("cannot parse bucket offset for field %q: %q", fieldName, bucketOffsetStr)
				}
				bf.bucketOffsetStr = bucketOffsetStr
				bf.bucketOffset = bucketOffset
			}
		}
		bfs = append(bfs, bf)
		switch {
		case lex.isKeyword(")"):
			lex.nextToken()
			return bfs, nil
		case lex.isKeyword(","):
		default:
			return nil, fmt.Errorf("unexpected token: %q; expecting ',' or ')'", lex.token)
		}
	}
}

// tryParseBucketOffset tries parsing bucket offset, which can have the following formats:
//
// - integer number: 12345
// - floating-point number: 1.2345
// - duration: 1.5s - it is converted to nanoseconds
// - bytes: 1.5KiB
func tryParseBucketOffset(s string) (float64, bool) {
	// Try parsing s as floating point number
	if f, ok := tryParseFloat64(s); ok {
		return f, true
	}

	// Try parsing s as duration (1s, 5m, etc.)
	if nsecs, ok := tryParseDuration(s); ok {
		return float64(nsecs), true
	}

	// Try parsing s as bytes (KiB, MB, etc.)
	if n, ok := tryParseBytes(s); ok {
		return float64(n), true
	}

	return 0, false
}

// tryParseBucketSize tries parsing bucket size, which can have the following formats:
//
// - integer number: 12345
// - floating-point number: 1.2345
// - duration: 1.5s - it is converted to nanoseconds
// - bytes: 1.5KiB
// - ipv4 mask: /24
func tryParseBucketSize(s string) (float64, bool) {
	switch s {
	case "nanosecond":
		return 1, true
	case "microsecond":
		return nsecsPerMicrosecond, true
	case "millisecond":
		return nsecsPerMillisecond, true
	case "second":
		return nsecsPerSecond, true
	case "minute":
		return nsecsPerMinute, true
	case "hour":
		return nsecsPerHour, true
	case "day":
		return nsecsPerDay, true
	case "week":
		return nsecsPerWeek, true
	}

	// Try parsing s as floating point number
	if f, ok := tryParseFloat64(s); ok {
		return f, true
	}

	// Try parsing s as duration (1s, 5m, etc.)
	if nsecs, ok := tryParseDuration(s); ok {
		return float64(nsecs), true
	}

	// Try parsing s as bytes (KiB, MB, etc.)
	if n, ok := tryParseBytes(s); ok {
		return float64(n), true
	}

	if n, ok := tryParseIPv4Mask(s); ok {
		return float64(n), true
	}

	return 0, false
}

func parseFieldNamesInParens(lex *lexer) ([]string, error) {
	if !lex.isKeyword("(") {
		return nil, fmt.Errorf("missing `(`")
	}
	var fields []string
	for {
		lex.nextToken()
		if lex.isKeyword(")") {
			lex.nextToken()
			return fields, nil
		}
		if lex.isKeyword(",") {
			return nil, fmt.Errorf("unexpected `,`")
		}
		field, err := parseFieldName(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse field name: %w", err)
		}
		fields = append(fields, field)
		switch {
		case lex.isKeyword(")"):
			lex.nextToken()
			return fields, nil
		case lex.isKeyword(","):
		default:
			return nil, fmt.Errorf("unexpected token: %q; expecting ',' or ')'", lex.token)
		}
	}
}

func parseFieldName(lex *lexer) (string, error) {
	fieldName, err := getCompoundToken(lex)
	if err != nil {
		return "", fmt.Errorf("cannot parse field name: %w", err)
	}
	fieldName = getCanonicalColumnName(fieldName)
	return fieldName, nil
}

func fieldNamesString(fields []string) string {
	a := make([]string, len(fields))
	for i, f := range fields {
		if f != "*" {
			f = quoteTokenIfNeeded(f)
		}
		a[i] = f
	}
	return strings.Join(a, ", ")
}

func areConstValues(values []string) bool {
	if len(values) == 0 {
		return false
	}
	v := values[0]
	for i := 1; i < len(values); i++ {
		if v != values[i] {
			return false
		}
	}
	return true
}
