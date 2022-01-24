package metrics

import (
	"context"
	"testing"
)

func TestMetrics(t *testing.T) {
	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req.(string) + "https://go-kratos.dev", nil
	}
	_, err := Server()(next)(context.Background(), "test:")
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}

	_, err = Client()(next)(context.Background(), "test:")
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
}
