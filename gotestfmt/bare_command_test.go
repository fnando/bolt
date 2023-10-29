package gotestfmt

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestShowHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"-h"}, []string{}, output)

	assert.Equal(t, 0, exitcode)
	assertEqualToFileContent(t, "./test/expected/help.txt", output.Stdout.String())
	assert.Empty(t, output.Stderr.String())
}
