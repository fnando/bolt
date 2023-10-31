package commands

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	c "github.com/fnando/bolt/common"
)

func Update(args []string, output *c.Output) int {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("bolt update", flag.ContinueOnError)
	flags.SetOutput(buffer)
	err := flags.Parse(args)

	usage := `
Download the latest binary matching the one being executed, and replace it.

  Usage: bolt update
`

	if err == flag.ErrHelp {
		fmt.Fprintln(output.Stdout, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	binary, _ := os.Executable()
	out, err := os.Create(binary)

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	defer out.Close()
	resp, err := http.Get(c.DownloadURL)

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	return 0
}
