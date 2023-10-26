package dotenvdefault

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestDefaultDotenvFile(t *testing.T) {
	assert.Equal(t, "", os.Getenv("GOTESTFMT_DOTENV_LOADED"))
}
