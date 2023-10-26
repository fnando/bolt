package dotenvdefault

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestCustomDotenvFile(t *testing.T) {
	assert.Equal(t, "👋", os.Getenv("GOTESTFMT_DOTENV_LOADED"))
}
