package circuitbreaker

import (
	"context"
	stderrors "errors"
	"testing"

	kratoserrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/transport"
)

type transportMock struct {
	kind      transport.Kind
	endpoint  string
	operation string
}

type circuitBreakerMock struct {
	err     error
	success int
	failed  int
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
func (c *circuitBreakerMock) MarkSuccess() { c.success++ }
func (c *circuitBreakerMock) MarkFailed()  { c.failed++ }

func TestWithBreakerFactory(t *testing.T) {
	var created int
	var o options

	WithBreakerFactory(func() CircuitBreaker {
		created++
		return &circuitBreakerMock{}
	})(&o)

	if o.group == nil {
		t.Fatal("breaker group must be updated")
	}
	first := o.group.Get("/foo")
	second := o.group.Get("/foo")
	third := o.group.Get("/bar")
	if first != second {
		t.Fatal("same operation must reuse the same circuit breaker")
	}
	if first == third {
		t.Fatal("different operations must use different circuit breakers")
	}
	if created != 2 {
		t.Fatalf("factory created %d breakers, want 2", created)
	}
}

func TestClient(t *testing.T) {
	breaker := &circuitBreakerMock{}
	ctx := transport.NewClientContext(context.Background(), &transportMock{operation: "/foo"})
	next := func(context.Context, any) (any, error) {
		return "Hello valid", nil
	}

	reply, err := Client(WithBreakerFactory(func() CircuitBreaker {
		return breaker
	}))(next)(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if reply != "Hello valid" {
		t.Fatalf("reply = %v, want Hello valid", reply)
	}
	if breaker.success != 1 || breaker.failed != 0 {
		t.Fatalf("breaker success=%d failed=%d, want success=1 failed=0", breaker.success, breaker.failed)
	}
}

func TestClientRejectsWhenBreakerIsOpen(t *testing.T) {
	breaker := &circuitBreakerMock{err: stderrors.New("circuitbreaker error")}
	ctx := transport.NewClientContext(context.Background(), &transportMock{operation: "/foo"})
	called := false
	next := func(context.Context, any) (any, error) {
		called = true
		return nil, nil
	}

	_, err := Client(WithBreakerFactory(func() CircuitBreaker {
		return breaker
	}))(next)(ctx, nil)

	if !stderrors.Is(err, ErrNotAllowed) {
		t.Fatalf("err = %v, want ErrNotAllowed", err)
	}
	if called {
		t.Fatal("handler must not be called when breaker is open")
	}
	if breaker.success != 0 || breaker.failed != 1 {
		t.Fatalf("breaker success=%d failed=%d, want success=0 failed=1", breaker.success, breaker.failed)
	}
}

func TestClientMarksServerErrorAsFailed(t *testing.T) {
	breaker := &circuitBreakerMock{}
	ctx := transport.NewClientContext(context.Background(), &transportMock{operation: "/foo"})
	next := func(context.Context, any) (any, error) {
		return nil, kratoserrors.InternalServer("", "")
	}

	_, _ = Client(WithBreakerFactory(func() CircuitBreaker {
		return breaker
	}))(next)(ctx, nil)

	if breaker.success != 0 || breaker.failed != 1 {
		t.Fatalf("breaker success=%d failed=%d, want success=0 failed=1", breaker.success, breaker.failed)
	}
}
