package gotestfmt

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
)

func BareCommand(args CommandArgs) (exitcode int) {
	var (
		version bool
	)

	var b bytes.Buffer
	buffer := bufio.NewWriter(&b)

	flags := flag.NewFlagSet("gotestfmt", flag.ContinueOnError)
	flags.SetOutput(buffer)
	flags.Usage = func() {}
	flags.BoolVar(&version, "version", false, "Show version")
	err := flags.Parse(args.Args)

	usage := fmt.Sprintf(`
gotestfmt is a golang test runner that has a nicer output.

  Usage: gotestfmt [command] [options]

  Options:
%s

  Commands:

    gotestfmt version                  Show gotestfmt version
    gotestfmt download-url             Output the latest binary download url
    gotestfmt run                      Run tests
    gotestfmt [command] --help         Display help on [command]


  Further information:
    https://github.com/fnando/gotestfmt
`, getOptionsDescription(flags))

	if err == flag.ErrHelp {
		fmt.Fprintln(args.Output.StdoutWriter, usage)
		return 0
	}

	if err != nil {
		fmt.Fprintf(args.Output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	if version {
		fmt.Fprintln(args.Output.StdoutWriter, Version)
		return 0
	}

	fmt.Fprintln(args.Output.StdoutWriter, usage)
	return 1
}
