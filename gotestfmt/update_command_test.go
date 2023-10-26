package gotestfmt

import (
	"syscall"
	"testing"

	"github.com/jarcoal/httpmock"
	assert "github.com/stretchr/testify/require"
)

func TestShowUpdateHelp(t *testing.T) {
	output := createBuffers()
	exitcode := Run([]string{"update", "-h"}, []string{}, output)

	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
	assertEqualToFileContent(t, "./test/expected/update-help.txt", output.Stdout.String())
}

func TestUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer func() {
		syscall.Unlink("./tmp/fake-bin")
	}()

	httpmock.RegisterResponder("GET", DownloadURL,
		httpmock.NewStringResponder(200, `downloaded bin`),
	)

	output := createBuffers()
	exitcode := Run([]string{"update"}, []string{"GOTESTFMT_BINARY=./tmp/fake-bin"}, output)

	assert.Equal(t, 0, exitcode)
	assert.Empty(t, output.Stderr.String())
	assert.FileExists(t, "./tmp/fake-bin")
	assertEqualToFileContent(t, "./tmp/fake-bin", "downloaded bin")
}
