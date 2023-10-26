package bye

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestBye(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, "Bye, John!", Bye("John"))
}
