package gotestfmt

import (
	"os"
	"time"
)

type Report struct {
	StartedAt   time.Time
	EndedAt     time.Time
	ElapsedTime time.Duration
	TestsCount  int
	FailCount   int
	PassCount   int
	SkipCount   int
	Tests       []Test
	Coverage    []Coverage
}

type Test struct {
	ReadableName string
	ErrorTrace   string
	Name         string
	Status       string
	Output       []string
	ElapsedTime  time.Duration
	Package      string
	StartedAt    time.Time
	EndedAt      time.Time
	Index        int
}

type Reporter interface {
	Progress(test Test, writer *os.File)
	Summary(report Report, writer *os.File)
	Coverage(list []Coverage, writer *os.File)
	Exit(report Report)
}

type Coverage struct {
	Package  string
	Coverage float64
}
