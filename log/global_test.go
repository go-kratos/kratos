package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestGlobalLog(t *testing.T) {
	defaultLogger := GetLogger()
	t.Cleanup(func() { SetLogger(defaultLogger) })

	buffer := &bytes.Buffer{}
	logger := NewStdLogger(buffer)
	SetLogger(logger)

	if GetLogger() != logger {
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

func TestGlobalContext(t *testing.T) {
	defaultLogger := GetLogger()
	t.Cleanup(func() { SetLogger(defaultLogger) })

	buffer := &bytes.Buffer{}
	SetLogger(NewStdLogger(buffer))
	Context(context.Background()).Info("111")
	if buffer.String() != "INFO msg=111\n" {
		t.Errorf("Expected:%s, got:%s", "INFO msg=111", buffer.String())
	}
}

type traceIdKey struct{}

func TestValuerUnderGlobalValue(t *testing.T) {
	defaultLogger := GetLogger()
	t.Cleanup(func() { SetLogger(defaultLogger) })

	var traceIdValuer Valuer = func(ctx context.Context) any {
		return ctx.Value(traceIdKey{})
	}

	var buf bytes.Buffer
	l1 := NewStdLogger(&buf)
	l2 := With(l1, "traceId", traceIdValuer)

	SetLogger(l2)
	l3 := GetLogger()

	ctx := context.WithValue(context.Background(), traceIdKey{}, "123")
	l4 := WithContext(ctx, l3)
	l4.Log(LevelInfo, "msg", "m")

	want := "INFO traceId=123 msg=m\n"
	if got := buf.String(); got != want {
		t.Errorf("Expected:%q, got:%q", want, got)
	}
}
