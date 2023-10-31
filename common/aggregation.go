package common

import (
	"cmp"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Aggregation struct {
	BenchmarksMap     map[string]*Benchmark
	CoverageCount     int
	CoverageMap       map[string]*Coverage
	CoverageThreshold float64
	OrphanOutput      []string
	SlowestCount      int
	SlowestThreshold  time.Duration
	TestsMap          map[string]*Test

	StartedAt time.Time
	EndedAt   time.Time
}

type Coverage struct {
	Package  string
	Coverage float64
}

func (agg Aggregation) Elapsed() time.Duration {
	return agg.EndedAt.Sub(agg.StartedAt)
}

func (agg Aggregation) TestsCount() int {
	return len(agg.Tests())
}

func (agg Aggregation) Benchmarks() []*Benchmark {
	benchmarks := maps.Values(agg.BenchmarksMap)

	slices.SortFunc(benchmarks, func(a, b *Benchmark) int {
		return cmp.Compare(a.Key, b.Key)
	})

	return benchmarks
}

func (agg Aggregation) Tests() []*Test {
	tests := maps.Values(agg.TestsMap)

	slices.SortFunc(tests, func(a, b *Test) int {
		if a.Package != b.Package {
			return cmp.Compare(a.Package, b.Package)
		}

		return cmp.Compare(a.Name, b.Name)
	})

	return tests
}

func (agg Aggregation) SlowestTests() []*Test {
	tests := []*Test{}

	for _, test := range agg.Tests() {
		if int64(test.Elapsed) > int64(agg.SlowestThreshold) {
			tests = append(tests, test)
		}
	}

	slices.SortFunc(tests, func(a, b *Test) int {
		return cmp.Compare(int64(a.Elapsed), int64(b.Elapsed))
	})

	slices.Reverse(tests)

	tests = tests[:min(len(tests), agg.SlowestCount)]

	return tests
}

func (agg Aggregation) Coverages() []*Coverage {
	coverages := []*Coverage{}

	for _, coverage := range maps.Values(agg.CoverageMap) {
		if coverage.Coverage < agg.CoverageThreshold {
			coverages = append(coverages, coverage)
		}
	}

	slices.SortFunc(coverages, func(a, b *Coverage) int {
		comparison := cmp.Compare(a.Coverage, b.Coverage)

		if comparison != 0 {
			return comparison
		}

		return cmp.Compare(a.Package, b.Package)
	})

	coverages = coverages[:min(len(coverages), agg.CoverageCount)]

	return coverages
}

func (agg Aggregation) CountBy(status string) int {
	count := 0

	for _, test := range agg.Tests() {
		if status == test.Status {
			count += 1
		}
	}

	return count
}

func (agg Aggregation) Status() string {
	if agg.CountBy("fail") > 0 {
		return "fail"
	} else if agg.CountBy("skip") > 0 {
		return "skip"
	}

	return "pass"
}
