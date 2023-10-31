package commands

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"

	c "github.com/fnando/bolt/common"
)

func Version(args []string, output *c.Output) int {
	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("bolt version", flag.ContinueOnError)
	flags.SetOutput(buffer)
	err := flags.Parse(args)

	usage := `
Show the version and commit hash.

  Usage: bolt version
`

	if err == flag.ErrHelp {
		fmt.Fprintln(output.Stdout, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(output.Stderr, "%s %v\n", c.Color.Fail("ERROR:"), err)
		return 1
	}

	fmt.Fprintf(output.Stdout, "bolt %s (%s)\n", c.Version, c.Commit)

	return 0
}
