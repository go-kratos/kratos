package logrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/go-kratos/kratos/v2/log"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		level     logrus.Level
		formatter logrus.Formatter
		logLevel  log.Level
		kvs       []interface{}
		want      string
	}{
		"json format": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelInfo,
			kvs:       []interface{}{"case", "json format", "msg", "1"},
			want:      `{"case":"json format","level":"info","msg":"1"`,
		},
		"level unmatch": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelDebug,
			kvs:       []interface{}{"case", "level unmatch", "msg", "1"},
			want:      "",
		},
		"fatal level": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelFatal,
			kvs:       []interface{}{"case", "json format", "msg", "1"},
			want:      `{"case":"json format","level":"fatal","msg":"1"`,
		},
		"no tags": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelInfo,
			kvs:       []interface{}{"msg", "1"},
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
			wrapped := NewLogger(logger)
			_ = wrapped.Log(test.logLevel, test.kvs...)

			if !strings.HasPrefix(output.String(), test.want) {
				t.Errorf("TestName(%s): %s has not prefix %s", name, output.String(), test.want)
			}
		})
	}
}
