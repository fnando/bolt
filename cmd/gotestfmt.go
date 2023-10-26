package main

import (
	"os"

	"github.com/fnando/gotestfmt/gotestfmt"
)

func main() {
	output := gotestfmt.OutputBuffers{StdoutWriter: os.Stdout, StderrWriter: os.Stderr}
	exitcode := gotestfmt.Run(os.Args[1:], os.Environ(), output)
	os.Exit(exitcode)
}
