package recovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

func TestOnce(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Error("fail")
		}
	}()

	next := func(context.Context, interface{}) (interface{}, error) {
		panic("panic reason")
	}
	_, e := Recovery(WithHandler(func(ctx context.Context, _, err interface{}) error {
		_, ok := ctx.Value(Latency{}).(float64)
		if !ok {
			t.Errorf("not latency")
		}
		return errors.InternalServer("RECOVERY", fmt.Sprintf("panic triggered: %v", err))
	}))(next)(context.Background(), "panic")
	t.Logf("succ and reason is %v", e)
}

func TestNotPanic(t *testing.T) {
	next := func(_ context.Context, req interface{}) (interface{}, error) {
		return req.(string) + "https://go-kratos.dev", nil
	}

	_, e := Recovery(WithHandler(func(_ context.Context, _ any, err any) error {
		return errors.InternalServer("RECOVERY", fmt.Sprintf("panic triggered: %v", err))
	}))(next)(context.Background(), "notPanic")
	if e != nil {
		t.Errorf("e isn't nil")
	}
}
