package gotestfmt

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ProgressReporter struct {
	Env           []string
	Output        *OutputBuffers
	ColorByStatus map[string]string
}

func (reporter ProgressReporter) PrintCoverage(coverages []Coverage) {
	if len(coverages) == 0 {
		return
	}

	fmt.Fprintf(reporter.Output.StdoutWriter, Color.Text("\nCoverage:\n\n"))

	for _, coverage := range coverages {
		line := fmt.Sprintf("[%.1f%%] %s\n", coverage.Coverage, coverage.Package)

		if coverage.Coverage < 50.0 {
			line = Color.Fail(line)
		} else if coverage.Coverage < 70.0 {
			line = Color.Skip(line)
		} else {
			line = Color.Pass(line)
		}

		fmt.Fprint(reporter.Output.StdoutWriter, line)
	}
}

func (reporter ProgressReporter) ExitCode(tests []*Test) int {
	stats := Stats{Tests: tests}

	return stats.FailedTestsCount()
}

func (reporter ProgressReporter) Progress(test *Test) {
	failSymbol := lookupEnvOrDefault(reporter.Env, "GOTESTFMT_FAIL_SYMBOL", "F")
	passSymbol := lookupEnvOrDefault(reporter.Env, "GOTESTFMT_PASS_SYMBOL", ".")
	skipSymbol := lookupEnvOrDefault(reporter.Env, "GOTESTFMT_SKIP_SYMBOL", "S")

	char := Color.Pass(passSymbol)

	if test.Status == TestStatus.Fail {
		char = Color.Fail(failSymbol)
	} else if test.Status == TestStatus.Skip {
		char = Color.Skip(skipSymbol)
	}

	fmt.Fprint(reporter.Output.StdoutWriter, char)
}

type ProgressBenchmarkGroup struct {
	Package           string
	Position          int
	Benchmarks        []*Benchmark
	MaxTitleSize      int
	MaxIterationsSize int
	MaxDurationSize   int
}

func (reporter ProgressReporter) PrintBenchmarks(benchmarks []*Benchmark) {

	if len(benchmarks) == 0 {
		return
	}

	fmt.Fprint(reporter.Output.StdoutWriter, Color.Text("Benchmarks:")+"\n\n")

	groups := map[string]ProgressBenchmarkGroup{}

	for _, benchmark := range benchmarks {
		group, exists := groups[benchmark.Package]

		if !exists {
			group = ProgressBenchmarkGroup{
				Package:    benchmark.Package,
				Benchmarks: []*Benchmark{},
			}
		}

		group.Benchmarks = append(group.Benchmarks, benchmark)
		group.MaxTitleSize = max(group.MaxTitleSize, len(benchmark.Name))
		group.MaxIterationsSize = max(group.MaxIterationsSize, len(fmt.Sprintf("%d", benchmark.IterationCount)))
		group.MaxDurationSize = max(group.MaxDurationSize, len(fmt.Sprintf("%d", len(benchmark.DurationPerOp.String()))))

		groups[benchmark.Package] = group
	}

	for index, group := range maps.Values(groups) {
		position := index + 1
		prefix := fmt.Sprintf("%d) ", position)
		indent := strings.Repeat(" ", len(prefix))

		fmt.Println(Color.Detail(prefix + group.Package))

		for _, benchmark := range group.Benchmarks {
			nameFormat := fmt.Sprintf("%%-%ds", group.MaxTitleSize+10)
			fmt.Fprintf(reporter.Output.StdoutWriter, indent+nameFormat, Color.Text(benchmark.Name))
			fmt.Fprintln(reporter.Output.StdoutWriter)
		}

		fmt.Fprintln(reporter.Output.StdoutWriter)
	}
}

func (reporter ProgressReporter) PrintTests(tests []*Test, statuses []string) {
	status := TestStatus.Pass

	fmt.Fprint(reporter.Output.StdoutWriter, "\n\n")

	position := 0

	for _, test := range tests {
		if !slices.Contains(statuses, test.Status) {
			continue
		}

		var output string

		position += 1

		if status == TestStatus.Pass && test.Status == TestStatus.Skip {
			status = test.Status
		} else if test.Status == TestStatus.Fail {
			status = test.Status
		}

		prefix := fmt.Sprintf("%d) ", position)
		indent := strings.Repeat(" ", len(prefix))

		color := reporter.ColorByStatus[test.Status]
		output += Color.apply(color, prefix+test.ReadableName) + "\n"

		if test.TestSource != "" {
			output += indent + Color.Detail(test.TestSource)
		} else if test.ErrorTrace != "" {
			output += indent + Color.Detail(test.ErrorTrace)
		}

		if test.Status == TestStatus.Skip {
			output += "\n\n" + indent + Color.Text(test.SkipMessage) + "\n"
		} else if len(test.Output) > 0 {
			output += "\n\n"

			lines := formatLines(test.Output)

			for _, line := range lines {
				output += indent + line + "\n"
			}

			if test.TestSource != "" && test.TestSource != test.ErrorTrace {
				output += "\n           " + Color.Fail(test.ErrorTrace) + "\n"
			}
		} else {
			output += "\n" + indent + Color.Text("[No output]") + "\n"
		}

		fmt.Fprint(reporter.Output.StdoutWriter, output+"\n")
	}

}

func (reporter ProgressReporter) PrintSummary(tests []*Test) {
	stats := Stats{Tests: tests}

	failedCount := stats.FailedTestsCount()
	skippedCount := stats.SkippedTestsCount()

	status := TestStatus.Pass

	if failedCount > 0 {
		status = TestStatus.Fail
	} else if skippedCount > 0 {
		status = TestStatus.Skip
	}

	fmt.Fprintf(
		reporter.Output.StdoutWriter,
		Color.apply(reporter.ColorByStatus[status], "Finished in %v, %s, %s, %s\n"),
		stats.Duration(),
		plural(stats.TestsCount(), "test"),
		plural(failedCount, "failure"),
		plural(skippedCount, "skip"),
	)
}

func (reporter ProgressReporter) Finish(tests []*Test, statuses []string, coverages []Coverage, benchmarks []*Benchmark) (int, error) {
	stats := Stats{Tests: tests}

	reporter.PrintTests(tests, statuses)
	reporter.PrintBenchmarks(benchmarks)
	reporter.PrintSummary(tests)
	reporter.PrintCoverage(coverages)

	return stats.CountByStatus(TestStatus.Fail), nil
}

func plural(count int, word string) string {
	if count == 0 {
		return "no " + word + "s"
	} else if count == 1 {
		return "1 " + word
	} else {
		return fmt.Sprintf("%d %s", count, word+"s")
	}
}

func formatLines(lines []string) []string {
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
		line = string(expected.ReplaceAllStringFunc(line, func(input string) string {
			matches := expected.FindStringSubmatch(input)
			return matches[1] + Color.Text(matches[2]) + matches[3] + Color.Pass(matches[4])
		}))

		line = string(actual.ReplaceAllStringFunc(line, func(input string) string {
			matches := actual.FindStringSubmatch(input)
			return matches[1] + Color.Text(matches[2]) + matches[3] + Color.Fail(matches[4])
		}))

		if isDiff {
			line = string(diffLocation.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffLocation.FindStringSubmatch(input)
				return matches[1] + Color.Detail(matches[2])
			}))

			line = string(diffExpected.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffExpected.FindStringSubmatch(input)
				return matches[1] + Color.Pass(matches[2])
			}))

			line = string(diffActual.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffActual.FindStringSubmatch(input)
				return matches[1] + Color.Fail(matches[2])
			}))

			line = string(diffExpectedChange.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffExpectedChange.FindStringSubmatch(input)
				return matches[1] + Color.Pass(matches[2])
			}))

			line = string(diffActualChange.ReplaceAllStringFunc(line, func(input string) string {
				matches := diffActualChange.FindStringSubmatch(input)
				return matches[1] + Color.Fail(matches[2])
			}))
		}

		lines[index] = line
	}

	return lines
}
