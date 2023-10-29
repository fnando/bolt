package gotestfmt

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
)

func VersionCommand(args CommandArgs) (exitcode int) {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("gotestfmt version", flag.ContinueOnError)
	flags.SetOutput(buffer)
	err := flags.Parse(args.Args)

	usage := `
Display the version

  Usage: gotestfmt version
`

	if err == flag.ErrHelp {
		fmt.Fprintln(args.Output.StdoutWriter, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	fmt.Fprintln(args.Output.StdoutWriter, Version)
	return 0
}
