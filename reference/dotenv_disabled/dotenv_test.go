package dotenvdefault

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestDefaultDotenvFile(t *testing.T) {
	value, exists := os.LookupEnv("GOTESTFMT_DOTENV_LOADED")
	assert.False(t, exists)
	assert.Equal(t, "", value)
}
