//go:build reference
// +build reference

package pass

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEqualStringPass(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, "hello", "hello")
}

func TestEqualNumberPass(t *testing.T) {
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, 1)
}
