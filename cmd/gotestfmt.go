package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	g "github.com/fnando/gotestfmt/gotestfmt"
	"golang.org/x/exp/slices"
)

type Data map[string]any

func countTestsByStatus(tests []g.Test, status string) int {
	count := 0

	for _, test := range tests {
		if test.Status == status {
			count = count + 1
		}
	}

	return count
}

func main() {
	var (
		reporterName string
		fastFail     bool
		showVersion  bool
	)

	flag.StringVar(&reporterName, "reporter", "dot", "Choose report type (dot, json)")
	flag.BoolVar(&fastFail, "fast-fail", false, "Fast fail")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Println("0.1.0")
		os.Exit(0)
	}

	var reporter g.Reporter

	if reporterName == "dot" {
		reporter = g.CreateDotReporter()
	} else if reporterName == "json" {
		reporter = g.JSONReporter{}
	} else {
		fmt.Fprintln(os.Stderr, "ERROR: expected reporter name to be one of [dot, json]; got", reporterName)
	}

	stat, err := os.Stdin.Stat()

	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: No data has been piped into gotestfmt")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	report := g.Report{
		StartedAt: time.Date(2100, time.January, 31, 23, 59, 59, 0, time.UTC),
		EndedAt:   time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		Tests:     []g.Test{},
	}
	allData := map[string]g.Test{}

	errorTraceRE := regexp.MustCompile("^\\s*Error Trace:\\s*(.+)")

	for scanner.Scan() {
		data := Data{}
		line := scanner.Bytes()
		err = json.Unmarshal(line, &data)

		if err != nil {
			fmt.Print(string(line))
			continue
		}

		time, _ := time.Parse(time.RFC3339, data["Time"].(string))
		action := data["Action"].(string)

		if action == "start" {
			if time.Before(report.StartedAt) {
				report.StartedAt = time
			}
		}

		if action == "run" {
			test := g.Test{
				Name:         data["Test"].(string),
				ReadableName: strings.Join(camelcase.Split(data["Test"].(string))[1:], " "),
				Package:      data["Package"].(string),
				Output:       []string{},
				StartedAt:    time,
				Index:        len(report.Tests),
			}

			allData[test.Package+":"+test.Name] = test
		}

		if slices.Contains([]string{"pass", "fail", "skip"}, action) {
			if time.After(report.EndedAt) {
				report.EndedAt = time
			}
		}

		if slices.Contains([]string{"pass", "fail", "skip"}, action) && data["Test"] != nil {
			test := allData[data["Package"].(string)+":"+data["Test"].(string)]
			test.Status = action
			test.EndedAt = time
			test.ElapsedTime = test.EndedAt.Sub(test.StartedAt)

			report.Tests = append(report.Tests, test)
			reporter.Progress(test, os.Stdout)

			if fastFail && test.Status == "fail" {
				break
			}
		}

		if action == "output" && data["Test"] != nil {
			key := data["Package"].(string) + ":" + data["Test"].(string)
			test := allData[key]

			base, _ := os.Getwd()
			output := strings.Replace(data["Output"].(string), base+"/", "", -1)

			errorTraceMatch := errorTraceRE.FindStringSubmatch(output)

			if len(errorTraceMatch) == 2 {
				test.ErrorTrace = errorTraceMatch[1]
			}

			ignore := strings.Contains(output, "=== RUN") ||
				strings.Contains(output, "--- PASS:") ||
				strings.Contains(output, "--- FAIL:") ||
				strings.Contains(output, "\tTest:") ||
				strings.Contains(output, "Error Trace:") ||
				strings.Contains(output, "Error Trace:") ||
				strings.Contains(output, filepath.Base(test.ErrorTrace))

			if !ignore {
				test.Output = append(test.Output, output)
			}

			allData[key] = test
		}
	}

	report.ElapsedTime = report.EndedAt.Sub(report.StartedAt)
	report.TestsCount = len(report.Tests)
	report.PassCount = countTestsByStatus(report.Tests, "pass")
	report.FailCount = countTestsByStatus(report.Tests, "fail")
	report.SkipCount = countTestsByStatus(report.Tests, "skip")

	reporter.Summary(report, os.Stdout)
	reporter.Exit(report)

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
