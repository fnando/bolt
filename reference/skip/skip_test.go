// +ignore
package skip

import (
	"testing"
	"time"
)

func TestSkipTestWithMessage(t *testing.T) {
	time.Sleep(10 * time.Millisecond)
	t.Skip("Skipping this test")
}

func TestSkipTestWithoutMessage(t *testing.T) {
	time.Sleep(20 * time.Millisecond)
	t.Skip()
}
