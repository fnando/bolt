//go:build reference
// +build reference

package letters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestA(t *testing.T) {
	assert.Equal(t, "A", A())
}

func TestB(t *testing.T) {
	assert.Equal(t, "B", B())
}
