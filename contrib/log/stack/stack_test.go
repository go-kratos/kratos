package stack

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func redirectStdout(fn func()) ([]byte, error) {
	fs, err := os.CreateTemp("", "stdout")
	if err != nil {
		return nil, err
	}
	defer os.Remove(fs.Name())
	defer fs.Close()

	org := os.Stdout
	os.Stdout = fs
	fn()
	os.Stdout = org

	return os.ReadFile(fs.Name())
}

type testLogger struct {
	name string
}

func (t *testLogger) Log(level log.Level, keyvals ...interface{}) error {
	fmt.Printf("name: %s, level: %v, keyvals: %v\n", t.name, level, keyvals)

	return nil
}

type errorLogger struct{}

func (t *errorLogger) Log(level log.Level, keyvals ...interface{}) error {
	return fmt.Errorf("error")
}

func TestStackLogger(t *testing.T) {
	logger := New([]log.Logger{
		&testLogger{"log1"},
		&testLogger{"log2"},
	})

	tests := []struct {
		name  string
		level log.Level
		want  string
	}{
		{
			name:  "debug",
			level: log.LevelDebug,
			want:  "name: log1, level: DEBUG, keyvals: [test test]\nname: log2, level: DEBUG, keyvals: [test test]\n",
		}, {
			name:  "info",
			level: log.LevelInfo,
			want:  "name: log1, level: INFO, keyvals: [test test]\nname: log2, level: INFO, keyvals: [test test]\n",
		}, {
			name:  "warn",
			level: log.LevelWarn,
			want:  "name: log1, level: WARN, keyvals: [test test]\nname: log2, level: WARN, keyvals: [test test]\n",
		},
		{
			name:  "error",
			level: log.LevelError,
			want:  "name: log1, level: ERROR, keyvals: [test test]\nname: log2, level: ERROR, keyvals: [test test]\n",
		},
		{
			name:  "fatal",
			level: log.LevelFatal,
			want:  "name: log1, level: FATAL, keyvals: [test test]\nname: log2, level: FATAL, keyvals: [test test]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := redirectStdout(func() {
				logger.Log(tt.level, "test", "test")
			})
			if err != nil {
				t.Fatal(err)
			}

			if string(buf) != tt.want {
				t.Fatal("unexpected output")
			}
		})
	}
}

func TestIgnoreExceptions(t *testing.T) {
	tests := []struct {
		name             string
		ignoreExceptions bool
		want             string
	}{
		{
			name:             "false",
			ignoreExceptions: false,
			want:             "",
		},
		{
			name:             "true",
			ignoreExceptions: true,
			want:             "name: log, level: DEBUG, keyvals: [test test]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []Option

			if tt.ignoreExceptions {
				opts = append(opts, IgnoreExceptions())
			}

			logger := New([]log.Logger{
				&errorLogger{},
				&testLogger{"log"},
			}, opts...)

			buf, err := redirectStdout(func() {
				logger.Log(log.LevelDebug, "test", "test")
			})
			if err != nil {
				t.Fatal(err)
			}

			if string(buf) != tt.want {
				t.Fatal("unexpected output")
			}
		})
	}
}
