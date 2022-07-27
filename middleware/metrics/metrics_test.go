package metrics

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type (
	mockCounter struct {
		lvs   []string
		value float64
	}
	mockObserver struct {
		lvs   []string
		value float64
	}
)

func (m *mockCounter) With(lvs ...string) metrics.Counter {
	return m
}

func (m *mockCounter) Inc() {
	m.value += 1.0
}

func (m *mockCounter) Add(delta float64) {
	m.value += delta
}

func (m *mockObserver) With(lvs ...string) metrics.Observer {
	return m
}

func (m *mockObserver) Observe(delta float64) {
	m.value += delta
}

func TestWithRequests(t *testing.T) {
	mock := mockCounter{
		lvs:   []string{"Initial"},
		value: 1.23,
	}
	o := options{
		requests: &mock,
	}

	WithRequests(&mock)(&o)

	if _, ok := o.requests.(*mockCounter); !ok {
		t.Errorf(`The type of the option requests property must be of "mockCounter", %T given.`, o.requests)
	}

	counter := o.requests.(*mockCounter)

	if len(counter.lvs) != 1 || counter.lvs[0] != "Initial" {
		t.Errorf(`The given counter lvs must have only one element equal to "Initial", %v given`, counter.lvs)
	}
	if counter.value != 1.23 {
		t.Errorf(`The given counter value must be equal to 1.23, %v given`, counter.value)
	}
}

func TestWithSeconds(t *testing.T) {
	mock := mockObserver{
		lvs:   []string{"Initial"},
		value: 1.23,
	}
	o := options{
		seconds: &mock,
	}

	WithSeconds(&mock)(&o)

	if _, ok := o.seconds.(*mockObserver); !ok {
		t.Errorf(`The type of the option requests property must be of "mockObserver", %T given.`, o.requests)
	}

	observer := o.seconds.(*mockObserver)

	if len(observer.lvs) != 1 || observer.lvs[0] != "Initial" {
		t.Errorf(`The given observer lvs must have only one element equal to "Initial", %v given`, observer.lvs)
	}
	if observer.value != 1.23 {
		t.Errorf(`The given observer value must be equal to 1.23, %v given`, observer.value)
	}
}

func TestServer(t *testing.T) {
	e := errors.New("got an error")
	nextError := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, e
	}
	nextValid := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "Hello valid", nil
	}

	_, err := Server()(nextError)(context.Background(), "test:")
	if err != e {
		t.Error("The given error mismatch the expected.")
	}

	res, err := Server(func(o *options) {
		o.requests = &mockCounter{
			lvs:   []string{"Initial"},
			value: 1.23,
		}
		o.seconds = &mockObserver{
			lvs:   []string{"Initial"},
			value: 1.23,
		}
	})(nextValid)(transport.NewServerContext(context.Background(), &http.Transport{}), "test:")
	if err != nil {
		t.Error("The server must not throw an error.")
	}
	if res != "Hello valid" {
		t.Error(`The server must return a "Hello valid" response.`)
	}
}

func TestClient(t *testing.T) {
	e := errors.New("got an error")
	nextError := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, e
	}
	nextValid := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "Hello valid", nil
	}

	_, err := Client()(nextError)(context.Background(), "test:")
	if err != e {
		t.Error("The given error mismatch the expected.")
	}

	res, err := Client(func(o *options) {
		o.requests = &mockCounter{
			lvs:   []string{"Initial"},
			value: 1.23,
		}
		o.seconds = &mockObserver{
			lvs:   []string{"Initial"},
			value: 1.23,
		}
	})(nextValid)(transport.NewClientContext(context.Background(), &http.Transport{}), "test:")
	if err != nil {
		t.Error("The server must not throw an error.")
	}
	if res != "Hello valid" {
		t.Error(`The server must return a "Hello valid" response.`)
	}
}
