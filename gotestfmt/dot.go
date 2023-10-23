package gotestfmt

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type DotReporter struct {
	FailColor         string
	PassColor         string
	SkipColor         string
	TextColor         string
	DetailColor       string
	ExpectedColor     string
	ActualColor       string
	CoverageBadColor  string
	CoverageGoodColor string
	CoverageOKColor   string
}

func (dot DotReporter) deindent(lines []string, prefixSize int) []string {
	if len(lines) == 0 {
		return lines
	}

	re := regexp.MustCompile("(?m)^(\\s+)Error:(\\s+)(.+)")
	matches := re.FindStringSubmatch(strings.Join(lines, ""))
	indent := ""

	if len(matches) > 0 {
		indent = matches[1] + matches[2] + "      "
		lines[0] = indent + matches[3] + "\n"
	}

	indentSize := len(indent)

	for index, line := range lines {
		content := ""

		if index == 0 {
			content = line[indentSize:]
		} else {
			content = line[indentSize:]
		}

		lines[index] = strings.Repeat(" ", prefixSize) + content
	}

	return lines
}

func CreateDotReporter() Reporter {
	reporter := DotReporter{
		FailColor:         "\u001b[31m",
		PassColor:         "\u001b[32m",
		SkipColor:         "\u001b[33m",
		TextColor:         "\u001b[30m",
		DetailColor:       "\u001b[34m",
		ExpectedColor:     "\u001b[32m",
		ActualColor:       "\u001b[31m",
		CoverageBadColor:  "\u001b[31m",
		CoverageGoodColor: "\u001b[32m",
		CoverageOKColor:   "\u001b[33m",
	}

	var (
		color  string
		exists bool
	)

	if color, exists = os.LookupEnv("GOTESTFMT_FAIL_COLOR"); exists {
		reporter.FailColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_PASS_COLOR"); exists {
		reporter.PassColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_SKIP_COLOR"); exists {
		reporter.SkipColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_TEXT_COLOR"); exists {
		reporter.TextColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_DETAIL_COLOR"); exists {
		reporter.DetailColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_COVERAGE_GOOD_COLOR"); exists {
		reporter.CoverageGoodColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_COVERAGE_BAD_COLOR"); exists {
		reporter.CoverageBadColor = "\u001b[" + color + "m"
	}

	if color, exists = os.LookupEnv("GOTESTFMT_COVERAGE_OK_COLOR"); exists {
		reporter.CoverageOKColor = "\u001b[" + color + "m"
	}

	return reporter
}

func (dot DotReporter) Progress(test Test, writer *os.File) {
	char := Color(dot.PassColor, ".")

	if test.Status == "skip" {
		char = Color(dot.SkipColor, "S")
	}

	if test.Status == "fail" {
		char = Color(dot.FailColor, "F")
	}

	fmt.Fprint(writer, char)
}

func (dot DotReporter) Summary(report Report, writer *os.File) {
	position := 0

	fmt.Fprint(writer, "\n\n")

	for _, test := range report.Tests {
		if test.Status != "fail" {
			continue
		}

		position += 1
		prefix := fmt.Sprintf("%d) ", position)
		prefixSize := len(prefix)
		indent := strings.Repeat(" ", prefixSize)

		fmt.Fprint(writer, Color(dot.FailColor, fmt.Sprintf("%s%s (%s)\n", prefix, test.ReadableName, test.Name)))

		if test.ErrorTrace != "" {
			fmt.Fprintf(writer, Color(dot.DetailColor, fmt.Sprintf("%s%s\n", indent, test.ErrorTrace)))
		}

		lines := dot.deindent(test.Output, prefixSize)
		content := strings.Join(lines, "")
		content = dot.formatOutput(content)

		fmt.Fprintf(writer, "\n%s\n\n", Color(dot.TextColor, content))
	}

	testsWord := "tests"

	if report.TestsCount == 1 {
		testsWord = "test"
	}

	summaryColor := dot.PassColor

	if report.FailCount > 0 {
		summaryColor = dot.FailColor
	} else if report.SkipCount > 0 {
		summaryColor = dot.SkipColor
	}

	fmt.Fprintf(
		writer,
		Color(
			summaryColor,
			fmt.Sprintf(
				"Finished in %v, %d %s, %d failed, %d skipped\n",
				report.ElapsedTime,
				report.TestsCount,
				testsWord,
				report.FailCount,
				report.SkipCount,
			),
		),
	)
}

func (_ DotReporter) Exit(report Report) {
	os.Exit(report.FailCount)
}

func (dot DotReporter) formatOutput(text string) string {
	var re *regexp.Regexp

	text = strings.Replace(text, "--- Expected", Color(dot.ExpectedColor, "--- Expected"), -1)
	text = strings.Replace(text, "+++ Actual", Color(dot.ActualColor, "+++ Actual"), -1)

	re = regexp.MustCompile("(@@.+@@)")
	text = string(re.ReplaceAllFunc([]byte(text), func(repl []byte) []byte {
		return []byte(Color(dot.DetailColor, string(repl)))
	}))

	re = regexp.MustCompile("(?m)^\\s+-(.+)")
	text = string(re.ReplaceAllFunc([]byte(text), func(repl []byte) []byte {
		return []byte(Color(dot.ExpectedColor, string(repl)))
	}))

	re = regexp.MustCompile("(?m)^\\s+\\+(.+)")
	text = string(re.ReplaceAllFunc([]byte(text), func(repl []byte) []byte {
		return []byte(Color(dot.ActualColor, string(repl)))
	}))

	re = regexp.MustCompile("(?m)^(\\s+expected\\s*:)(.+)")
	text = string(re.ReplaceAllFunc([]byte(text), func(repl []byte) []byte {
		matches := re.FindStringSubmatch(string(repl))

		return []byte(Color(dot.TextColor, matches[1]) + Color(dot.ExpectedColor, matches[2]))
	}))

	re = regexp.MustCompile("(?m)^(\\s+actual\\s*:)(.+)")
	text = string(re.ReplaceAllFunc([]byte(text), func(repl []byte) []byte {
		matches := re.FindStringSubmatch(string(repl))

		return []byte(Color(dot.TextColor, matches[1]) + Color(dot.ActualColor, matches[2]))
	}))

	return text
}

func (dot DotReporter) Coverage(list []Coverage, writer *os.File) {
	count := len(list)

	if count == 0 {
		return
	}

	fmt.Fprintf(
		writer,
		fmt.Sprintf("\n%s\n\n", Color(dot.TextColor, "Coverage:")),
	)

	for _, coverage := range list {
		color := dot.CoverageGoodColor

		if coverage.Coverage < 70.0 {
			color = dot.CoverageOKColor
		}

		if coverage.Coverage < 50.0 {
			color = dot.CoverageBadColor
		}

		cov := fmt.Sprintf("%.1f", coverage.Coverage)
		line := fmt.Sprintf("[%+4s%s] %s\n", cov, "%%", coverage.Package)

		fmt.Fprintf(writer, Color(color, line))
	}

	fmt.Println()
}
