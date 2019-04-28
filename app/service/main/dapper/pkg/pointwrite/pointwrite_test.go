package pointwrite

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/dapper/model"
)

func TestPointWrite(t *testing.T) {
	var data []*model.SpanPoint
	mockFn := func(ctx context.Context, points []*model.SpanPoint) error {
		data = append(data, points...)
		return nil
	}
	pw := &pointwriter{
		fn:        mockFn,
		current:   make(map[string]*model.SpanPoint),
		precision: 5,
		timeout:   time.Second,
		tk:        time.NewTicker(time.Second * time.Duration(5)),
	}
	spans := []*model.Span{
		&model.Span{
			ServiceName: "test1",
			StartTime:   time.Unix(100, 0),
		},
		&model.Span{
			ServiceName: "test1",
			StartTime:   time.Unix(110, 0),
		},
	}
	for _, span := range spans {
		if err := pw.WriteSpan(span); err != nil {
			t.Error(err)
		}
	}
	if len(pw.current) != 2 {
		t.Errorf("expect 2 point get %d", len(pw.current))
	}
	pw.flush()
	if len(data) != 2 {
		t.Errorf("expect 2 point get %d", len(data))
	}
}

func TestPointWriteFlush(t *testing.T) {
	var data []*model.SpanPoint
	wait := make(chan bool, 1)
	mockFn := func(ctx context.Context, points []*model.SpanPoint) error {
		data = append(data, points...)
		wait <- true
		return nil
	}
	pw := New(mockFn, 1, time.Second)
	spans := []*model.Span{
		&model.Span{
			ServiceName: "test1",
			StartTime:   time.Unix(100, 0),
		},
		&model.Span{
			ServiceName: "test1",
			StartTime:   time.Unix(110, 0),
		},
	}
	for _, span := range spans {
		if err := pw.WriteSpan(span); err != nil {
			t.Error(err)
		}
	}
	<-wait
	if len(data) != 2 {
		t.Errorf("expect 2 point get %d", len(data))
	}
}
