package reporters

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	c "github.com/fnando/bolt/common"
)

type PostRunCommandReporter struct {
	Command string
	Output  *c.Output
}

func (reporter PostRunCommandReporter) Name() string {
	return "post-command"
}

func (reporter PostRunCommandReporter) OnFinished(options ReporterFinishedOptions) {
	if reporter.Command == "" {
		return
	}

	total := options.Aggregation.TestsCount()
	fail := options.Aggregation.CountBy("fail")
	pass := options.Aggregation.CountBy("pass")
	skip := options.Aggregation.CountBy("skip")
	benchmarks := len(options.Aggregation.Benchmarks())
	elapsed := options.Aggregation.Elapsed()
	elapsedNS := int(elapsed)
	title := "Passed!"

	if fail > 0 {
		title = "Failed!"
	}

	env := os.Environ()
	env = append(
		env,
		fmt.Sprintf(
			"BOLT_SUMMARY=Finished in %s, %d tests, %d fails, %d skips, %d benchmarks",
			formatDuration(elapsed, 2),
			total,
			fail,
			skip,
			benchmarks,
		),
		fmt.Sprintf("BOLT_TEST_COUNT=%d", total),
		fmt.Sprintf("BOLT_FAIL_COUNT=%d", fail),
		fmt.Sprintf("BOLT_PASS_COUNT=%d", pass),
		fmt.Sprintf("BOLT_SKIP_COUNT=%d", skip),
		fmt.Sprintf("BOLT_BENCHMARK_COUNT=%d", benchmarks),
		fmt.Sprintf("BOLT_ELAPSED_NANOSECONDS=%d", elapsedNS),
		fmt.Sprintf("BOLT_ELAPSED=%s", formatDuration(elapsed, 2)),
		fmt.Sprintf("BOLT_TITLE=%s", title),
	)

	var buffer bytes.Buffer
	out := io.Writer(&buffer)

	dir, _ := os.Getwd()

	cmd := exec.Command("bash", "-c", reporter.Command)
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Dir = dir
	cmd.Env = env
	err := cmd.Start()

	if options.Debug && err != nil {
		fmt.Fprintln(
			reporter.Output.Stderr,
			"\n",
			c.Color.Detail("⚡️")+" failed to run post run command:",
			err,
			buffer.String(),
		)
	}

	if err != nil {
		return
	}

	err = cmd.Wait()

	if options.Debug && err != nil {
		fmt.Fprintln(
			reporter.Output.Stderr,
			"\n",
			c.Color.Detail("⚡️")+" failed to run post run command:",
			err,
		)

		fmt.Fprintln(
			reporter.Output.Stderr,
			"output:",
			buffer.String(),
		)
	}

	if err != nil {
		return
	}
}

func (reporter PostRunCommandReporter) OnProgress(test c.Test) {
}

func (reporter PostRunCommandReporter) OnData(line string) {
}
