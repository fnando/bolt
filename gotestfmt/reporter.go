package gotestfmt

import (
	"math"
	"time"
)

type Reporter interface {
	Progress(test *Test)
	Finish(tests []*Test, statuses []string, coverages []Coverage, benchmarks []*Benchmark) (int, error)
}

type Stats struct {
	Tests []*Test
}

func (stats Stats) TestsCount() int {
	return len(stats.Tests)
}

func (stats Stats) FailedTestsCount() int {
	return stats.CountByStatus(TestStatus.Fail)
}

func (stats Stats) PassedTestsCount() int {
	return stats.CountByStatus(TestStatus.Pass)
}

func (stats Stats) SkippedTestsCount() int {
	return stats.CountByStatus(TestStatus.Skip)
}

func (stats Stats) Duration() time.Duration {
	total := 0.0

	for _, test := range stats.Tests {
		total += test.Elapsed
	}

	return time.Duration(math.Round(total*1000)) * time.Millisecond
}

func (stats Stats) CountByStatus(status string) int {
	count := 0

	for _, test := range stats.Tests {
		if test.Status == status {
			count += 1
		}
	}

	return count
}
