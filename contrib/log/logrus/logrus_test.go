package logrus

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		level     logrus.Level
		formatter logrus.Formatter
		logLevel  slog.Level
		msg       string
		kvs       []any
		want      string
	}{
		"json format": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  slog.LevelInfo,
			msg:       "1",
			kvs:       []any{"case", "json format"},
			want:      `{"case":"json format","level":"info","msg":"1"`,
		},
		"level unmatch": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  slog.LevelDebug,
			msg:       "1",
			kvs:       []any{"case", "level unmatch"},
			want:      "",
		},
		"fatal level": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  slog.LevelError + 4,
			msg:       "1",
			kvs:       []any{"case", "json format"},
			want:      `{"case":"json format","level":"fatal","msg":"1"`,
		},
		"no tags": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  slog.LevelInfo,
			msg:       "1",
			want:      `{"level":"info","msg":"1"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := new(bytes.Buffer)
			logger := logrus.New()
			logger.Level = test.level
			logger.Out = output
			logger.Formatter = test.formatter
			logger.ExitFunc = func(int) {}
			wrapped := NewLogger(logger)
			wrapped.Log(context.Background(), test.logLevel, test.msg, test.kvs...)

			if !strings.HasPrefix(output.String(), test.want) {
				t.Errorf("TestName(%s): %s has not prefix %s", name, output.String(), test.want)
			}
		})
	}
}
