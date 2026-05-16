package log

import (
	"context"
	"log/slog"
	"testing"
)

type benchmarkHandler struct {
	enabled bool
}

func (h benchmarkHandler) Enabled(context.Context, slog.Level) bool {
	return h.enabled
}

func (h benchmarkHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h benchmarkHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h benchmarkHandler) WithGroup(string) slog.Handler {
	return h
}

func benchmarkWithDefault(b *testing.B, h slog.Handler) {
	b.Helper()
	old := Default()
	SetDefault(slog.New(h))
	b.Cleanup(func() {
		SetDefault(old)
	})
	b.ReportAllocs()
}

func BenchmarkInfoDisabled(b *testing.B) {
	benchmarkWithDefault(b, benchmarkHandler{})
	for i := 0; i < b.N; i++ {
		Info("hello", "kind", "test")
	}
}

func BenchmarkInfoEnabled(b *testing.B) {
	benchmarkWithDefault(b, benchmarkHandler{enabled: true})
	for i := 0; i < b.N; i++ {
		Info("hello", "kind", "test")
	}
}

func BenchmarkLogAttrsEnabled(b *testing.B) {
	benchmarkWithDefault(b, benchmarkHandler{enabled: true})
	attr := slog.String("kind", "test")
	for i := 0; i < b.N; i++ {
		LogAttrs(context.Background(), LevelInfo, "hello", attr)
	}
}

func BenchmarkSlogInfoEnabled(b *testing.B) {
	benchmarkWithDefault(b, benchmarkHandler{enabled: true})
	for i := 0; i < b.N; i++ {
		slog.Info("hello", "kind", "test")
	}
}
