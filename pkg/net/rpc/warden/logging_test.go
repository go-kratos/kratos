package warden

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/bilibili/kratos/pkg/log"
)

func Test_logFn(t *testing.T) {
	type args struct {
		code int
		dt   time.Duration
	}
	tests := []struct {
		name string
		args args
		want func(context.Context, ...log.D)
	}{
		{
			name: "ok",
			args: args{code: 0, dt: time.Millisecond},
			want: log.Infov,
		},
		{
			name: "slowlog",
			args: args{code: 0, dt: time.Second},
			want: log.Warnv,
		},
		{
			name: "business error",
			args: args{code: 2233, dt: time.Millisecond},
			want: log.Warnv,
		},
		{
			name: "system error",
			args: args{code: -1, dt: 0},
			want: log.Errorv,
		},
		{
			name: "system error and slowlog",
			args: args{code: -1, dt: time.Second},
			want: log.Errorv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logFn(tt.args.code, tt.args.dt); reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("unexpect log function!")
			}
		})
	}
}

func callInterceptor(err error, interceptor grpc.UnaryClientInterceptor, opts ...grpc.CallOption) {
	interceptor(context.Background(),
		"test-method",
		bytes.NewBufferString("test-req"),
		"test_reply",
		&grpc.ClientConn{},
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return err
		}, opts...)
}

func TestClientLog(t *testing.T) {
	stderr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = stderr
	t.Logf("capture stderr file: %s", stderr.Name())

	t.Run("test no option", func(t *testing.T) {
		callInterceptor(nil, clientLogging())

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-req")
		assert.Contains(t, string(data), "path")
		assert.Contains(t, string(data), "ret")
		assert.Contains(t, string(data), "ts")
		assert.Contains(t, string(data), "grpc-access-log")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test disable args", func(t *testing.T) {
		callInterceptor(nil, clientLogging(WithDialLogFlag(LogFlagDisableArgs)))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.NotContains(t, string(data), "test-req")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test disable args and disable info", func(t *testing.T) {
		callInterceptor(nil, clientLogging(WithDialLogFlag(LogFlagDisableArgs|LogFlagDisableInfo)))
		callInterceptor(errors.New("test-error"), clientLogging(WithDialLogFlag(LogFlagDisableArgs|LogFlagDisableInfo)))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-error")
		assert.NotContains(t, string(data), "INFO")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test call option", func(t *testing.T) {
		callInterceptor(nil, clientLogging(), WithLogFlag(LogFlagDisableArgs))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.NotContains(t, string(data), "test-req")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test combine option", func(t *testing.T) {
		interceptor := clientLogging(WithDialLogFlag(LogFlagDisableInfo))
		callInterceptor(nil, interceptor, WithLogFlag(LogFlagDisableArgs))
		callInterceptor(errors.New("test-error"), interceptor, WithLogFlag(LogFlagDisableArgs))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-error")
		assert.NotContains(t, string(data), "INFO")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})
	t.Run("test no log", func(t *testing.T) {
		callInterceptor(errors.New("test error"), clientLogging(WithDialLogFlag(LogFlagDisable)))
		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Empty(t, data)

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test multi flag", func(t *testing.T) {
		interceptor := clientLogging(WithDialLogFlag(LogFlagDisableInfo | LogFlagDisableArgs))
		callInterceptor(nil, interceptor)
		callInterceptor(errors.New("test-error"), interceptor)

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-error")
		assert.NotContains(t, string(data), "INFO")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})
	os.Stderr = old
}

func callServerInterceptor(err error, interceptor grpc.UnaryServerInterceptor) {
	interceptor(context.Background(),
		bytes.NewBufferString("test-req"),
		&grpc.UnaryServerInfo{
			FullMethod: "test-method",
		},
		func(ctx context.Context, req interface{}) (interface{}, error) { return nil, err })
}

func TestServerLog(t *testing.T) {
	stderr, err := ioutil.TempFile(os.TempDir(), "stderr")
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = stderr
	t.Logf("capture stderr file: %s", stderr.Name())

	t.Run("test no option", func(t *testing.T) {
		callServerInterceptor(nil, serverLogging(0))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-req")
		assert.Contains(t, string(data), "path")
		assert.Contains(t, string(data), "ret")
		assert.Contains(t, string(data), "ts")
		assert.Contains(t, string(data), "grpc-access-log")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test disable args", func(t *testing.T) {
		callServerInterceptor(nil, serverLogging(LogFlagDisableArgs))

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.NotContains(t, string(data), "test-req")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test no log", func(t *testing.T) {
		callServerInterceptor(errors.New("test error"), serverLogging(LogFlagDisable))
		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Empty(t, data)

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})

	t.Run("test multi flag", func(t *testing.T) {
		interceptor := serverLogging(LogFlagDisableInfo | LogFlagDisableArgs)
		callServerInterceptor(nil, interceptor)
		callServerInterceptor(errors.New("test-error"), interceptor)

		stderr.Seek(0, os.SEEK_SET)

		data, err := ioutil.ReadAll(stderr)
		if err != nil {
			t.Error(err)
		}
		assert.Contains(t, string(data), "test-method")
		assert.Contains(t, string(data), "test-error")
		assert.NotContains(t, string(data), "INFO")

		stderr.Seek(0, os.SEEK_SET)
		stderr.Truncate(0)
	})
	os.Stderr = old
}
