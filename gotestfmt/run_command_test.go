package gotestfmt

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

var env []string = os.Environ()

func TestRunHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "-h"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-help.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunMissingReplayFile(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./cli/replays/doesnt-exist.txt"}, env, output)

	assert.Equal(t, "ERROR: replay file doesn't exist\n", output.Stderr.String())
	assert.Empty(t, output.Stdout.String())
	assert.Equal(t, 1, exitcode)
}

func TestRunDirectoryAsReplayFile(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "."}, env, output)

	assert.Equal(t, "ERROR: replay file can't be a directory\n", output.Stderr.String())
	assert.Empty(t, output.Stdout.String())
	assert.Equal(t, 1, exitcode)
}

func TestRunPassedReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/passed.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-passed.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunWithShowAllReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--all", "--replay-file", "./test/replays/passed.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-passed-all.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunFailedReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/failed.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-failed.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 2, exitcode)
}

func TestRunSkippedReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/skipped.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-skipped.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunErrorsReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/error.txt"}, env, output)

	want := `# github.com/fnando/gotestfmt/reference/fail [github.com/fnando/gotestfmt/reference/fail.test]
reference/fail/fail_test.go:5:2: "os" imported and not used
FAIL	github.com/fnando/gotestfmt/reference/fail [build failed]
`

	assert.Equal(t, want, output.Stdout.String())
	assert.Equal(t, 1, exitcode)
}

func TestRunMixedReplayFileWithCustomSymbolsOnDotsReporter(t *testing.T) {
	output := createBuffers()
	extraEnv := []string{
		"GOTESTFMT_FAIL_SYMBOL=❌",
		"GOTESTFMT_PASS_SYMBOL=🔥",
		"GOTESTFMT_SKIP_SYMBOL=😴",
	}
	env := append(os.Environ(), extraEnv...)
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/mixed.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/dots-custom-symbol.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, 3, exitcode)
}

func TestRunCoverageReplayFileWithDotsReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--no-color", "--replay-file", "./test/replays/coverage.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-coverage.txt", output.Stdout.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunPassedReplayFileWithJSONReporter(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--reporter", "json", "--no-color", "--replay-file", "./test/replays/coverage.txt"}, env, output)

	assertEqualToFileContent(t, "./test/expected/run-passed.json", output.Stdout.String())
	assert.Equal(t, 0, exitcode)
}

func TestRunLoadDefaultDotenvFile(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--reporter", "json", "--no-color", "./reference/dotenv_default"}, env, output)

	assert.Equal(t, 0, exitcode)
}

func TestRunCustomDotenvFile(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--reporter", "json", "--no-color", "--dotenv", "test/dotenv", "./reference/dotenv_custom"}, env, output)

	assert.Equal(t, 0, exitcode)
}

func TestRunDisableDefaultDotenvFile(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"run", "--reporter", "json", "--no-color", "--dotenv", "false", "./reference/dotenv_disabled"}, env, output)

	assert.Equal(t, 0, exitcode)
}
