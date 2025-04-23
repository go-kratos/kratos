package circuitbreaker

import (
	"context"
	"errors"
	"testing"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/group"
	"github.com/go-kratos/kratos/v2/transport"
)

type transportMock struct {
	kind      transport.Kind
	endpoint  string
	operation string
}

type circuitBreakerMock struct {
	err error
}

func (tr *transportMock) Kind() transport.Kind {
	return tr.kind
}

func (tr *transportMock) Endpoint() string {
	return tr.endpoint
}

func (tr *transportMock) Operation() string {
	return tr.operation
}

func (tr *transportMock) RequestHeader() transport.Header {
	return nil
}

func (tr *transportMock) ReplyHeader() transport.Header {
	return nil
}

func (c *circuitBreakerMock) Allow() error { return c.err }
func (c *circuitBreakerMock) MarkSuccess() {}
func (c *circuitBreakerMock) MarkFailed()  {}

func Test_WithGroup(t *testing.T) {
	o := options{
		group: group.NewGroup(func() any {
			return ""
		}),
	}

	WithGroup(nil)(&o)
	if o.group != nil {
		t.Error("The group property must be updated to nil.")
	}
}

func TestServer(_ *testing.T) {
	nextValid := func(context.Context, any) (any, error) {
		return "Hello valid", nil
	}
	nextInvalid := func(context.Context, any) (any, error) {
		return nil, kratoserrors.InternalServer("", "")
	}

	ctx := transport.NewClientContext(context.Background(), &transportMock{})

	_, _ = Client(func(o *options) {
		o.group = group.NewGroup(func() any {
			return &circuitBreakerMock{err: errors.New("circuitbreaker error")}
		})
	})(nextValid)(ctx, nil)

	_, _ = Client(func(_ *options) {})(nextValid)(ctx, nil)

	_, _ = Client(func(_ *options) {})(nextInvalid)(ctx, nil)
}
