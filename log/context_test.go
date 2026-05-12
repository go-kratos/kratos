package log

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestContextWithAttrsAppendsAttrs(t *testing.T) {
	ctx := ContextWithAttrs(context.Background(), slog.String("a", "1"))
	ctx = ContextWithAttrs(ctx, slog.String("b", "2"))

	attrs := AttrsFromContext(ctx)
	if len(attrs) != 2 {
		t.Fatalf("len(attrs) = %d, want 2", len(attrs))
	}
	if attrs[0].Key != "a" || attrs[0].Value.String() != "1" {
		t.Fatalf("attrs[0] = %v", attrs[0])
	}
	if attrs[1].Key != "b" || attrs[1].Value.String() != "2" {
		t.Fatalf("attrs[1] = %v", attrs[1])
	}
}

func TestContextWithAttrsNilCtx(t *testing.T) {
	ctx := ContextWithAttrs(nil, slog.Int("n", 1)) //nolint:staticcheck
	if attrs := AttrsFromContext(ctx); len(attrs) != 1 {
		t.Fatalf("len(attrs) = %d, want 1", len(attrs))
	}
}

func TestContextHandlerMerges(t *testing.T) {
	cap := &captureHandler{}
	h := newContextHandler(cap, AttrsFromContext)

	ctx := ContextWithAttrs(context.Background(), slog.String("trace_id", "t1"))
	record := slog.NewRecord(now(), LevelInfo, "msg", 0)
	record.AddAttrs(slog.String("k", "v"))
	if err := h.Handle(ctx, record); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if got := cap.attrs[0]["trace_id"]; got != "t1" {
		t.Fatalf("trace_id = %v", got)
	}
	if got := cap.attrs[0]["k"]; got != "v" {
		t.Fatalf("k = %v", got)
	}
}

func TestContextHandlerNoAttrs(t *testing.T) {
	cap := &captureHandler{}
	h := newContextHandler(cap, AttrsFromContext)
	record := slog.NewRecord(now(), LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), record); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if len(cap.records) != 1 {
		t.Fatalf("records = %d", len(cap.records))
	}
}

func now() time.Time { return time.Now() }
