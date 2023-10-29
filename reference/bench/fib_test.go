package bench

import "testing"

func BenchmarkFib10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10)
	}
}

func BenchmarkFib5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(5)
	}
}
