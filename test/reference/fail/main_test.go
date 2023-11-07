//go:build reference
// +build reference

package fail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func failHelper(t *testing.T) {
	assert.Equal(t, 1, 2)
}

func TestEqualNumberFail(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, 2)
}

func TestFailedThroughHelper(t *testing.T) {
	time.Sleep(20 * time.Millisecond)
	failHelper(t)
}

func TestEqualStructFail(t *testing.T) {
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, map[string]any{"a": 1, "b": 2, "c": 3}, map[string]any{"a": 1, "b": 3, "c": 2})
}
