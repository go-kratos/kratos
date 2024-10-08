package logging

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Transporter = (*Transport)(nil)

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return nil
}

func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

func TestHTTP(t *testing.T) {
	err := errors.New("reply.error")
	bf := bytes.NewBuffer(nil)
	logger := log.NewStdLogger(bf)

	tests := []struct {
		name string
		kind func(logger log.Logger) middleware.Middleware
		err  error
		ctx  context.Context
	}{
		{
			"http-server@fail",
			Server,
			err,
			func() context.Context {
				return transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"})
			}(),
		},
		{
			"http-server@succ",
			Server,
			nil,
			func() context.Context {
				return transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"})
			}(),
		},
		{
			"http-client@succ",
			Client,
			nil,
			func() context.Context {
				return transport.NewClientContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"})
			}(),
		},
		{
			"http-client@fail",
			Client,
			err,
			func() context.Context {
				return transport.NewClientContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"})
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bf.Reset()
			next := func(context.Context, interface{}) (interface{}, error) {
				return "reply", test.err
			}
			next = test.kind(logger)(next)
			v, e := next(test.ctx, "req.args")
			t.Logf("[%s]reply: %v, error: %v", test.name, v, e)
			t.Logf("[%s]log:%s", test.name, bf.String())
		})
	}
}

type (
	dummy struct {
		field string
	}
	dummyStringer struct {
		field string
	}
	dummyStringerRedacter struct {
		field string
	}
)

func (d *dummyStringer) String() string {
	return "my value"
}

func (d *dummyStringerRedacter) String() string {
	return "my value"
}

func (d *dummyStringerRedacter) Redact() string {
	return "my value redacted"
}

func TestExtractArgs(t *testing.T) {
	tests := []struct {
		name     string
		req      interface{}
		expected string
	}{
		{
			name:     "dummyStringer",
			req:      &dummyStringer{field: ""},
			expected: "my value",
		}, {
			name:     "dummy",
			req:      &dummy{field: "value"},
			expected: "&{field:value}",
		}, {
			name:     "dummyStringerRedacter",
			req:      &dummyStringerRedacter{field: ""},
			expected: "my value redacted",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if value := extractArgs(test.req); value != test.expected {
				t.Errorf(`The stringified %s structure must be equal to "%s", %v given`, test.name, test.expected, value)
			}
		})
	}
}

func TestExtractError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantLevel  log.Level
		wantErrStr string
	}{
		{
			"no error", nil, log.LevelInfo, "",
		},
		{
			"error", errors.New("test error"), log.LevelError, "test error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			level, errStr := extractError(test.err)
			if level != test.wantLevel {
				t.Errorf("want: %d, got: %d", test.wantLevel, level)
			}
			if errStr != test.wantErrStr {
				t.Errorf("want: %s, got: %s", test.wantErrStr, errStr)
			}
		})
	}
}

type extractKeyValues [][]any

func (l *extractKeyValues) Log(_ log.Level, kv ...any) error { *l = append(*l, kv); return nil }

func TestServer_CallerPath(t *testing.T) {
	var a extractKeyValues
	logger := log.With(&a, "caller", log.Caller(5)) // report where the helper was called

	// make sure the caller is same
	sameCaller := func(fn middleware.Handler) { _, _ = fn(context.Background(), nil) }

	// caller: [... log inside middleware, fn(context.Background(), nil)]
	h := func(context.Context, any) (a any, e error) { return }
	h = Server(logger)(h)
	sameCaller(h)

	// caller: [... helper.Info("foo"), fn(context.Background(), nil)]
	helper := log.NewHelper(logger)
	sameCaller(func(context.Context, any) (a any, e error) { helper.Info("foo"); return })

	t.Log(a[0])
	t.Log(a[1])
	if a[0][0] != "caller" || a[1][0] != "caller" {
		t.Fatal("caller not found")
	}
	if a[0][1] != a[1][1] {
		t.Fatalf("middleware should have the same caller as log.Helper. middleware: %s, helper: %s", a[0][1], a[1][1])
	}
}

func TestClient_CallerPath(t *testing.T) {
	var a extractKeyValues
	logger := log.With(&a, "caller", log.Caller(5)) // report where the helper was called

	// make sure the caller is same
	sameCaller := func(fn middleware.Handler) { _, _ = fn(context.Background(), nil) }

	// caller: [... log inside middleware, fn(context.Background(), nil)]
	h := func(context.Context, any) (a any, e error) { return }
	h = Client(logger)(h)
	sameCaller(h)

	// caller: [... helper.Info("foo"), fn(context.Background(), nil)]
	helper := log.NewHelper(logger)
	sameCaller(func(context.Context, any) (a any, e error) { helper.Info("foo"); return })

	t.Log(a[0])
	t.Log(a[1])
	if a[0][0] != "caller" || a[1][0] != "caller" {
		t.Fatal("caller not found")
	}
	if a[0][1] != a[1][1] {
		t.Fatalf("middleware should have the same caller as log.Helper. middleware: %s, helper: %s", a[0][1], a[1][1])
	}
}
