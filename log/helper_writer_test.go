package log

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestWriterWrapper(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStdLogger(&buf)
	content := "ThisIsSomeTestLogMessage"
	testCases := []struct {
		w                io.Writer
		acceptLevel      Level
		acceptMessageKey string
	}{
		{
			w:                NewWriter(logger),
			acceptLevel:      LevelInfo, // default level
			acceptMessageKey: DefaultMessageKey,
		},
		{
			w:                NewWriter(logger, WithWriterLevel(LevelDebug)),
			acceptLevel:      LevelDebug,
			acceptMessageKey: DefaultMessageKey,
		},
		{
			w:                NewWriter(logger, WithWriteMessageKey("XxXxX")),
			acceptLevel:      LevelInfo, // default level
			acceptMessageKey: "XxXxX",
		},
		{
			w:                NewWriter(logger, WithWriterLevel(LevelError), WithWriteMessageKey("XxXxX")),
			acceptLevel:      LevelError,
			acceptMessageKey: "XxXxX",
		},
	}
	for _, tc := range testCases {
		_, err := tc.w.Write([]byte(content))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), tc.acceptLevel.String()) {
			t.Errorf("expected level: %s, got: %s", tc.acceptLevel, buf.String())
		}
		if !strings.Contains(buf.String(), tc.acceptMessageKey) {
			t.Errorf("expected message key: %s, got: %s", tc.acceptMessageKey, buf.String())
		}
	}
}
