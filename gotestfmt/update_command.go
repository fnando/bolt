package gotestfmt

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func UpdateCommand(args CommandArgs) (exitcode int) {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("gotestfmt update", flag.ContinueOnError)
	flags.SetOutput(buffer)
	err := flags.Parse(args.Args)

	usage := `
Download the latest binary matching the one being executed, and replace it.

  Usage: gotestfmt update
`

	if err == flag.ErrHelp {
		fmt.Fprintln(args.Output.StdoutWriter, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	out, err := os.Create(args.Binary)

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	defer out.Close()
	resp, err := http.Get(DownloadURL)

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	return 0
}
