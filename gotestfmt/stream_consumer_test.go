package gotestfmt

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func createFileStreamConsumer(path string, setup func(consumer *StreamConsumer)) *StreamConsumer {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	output := createBuffers()

	consumer := CreateStreamConsumer(CreateStreamConsumerOptions{
		WorkingDir: "/home/test/gotestfmt",
		HomeDir:    "/home/test",
		Scanner:    scanner,
		Output:     &output,
	})
	setup(consumer)
	consumer.Run()

	return consumer
}

func TestConsumePassStreaming(t *testing.T) {
	consumer := createFileStreamConsumer("./test/replays/passed.txt", func(consumer *StreamConsumer) {})

	test := consumer.Tests["github.com/fnando/gotestfmt/reference/pass:TestEqualNumberPass"]

	assert.NotNil(t, test)
	assert.Equal(t, "TestEqualNumberPass", test.Name)
	assert.Equal(t, "Equal Number Pass", test.ReadableName)
	assert.Equal(t, TestStatus.Pass, test.Status)
	assert.Equal(t, 0.02, test.Elapsed)
	assert.Equal(t, "2023-10-27T01:26:18Z", test.StartedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, "2023-10-27T01:26:18Z", test.EndedAt.UTC().Format(time.RFC3339))
	assert.Empty(t, test.Output)
}

func TestConsumeSkipStreaming(t *testing.T) {
	consumer := createFileStreamConsumer("./test/replays/skipped.txt", func(consumer *StreamConsumer) {})

	t.Run("WithMessage", func(t *testing.T) {
		test := consumer.Tests["github.com/fnando/gotestfmt/reference/skip:TestSkipTestWithMessage"]

		assert.Equal(t, "TestSkipTestWithMessage", test.Name)
		assert.Equal(t, "Skip Test With Message", test.ReadableName)
		assert.Equal(t, "reference/skip/skip_test.go:10", test.TestSource)
		assert.Equal(t, "Skipping this test", test.SkipMessage)
		assert.Equal(t, TestStatus.Skip, test.Status)
		assert.Equal(t, 0.01, test.Elapsed)
		assert.Equal(t, "2023-10-26T21:52:15Z", test.StartedAt.UTC().Format(time.RFC3339))
		assert.Equal(t, "2023-10-26T21:52:15Z", test.EndedAt.UTC().Format(time.RFC3339))
		assert.Empty(t, test.Output)
	})

	t.Run("WithoutMessage", func(t *testing.T) {
		test := consumer.Tests["github.com/fnando/gotestfmt/reference/skip:TestSkipTestWithoutMessage"]

		assert.Equal(t, "TestSkipTestWithoutMessage", test.Name)
		assert.Equal(t, "Skip Test Without Message", test.ReadableName)
		assert.Equal(t, "reference/skip/skip_test.go:15", test.TestSource)
		assert.Equal(t, "Skipped", test.SkipMessage)
		assert.Equal(t, TestStatus.Skip, test.Status)
		assert.Equal(t, 0.02, test.Elapsed)
		assert.Equal(t, "2023-10-26T21:52:15Z", test.StartedAt.UTC().Format(time.RFC3339))
		assert.Equal(t, "2023-10-26T21:52:15Z", test.EndedAt.UTC().Format(time.RFC3339))
		assert.Empty(t, test.Output)
	})
}

func TestNotifyWhenTestFinishes(t *testing.T) {
	want := "github.com/fnando/gotestfmt/reference/pass:TestEqualNumberPass"
	var got string

	createFileStreamConsumer("./test/replays/passed.txt", func(c *StreamConsumer) {
		c.OnNotifyTestFinish = func(tst *Test) {
			if tst.Key == want {
				got = tst.Key
			}
		}
	})

	assert.Equal(t, want, got)
}

func TestNotifyWhenStreamFinishes(t *testing.T) {
	var got bool

	createFileStreamConsumer("./test/replays/passed.txt", func(c *StreamConsumer) {
		c.OnFinish = func(tests []*Test, coverage []Coverage, benchmarks []*Benchmark) {
			got = true
		}
	})

	assert.True(t, got)
}

func TestConsumeFailStreaming(t *testing.T) {
	consumer := createFileStreamConsumer("./test/replays/failed.txt", func(consumer *StreamConsumer) {})

	test := consumer.Tests["github.com/fnando/gotestfmt/reference/fail:TestEqualNumberFail"]

	assert.NotNil(t, test)

	want := `Error:  Not equal:
        expected: 1
        actual  : 2`

	got := strings.Join(test.Output, "\n")

	assert.Equal(t, "TestEqualNumberFail", test.Name)
	assert.Equal(t, "Equal Number Fail", test.ReadableName)
	assert.Equal(t, TestStatus.Fail, test.Status)
	assert.Equal(t, 0.01, test.Elapsed)
	assert.Equal(t, "reference/failed/failed_test.go:16", test.ErrorTrace)
	assert.Equal(t, "", test.TestSource)
	assert.Equal(t, "2023-10-26T21:41:41Z", test.StartedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, "2023-10-26T21:41:41Z", test.EndedAt.UTC().Format(time.RFC3339))
	assert.NotEmpty(t, test.Output)
	assert.Equal(t, want, got)
}

func TestConsumeFailStreamingThroughIndirectErrorTrace(t *testing.T) {
	consumer := createFileStreamConsumer("./test/replays/failed.txt", func(consumer *StreamConsumer) {})

	test := consumer.Tests["github.com/fnando/gotestfmt/reference/fail:TestFailedThroughHelper"]

	assert.NotNil(t, test)

	want := `Error:  Not equal:
        expected: 1
        actual  : 2`

	got := strings.Join(test.Output, "\n")

	assert.Equal(t, "TestFailedThroughHelper", test.Name)
	assert.Equal(t, "Failed Through Helper", test.ReadableName)
	assert.Equal(t, TestStatus.Fail, test.Status)
	assert.Equal(t, 0.02, test.Elapsed)
	assert.Equal(t, "reference/failed/failed_test.go:11", test.ErrorTrace)
	assert.Equal(t, "reference/failed/failed_test.go:21", test.TestSource)
	assert.Equal(t, "2023-10-26T21:41:41Z", test.StartedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, "2023-10-26T21:41:41Z", test.EndedAt.UTC().Format(time.RFC3339))
	assert.Equal(t, want, got)
}
