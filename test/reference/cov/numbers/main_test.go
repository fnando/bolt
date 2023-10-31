//go:build reference
// +build reference

package numbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	assert.Equal(t, 1, One())
}
