package reporters

import (
	"encoding/json"
	"fmt"

	c "github.com/fnando/bolt/common"
)

type JSONReporter struct {
	Output *c.Output
}

type JSONData struct {
	Coverage   []*c.Coverage
	Tests      []*c.Test
	Benchmarks []*c.Benchmark
	Elapsed    float64
}

func (reporter JSONReporter) Name() string {
	return "json"
}

func (reporter JSONReporter) OnFinished(options ReporterFinishedOptions) {
	data := JSONData{
		Coverage:   options.Aggregation.Coverages(),
		Tests:      options.Aggregation.Tests(),
		Benchmarks: options.Aggregation.Benchmarks(),
		Elapsed:    float64(options.Aggregation.Elapsed()),
	}
	contents, _ := json.MarshalIndent(data, "", "  ")
	fmt.Fprintln(reporter.Output.Stdout, string(contents))
}

func (reporter JSONReporter) OnProgress(test c.Test) {
}

func (reporter JSONReporter) OnData(line string) {
}
