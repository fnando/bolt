package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	g "github.com/fnando/gotestfmt/gotestfmt"
	"golang.org/x/exp/slices"
)

var cliVersion = "0.1.7"

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
		reporterName      string
		showCoverage      bool
		coverageThreshold float64
		coverageCount     int
	)

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "\ngotestfmt is a tool that generates a better output format for golang tests.\n\n")
		fmt.Fprint(w, "Usage: gotestfmt [OPTIONS]\n\n")
		flag.PrintDefaults()
		fmt.Fprint(w, "\nOther commands:\n\n")
		fmt.Fprint(w, "  gotestfmt download-url\n      display the latest binary download url\n\n")
		fmt.Fprint(w, "  gotestfmt update\n      download the latest binary and replace the running one\n\n")
		fmt.Fprint(w, "  gotestfmt version\n      display the version\n\n")
		fmt.Fprint(w, "  gotestfmt help\n      display this help\n\n")
		fmt.Fprintln(w, "For more info, visit https://github.com/fnando/gotestfmt")
	}

	flag.StringVar(&reporterName, "reporter", "dot", "Choose report type (dot, json)")
	flag.BoolVar(&showCoverage, "cover", true, "Show module coverage")
	flag.Float64Var(&coverageThreshold, "cover-threshold", 100.0, "Only show module coverage below this threshold")
	flag.IntVar(&coverageCount, "cover-count", 10, "Number of coverage items to display")
	flag.Parse()

	tailArgs := flag.Args()

	if len(tailArgs) > 0 {
		cmd := tailArgs[0]

		if cmd == "download-url" {
			fmt.Print(downloadUrl())
			os.Exit(0)
		} else if cmd == "update" {
			update()
			os.Exit(0)
		} else if cmd == "version" {
			fmt.Println(cliVersion)
			os.Exit(0)
		} else if cmd == "help" {
			flag.Usage()
			os.Exit(0)
		} else {
			flag.Usage()
			os.Exit(1)
		}
	}

	var reporter g.Reporter

	if reporterName == "dot" {
		reporter = g.CreateDotReporter()
	} else if reporterName == "json" {
		reporter = g.JSONReporter{}
	} else {
		fmt.Fprintln(os.Stderr, "ERROR: expected reporter name to be one of [dot, json]; got", reporterName)
		flag.Usage()
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()

	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: No data has been piped into gotestfmt")
		flag.Usage()
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	report := g.Report{
		StartedAt: time.Date(2100, time.January, 31, 23, 59, 59, 0, time.UTC),
		EndedAt:   time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		Tests:     []g.Test{},
		Coverage:  []g.Coverage{},
	}
	allData := map[string]g.Test{}

	errorTraceRE := regexp.MustCompile("^\\s*Error Trace:\\s*(.+)")
	coverageRE := regexp.MustCompile("^coverage: ([\\d.]+)% of statements")
	var hasFailed bool

	for scanner.Scan() {
		data := Data{}
		lineStr := string(scanner.Bytes())

		base, _ := os.Getwd()
		lineStr = strings.Replace(lineStr, base+"/", "", -1)
		line := []byte(lineStr)

		err = json.Unmarshal(line, &data)

		hasFailed = hasFailed ||
			strings.HasPrefix(lineStr, "FAIL\t") ||
			strings.Contains(lineStr, "panic:") ||
			strings.Contains(lineStr, "[build failed]")

		if err != nil {
			re, _ := regexp.Compile("(?m)^FAIL")

			fmt.Println(lineStr, err)

			if re.Match(line) {
				os.Exit(1)
			}

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
		}

		if action == "output" && data["Test"] == nil {
			output := data["Output"].(string)
			noise := strings.HasPrefix(output, "FAIL\t") ||
				strings.Contains(output, "[no test files]")

			if noise {
				continue
			}

			coverageMatch := coverageRE.FindStringSubmatch(output)

			if len(coverageMatch) == 2 {
				val, _ := strconv.ParseFloat(coverageMatch[1], 32)

				report.Coverage = append(report.Coverage, g.Coverage{
					Package:  data["Package"].(string),
					Coverage: val,
				})
			} else {
				fmt.Println(output)
			}
		}

		if action == "output" && data["Test"] != nil {
			key := data["Package"].(string) + ":" + data["Test"].(string)
			test := allData[key]
			output := data["Output"].(string)

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

	filteredCoverage := prepareCoverage(report, coverageThreshold)

	reporter.Summary(report, os.Stdout)
	reporter.Coverage(filteredCoverage[:min(coverageCount, len(filteredCoverage))], os.Stdout)

	exitcode := 0

	if hasFailed {
		exitcode = 1
	}

	os.Exit(max(exitcode, reporter.Exit(report)))

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func prepareCoverage(report g.Report, threshold float64) []g.Coverage {
	coverageList := report.Coverage
	sort.SliceStable(coverageList, func(i, j int) bool {
		return coverageList[i].Coverage < coverageList[j].Coverage
	})

	filteredCoverage := []g.Coverage{}

	for _, coverage := range coverageList {
		if coverage.Coverage < threshold {
			filteredCoverage = append(filteredCoverage, coverage)
		}
	}

	return filteredCoverage
}

func downloadUrl() string {
	return fmt.Sprintf(
		"https://github.com/fnando/gotestfmt/releases/latest/download/gotestfmt-%s",
		g.Arch,
	)
}

func update() {
	url := downloadUrl()

	bin, err := os.Executable()
	if err != nil {
		panic(err)
	}

	out, err := os.Create(bin)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
