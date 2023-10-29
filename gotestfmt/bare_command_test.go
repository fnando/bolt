package gotestfmt

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestShowHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"-h"}, []string{}, output)

	assertEqualToFileContent(t, "./test/expected/help.txt", output.Stdout.String())
	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
}
