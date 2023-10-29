package gotestfmt

import (
	"bytes"
	"io"
)

type OutputBuffers struct {
	Stdout       *bytes.Buffer
	Stderr       *bytes.Buffer
	StdoutWriter io.Writer
	StderrWriter io.Writer
}
