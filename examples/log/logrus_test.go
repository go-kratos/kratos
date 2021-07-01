package log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		level     logrus.Level
		formatter logrus.Formatter
		logLevel  log.Level
		kvs       []interface{}
		want      string
	}{
		"json format":   {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelInfo, kvs: []interface{}{"case", "json format", "msg", "1"}, want: `{"case":"json format","level":"info","msg":"1"`},
		"level unmatch": {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelDebug, kvs: []interface{}{"case", "level unmatch", "msg", "1"}, want: ""},
		"no tags":       {level: logrus.InfoLevel, formatter: &logrus.JSONFormatter{}, logLevel: log.LevelInfo, kvs: []interface{}{"msg", "1"}, want: `{"level":"info","msg":"1"`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := new(bytes.Buffer)
			logger := NewLogrusLogger(Level(test.level), Formatter(test.formatter), Output(output))
			_ = logger.Log(test.logLevel, test.kvs...)

			assert.True(t, strings.HasPrefix(output.String(), test.want))
		})
	}
}
