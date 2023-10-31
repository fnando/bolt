package bench

import "testing"

func BenchmarkFib1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(1)
	}
}

func BenchmarkFib2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(2)
	}
}

func BenchmarkFib3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(3)
	}
}

func BenchmarkFib4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(4)
	}
}

func BenchmarkFib5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(5)
	}
}

func BenchmarkFib6(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(6)
	}
}

func BenchmarkFib7(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(7)
	}
}

func BenchmarkFib8(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(8)
	}
}

func BenchmarkFib9(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(9)
	}
}

func BenchmarkFib10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10)
	}
}
