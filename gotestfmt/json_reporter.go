package gotestfmt

import (
	"encoding/json"
	"fmt"
)

type JSONReporter struct {
	Output *OutputBuffers
}

func (reporter JSONReporter) Progress(test *Test) {}

func (reporter JSONReporter) Finish(tests []*Test, statuses []string, coverages []Coverage, benchmarks []*Benchmark) (int, error) {
	data := map[string]any{"Tests": tests, "Coverage": coverages, "Benchmarks": benchmarks}

	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return 1, err
	}

	fmt.Fprintln(reporter.Output.StdoutWriter, string(output))
	return 0, nil
}
