package gotestfmt

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestShowDownloadURL(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"download-url"}, []string{}, output)

	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, "https://github.com/fnando/gotestfmt/releases/latest/download/gotestfmt-unknown\n", output.Stdout.String())
	assert.Equal(t, 0, exitcode)
}

func TestShowDownloadURLHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"download-url", "-h"}, []string{}, output)

	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
	assertEqualToFileContent(t, "./test/expected/download-url-help.txt", output.Stdout.String())
}
