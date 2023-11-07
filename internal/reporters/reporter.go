package reporters

import c "github.com/fnando/bolt/common"

type Reporter interface {
	Name() string
	OnData(line string)
	OnProgress(test c.Test)
	OnFinished(options ReporterFinishedOptions)
}

type ReporterFinishedOptions struct {
	Aggregation  *c.Aggregation
	HideCoverage bool
	HideSlowest  bool
	Debug        bool
}
