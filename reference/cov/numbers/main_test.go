package numbers

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestOne(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, One())
}

func TestTwo(t *testing.T) {
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 2, Two())
}

func TestThree(t *testing.T) {
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 3, Three())
}

func TestFour(t *testing.T) {
	time.Sleep(40 * time.Millisecond)
	assert.Equal(t, 4, Four())
}
