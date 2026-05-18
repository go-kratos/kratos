package log

import (
	"log/slog"
	"testing"
)

func TestLevelAliases(t *testing.T) {
	if LevelDebug != slog.LevelDebug {
		t.Fatalf("LevelDebug = %v, want %v", LevelDebug, slog.LevelDebug)
	}
	if LevelFatal != slog.LevelError+4 {
		t.Fatalf("LevelFatal = %v, want %v", LevelFatal, slog.LevelError+4)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Level
	}{
		{name: "debug", in: "debug", want: LevelDebug},
		{name: "info", in: "info", want: LevelInfo},
		{name: "warn", in: "warn", want: LevelWarn},
		{name: "error", in: "error", want: LevelError},
		{name: "fatal", in: "fatal", want: LevelFatal},
		{name: "custom", in: "INFO+1", want: LevelInfo + 1},
		{name: "default", in: "unknown", want: LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLevel(tt.in); got != tt.want {
				t.Fatalf("ParseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
