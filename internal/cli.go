package internal

import (
	"fmt"
	"os"

	"github.com/fnando/bolt/common"
	"github.com/fnando/bolt/internal/commands"
	"golang.org/x/exp/slices"
)

var availableCommands = []string{"run", "update", "version"}

var usage string = `
bolt is a golang test runner that has a nicer output.

  Usage: bolt [command] [options]

  Commands:

    bolt version                  Show bolt version
    bolt run                      Run tests
    bolt update                   Update to the latest released version
    bolt [command] --help         Display help on [command]


  Further information:
    https://github.com/fnando/bolt
`

func Run(workingDir string, homeDir string, output common.Output) (exitcode int) {
	args := os.Args[1:]
	cmd := ""

	common.Color.Disabled = slices.Contains(args, "-no-color") ||
		slices.Contains(args, "--no-color") ||
		os.Getenv("NO_COLOR") == "1"

	if len(args) >= 1 && slices.Contains(availableCommands, args[0]) {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "update":
		return commands.Update(args, &output)

	case "version":
		return commands.Version(args, &output)

	case "run":
		return commands.Run(
			args,
			commands.RunArgs{HomeDir: homeDir, WorkingDir: workingDir},
			&output,
		)

	default:
		fmt.Fprint(output.Stdout, usage)
		return 1
	}
}
