package t1ha

import (
	"strconv"
	"testing"
)

var buf [8192]byte

func BenchmarkHash(b *testing.B) {
	var sizes = []int{1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1024, 8192}
	for _, n := range sizes {
		b.Run(strconv.Itoa(n), func(b *testing.B) { benchmarkHashn(b, int64(n), Sum64) })
	}
}

var total uint64

func benchmarkHashn(b *testing.B, size int64, h func([]byte, uint64) uint64) {
	b.SetBytes(size)
	for i := 0; i < b.N; i++ {
		total += h(buf[:size], 0)
	}
}
