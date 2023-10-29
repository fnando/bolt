package gotestfmt

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
)

func DownloadURLCommand(args CommandArgs) (exitcode int) {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("gotestfmt download-url", flag.ContinueOnError)
	flags.SetOutput(buffer)
	err := flags.Parse(args.Args)

	usage := `
Return the latest binary url, matching  the one being executed.

  Usage: gotestfmt download-url
`

	if err == flag.ErrHelp {
		fmt.Fprintln(args.Output.StdoutWriter, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	fmt.Fprintln(args.Output.StdoutWriter, DownloadURL)
	return 0
}
