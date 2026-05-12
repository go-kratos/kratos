package logging

import (
	"context"
	"errors"
	"log/slog"
	"testing"

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
	handler := &captureHandler{}
	logger := slog.New(handler)

	tests := []struct {
		name string
		kind func(*slog.Logger) middleware.Middleware
		err  error
		ctx  context.Context
		want slog.Level
	}{
		{
			name: "http-server@fail",
			kind: Server,
			err:  err,
			ctx:  transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"}),
			want: slog.LevelError,
		},
		{
			name: "http-server@succ",
			kind: Server,
			ctx:  transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"}),
			want: slog.LevelInfo,
		},
		{
			name: "http-client@succ",
			kind: Client,
			ctx:  transport.NewClientContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"}),
			want: slog.LevelInfo,
		},
		{
			name: "http-client@fail",
			kind: Client,
			err:  err,
			ctx:  transport.NewClientContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/package.service/method"}),
			want: slog.LevelError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.reset()
			next := func(context.Context, any) (any, error) {
				return "reply", test.err
			}
			next = test.kind(logger)(next)
			reply, gotErr := next(test.ctx, "req.args")
			if reply != "reply" {
				t.Fatalf("reply = %v, want %q", reply, "reply")
			}
			if gotErr != test.err {
				t.Fatalf("err = %v, want %v", gotErr, test.err)
			}
			if len(handler.records) != 1 {
				t.Fatalf("records len = %d, want 1", len(handler.records))
			}
			if handler.records[0].Level != test.want {
				t.Fatalf("level = %v, want %v", handler.records[0].Level, test.want)
			}
			if got := handler.attrs[0]["component"]; got != "http" {
				t.Fatalf("component = %v, want %q", got, "http")
			}
			if got := handler.attrs[0]["operation"]; got != "/package.service/method" {
				t.Fatalf("operation = %v, want %q", got, "/package.service/method")
			}
			if got := handler.attrs[0]["args"]; got != "req.args" {
				t.Fatalf("args = %v, want %q", got, "req.args")
			}
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
		req      any
		expected string
	}{
		{name: "dummyStringer", req: &dummyStringer{field: ""}, expected: "my value"},
		{name: "dummy", req: &dummy{field: "value"}, expected: "&{field:value}"},
		{name: "dummyStringerRedacter", req: &dummyStringerRedacter{field: ""}, expected: "my value redacted"},
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
		wantLevel  slog.Level
		wantErrStr string
	}{
		{name: "no error", err: nil, wantLevel: slog.LevelInfo, wantErrStr: ""},
		{name: "error", err: errors.New("test error"), wantLevel: slog.LevelError, wantErrStr: "test error"},
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

type captureHandler struct {
	records []slog.Record
	attrs   []map[string]any
}

func (h *captureHandler) reset() {
	h.records = nil
	h.attrs = nil
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})
	h.records = append(h.records, record.Clone())
	h.attrs = append(h.attrs, attrs)
	return nil
}

func (h *captureHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *captureHandler) WithGroup(string) slog.Handler {
	return h
}
