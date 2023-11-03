package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"testing"
	"time"

	c "github.com/fnando/bolt/common"
	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)

	if err != nil {
		panic(err)
	}
}

type execResult struct {
	exitcode int
	stdout   string
	stderr   string
}

func read(path string) string {
	contents, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return string(contents)
}

func replaceSummary(input string) string {
	re := regexp.MustCompile(`Finished in ([^,]+)`)
	return re.ReplaceAllString(input, "Finished in 0s")
}

func run(args []string, env []string) (execResult, error) {
	args = append([]string{"run", "./cmd/bolt.go"}, args...)
	stdout := bytes.NewBufferString("")
	stderr := bytes.NewBufferString("")
	dir, _ := os.Getwd()
	cmd := exec.Command("go", args...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = dir
	err := cmd.Start()

	if err != nil {
		return execResult{
			stdout:   stdout.String(),
			stderr:   stderr.String(),
			exitcode: -1,
		}, err
	}

	err = cmd.Wait()

	result := execResult{
		stdout:   stdout.String(),
		stderr:   stderr.String(),
		exitcode: 0,
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		result.exitcode = exiterr.ExitCode()

		return result, nil
	}

	if err != nil {
		return execResult{
			stdout:   stdout.String(),
			stderr:   stderr.String(),
			exitcode: -1,
		}, err
	}

	return result, err
}

func TestMain(m *testing.M) {
	c.Clock.Now = time.Now
	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestRunCommand(t *testing.T) {
	require.NoError(t, nil)

	t.Run("Help", func(t *testing.T) {
		result, err := run([]string{"run", "--help"}, []string{})

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-help.txt"), result.stdout)
		require.Equal(t, 0, result.exitcode)
	})

	t.Run("SuccessReplayFile", func(t *testing.T) {
		c.Clock.Now = func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2023-10-31T12:15:22.773212-07:00")
			return t
		}
		result, err := run(
			[]string{"run", "--no-color", "--replay", "test/replays/run-pass.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-pass.txt"), replaceSummary(result.stdout))
		require.Equal(t, 0, result.exitcode)
	})

	t.Run("FailReplayFile", func(t *testing.T) {
		c.Clock.Now = func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2023-10-26T14:41:41.05297-07:00")
			return t
		}
		result, err := run(
			[]string{"run", "--no-color", "--replay", "test/replays/run-fail.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-fail.txt"), replaceSummary(result.stdout))
		require.Contains(t, result.stderr, "exit status 3")
		require.Equal(t, 1, result.exitcode)
	})

	t.Run("SkipReplayFile", func(t *testing.T) {
		c.Clock.Now = func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2023-10-26T14:41:41.05297-07:00")
			return t
		}
		result, err := run(
			[]string{"run", "--no-color", "--replay", "test/replays/run-skip.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-skip.txt"), replaceSummary(result.stdout))
		require.Equal(t, 0, result.exitcode)
	})

	t.Run("CustomSymbols", func(t *testing.T) {
		result, err := run(
			[]string{"run", "--no-color", "--replay", "test/replays/run-mixed.txt"},
			[]string{"BOLT_FAIL_SYMBOL=❌", "BOLT_PASS_SYMBOL=⚡️", "BOLT_SKIP_SYMBOL=😴"},
		)

		require.NoError(t, err)
		require.Contains(t, result.stdout, "⚡️⚡️❌❌❌⚡️⚡️⚡️😴😴")
	})

	t.Run("ColorOutput", func(t *testing.T) {
		result, err := run(
			[]string{"run", "--replay", "test/replays/run-mixed.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(
			t,
			read("test/expected/run-mixed-color.txt"),
			replaceSummary(result.stdout),
		)
		require.Contains(t, result.stderr, "exit status 3")
		require.Equal(t, 1, result.exitcode)
	})

	t.Run("CustomColorOutput", func(t *testing.T) {
		result, err := run(
			[]string{"run", "--replay", "test/replays/run-mixed.txt"},
			[]string{
				"BOLT_TEXT_COLOR=31",
				"BOLT_FAIL_COLOR=32",
				"BOLT_PASS_COLOR=33",
				"BOLT_SKIP_COLOR=34",
				"BOLT_DETAIL_COLOR=35",
			},
		)

		require.NoError(t, err)
		require.Equal(
			t,
			read("test/expected/run-mixed-custom-color.txt"),
			replaceSummary(result.stdout),
		)
		require.Contains(t, result.stderr, "exit status 3")
		require.Equal(t, 1, result.exitcode)
	})

	t.Run("ReplayError", func(t *testing.T) {
		result, err := run(
			[]string{"run", "--replay", "test/replays/run-error.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-error.txt"), result.stdout)
		require.Equal(t, 1, result.exitcode)
	})

	t.Run("ReplayNoTests", func(t *testing.T) {
		result, err := run(
			[]string{"run", "--no-color", "--replay", "test/replays/run-no-tests.txt"},
			[]string{},
		)

		require.NoError(t, err)
		require.Equal(t, read("test/expected/run-no-tests.txt"), replaceSummary(result.stdout))
		require.Equal(t, 0, result.exitcode)
	})

	t.Run("JSON", func(t *testing.T) {
		c.Clock.Now = func() time.Time {
			t, _ := time.Parse(time.RFC3339, "2023-10-26T14:41:41.05297-07:00")
			return t
		}

		result, err := run(
			[]string{"run", "--reporter", "json", "--replay", "test/replays/run-mixed.txt"},
			[]string{},
		)

		require.NoError(t, err)

		var data any
		err = json.Unmarshal([]byte(result.stdout), &data)
		require.NoError(t, err)

		require.Contains(t, result.stderr, "exit status 3")
		require.Equal(t, 1, result.exitcode)
	})
}
