package gotestfmt

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONReporter struct {
}

func (_ JSONReporter) Setup() {
	// noop
}

func (_ JSONReporter) Progress(test Test, writer *os.File) {
	// noop
}

func (_ JSONReporter) Coverage(list []Coverage, writer *os.File) {
	// noop
}

func (_ JSONReporter) Summary(report Report, writer *os.File) {
	output, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintln(writer, string(output))
}

func (_ JSONReporter) Exit(report Report) {
	os.Exit(0)
}
