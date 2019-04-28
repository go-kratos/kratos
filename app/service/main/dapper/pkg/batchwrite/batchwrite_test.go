package batchwrite

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/service/main/dapper/model"
)

var (
	emptyspan = &model.Span{}
)

func TestRawDataBatchWriter(t *testing.T) {
	storage := make(map[string]map[string][]byte)
	writeFunc := func(ctx context.Context, traceID string, data map[string][]byte) error {
		if _, ok := storage[traceID]; !ok {
			storage[traceID] = make(map[string][]byte)
		}
		for k, v := range data {
			storage[traceID][k] = v
		}
		return nil
	}
	rbw := NewRawDataBatchWriter(writeFunc, 16, 2, 2, 0)
	spans := []*model.Span{
		&model.Span{
			TraceID: 1,
			SpanID:  11,
		},
		&model.Span{
			TraceID: 1,
			SpanID:  12,
		},
		&model.Span{
			TraceID: 2,
			SpanID:  21,
		},
		&model.Span{
			TraceID: 2,
			SpanID:  22,
		},
	}
	for _, span := range spans {
		if err := rbw.WriteSpan(span); err != nil {
			t.Error(err)
		}
	}
	rbw.Close()
	if len(storage) != 2 {
		t.Errorf("expect get 2 trace data, get %v", storage)
	}
	if len(storage["1"]) != 2 {
		t.Errorf("expect get 2 span data, get %v", storage["1"])
	}
	t.Logf("%v\n", storage)
}

func TestBatchWriterClosed(t *testing.T) {
	writeFunc2 := func(ctx context.Context, traceID string, data map[string][]byte) error {
		return nil
	}
	rbw := NewRawDataBatchWriter(writeFunc2, 1024*1024, 2, 2, 0)
	rbw.Close()
	if err := rbw.WriteSpan(emptyspan); err != ErrClosed {
		t.Errorf("expect err == ErrClosed get: %v", err)
	}
}

func randSpan() *model.Span {
	return &model.Span{
		TraceID: rand.Uint64() % 128,
		SpanID:  rand.Uint64() % 16,
	}
}

func BenchmarkRawDataWriter(b *testing.B) {
	writeFunc := func(ctx context.Context, traceID string, data map[string][]byte) error {
		return nil
	}
	rbw := NewRawDataBatchWriter(writeFunc, 1024*1024, 2, 2, 0)
	for i := 0; i < b.N; i++ {
		if err := rbw.WriteSpan(randSpan()); err != nil {
			b.Error(err)
		}
	}
	rbw.Close()
}
