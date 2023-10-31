package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"
)

type StreamConsumer struct {
	Aggregation *Aggregation
	OnData      func(data string)
	OnProgress  func(test Test)
	OnFinished  func(aggregation *Aggregation)
}

type Stream struct {
	Action  string
	Elapsed float64
	Output  string
	Package string
	Test    string
	Time    string
}

type Test struct {
	Key             string `json:"-"`
	ErrorTrace      string
	ErrorTraceIndex int `json:"-"`
	Source          string
	ReadableName    string
	Name            string
	StartedAt       time.Time
	EndedAt         time.Time
	Elapsed         time.Duration
	Output          []string
	Status          string
	SkipMessage     string
	Package         string
}

type Benchmark struct {
	Name                 string
	Package              string
	Key                  string `json:"-"`
	Processors           int
	Iterations           int
	DurationPerOperation time.Duration
}

func (consumer StreamConsumer) Ingest(scanner *bufio.Scanner) {
	consumer.Aggregation.StartedAt = time.Now()

	for scanner.Scan() {
		var stream Stream
		line := scanner.Bytes()
		lineStr := string(line)

		consumer.OnData(lineStr)

		err := json.Unmarshal(line, &stream)

		if err != nil {
			fmt.Println(lineStr)

			if strings.Contains(lineStr, "[build failed]") {
				os.Exit(1)
			}
		} else {
			consumer.process(stream)
		}
	}

	consumer.Aggregation.EndedAt = time.Now()
	consumer.OnFinished(consumer.Aggregation)
}

func (consumer StreamConsumer) process(stream Stream) {
	switch stream.Action {
	case "start":
		// Suite started running

	case "run":
		// A test/benchmark just started running.

		if strings.HasPrefix(stream.Test, "Test") {
			test := Test{
				Name:            stream.Test,
				Package:         stream.Package,
				ErrorTraceIndex: -1,
				ReadableName:    strings.ReplaceAll(strings.Join(camelcase.Split(stream.Test)[1:], " "), " _ ", " "),
				Key:             stream.Package + ":" + stream.Test,
				StartedAt:       time.Now(),
			}

			coverage := Coverage{Package: test.Package}

			consumer.Aggregation.TestsMap[test.Key] = &test
			consumer.Aggregation.CoverageMap[coverage.Package] = &coverage
		} else if strings.HasPrefix(stream.Test, "Benchmark") {
			benchmark := Benchmark{
				Name:    stream.Test,
				Package: stream.Package,
				Key:     stream.Package + ":" + stream.Test,
			}

			consumer.Aggregation.BenchmarksMap[benchmark.Key] = &benchmark
		}

	case "output":
		// Something was printed to the console.
		key := stream.Package + ":" + stream.Test

		if strings.HasPrefix(stream.Test, "Benchmark") {
			output := strings.TrimSpace(stream.Output)
			re := regexp.MustCompile(`^(?:.*?-(\d+))\s+(\d+)\s+(.*?)/op$`)
			matches := re.FindStringSubmatch(output)

			if matches != nil {
				benchmark := consumer.Aggregation.BenchmarksMap[key]
				procs, _ := strconv.Atoi(matches[1])
				iters, _ := strconv.Atoi(matches[2])
				dur, _ := time.ParseDuration(strings.ReplaceAll(matches[3], " ", ""))
				benchmark.Processors = procs
				benchmark.Iterations = iters
				benchmark.DurationPerOperation = dur
			}

			return
		}

		if stream.Test == "" {
			re := regexp.MustCompile(`coverage: ([\d.]+)% of statements`)
			matches := re.FindStringSubmatch(stream.Output)

			if matches != nil {
				percent, _ := strconv.ParseFloat(matches[1], 64)
				consumer.Aggregation.CoverageMap[stream.Package].Coverage = percent
			}

			return
		}

		output := strings.TrimRight(stream.Output, "\r\n")
		test := consumer.Aggregation.TestsMap[key]
		index := len(test.Output)
		errorTrace := findErrorTrace(output)
		shouldAppend := true

		if errorTrace != "" {
			test.ErrorTraceIndex = index
			test.ErrorTrace = errorTrace
			shouldAppend = false
		} else if index == test.ErrorTraceIndex {
			source := findSource(output)

			if source != "" {
				test.Source = test.ErrorTrace
				test.ErrorTrace = source
				shouldAppend = false
			}
		}

		if shouldAppend {
			test.Output = append(test.Output, output)
		}

	case "fail":
		fallthrough
	case "skip":
		fallthrough
	case "pass":
		// Test/benchmark has finished running.
		if stream.Test == "" {
			return
		}

		key := stream.Package + ":" + stream.Test
		test := consumer.Aggregation.TestsMap[key]
		test.EndedAt = time.Now()
		test.Elapsed = test.EndedAt.Sub(test.StartedAt)
		test.Status = stream.Action
		consumer.OnProgress(*test)

		if test.Status == "skip" {
			// let's extract the error trace and message
			index := len(test.Output) - 2
			line := test.Output[index]
			re := regexp.MustCompile(`^(\s*)(.*?\.go:\d+):(?:\s+(.*?))?$`)
			matches := re.FindStringSubmatch(line)
			test.ErrorTrace = matches[2]

			if matches[3] != "" {
				test.Output[index] = matches[3]
			} else {
				test.Output[index] = "[No message]"
			}
		}

	default:
		panic("here")
	}
}

func findErrorTrace(line string) string {
	re := regexp.MustCompile(`^Error Trace:\s*(.*?)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(line))

	if matches == nil {
		return ""
	}

	return matches[1]
}

func findSource(line string) string {
	re := regexp.MustCompile(`^(.*?\.go:\d+):?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(line))

	if matches == nil {
		return ""
	}

	return matches[1]
}
