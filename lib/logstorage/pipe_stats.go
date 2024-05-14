package logstorage

import (
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/bytesutil"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/encoding"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/memory"
)

// pipeStats processes '| stats ...' queries.
//
// See https://docs.victoriametrics.com/victorialogs/logsql/#stats-pipe
type pipeStats struct {
	// byFields contains field names with optional buckets from 'by(...)' clause.
	byFields []*byStatsField

	// resultNames contains names of output results generated by funcs.
	resultNames []string

	// funcs contains stats functions to execute.
	funcs []statsFunc
}

type statsFunc interface {
	// String returns string representation of statsFunc
	String() string

	// neededFields returns the needed fields for calculating the given stats
	neededFields() []string

	// newStatsProcessor must create new statsProcessor for calculating stats for the given statsFunc.
	//
	// It also must return the size in bytes of the returned statsProcessor.
	newStatsProcessor() (statsProcessor, int)
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

	// finalizeStats must return the collected stats result from statsProcessor.
	finalizeStats() string
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
		a[i] = f.String() + " as " + quoteTokenIfNeeded(ps.resultNames[i])
	}
	s += strings.Join(a, ", ")
	return s
}

func (ps *pipeStats) updateNeededFields(neededFields, unneededFields fieldsSet) {
	neededFieldsOrig := neededFields.clone()
	neededFields.reset()

	byFields := make([]string, len(ps.byFields))
	for i, bf := range ps.byFields {
		byFields[i] = bf.name
	}

	for _, f := range byFields {
		if neededFieldsOrig.contains(f) && !unneededFields.contains(f) {
			neededFields.addAll(byFields)
		}
	}

	for i, resultName := range ps.resultNames {
		if neededFieldsOrig.contains(resultName) && !unneededFields.contains(resultName) {
			funcFields := ps.funcs[i].neededFields()
			neededFields.addAll(byFields)
			neededFields.addAll(funcFields)
		}
	}

	unneededFields.reset()
}

const stateSizeBudgetChunk = 1 << 20

func (ps *pipeStats) newPipeProcessor(workersCount int, stopCh <-chan struct{}, cancel func(), ppBase pipeProcessor) pipeProcessor {
	maxStateSize := int64(float64(memory.Allowed()) * 0.3)

	shards := make([]pipeStatsProcessorShard, workersCount)
	for i := range shards {
		shard := &shards[i]
		shard.ps = ps
		shard.m = make(map[string]*pipeStatsGroup)
		shard.stateSizeBudget = stateSizeBudgetChunk
		maxStateSize -= stateSizeBudgetChunk
	}

	psp := &pipeStatsProcessor{
		ps:     ps,
		stopCh: stopCh,
		cancel: cancel,
		ppBase: ppBase,

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
	ppBase pipeProcessor

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
	m  map[string]*pipeStatsGroup

	columnValues [][]string
	keyBuf       []byte

	stateSizeBudget int
}

func (shard *pipeStatsProcessorShard) writeBlock(br *blockResult) {
	byFields := shard.ps.byFields

	if len(byFields) == 0 {
		// Fast path - pass all the rows to a single group with empty key.
		psg := shard.getPipeStatsGroup(nil)
		shard.stateSizeBudget -= psg.updateStatsForAllRows(br)
		return
	}
	if len(byFields) == 1 {
		// Special case for grouping by a single column.
		bf := byFields[0]
		c := br.getColumnByName(bf.name)
		if c.isConst {
			// Fast path for column with constant value.
			v := br.getBucketedValue(c.encodedValues[0], bf)
			shard.keyBuf = encoding.MarshalBytes(shard.keyBuf[:0], bytesutil.ToUnsafeBytes(v))
			psg := shard.getPipeStatsGroup(shard.keyBuf)
			shard.stateSizeBudget -= psg.updateStatsForAllRows(br)
			return
		}

		values := c.getBucketedValues(br, bf)
		if areConstValues(values) {
			// Fast path for column with constant values.
			shard.keyBuf = encoding.MarshalBytes(shard.keyBuf[:0], bytesutil.ToUnsafeBytes(values[0]))
			psg := shard.getPipeStatsGroup(shard.keyBuf)
			shard.stateSizeBudget -= psg.updateStatsForAllRows(br)
			return
		}

		// Slower generic path for a column with different values.
		var psg *pipeStatsGroup
		keyBuf := shard.keyBuf[:0]
		for i := range br.timestamps {
			if i <= 0 || values[i-1] != values[i] {
				keyBuf = encoding.MarshalBytes(keyBuf[:0], bytesutil.ToUnsafeBytes(values[i]))
				psg = shard.getPipeStatsGroup(keyBuf)
			}
			shard.stateSizeBudget -= psg.updateStatsForRow(br, i)
		}
		shard.keyBuf = keyBuf
		return
	}

	// Obtain columns for byFields
	columnValues := shard.columnValues[:0]
	for _, bf := range byFields {
		c := br.getColumnByName(bf.name)
		values := c.getBucketedValues(br, bf)
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
		psg := shard.getPipeStatsGroup(keyBuf)
		shard.stateSizeBudget -= psg.updateStatsForAllRows(br)
		shard.keyBuf = keyBuf
		return
	}

	// The slowest path - group by multiple columns with different values across rows.
	var psg *pipeStatsGroup
	keyBuf := shard.keyBuf[:0]
	for i := range br.timestamps {
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
			psg = shard.getPipeStatsGroup(keyBuf)
		}
		shard.stateSizeBudget -= psg.updateStatsForRow(br, i)
	}
	shard.keyBuf = keyBuf
}

func (shard *pipeStatsProcessorShard) getPipeStatsGroup(key []byte) *pipeStatsGroup {
	psg := shard.m[string(key)]
	if psg != nil {
		return psg
	}

	sfps := make([]statsProcessor, len(shard.ps.funcs))
	for i, f := range shard.ps.funcs {
		sfp, stateSize := f.newStatsProcessor()
		sfps[i] = sfp
		shard.stateSizeBudget -= stateSize
	}
	psg = &pipeStatsGroup{
		sfps: sfps,
	}
	shard.m[string(key)] = psg
	shard.stateSizeBudget -= len(key) + int(unsafe.Sizeof("")+unsafe.Sizeof(psg)+unsafe.Sizeof(sfps[0])*uintptr(len(sfps)))

	return psg
}

type pipeStatsGroup struct {
	sfps []statsProcessor
}

func (psg *pipeStatsGroup) updateStatsForAllRows(br *blockResult) int {
	n := 0
	for _, sfp := range psg.sfps {
		n += sfp.updateStatsForAllRows(br)
	}
	return n
}

func (psg *pipeStatsGroup) updateStatsForRow(br *blockResult, rowIdx int) int {
	n := 0
	for _, sfp := range psg.sfps {
		n += sfp.updateStatsForRow(br, rowIdx)
	}
	return n
}

func (psp *pipeStatsProcessor) writeBlock(workerID uint, br *blockResult) {
	if len(br.timestamps) == 0 {
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

	// Merge states across shards
	shards := psp.shards
	m := shards[0].m
	shards = shards[1:]
	for i := range shards {
		shard := &shards[i]
		for key, psg := range shard.m {
			// shard.m may be quite big, so this loop can take a lot of time and CPU.
			// Stop processing data as soon as stopCh is closed without wasting additional CPU time.
			select {
			case <-psp.stopCh:
				return nil
			default:
			}

			spgBase := m[key]
			if spgBase == nil {
				m[key] = psg
			} else {
				for i, sfp := range spgBase.sfps {
					sfp.mergeState(psg.sfps[i])
				}
			}
		}
	}

	// Write per-group states to ppBase
	byFields := psp.ps.byFields
	if len(byFields) == 0 && len(m) == 0 {
		// Special case - zero matching rows.
		_ = shards[0].getPipeStatsGroup(nil)
		m = shards[0].m
	}

	rcs := make([]resultColumn, 0, len(byFields)+len(psp.ps.resultNames))
	for _, bf := range byFields {
		rcs = append(rcs, resultColumn{
			name: bf.name,
		})
	}
	for _, resultName := range psp.ps.resultNames {
		rcs = append(rcs, resultColumn{
			name: resultName,
		})
	}
	var br blockResult

	var values []string
	valuesLen := 0
	for key, psg := range m {
		// m may be quite big, so this loop can take a lot of time and CPU.
		// Stop processing data as soon as stopCh is closed without wasting additional CPU time.
		select {
		case <-psp.stopCh:
			return nil
		default:
		}

		// Unmarshal values for byFields from key.
		values = values[:0]
		keyBuf := bytesutil.ToUnsafeBytes(key)
		for len(keyBuf) > 0 {
			v, nSize := encoding.UnmarshalBytes(keyBuf)
			if nSize <= 0 {
				logger.Panicf("BUG: cannot unmarshal value from keyBuf=%q", keyBuf)
			}
			keyBuf = keyBuf[nSize:]
			values = append(values, bytesutil.ToUnsafeString(v))
		}
		if len(values) != len(byFields) {
			logger.Panicf("BUG: unexpected number of values decoded from keyBuf; got %d; want %d", len(values), len(byFields))
		}

		// calculate values for stats functions
		for _, sfp := range psg.sfps {
			value := sfp.finalizeStats()
			values = append(values, value)
		}

		if len(values) != len(rcs) {
			logger.Panicf("BUG: len(values)=%d must be equal to len(rcs)=%d", len(values), len(rcs))
		}
		for i, v := range values {
			rcs[i].addValue(v)
			valuesLen += len(v)
		}
		if valuesLen >= 1_000_000 {
			br.setResultColumns(rcs)
			psp.ppBase.writeBlock(0, &br)
			br.reset()
			for i := range rcs {
				rcs[i].resetKeepName()
			}
			valuesLen = 0
		}
	}

	br.setResultColumns(rcs)
	psp.ppBase.writeBlock(0, &br)

	return nil
}

func parsePipeStats(lex *lexer) (*pipeStats, error) {
	if !lex.isKeyword("stats") {
		return nil, fmt.Errorf("expecting 'stats'; got %q", lex.token)
	}

	lex.nextToken()

	var ps pipeStats
	if lex.isKeyword("by") {
		lex.nextToken()
		bfs, err := parseByStatsFields(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse 'by' clause: %w", err)
		}
		ps.byFields = bfs
	}

	var resultNames []string
	var funcs []statsFunc
	for {
		sf, resultName, err := parseStatsFunc(lex)
		if err != nil {
			return nil, err
		}
		resultNames = append(resultNames, resultName)
		funcs = append(funcs, sf)
		if lex.isKeyword("|", ")", "") {
			ps.resultNames = resultNames
			ps.funcs = funcs
			return &ps, nil
		}
		if !lex.isKeyword(",") {
			return nil, fmt.Errorf("unexpected token %q; want ',', '|' or ')'", lex.token)
		}
		lex.nextToken()
	}
}

func parseStatsFunc(lex *lexer) (statsFunc, string, error) {
	var sf statsFunc
	switch {
	case lex.isKeyword("count"):
		scs, err := parseStatsCount(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'count' func: %w", err)
		}
		sf = scs
	case lex.isKeyword("count_empty"):
		scs, err := parseStatsCountEmpty(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'count_empty' func: %w", err)
		}
		sf = scs
	case lex.isKeyword("count_uniq"):
		sus, err := parseStatsCountUniq(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'count_uniq' func: %w", err)
		}
		sf = sus
	case lex.isKeyword("sum"):
		sss, err := parseStatsSum(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'sum' func: %w", err)
		}
		sf = sss
	case lex.isKeyword("max"):
		sms, err := parseStatsMax(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'max' func: %w", err)
		}
		sf = sms
	case lex.isKeyword("min"):
		sms, err := parseStatsMin(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'min' func: %w", err)
		}
		sf = sms
	case lex.isKeyword("avg"):
		sas, err := parseStatsAvg(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'avg' func: %w", err)
		}
		sf = sas
	case lex.isKeyword("uniq_values"):
		sus, err := parseStatsUniqValues(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'uniq_values' func: %w", err)
		}
		sf = sus
	case lex.isKeyword("values"):
		svs, err := parseStatsValues(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'values' func: %w", err)
		}
		sf = svs
	case lex.isKeyword("sum_len"):
		sss, err := parseStatsSumLen(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'sum_len' func: %w", err)
		}
		sf = sss
	case lex.isKeyword("quantile"):
		sqs, err := parseStatsQuantile(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'quantile' func: %w", err)
		}
		sf = sqs
	case lex.isKeyword("median"):
		sms, err := parseStatsMedian(lex)
		if err != nil {
			return nil, "", fmt.Errorf("cannot parse 'median' func: %w", err)
		}
		sf = sms
	default:
		return nil, "", fmt.Errorf("unknown stats func %q", lex.token)
	}

	resultName, err := parseResultName(lex)
	if err != nil {
		return nil, "", fmt.Errorf("cannot parse result name for %s: %w", sf, err)
	}
	return sf, resultName, nil
}

func parseResultName(lex *lexer) (string, error) {
	if lex.isKeyword("as") {
		lex.nextToken()
	}
	resultName, err := parseFieldName(lex)
	if err != nil {
		return "", err
	}
	return resultName, nil
}

var zeroByStatsField = &byStatsField{}

// byStatsField represents 'by (...)' part of the pipeStats.
//
// It can have either 'name' representation or 'name:bucket' or 'name:buket offset off' representation,
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
		fieldName, err := parseFieldName(lex)
		if err != nil {
			return nil, fmt.Errorf("cannot parse field name: %w", err)
		}
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

// parseFieldNamesForStatsFunc parses field names for statsFunc.
//
// It returns ["*"] if the fields names list is empty or if it contains "*" field.
func parseFieldNamesForStatsFunc(lex *lexer, funcName string) ([]string, error) {
	if !lex.isKeyword(funcName) {
		return nil, fmt.Errorf("unexpected func; got %q; want %q", lex.token, funcName)
	}
	lex.nextToken()
	fields, err := parseFieldNamesInParens(lex)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q args: %w", funcName, err)
	}
	if len(fields) == 0 || slices.Contains(fields, "*") {
		fields = []string{"*"}
	}
	return fields, nil
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
	if lex.isKeyword(",", "(", ")", "[", "]", "|", ":", "") {
		return "", fmt.Errorf("unexpected token: %q", lex.token)
	}
	fieldName := getCompoundPhrase(lex, false)
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
