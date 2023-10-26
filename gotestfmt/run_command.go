package gotestfmt

import (
	"bufio"
	"bytes"
	"cmp"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/exp/slices"
)

type RunCommandArgs struct {
	Reporter          string
	ReplayFile        string
	Dotenv            string
	NoColor           bool
	CoverageCount     int
	CoverageThreshold float64
	ShowAll           bool
	Args              []string
	Env               []string
	HomeDir           string
	WorkingDir        string
	Output            *OutputBuffers
	Binary            string
}

func RunCommand(args CommandArgs) (exitcode int) {
	var cmdArgs RunCommandArgs = RunCommandArgs{
		Env:        args.Env,
		Output:     args.Output,
		Binary:     args.Binary,
		WorkingDir: args.WorkingDir,
		HomeDir:    args.HomeDir,
	}

	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("gotestfmt run", flag.ContinueOnError)
	flags.SetOutput(buffer)
	flags.StringVar(&cmdArgs.Reporter, "reporter", "progress", "Set the reporter type")
	flags.StringVar(&cmdArgs.ReplayFile, "replay-file", "", "Use a replay file instead of running tests")
	flags.StringVar(&cmdArgs.Dotenv, "dotenv", ".env.test", "Load an env var file before running tests. To disable it, set it to false")
	flags.IntVar(&cmdArgs.CoverageCount, "coverage-count", 10, "The number of coverage items to display")
	flags.Float64Var(&cmdArgs.CoverageThreshold, "coverage-threshold", 100, "The coverage threshold")
	flags.BoolVar(&cmdArgs.NoColor, "no-color", false, "Disable colored output. When unset, respects the NO_COLOR=1 env var")
	flags.BoolVar(&cmdArgs.ShowAll, "all", false, "Show all tests output, including passed ones")
	err := flags.Parse(args.Args)
	cmdArgs.Args = flags.Args()

	var reporter Reporter

	if cmdArgs.Reporter == "progress" {
		reporter = ProgressReporter{
			Env:    args.Env,
			Output: args.Output,
			ColorByStatus: map[string]string{
				"fail": Color.FailColor,
				"skip": Color.SkipColor,
				"pass": Color.PassColor,
			},
		}
	} else if cmdArgs.Reporter == "json" {
		reporter = JSONReporter{
			Output: args.Output,
		}
	} else {
		err = errors.New(Color.Fail("Invalid reporter"))
	}

	Color.Disabled = Color.Disabled || cmdArgs.NoColor

	usage := fmt.Sprintf(`
Run tests by wrapping "go tests".

  Usage: gotestfmt run [options] [packages...] -- [additional "go test" arguments]

  Options:
%s

  Available reporters:
    progress
      Print a character for each test, with a test summary and list of
      failed/skipped tests.

    json
      Print a JSON representation of the gotestfmt state.


  How it works:
    This is what gotestfmt runs if you execute "gotestfmt run ./...":

    $ go test -cover -json -fullpath ./...

    You can pass additional arguments to the "go test" command like this:

    $ gotestfmt run ./... -- -run TestExample

    These arguments will be appended to the default arguments used by gotestfmt.
    The example above would be execute like this:

    $ go test -cover -json -fullpath -run TestExample ./...

    To execute a raw "go test" command, use the switch --raw. This will avoid
    default arguments from being added to the final execution. In practice, it
    means you'll need to run the whole command:

    $ gotestfmt run --raw -- ./some_module -run TestExample


  Replaying files:
    To replay a file, you need to save the output by running something like
    "go test -json -fullpath &> replay.txt". Then you can replay it using
    "gotestfmt run -replay-file replay.txt".


  Color:
    gotestfmt will output colored text based on ANSI colors. By default, the
    following env vars will be used and you can override any of them to set
    a custom color:

    export GOTESTFMT_TEXT_COLOR="30"
    export GOTESTFMT_FAIL_COLOR="31"
    export GOTESTFMT_PASS_COLOR="32"
    export GOTESTFMT_SKIP_COLOR="33"
    export GOTESTFMT_DETAIL_COLOR="34"

    To disable colored output you can use "--no-color" or
    set the env var NO_COLOR=1.


  Progress reporter:
    You can override the default progress symbols by setting env vars. The
    following example shows how to use emojis instead:

    export GOTESTFMT_FAIL_SYMBOL=❌
    export GOTESTFMT_PASS_SYMBOL=🔥
    export GOTESTFMT_SKIP_SYMBOL=😴

`, getOptionsDescription(flags))

	if err == flag.ErrHelp {
		fmt.Fprintln(args.Output.StdoutWriter, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	if cmdArgs.ReplayFile != "" {
		return ReplayFile(&cmdArgs, reporter)
	}

	return Exec(&cmdArgs, reporter)
}

func Exec(args *RunCommandArgs, reporter Reporter) int {
	var (
		packages []string
		extra    []string
		err      error
	)

	if len(args.Args) == 0 {
		packages = []string{"./..."}
	} else if slices.Contains(args.Args, "--") {
		index := slices.Index(args.Args, "--")

		if index+1 < len(args.Args) {
			extra = args.Args[index+1:]
		}

		if index > 0 {
			packages = args.Args[:index]
		}

	} else {
		packages = args.Args
	}

	cmdArgs := append([]string{"test", "-cover", "-json", "-fullpath"}, packages...)
	cmdArgs = append(cmdArgs, extra...)
	cmd := exec.Command("go", cmdArgs...)

	cmd.Stderr = cmd.Stdout
	cmd.Env = args.Env
	out, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	exitcode := 0

	consumer := CreateStreamConsumer(CreateStreamConsumerOptions{
		HomeDir:    args.HomeDir,
		WorkingDir: args.WorkingDir,
		Scanner:    scanner,
		Output:     args.Output,
	})
	consumer.OnNotifyTestFinish = func(test *Test) {
		reporter.Progress(test)
	}

	consumer.OnError = func() {
		exitcode = max(1, exitcode)
	}

	consumer.OnFinish = func(tests []*Test, coverage []Coverage, benchmarks []*Benchmark) {
		statuses := []string{TestStatus.Fail, TestStatus.Skip}

		if args.ShowAll {
			statuses = append(statuses, TestStatus.Pass)
		}

		exitcode, err = reporter.Finish(
			tests,
			statuses,
			prepareCoverage(coverage, args.CoverageCount, args.CoverageThreshold),
			benchmarks,
		)
	}

	cmd.Start()
	consumer.Run()

	if err != nil {
		panic(err)
	}

	return exitcode
}

func ReplayFile(args *RunCommandArgs, reporter Reporter) int {
	stat, err := os.Stat(args.ReplayFile)

	if os.IsNotExist(err) {
		fmt.Fprintf(args.Output.StderrWriter, "%s replay file doesn't exist\n", Color.Fail("ERROR:"))
		return 1
	}

	if stat.IsDir() {
		fmt.Fprintf(args.Output.StderrWriter, "%s replay file can't be a directory\n", Color.Fail("ERROR:"))
		return 1
	}

	file, err := os.Open(args.ReplayFile)

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	exitcode := 0

	consumer := CreateStreamConsumer(CreateStreamConsumerOptions{
		HomeDir:    args.HomeDir,
		WorkingDir: args.WorkingDir,
		Scanner:    scanner,
		Output:     args.Output,
	})

	consumer.OnNotifyTestFinish = func(test *Test) {
		reporter.Progress(test)
	}

	consumer.OnError = func() {
		exitcode = max(1, exitcode)
	}

	consumer.OnFinish = func(tests []*Test, coverage []Coverage, benchmarks []*Benchmark) {
		statuses := []string{TestStatus.Fail, TestStatus.Skip}

		if args.ShowAll {
			statuses = append(statuses, TestStatus.Pass)
		}

		exitcode, err = reporter.Finish(
			tests,
			statuses,
			prepareCoverage(coverage, args.CoverageCount, args.CoverageThreshold),
			benchmarks,
		)
	}

	consumer.Run()

	if err != nil {
		panic(err)
	}

	return exitcode
}

func prepareCoverage(coverages []Coverage, count int, threshold float64) []Coverage {
	list := []Coverage{}

	for _, coverage := range coverages {
		if coverage.Coverage < threshold {
			list = append(list, coverage)
		}

		if len(list) == count {
			break
		}
	}

	slices.SortFunc(list, func(a, b Coverage) int {
		return cmp.Compare(a.Coverage, b.Coverage)
	})

	return list
}
