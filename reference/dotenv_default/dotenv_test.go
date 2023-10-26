package dotenvdefault

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestDefaultDotenvFile(t *testing.T) {
	assert.Equal(t, "1", os.Getenv("GOTESTFMT_DOTENV_LOADED"))
}
