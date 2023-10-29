package gotestfmt

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"runtime"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func init() {
	var err error

	err = os.Setenv("GOTESTFMT_HOME_DIR", "/home/test")

	if err != nil {
		panic(err)
	}

	err = os.Setenv("GOTESTFMT_WORKING_DIR", "/home/test/gotestfmt")

	if err != nil {
		panic(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err = os.Chdir(dir)

	if err != nil {
		panic(err)
	}
}

func createBuffers() (output OutputBuffers) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	stdoutw := bufio.NewWriterSize(&stdout, 1)
	stderrw := bufio.NewWriterSize(&stderr, 1)

	return OutputBuffers{
		Stdout:       &stdout,
		Stderr:       &stderr,
		StdoutWriter: stdoutw,
		StderrWriter: stderrw,
	}
}

func assertEqualToFileContent(t *testing.T, path string, actual string) {
	assert.FileExists(t, path)

	expected, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), actual)
}
