package reporters

import (
	"encoding/json"
	"fmt"
	"strings"

	c "github.com/fnando/bolt/common"
)

type StandardReporter struct {
	Output *c.Output
}

func (reporter StandardReporter) Name() string {
	return "standard"
}

func (reporter StandardReporter) OnFinished(options ReporterFinishedOptions) {

}

func (reporter StandardReporter) OnProgress(test c.Test) {

}

func (reporter StandardReporter) OnData(line string) {
	var data c.Stream
	err := json.Unmarshal([]byte(line), &data)

	if err != nil {
		fmt.Fprint(reporter.Output.Stdout, line)
	} else {
		fmt.Fprint(reporter.Output.Stdout, reporter.formatLine(data.Output))
	}

}

func (reporter StandardReporter) formatLine(line string) string {
	line = strings.Replace(line, "=== RUN", c.Color.Detail("=== RUN"), 1)
	line = strings.Replace(line, "--- PASS:", c.Color.Pass("--- PASS:"), 1)
	line = strings.Replace(line, "--- SKIP:", c.Color.Skip("--- SKIP:"), 1)
	line = strings.Replace(line, "--- FAIL:", c.Color.Fail("--- FAIL:"), 1)

	return line
}
