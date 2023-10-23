package equal

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	g "github.com/fnando/gotestfmt/gotestfmt"
	"github.com/stretchr/testify/assert"
)

func TestEqualStructFail(t *testing.T) {
	assert.Equal(t, map[string]any{"a": 1, "b": 2, "c": 3}, map[string]any{"a": 1, "b": 3, "c": 2})
}

func TestEqualNumberFail(t *testing.T) {
	assert.Equal(t, 1, 2)
}

func TestEqualStringFail(t *testing.T) {
	assert.Equal(t, "hello", "hi")
}

func TestEqualBoolFail(t *testing.T) {
	assert.Equal(t, true, false)
}

func TestEqualErrorFail(t *testing.T) {
	assert.Equal(t, errors.New("oh noes"), errors.New("derp"))
}

func TestContainStringFail(t *testing.T) {
	assert.Contains(t, "hello", "bye")
}

func TestContainSliceFail(t *testing.T) {
	assert.Contains(t, []string{"hello"}, "bye")
}

func TestDirExistsFail(t *testing.T) {
	assert.DirExists(t, "/invalid/path")
}

func TestElementsMatchFail(t *testing.T) {
	assert.ElementsMatch(t, []int{1, 2}, []int{1, 3})
}

func TestEmptyStringFail(t *testing.T) {
	assert.Empty(t, "hello")
}

func TestEmptyListFail(t *testing.T) {
	assert.Empty(t, []int{1, 2, 3})
}

func TestAnotherEqualErrorFail(t *testing.T) {
	assert.EqualError(t, errors.New("oh noes"), "derp")
}

func TestExactlyIntFail(t *testing.T) {
	assert.Exactly(t, int16(1), int32(1))
}

func TestFail(t *testing.T) {
	assert.Fail(t, "Failing this test")
}

func TestFailNow(t *testing.T) {
	assert.FailNow(t, "Failing this test now")
}

func TestFalseFail(t *testing.T) {
	assert.False(t, true)
}

func TestFileExistsFail(t *testing.T) {
	assert.FileExists(t, "/invalid/file.txt")
}

func TestGreaterNumberFail(t *testing.T) {
	assert.Greater(t, 1, 2)
}

func TestGreaterStringFail(t *testing.T) {
	assert.Greater(t, "a", "b")
}

func TestGreaterOrEqualNumberFail(t *testing.T) {
	assert.GreaterOrEqual(t, 1, 2)
}

func TestGreaterOrEqualStringFail(t *testing.T) {
	assert.GreaterOrEqual(t, "a", "b")
}

func TestHTTPBodyContainsFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPBodyContains(t, handler, "GET", "https://example.com/", nil, "Hello")
}

func TestHTTPBodyNotContainsFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPBodyNotContains(t, handler, "GET", "https://example.com/", nil, "👋")
}

func TestHTTPErrorFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPError(t, handler, "GET", "https://example.com/", nil)
}

func TestHTTPRedirectFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPRedirect(t, handler, "GET", "https://example.com/", nil)
}

func TestHTTPStatusCodeFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPStatusCode(t, handler, "GET", "https://example.com/", nil, 404)
}

func TestHTTPSuccessFail(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = fmt.Fprintf(w, "👋")
	}
	assert.HTTPSuccess(t, handler, "GET", "https://example.com/", nil)
}

func TestImplementsFail(t *testing.T) {
	assert.Implements(t, (*g.Reporter)(nil), []string{})
}

func TestInDeltaFail(t *testing.T) {
	assert.InDelta(t, 2.0, 1.0, 0.01)
}

func TestIsDecreasingFail(t *testing.T) {
	assert.IsDecreasing(t, []int{1, 2, 3})
}

func TestIsIncreasingFail(t *testing.T) {
	assert.IsIncreasing(t, []int{3, 2, 1})
}

func TestIsTypeFail(t *testing.T) {
	assert.IsType(t, (*int)(nil), "hello")
}
