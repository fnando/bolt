package gotestfmt

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type StreamConsumer struct {
	WorkingDir         string
	HomeDir            string
	Tests              map[string]*Test
	Benchmarks         map[string]*Benchmark
	Coverage           []Coverage
	Scanner            *bufio.Scanner
	OnNotifyTestFinish func(test *Test)
	OnFinish           func(test []*Test, coverage []Coverage, benchmarks []*Benchmark)
	OnError            func()
	Output             *OutputBuffers
}

type CreateStreamConsumerOptions struct {
	HomeDir    string
	WorkingDir string
	Scanner    *bufio.Scanner
	Output     *OutputBuffers
}

type testStatus struct {
	Pass string
	Fail string
	Skip string
}

var TestStatus testStatus = testStatus{
	Pass: "pass",
	Fail: "fail",
	Skip: "skip",
}

type Test struct {
	Key          string
	Name         string
	Package      string
	StartedAt    time.Time
	EndedAt      time.Time
	Elapsed      float64
	ReadableName string
	Status       string
	Output       []string
	ErrorTrace   string
	TestSource   string
	SkipMessage  string
}

type Benchmark struct {
	Key            string
	Name           string
	Package        string
	StartedAt      time.Time
	EndedAt        time.Time
	Output         []string
	Elapsed        float64
	DurationPerOp  time.Duration
	IterationCount int
	MaxProcessors  int
}

type Coverage struct {
	Package  string
	Coverage float64
}

type Data struct {
	Action  string
	Elapsed float64
	Output  string
	Package string
	Test    string
	Time    string
}

func (s *StreamConsumer) Run() {
	var err error

	for s.Scanner.Scan() {
		err = s.ingest(s.Scanner.Bytes())

		if err != nil {
			break
		}
	}

	tests := maps.Values(s.Tests)
	slices.SortFunc(tests, func(a, b *Test) int {
		if n := cmp.Compare(a.StartedAt.Unix(), b.StartedAt.Unix()); n != 0 {
			return n
		}

		return cmp.Compare(a.Name, b.Name)
	})

	if err != nil {
		s.OnError()
	} else {
		s.OnFinish(tests, s.Coverage, maps.Values(s.Benchmarks))
	}
}

func (s *StreamConsumer) ingest(lineBytes []byte) error {
	var data Data
	err := json.Unmarshal(lineBytes, &data)

	if err != nil {
		// These are usually build errors, package installing lines and similar.
		// We can just output these to stderr, because they're likely informing of an
		// issue.
		line := string(lineBytes)
		hasFailed := strings.Contains(line, "[build failed]")

		fmt.Fprintln(s.Output.StdoutWriter, line)

		if hasFailed {
			return errors.New("Build has failed")
		}
	}

	// JSON lines were spitted out by `go test`, so we can safely process it.

	isBenchmark := strings.HasPrefix(data.Test, "Benchmark")

	switch data.Action {
	case "start":
		s.processStart(&data)
	case "run":
		if isBenchmark {
			s.processBenchmark(&data)
		} else {
			s.processTest(&data)
		}

	case "output":
		if isBenchmark {
			s.processBenchmarkOutput(&data)
		} else if strings.HasPrefix(data.Output, "coverage:") {
			s.processCoverage(&data)
		} else {
			s.processTestOutput(&data)
		}

	case TestStatus.Skip:
		fallthrough
	case TestStatus.Fail:
		fallthrough
	case TestStatus.Pass:
		if !isBenchmark {
			s.processTestFinished(&data)
		}
	}

	return nil
}

func (s *StreamConsumer) processStart(data *Data) {

}

func (s *StreamConsumer) processBenchmark(data *Data) {
	key := data.Package + ":" + data.Test
	startedAt, _ := time.Parse(time.RFC3339, data.Time)

	s.Benchmarks[key] = &Benchmark{
		Key:       key,
		Package:   data.Package,
		Name:      data.Test,
		StartedAt: startedAt,
	}
}

func (s *StreamConsumer) processTest(data *Data) {
	key := data.Package + ":" + data.Test
	startedAt, _ := time.Parse(time.RFC3339, data.Time)

	s.Tests[key] = &Test{
		Key:          key,
		Package:      data.Package,
		Name:         data.Test,
		ReadableName: strings.Join(camelcase.Split(data.Test)[1:], " "),
		StartedAt:    startedAt,
	}
}

func (s *StreamConsumer) processCoverage(data *Data) {
	re := regexp.MustCompile(`coverage: ([\d.]+)`)
	matches := re.FindStringSubmatch(data.Output)

	if matches == nil {
		return
	}

	value, err := strconv.ParseFloat(matches[1], 64)

	if err != nil {
		return
	}

	s.Coverage = append(s.Coverage, Coverage{
		Package:  data.Package,
		Coverage: value,
	})
}

func (s *StreamConsumer) processTestOutput(data *Data) {
	// Receiving output for package message that we can ignore.
	if data.Test == "" {
		return
	}

	output := s.normalizePaths(data.Output)

	key := data.Package + ":" + data.Test
	test := s.Tests[key]
	test.Output = append(test.Output, output)
}

func (s *StreamConsumer) processBenchmarkOutput(data *Data) {
	// Receiving output for package message that we can ignore.
	if data.Test == "" {
		return
	}

	key := data.Package + ":" + data.Test
	benchmark := s.Benchmarks[key]
	output := s.normalizePaths(data.Output)
	resultRE := regexp.MustCompile(`^.*?-(\d+)\s+(\d+)\s+([\d.]+ [a-z]+)/op$`)
	matches := resultRE.FindStringSubmatch(output)

	if matches != nil {
		// Benchmarks don't have a skip/fail/pass action like regular tests.
		// The best we can do is acting whenever results are ready.
		benchmark.EndedAt, _ = time.Parse(data.Test, time.RFC3339)
		durationPerOp, _ := time.ParseDuration(strings.ReplaceAll(matches[3], " ", ""))
		iterCount, _ := strconv.Atoi(matches[2])
		maxProcs, _ := strconv.Atoi(matches[1])

		benchmark.MaxProcessors = maxProcs
		benchmark.IterationCount = iterCount
		benchmark.DurationPerOp = durationPerOp
	}

	benchmark.Output = append(benchmark.Output, output)
}

func (s *StreamConsumer) processTestFinished(data *Data) {
	// Receiving a package PASS/SKIP/FAIL message that we can ignore.
	if data.Test == "" {
		return
	}

	key := data.Package + ":" + data.Test
	endedAt, _ := time.Parse(time.RFC3339, data.Time)
	test := s.Tests[key]
	test.EndedAt = endedAt
	test.Status = data.Action
	test.Elapsed = data.Elapsed
	test.ErrorTrace = findErrorTrace(test.Output)
	test.TestSource = findTestSource(test.Output)

	// A skip message may have an output that follows the `file.go:line:( <text>)?`
	// format on its last line. If that happens, let's use the specified path as the test
	// source.
	if test.Status == TestStatus.Skip {
		output := filterOutput(test)
		line := output[len(output)-1]
		re := regexp.MustCompile(`^\s*(.*?:\d+):(?:\s+(.*?))?$`)
		matches := re.FindStringSubmatch(line)
		test.TestSource = matches[1]
		skipMessage := matches[2]

		if skipMessage == "" {
			skipMessage = "Skipped"
		}

		test.SkipMessage = skipMessage

		// remove the last line already, because we already stored the info out of it.
		test.Output = output[:len(output)-1]
	}

	test.Output = filterOutput(test)
	test.Output = deindentOutput(test)

	s.OnNotifyTestFinish(test)
}

func (s *StreamConsumer) normalizePaths(input string) string {
	output := strings.TrimRight(input, "\r\n\t ")

	if s.WorkingDir != "" {
		output = strings.ReplaceAll(output, s.WorkingDir+"/", "")
	}

	if s.HomeDir != "" {
		output = strings.ReplaceAll(output, s.HomeDir+"/", "~/")
	}

	return output
}

func CreateStreamConsumer(options CreateStreamConsumerOptions) *StreamConsumer {
	return &StreamConsumer{
		WorkingDir:         options.WorkingDir,
		HomeDir:            options.HomeDir,
		Tests:              map[string]*Test{},
		Benchmarks:         map[string]*Benchmark{},
		Coverage:           []Coverage{},
		Scanner:            options.Scanner,
		OnNotifyTestFinish: func(test *Test) {},
		OnFinish:           func(tests []*Test, coverage []Coverage, benchmarks []*Benchmark) {},
		OnError:            func() {},
		Output:             options.Output,
	}
}

func findErrorTrace(lines []string) string {
	re := regexp.MustCompile(`^\s+Error Trace:\s+(.*?\.go:\d+)$`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)

		if matches != nil {
			return matches[1]
		}
	}

	return ""
}

func findTestSource(lines []string) string {
	reTrace := regexp.MustCompile(`^\s+Error Trace:\s+(.*?\.go:\d+)$`)
	rePath := regexp.MustCompile(`(?i)^(.*?\.go:\d+)$`)
	errTraceLineNo := -1

	for index, line := range lines {
		matches := reTrace.FindStringSubmatch(line)

		if matches != nil {
			errTraceLineNo = index
			break
		}
	}

	if errTraceLineNo != -1 && len(lines) > errTraceLineNo {
		line := strings.TrimSpace(lines[errTraceLineNo+1])
		matches := rePath.FindStringSubmatch(line)

		if matches != nil {
			return matches[1]
		}
	}

	return ""
}

func filterOutput(test *Test) []string {
	filtered := []string{}
	re := regexp.MustCompile(`^\s+.*?\.go:\d+:?$`)

	for _, line := range test.Output {
		ignore := (test.TestSource != "" && re.Match([]byte(line))) ||
			strings.TrimSpace(line) == test.ErrorTrace+":" ||
			strings.TrimSpace(line) == filepath.Base(test.ErrorTrace)+":" ||
			strings.Contains(line, "Error Trace:") ||
			strings.Contains(line, test.Name) ||
			strings.HasPrefix(line, "coverage:")

		if ignore {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}

func deindentOutput(test *Test) []string {
	lines := []string{}
	indent := ""
	re := regexp.MustCompile(`^(\s+)Error:(\s+)(.+)$`)
	errorIndex := -1
	errorLabel := "Error:  "
	errorSpacing := strings.Repeat(" ", len(errorLabel))

	for index, line := range test.Output {
		matches := re.FindStringSubmatch(line)

		if matches != nil {
			indent = matches[1] + errorSpacing + matches[2]
			test.Output[index] = errorLabel + matches[3]
			errorIndex = index

			break
		}
	}

	if indent == "" {
		return test.Output
	}

	for index, line := range test.Output {
		if errorIndex != index {
			line = errorSpacing + strings.TrimLeft(line, "\r\n\t ")
		}

		lines = append(lines, line)
	}

	return lines
}
