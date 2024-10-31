package log

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGlobalLog(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := NewStdLogger(buffer)
	SetLogger(logger)

	if global.Logger != logger {
		t.Error("GetLogger() is not equal to logger")
	}

	testCases := []struct {
		level   Level
		content []interface{}
	}{
		{
			LevelDebug,
			[]interface{}{"test debug"},
		},
		{
			LevelInfo,
			[]interface{}{"test info"},
		},
		{
			LevelInfo,
			[]interface{}{"test %s", "info"},
		},
		{
			LevelWarn,
			[]interface{}{"test warn"},
		},
		{
			LevelError,
			[]interface{}{"test error"},
		},
		{
			LevelError,
			[]interface{}{"test %s", "error"},
		},
	}

	var expected []string
	for _, tc := range testCases {
		msg := fmt.Sprintf(tc.content[0].(string), tc.content[1:]...)
		switch tc.level {
		case LevelDebug:
			Debug(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "DEBUG", msg))
			Debugf(tc.content[0].(string), tc.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "DEBUG", msg))
			Debugw("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "DEBUG", msg))
		case LevelInfo:
			Info(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			Infof(tc.content[0].(string), tc.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", msg))
			Infow("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "INFO", msg))
		case LevelWarn:
			Warn(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			Warnf(tc.content[0].(string), tc.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "WARN", msg))
			Warnw("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "WARN", msg))
		case LevelError:
			Error(msg)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			Errorf(tc.content[0].(string), tc.content[1:]...)
			expected = append(expected, fmt.Sprintf("%s msg=%s", "ERROR", msg))
			Errorw("log", msg)
			expected = append(expected, fmt.Sprintf("%s log=%s", "ERROR", msg))
		}
	}
	Log(LevelInfo, DefaultMessageKey, "test log")
	expected = append(expected, fmt.Sprintf("%s msg=%s", "INFO", "test log"))

	expected = append(expected, "")

	t.Logf("Content: %s", buffer.String())
	if buffer.String() != strings.Join(expected, "\n") {
		t.Errorf("Expected: %s, got: %s", strings.Join(expected, "\n"), buffer.String())
	}
}

func TestGlobalLogUpdate(t *testing.T) {
	l := &loggerAppliance{}
	l.SetLogger(NewStdLogger(os.Stdout))
	LOG := NewHelper(l)
	LOG.Info("Log to stdout")

	buffer := &bytes.Buffer{}
	l.SetLogger(NewStdLogger(buffer))
	LOG.Info("Log to buffer")

	expected := "INFO msg=Log to buffer\n"
	if buffer.String() != expected {
		t.Errorf("Expected: %s, got: %s", expected, buffer.String())
	}
}

func TestGlobalContext(t *testing.T) {
	buffer := &bytes.Buffer{}
	SetLogger(NewStdLogger(buffer))
	Context(context.Background()).Infof("111")
	if buffer.String() != "INFO msg=111\n" {
		t.Errorf("Expected:%s, got:%s", "INFO msg=111", buffer.String())
	}
}

func TestContextWithGlobalLog(t *testing.T) {
	buffer := &bytes.Buffer{}

	type traceKey struct{}
	// set "trace-id" Valuer
	newLogger := With(NewStdLogger(buffer), "trace-id", Valuer(func(ctx context.Context) interface{} {
		return ctx.Value(traceKey{})
	}))

	SetLogger(newLogger)

	// add value to ctx
	ctx := context.WithValue(context.Background(), traceKey{}, "test-trace-id")

	_ = WithContext(ctx, GetLogger()).Log(LevelInfo)
	if buffer.String() != "INFO trace-id=test-trace-id\n" {
		t.Errorf("Expected:%s, got:%s", "INFO trace-id=test-trace-id", buffer.String())
	}
}
