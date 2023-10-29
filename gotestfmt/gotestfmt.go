package gotestfmt

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type CommandArgs struct {
	Env        []string
	HomeDir    string
	WorkingDir string
	Binary     string
	Args       []string
	Output     *OutputBuffers
}

func lookupEnvOrDefault(env []string, name string, defaultValue string) string {
	found := -1

	for index, item := range env {
		if strings.HasPrefix(item, name+"=") {
			found = index
			break
		}
	}

	if found == -1 {
		return defaultValue
	}

	parts := strings.Split(env[found], "=")

	if len(parts) == 1 {
		return defaultValue
	}

	return parts[1]
}

func getOptionsDescription(flags *flag.FlagSet) (out string) {
	flags.VisitAll(func(flag *flag.Flag) {
		defaultValue := ""
		ignoreDefaultValue := slices.Contains([]string{"version", "help"}, flag.Name)

		if flag.DefValue != "" && !ignoreDefaultValue {
			defaultValue = fmt.Sprintf(" (default to %v)", flag.DefValue)
		}

		var flagStr string = "--" + flag.Name

		if reflect.TypeOf(flag.Value).String() != "*flag.boolValue" {
			parts := strings.Split(strings.ToUpper(flag.Name), "-")
			flagStr += "=" + parts[len(parts)-1]
		}

		out += fmt.Sprintf("    %-35s%s%s\n", flagStr, flag.Usage, defaultValue)
	})

	return out
}

func Run(args []string, env []string, output OutputBuffers) (exitcode int) {
	var err error
	commands := []string{"update", "version", "run", "download-url"}

	cmd := ""
	if len(args) > 0 && slices.Contains(commands, args[0]) {
		cmd = args[0]
		args = args[1:]
	}

	bin, err := os.Executable()
	binary := lookupEnvOrDefault(env, "GOTESTFMT_BINARY", bin)

	if err != nil {
		fmt.Fprintf(output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	dir, err := os.UserHomeDir()
	homeDir := lookupEnvOrDefault(env, "GOTESTFMT_HOME_DIR", dir)

	if err != nil {
		fmt.Fprintf(output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	dir, err = os.Getwd()
	workingDir := lookupEnvOrDefault(env, "GOTESTFMT_WORKING_DIR", dir)

	if err != nil {
		fmt.Fprintf(output.StderrWriter, "%s %v\n", Color.Fail("ERROR:"), err)
		return 1
	}

	cmdArgs := CommandArgs{
		HomeDir:    homeDir,
		WorkingDir: workingDir,
		Binary:     binary,
		Env:        env,
		Output:     &output,
		Args:       args,
	}

	switch cmd {
	case "run":
		return RunCommand(cmdArgs)

	case "version":
		return VersionCommand(cmdArgs)

	case "download-url":
		return DownloadURLCommand(cmdArgs)

	case "update":
		return UpdateCommand(cmdArgs)

	default:
		return BareCommand(cmdArgs)
	}
}

func loadDotenvFile(environ []string, path string) ([]string, error) {
	contents, err := os.ReadFile(path)

	if err != nil {
		return environ, err
	}

	env := map[string]string{}

	// Convert all existing entries from environ (slice) into a map.
	for _, entry := range environ {
		parsedEntry, _ := godotenv.Unmarshal(entry)

		for key, value := range parsedEntry {
			env[key] = value
		}
	}

	loadedEnv, err := godotenv.Unmarshal(string(contents))

	// Then set all loaded env vars from dotenv file and assign to the same map.
	// This will ensure any overriding value will actually be set.
	for key, value := range loadedEnv {
		env[key] = value

		// Actually set the loaded env vars.
		os.Setenv(key, value)
	}

	environ = []string{}

	// Finally, convert the map into an array, so we can keep using it.
	for key, value := range env {
		environ = append(environ, key+"="+value)
	}

	return environ, nil
}
