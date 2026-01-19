package monitoring

import (
	"runtime"
	"testing"
)

// BenchmarkCollectGoroutineCount benchmarks goroutine count collection
func BenchmarkCollectGoroutineCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = runtime.NumGoroutine()
	}
}

// BenchmarkCollectMemStats benchmarks memory stats collection
func BenchmarkCollectMemStats(b *testing.B) {
	var m runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&m)
	}
}

// BenchmarkCollectMemStatsIndividualAccess benchmarks accessing individual memory stats
func BenchmarkCollectMemStatsIndividualAccess(b *testing.B) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Alloc
		_ = m.TotalAlloc
		_ = m.Sys
		_ = m.HeapAlloc
		_ = m.HeapSys
	}
}

// BenchmarkCompleteCollection simulates the full collection process
func BenchmarkCompleteCollection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Goroutine count
		_ = runtime.NumGoroutine()

		// Memory stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		_ = m.Alloc
		_ = m.TotalAlloc
		_ = m.Sys
		_ = m.HeapAlloc
		_ = m.HeapSys
	}
}
