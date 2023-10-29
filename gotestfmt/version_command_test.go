package gotestfmt

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestShowVersion(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"version"}, []string{}, output)

	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
	assert.Equal(t, Version+"\n", output.Stdout.String())
}

func TestShowVersionHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"version", "-h"}, []string{}, output)

	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
	assertEqualToFileContent(t, "./test/expected/version-help.txt", output.Stdout.String())
}
