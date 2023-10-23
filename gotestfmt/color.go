package gotestfmt

import (
	"os"
)

func Color(color string, text string) string {
	if os.Getenv("NO_COLOR") == "1" {
		return text
	}

	return color + text + "\033[0m"
}
