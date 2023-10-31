package commands

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

func getFlagsUsage(flags *flag.FlagSet) (out string) {
	flags.VisitAll(func(flag *flag.Flag) {
		if flag.Usage == "" {
			return
		}

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
