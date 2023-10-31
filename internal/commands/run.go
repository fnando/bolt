package commands

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	c "github.com/fnando/bolt/common"
	"github.com/fnando/bolt/internal/reporters"
	"github.com/joho/godotenv"
)

type RunArgs struct {
	Compat            bool
	CoverageCount     int
	CoverageThreshold float64
	Debug             bool
	Dotenv            string
	HideCoverage      bool
	HideSlowest       bool
	HomeDir           string
	NoColor           bool
	Raw               bool
	Replay            string
	Reporter          string
	SlowestCount      int
	SlowestThreshold  string
	WorkingDir        string
}

var usage string = `
Run tests by wrapping "go tests".

  Usage: bolt [options] [packages...] -- [additional "go test" arguments]

  Options:
%s

  Available reporters:
    progress
      Print a character for each test, with a test summary and list of
      failed/skipped tests.

    json
      Print a JSON representation of the bolt state.


  How it works:
    This is what bolt runs if you execute "bolt ./...":

    $ go test ./... -cover -json -fullpath

    You can pass additional arguments to the "go test" command like this:

    $ bolt ./... -- -run TestExample

    These arguments will be appended to the default arguments used by bolt.
    The example above would be executed like this:

    $ go test -cover -json -fullpath -run TestExample ./...

    To execute a raw "go test" command, use the switch --raw. This will avoid
    default arguments from being added to the final execution. In practice, it
    means you'll need to run the whole command:

    $ bolt --raw -- ./some_module -run TestExample

    Note: -fullpath was introduced on go 1.21. If you're using an older
    version, you can use --compat or manually set arguments by using --raw.


  Env files:
    bolt will load .env.test by default. You can also set it to a
    different file by using --env. If you want to disable env files
    completely, use --env=false.


  Color:
    bolt will output colored text based on ANSI colors. By default, the
    following env vars will be used and you can override any of them to set
    a custom color:

    export BOLT_TEXT_COLOR="30"
    export BOLT_FAIL_COLOR="31"
    export BOLT_PASS_COLOR="32"
    export BOLT_SKIP_COLOR="33"
    export BOLT_DETAIL_COLOR="34"

    To disable colored output you can use "--no-color" or
    set the env var NO_COLOR=1.


  Progress reporter:
    You can override the default progress symbols by setting env vars. The
    following example shows how to use emojis instead:

    export BOLT_FAIL_SYMBOL=❌
    export BOLT_PASS_SYMBOL=⚡️
    export BOLT_SKIP_SYMBOL=😴
`

func Run(args []string, options RunArgs, output *c.Output) int {
	flags := flag.NewFlagSet("bolt", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(
		&options.NoColor,
		"no-color",
		false,
		"Disable colored output. When unset, respects the NO_COLOR=1 env var",
	)

	flags.BoolVar(&options.Raw, "raw", false, "Don't append arguments to `go test`")
	flags.BoolVar(&options.Compat, "compat", false, "Don't append -fullpath, available on go 1.21 or new")
	flags.BoolVar(&options.HideCoverage, "hide-coverage", false, "Don't display the coverage section")
	flags.BoolVar(&options.HideSlowest, "hide-slowest", false, "Don't display the slowest tests section")
	flags.StringVar(&options.Dotenv, "env", ".env.test", "Load env file")
	flags.IntVar(&options.CoverageCount, "coverage-count", 10, "Number of coverate items to show")
	flags.Float64Var(&options.CoverageThreshold, "coverage-threshold", 100.0, "Anything below this threshold will be listed")
	flags.StringVar(&options.SlowestThreshold, "slowest-threshold", "1s", "Anything above this threshold will be listed. Must be a valid duration string")
	flags.IntVar(&options.SlowestCount, "slowest-count", 10, "Number of slowest tests to show")

	flags.BoolVar(&options.Debug, "debug", false, "")
	flags.StringVar(&options.Replay, "replay", "", "")
	flags.StringVar(&options.Reporter, "reporter", "progress", "")

	flags.SetOutput(bufio.NewWriter(&bytes.Buffer{}))
	err := flags.Parse(args)

	if options.Dotenv != "false" {
		dotenvErr := godotenv.Load(options.Dotenv)

		if dotenvErr != nil {
			ignore := strings.Contains(dotenvErr.Error(), "no such file") &&
				options.Dotenv == ".env.test"

			if !ignore {
				err = dotenvErr
			}
		}
	}

	if options.Debug {
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" version:", c.Version)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" arch:", c.Arch)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" commit:", c.Commit)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" working dir:", options.WorkingDir)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" home dir:", options.HomeDir)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" reporter:", options.Reporter)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" env file:", options.Dotenv)
		fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" compat:", options.Compat)

		if options.Replay != "" {
			fmt.Fprintln(output.Stdout, c.Color.Detail("⚡️")+" replay file:", options.Replay)
		}
	}

	if err == flag.ErrHelp {
		fmt.Fprintf(output.Stdout, usage, getFlagsUsage(flags))
		return 0
	} else if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	slowestThreshold, err := time.ParseDuration(options.SlowestThreshold)

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	exitcode := 1
	consumer := c.StreamConsumer{
		Aggregation: &c.Aggregation{
			TestsMap:          map[string]*c.Test{},
			CoverageMap:       map[string]*c.Coverage{},
			BenchmarksMap:     map[string]*c.Benchmark{},
			CoverageThreshold: options.CoverageThreshold,
			CoverageCount:     options.CoverageCount,
			SlowestThreshold:  slowestThreshold,
			SlowestCount:      options.SlowestCount,
		},
	}
	var reporter reporters.Reporter

	if options.Reporter == "progress" {
		reporter = reporters.ProgressReporter{Output: output}
	} else if options.Reporter == "standard" {
		reporter = reporters.StandardReporter{Output: output}
	} else if options.Reporter == "json" {
		reporter = reporters.JSONReporter{Output: output}
	} else {
		fmt.Fprintf(output.Stderr, "%s %s\n", c.Color.Fail("ERROR:"), "Invalid reporter")
		return 1
	}

	consumer.OnData = func(line string) { reporter.OnData(line) }
	consumer.OnProgress = func(test c.Test) { reporter.OnProgress(test) }
	consumer.OnFinished = func(aggregation *c.Aggregation) {
		reporter.OnFinished(
			reporters.ReporterFinishedOptions{
				Aggregation:  aggregation,
				HideCoverage: options.HideCoverage,
				HideSlowest:  options.HideSlowest,
			},
		)
	}

	if options.Replay == "" {
		execArgs := append([]string{"-json", "-cover"})

		if !options.Compat {
			execArgs = append(execArgs, "-fullpath")
		}

		extraArgs := []string{}

		for _, arg := range flags.Args() {
			if arg != "--" {
				extraArgs = append(extraArgs, arg)
			}
		}

		execArgs = append(execArgs, extraArgs...)

		if options.Raw {
			execArgs = flags.Args()
		}

		if options.Debug {
			fmt.Fprintln(
				output.Stdout,
				c.Color.Detail("⚡️"),
				"command:",
				"go test",
				strings.Join(execArgs, " "),
			)
		}

		exitcode, err = Exec(&consumer, output, execArgs)
	} else {
		exitcode, err = Replay(&consumer, &options)
	}

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	return exitcode
}

func Replay(consumer *c.StreamConsumer, options *RunArgs) (int, error) {
	stat, err := os.Stat(options.Replay)

	if os.IsNotExist(err) {
		return 1, errors.New("replay file doesn't exist")
	}

	if stat.IsDir() {
		return 1, errors.New("can't read directory (" + options.Replay + ")")
	}

	file, err := os.Open(options.Replay)

	if err != nil {
		return 1, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	consumer.Ingest(scanner)

	return consumer.Aggregation.CountBy("fail"), nil
}

func Exec(consumer *c.StreamConsumer, output *c.Output, args []string) (int, error) {
	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = cmd.Stdout
	out, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)

	err := cmd.Start()

	if err != nil {
		return 1, err
	}

	consumer.Ingest(scanner)

	err = cmd.Wait()

	if err != nil {
		return 1, err
	}

	exitcode := 0

	if exiterr, ok := err.(*exec.ExitError); ok {
		exitcode = exiterr.ExitCode()
		return exitcode, nil
	}

	return exitcode, err
}
