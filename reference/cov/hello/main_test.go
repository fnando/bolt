package hello

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, "Hello, John!", Hello("John"))
}
