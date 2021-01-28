package ulog

import (
	"testing"
)

func withProductionLogger(b *testing.B, f func(logger Logger)) {
	loggerConfig := pkgConfig{
		DefaultLevel: "info",
		Console:      false,
		Directory:    "/tmp",
		Name:         "myApp",
	}
	pkgCfg = loggerConfig
	logger := NewLogger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func BenchmarkProduction(b *testing.B) {
	withProductionLogger(b, func(log Logger) {
		log.Infow("message with kvs", "k1", "v1", "k2", "v2")
	})
}

func withDevLogger(b *testing.B, f func(logger Logger)) {
	logger := NewLogger()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func BenchmarkDevelopment(b *testing.B) {
	withDevLogger(b, func(log Logger) {
		log.Infow("message with kvs", "k1", "v1", "k2", "v2")
	})
}
