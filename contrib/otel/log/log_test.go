package log

import (
	"context"
	"log/slog"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestTraceAttrs(t *testing.T) {
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	attrs := map[string]string{}
	for _, attr := range TraceAttrs(ctx) {
		if attr.Value.Kind() == slog.KindString {
			attrs[attr.Key] = attr.Value.String()
		}
	}
	if attrs["trace_id"] != traceID.String() {
		t.Fatalf("trace_id = %q, want %q", attrs["trace_id"], traceID.String())
	}
	if attrs["span_id"] != spanID.String() {
		t.Fatalf("span_id = %q, want %q", attrs["span_id"], spanID.String())
	}
	if attrs["trace_flags"] != trace.FlagsSampled.String() {
		t.Fatalf("trace_flags = %q, want %q", attrs["trace_flags"], trace.FlagsSampled.String())
	}
}
