package storage

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"sort"
	"sync"
	"testing"
	"testing/quick"
	"time"
)

func TestSearchQueryMarshalUnmarshal(t *testing.T) {
	rnd := rand.New(rand.NewSource(0))
	typ := reflect.TypeOf(&SearchQuery{})
	var buf []byte
	var sq2 SearchQuery

	for i := 0; i < 1000; i++ {
		v, ok := quick.Value(typ, rnd)
		if !ok {
			t.Fatalf("cannot create random SearchQuery via testing/quick.Value")
		}
		sq1 := v.Interface().(*SearchQuery)
		if sq1 == nil {
			// Skip nil sq1.
			continue
		}
		buf = sq1.Marshal(buf[:0])

		tail, err := sq2.Unmarshal(buf)
		if err != nil {
			t.Fatalf("cannot unmarshal SearchQuery: %s", err)
		}
		if len(tail) > 0 {
			t.Fatalf("unexpected tail left after SearchQuery unmarshaling; tail (len=%d): %q", len(tail), tail)
		}
		if sq1.MinTimestamp != sq2.MinTimestamp {
			t.Fatalf("unexpected MinTimestamp; got %d; want %d", sq2.MinTimestamp, sq1.MinTimestamp)
		}
		if sq1.MaxTimestamp != sq2.MaxTimestamp {
			t.Fatalf("unexpected MaxTimestamp; got %d; want %d", sq2.MaxTimestamp, sq1.MaxTimestamp)
		}
		if len(sq1.TagFilterss) != len(sq2.TagFilterss) {
			t.Fatalf("unexpected TagFilterss len; got %d; want %d", len(sq2.TagFilterss), len(sq1.TagFilterss))
		}
		for ii := range sq1.TagFilterss {
			tagFilters1 := sq1.TagFilterss[ii]
			tagFilters2 := sq2.TagFilterss[ii]
			for j := range tagFilters1 {
				tf1 := &tagFilters1[j]
				tf2 := &tagFilters2[j]
				if string(tf1.Key) != string(tf2.Key) {
					t.Fatalf("unexpected Key on iteration %d,%d; got %X; want %X", i, j, tf2.Key, tf1.Key)
				}
				if string(tf1.Value) != string(tf2.Value) {
					t.Fatalf("unexpected Value on iteration %d,%d; got %X; want %X", i, j, tf2.Value, tf1.Value)
				}
				if tf1.IsNegative != tf2.IsNegative {
					t.Fatalf("unexpected IsNegative on iteration %d,%d; got %v; want %v", i, j, tf2.IsNegative, tf1.IsNegative)
				}
				if tf1.IsRegexp != tf2.IsRegexp {
					t.Fatalf("unexpected IsRegexp on iteration %d,%d; got %v; want %v", i, j, tf2.IsRegexp, tf1.IsRegexp)
				}
			}
		}
	}
}

func TestSearch(t *testing.T) {
	path := "TestSearch"
	st := MustOpenStorage(path, OpenOptions{})
	defer func() {
		st.MustClose()
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("cannot remove storage %q: %s", path, err)
		}
	}()

	// Add rows to storage.
	const rowsCount = 2e4
	const rowsPerBlock = 1e3
	const metricGroupsCount = rowsCount / 5

	mrs := make([]MetricRow, rowsCount)
	var mn MetricName
	mn.Tags = []Tag{
		{[]byte("job"), []byte("super-service")},
		{[]byte("instance"), []byte("8.8.8.8:1234")},
	}
	startTimestamp := timestampFromTime(time.Now())
	startTimestamp -= startTimestamp % (1e3 * 60 * 30)
	blockRowsCount := 0
	for i := 0; i < rowsCount; i++ {
		mn.MetricGroup = []byte(fmt.Sprintf("metric_%d", i%metricGroupsCount))

		mr := &mrs[i]
		mr.MetricNameRaw = mn.marshalRaw(nil)
		mr.Timestamp = startTimestamp + int64(i)
		mr.Value = float64(i)

		blockRowsCount++
		if blockRowsCount == rowsPerBlock {
			st.AddRows(mrs[i-blockRowsCount+1:i+1], defaultPrecisionBits)
			blockRowsCount = 0
		}
	}
	st.AddRows(mrs[rowsCount-blockRowsCount:], defaultPrecisionBits)
	endTimestamp := mrs[len(mrs)-1].Timestamp

	// Re-open the storage in order to flush all the pending cached data.
	st.MustClose()
	st = MustOpenStorage(path, OpenOptions{})

	// Run search.
	tr := TimeRange{
		MinTimestamp: startTimestamp + int64(rowsCount)/3,
		MaxTimestamp: endTimestamp - int64(rowsCount)/3,
	}

	t.Run("serial", func(t *testing.T) {
		if err := testSearchInternal(st, tr, mrs); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		ch := make(chan error, 3)
		for i := 0; i < cap(ch); i++ {
			go func() {
				ch <- testSearchInternal(st, tr, mrs)
			}()
		}
		var firstError error
		for i := 0; i < cap(ch); i++ {
			select {
			case err := <-ch:
				if err != nil && firstError == nil {
					firstError = err
				}
			case <-time.After(10 * time.Second):
				t.Fatalf("timeout")
			}
		}
		if firstError != nil {
			t.Fatalf("unexpected error: %s", firstError)
		}
	})
}

// The test ensures that dropping partitions outside the retention period
// concurrently with search requests happens gracefully. I.e. there are no
// panics that cause the process to exit.
func TestSearch_PartitionsDroppedConcurrently(t *testing.T) {
	defer testRemoveAll(t)

	const (
		days = 1000
	)
	mrs := make([]MetricRow, days)
	name := fmt.Sprintf("metric")
	mn := &MetricName{
		MetricGroup: []byte(name),
	}
	metricNameRaw := mn.marshalRaw(nil)
	for day := 1; day <= days; day++ {
		mr := &mrs[day-1]
		mr.MetricNameRaw = metricNameRaw
		mr.Timestamp = time.Now().Add(-time.Duration(day) * 24 * time.Hour).UnixMilli()
		mr.Value = float64(rand.Intn(1000))
	}

	s := MustOpenStorage(t.Name(), OpenOptions{
		Retention: days * 24 * time.Hour,
	})
	s.AddRows(mrs, defaultPrecisionBits)
	s.DebugFlush()
	defer s.MustClose()

	tfs := NewTagFilters()
	if err := tfs.Add(nil, []byte("metric"), false, false); err != nil {
		t.Fatalf("unexpected error in TagFilters.Add: %v", err)
	}
	tr := TimeRange{
		MinTimestamp: time.Now().Add(-days * 24 * time.Hour).UnixMilli(),
		MaxTimestamp: time.Now().UnixMilli(),
	}

	// Start several readers that query and read the data over entire time
	// range.
	var wg sync.WaitGroup
	stop := make(chan bool)
	for range 100 {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-stop:
					wg.Done()
					return
				default:
				}

				var search Search
				search.Init(nil, s, []*TagFilters{tfs}, tr, 1e5, noDeadline)
				for search.NextMetricBlock() {
					var b Block
					search.MetricBlockRef.BlockRef.MustReadBlock(&b)
				}
				if err := search.Error(); err != nil {
					panic(fmt.Sprintf("search error: %v", err))
				}
				search.MustClose()
			}
		}()
	}

	// Drop partitions one by one concurrently with readers.
	for month := range int(days / 31) {
		d := time.Duration(month) * 30 * 24 * time.Hour
		currentMillis := time.Now().Add(d).UnixMilli()
		s.tb.dropPartitionsOutsideRetentionPeriod(currentMillis)
	}

	close(stop)
	wg.Wait()
}

func testSearchInternal(s *Storage, tr TimeRange, mrs []MetricRow) error {
	for i := 0; i < 10; i++ {
		// Prepare TagFilters for search.
		tfs := NewTagFilters()
		metricGroupRe := fmt.Sprintf(`metric_\d*%d%d`, i, i)
		if err := tfs.Add(nil, []byte(metricGroupRe), false, true); err != nil {
			return fmt.Errorf("cannot add metricGroupRe=%q: %w", metricGroupRe, err)
		}
		if err := tfs.Add([]byte("job"), []byte("nonexisting-service"), true, false); err != nil {
			return fmt.Errorf("cannot add tag filter %q=%q: %w", "job", "nonexsitsing-service", err)
		}
		if err := tfs.Add([]byte("instance"), []byte(".*"), false, true); err != nil {
			return fmt.Errorf("cannot add tag filter %q=%q: %w", "instance", ".*", err)
		}

		// Build extectedMrs.
		var expectedMrs []MetricRow
		metricGroupRegexp := regexp.MustCompile(fmt.Sprintf("^%s$", metricGroupRe))
		var mn MetricName
		for j := range mrs {
			mr := &mrs[j]
			if mr.Timestamp < tr.MinTimestamp || mr.Timestamp > tr.MaxTimestamp {
				continue
			}
			if err := mn.UnmarshalRaw(mr.MetricNameRaw); err != nil {
				return fmt.Errorf("cannot unmarshal MetricName: %w", err)
			}
			if !metricGroupRegexp.Match(mn.MetricGroup) {
				continue
			}
			expectedMrs = append(expectedMrs, *mr)
		}

		if err := testAssertSearchResult(s, tr, tfs, expectedMrs); err != nil {
			return err
		}
	}
	return nil
}

func testAssertSearchResult(st *Storage, tr TimeRange, tfs *TagFilters, want []MetricRow) error {
	type metricBlock struct {
		MetricName []byte
		Block      *Block
	}

	var s Search
	s.Init(nil, st, []*TagFilters{tfs}, tr, 1e5, noDeadline)
	var mbs []metricBlock
	for s.NextMetricBlock() {
		var b Block
		s.MetricBlockRef.BlockRef.MustReadBlock(&b)

		var mb metricBlock
		mb.MetricName = append(mb.MetricName, s.MetricBlockRef.MetricName...)
		mb.Block = &b
		mbs = append(mbs, mb)
	}
	if err := s.Error(); err != nil {
		return fmt.Errorf("search error: %w", err)
	}
	s.MustClose()

	var got []MetricRow
	var mn MetricName
	for _, mb := range mbs {
		rb := newTestRawBlock(mb.Block, tr)
		if err := mn.Unmarshal(mb.MetricName); err != nil {
			return fmt.Errorf("cannot unmarshal MetricName: %w", err)
		}
		metricNameRaw := mn.marshalRaw(nil)
		for i, timestamp := range rb.Timestamps {
			mr := MetricRow{
				MetricNameRaw: metricNameRaw,
				Timestamp:     timestamp,
				Value:         rb.Values[i],
			}
			got = append(got, mr)
		}
	}

	testSortMetricRows(got)
	testSortMetricRows(want)
	if !reflect.DeepEqual(got, want) {
		return fmt.Errorf("unexpected rows found;\ngot\n%s\nwant\n%s", mrsToString(got), mrsToString(want))
	}

	return nil
}

func testSortMetricRows(mrs []MetricRow) {
	sort.Slice(mrs, func(i, j int) bool {
		a, b := &mrs[i], &mrs[j]
		cmp := bytes.Compare(a.MetricNameRaw, b.MetricNameRaw)
		if cmp < 0 {
			return true
		}
		if cmp > 0 {
			return false
		}
		return a.Timestamp < b.Timestamp
	})
}

func mrsToString(mrs []MetricRow) string {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "len=%d\n", len(mrs))
	for i := range mrs {
		mr := &mrs[i]
		fmt.Fprintf(&bb, "[%q, Timestamp=%d, Value=%f]\n", mr.MetricNameRaw, mr.Timestamp, mr.Value)
	}
	return bb.String()
}
