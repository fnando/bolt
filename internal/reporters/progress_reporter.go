package reporters

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	h "github.com/dustin/go-humanize"
	c "github.com/fnando/bolt/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type ProgressReporter struct {
	Output *c.Output
}

func (reporter ProgressReporter) Name() string {
	return "progress"
}

func (reporter ProgressReporter) OnFinished(options ReporterFinishedOptions) {
	reporter.PrintTests(options.Aggregation)
	reporter.PrintBenchmarks(options.Aggregation)
	reporter.PrintSummary(options.Aggregation)

	if options.Aggregation.CountBy("failed") == 0 {
		if !options.HideCoverage {
			reporter.PrintCoverage(options.Aggregation)
		}

		if !options.HideSlowest {
			reporter.PrintSlowestTests(options.Aggregation)
		}
	}
}

func (reporter ProgressReporter) PrintSummary(aggregation *c.Aggregation) {
	testsCount := aggregation.TestsCount()
	failCount := aggregation.CountBy("fail")
	skipCount := aggregation.CountBy("skip")
	benchmarksCount := len(aggregation.Benchmarks())

	summary := fmt.Sprintf(
		"\nFinished in %s, %d tests, %d failures, %d skips, %d benchmarks\n",
		formatDuration(aggregation.Elapsed(), 0),
		testsCount,
		failCount,
		skipCount,
		benchmarksCount,
	)

	fmt.Fprintf(
		reporter.Output.Stdout,
		c.Color.Apply(c.Color.Color(aggregation.Status()), summary),
	)
}

func (reporter ProgressReporter) PrintBenchmarks(aggregation *c.Aggregation) {
	benchmarks := aggregation.Benchmarks()

	if len(benchmarks) == 0 {
		return
	}

	fmt.Fprint(reporter.Output.Stdout, "\n"+c.Color.Text("Benchmarks:")+"\n\n")

	t := table.NewWriter()
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Header = text.FormatTitle
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1},
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignRight},
	})
	t.SetOutputMirror(reporter.Output.Stdout)
	t.AppendHeader(table.Row{"Name", "Iterations", "Time/op"})

	for _, benchmark := range benchmarks {
		t.AppendRow([]interface{}{
			fmt.Sprintf("%s-%d", benchmark.Name, benchmark.Processors),
			h.Comma(int64(benchmark.Iterations)),
			formatDuration(benchmark.DurationPerOperation, 2),
		})
	}

	t.Render()
}

func (reporter ProgressReporter) PrintTests(aggregation *c.Aggregation) {
	fmt.Fprintln(reporter.Output.Stdout)

	position := 0

	if aggregation.TestsCount() == 0 {
		return
	}

	for _, test := range aggregation.Tests() {
		if test.Status == "pass" {
			continue
		}

		position += 1
		output := "\n"
		prefix := fmt.Sprintf("%d) ", position)
		indent := strings.Repeat(" ", len(prefix))
		output += c.Color.Apply(c.Color.Color(test.Status), prefix+test.ReadableName) + "\n"

		if test.ErrorTrace != "" {
			output += indent + c.Color.Detail(test.ErrorTrace) + "\n\n"
		} else {
			output += "\n"
		}

		lines := reporter.formatLines(reporter.deindentOutput(test.Output))

		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)

			ignore := strings.HasPrefix(trimmedLine, "=== RUN") ||
				strings.HasPrefix(trimmedLine, "--- FAIL:") ||
				strings.HasPrefix(trimmedLine, "Error Trace:") ||
				trimmedLine == test.ErrorTrace+":" ||
				trimmedLine == test.Source+":" ||
				strings.HasPrefix(trimmedLine, "Test:") ||
				strings.Contains(line, test.Name)

			if ignore {
				continue
			}

			if trimmedLine != "" {
				output += indent + c.Color.Text(line) + "\n"
			} else {
				output += "\n"
			}
		}

		if test.Source != "" {
			output += "\n" + indent + "        " + c.Color.Fail(test.Source) + "\n"
		}

		fmt.Fprint(reporter.Output.Stdout, output)
	}
}

func (reporter ProgressReporter) PrintCoverage(aggregation *c.Aggregation) {
	coverages := aggregation.Coverages()

	if len(coverages) == 0 {
		return
	}

	fmt.Fprint(reporter.Output.Stdout, "\n"+c.Color.Text("Coverage:")+"\n\n")

	for _, coverage := range coverages {
		line := fmt.Sprintf("[%.1f%%] %s", coverage.Coverage, coverage.Package)

		if coverage.Coverage < 50.0 {
			line = c.Color.Fail(line)
		} else if coverage.Coverage < 70.0 {
			line = c.Color.Skip(line)
		} else {
			line = c.Color.Pass(line)
		}

		fmt.Fprint(reporter.Output.Stdout, line+"\n")
	}
}

func (reporter ProgressReporter) PrintSlowestTests(aggregation *c.Aggregation) {
	tests := aggregation.SlowestTests()

	if len(tests) == 0 {
		return
	}
	durationSize := 0
	var totalTime int64 = 0

	for _, test := range tests {
		totalTime += test.Elapsed.Nanoseconds()
		durationSize = max(durationSize, len(formatDuration(test.Elapsed, 2)))
	}

	fmt.Fprintf(
		reporter.Output.Stdout,
		"\nTop %d %s (%s, %.2f%% of total time):\n\n",
		aggregation.SlowestCount,
		"slowest tests",
		formatDuration(time.Duration(totalTime), 0),
		min(100.0, float64(totalTime)/float64(aggregation.Elapsed())*100),
	)

	for _, test := range tests {
		elapsed := formatDuration(test.Elapsed, 2)
		padding := strings.Repeat(" ", max(0, durationSize-len(elapsed)))

		fmt.Fprintf(
			reporter.Output.Stdout,
			"%s %s\n%s\n",
			c.Color.Detail(padding+elapsed),
			c.Color.Fail(test.Name),
			strings.Repeat(" ", durationSize)+" "+c.Color.Text(test.Package),
		)
	}
}

func (reporter ProgressReporter) OnProgress(test c.Test) {
	env := func(name, defaultVal string) string {
		val := os.Getenv(name)

		if val != "" {
			return val
		}

		return defaultVal
	}

	symbols := map[string]string{
		"fail": env("BOLT_FAIL_SYMBOL", "F"),
		"pass": env("BOLT_PASS_SYMBOL", "."),
		"skip": env("BOLT_SKIP_SYMBOL", "S"),
	}

	fmt.Fprint(
		reporter.Output.Stdout,
		c.Color.Apply(c.Color.Color(test.Status), symbols[test.Status]),
	)
}

func (reporter ProgressReporter) OnData(line string) {

}

func (reporter ProgressReporter) formatLines(lines []string) []string {
	expected := regexp.MustCompile(`^(\s*)(expected:)(\s*)(.*?)$`)
	actual := regexp.MustCompile(`^(\s*)(actual\s*:)(\s*)(.*?)$`)
	diffLocation := regexp.MustCompile(`^(\s*)(@@.*?@@)$`)
	diffExpected := regexp.MustCompile(`^(\s*)(--- Expected)$`)
	diffActual := regexp.MustCompile(`^(\s*)(\+\+\+ Actual)$`)
	diffExpectedChange := regexp.MustCompile(`^(\s*)(\-.*?)$`)
	diffActualChange := regexp.MustCompile(`^(\s*)(\+.*?)$`)
	diff := regexp.MustCompile(`(?m)^\s+Diff:`)

	isDiff := diff.MatchString(strings.Join(lines, "\n"))

	for index, line := range lines {
		line = strings.TrimRight(line, " \t\r\n")

		line = string(expected.ReplaceAllStringFunc(line, func(input string) string {
			matches := expected.FindStringSubmatch(input)
			return matches[1] + c.Color.Text(matches[2]) + matches[3] + c.Color.Pass(matches[4])
		}))

		line = string(actual.ReplaceAllStringFunc(line, func(input string) string {
			matches := actual.FindStringSubmatch(input)
			return matches[1] + c.Color.Text(matches[2]) + matches[3] + c.Color.Fail(matches[4])
		}))

		if isDiff {
			line = string(diffLocation.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffLocation.FindStringSubmatch(input)
				return matches[1] + c.Color.Detail(matches[2])
			}))

			line = string(diffExpected.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffExpected.FindStringSubmatch(input)
				return matches[1] + c.Color.Pass(matches[2])
			}))

			line = string(diffActual.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffActual.FindStringSubmatch(input)
				return matches[1] + c.Color.Fail(matches[2])
			}))

			line = string(diffExpectedChange.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffExpectedChange.FindStringSubmatch(input)
				return matches[1] + c.Color.Pass(matches[2])
			}))

			line = string(diffActualChange.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffActualChange.FindStringSubmatch(input)
				return matches[1] + c.Color.Fail(matches[2])
			}))
		}

		lines[index] = line
	}

	return lines
}

func (reporter ProgressReporter) deindentOutput(output []string) []string {
	lines := []string{}
	indent := ""
	re := regexp.MustCompile(`^(\s+)Error:(\s+)(.+)$`)
	errorIndex := -1
	errorLabel := "Error:  "
	errorSpacing := strings.Repeat(" ", len(errorLabel))

	for index, line := range output {
		matches := re.FindStringSubmatch(line)

		if matches != nil {
			indent = matches[1] + errorSpacing + matches[2]
			output[index] = errorLabel + matches[3]
			errorIndex = index

			break
		}
	}

	if indent == "" {
		return output
	}

	for index, line := range output {
		if errorIndex != index {

			trimmed := strings.TrimLeft(line, "\r\n\t ")

			if trimmed == "" {
				line = ""
			} else {
				line = errorSpacing + trimmed
			}
		}

		lines = append(lines, line)
	}

	return lines
}

func formatDuration(duration time.Duration, places int) string {
	result := duration.String()
	re := regexp.MustCompile(`(?:(\d+(?:\.\d+)?)([^\d]+))`)
	matches := re.FindAllStringSubmatch(result, -1)

	if matches == nil {
		return result
	}

	result = ""

	for _, pair := range matches {
		if strings.Contains(pair[0], ".") {
			numeric, _ := strconv.ParseFloat(pair[1], 64)
			format := "%." + strconv.Itoa(places) + "f"
			result += fmt.Sprintf(format, numeric) + pair[2]
		} else {
			result += pair[0]
		}
	}

	return result
}
