package log

import (
	"context"
	"log/slog"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
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

func TestNewLoggerAppliesLogBuilderOptions(t *testing.T) {
	recorder := logtest.NewRecorder()
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	ctx := klog.ContextWithAttrs(context.Background(), slog.String("request_id", "req-1"))
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	logger := NewLogger("helloworld",
		WithLoggerProvider(recorder),
		WithVersion("v1.2.3"),
		WithSchemaURL("https://example.test/schema"),
		WithLogOptions(klog.WithAttrs(slog.String("service.name", "helloworld"))),
		WithFilter(klog.FilterKey("password")),
	)
	logger.InfoContext(ctx, "user created", "user_id", "42", "password", "secret")

	results := recorder.Result()
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Name != "helloworld" {
		t.Fatalf("scope name = %q, want %q", results[0].Name, "helloworld")
	}
	if results[0].Version != "v1.2.3" {
		t.Fatalf("scope version = %q, want %q", results[0].Version, "v1.2.3")
	}
	if results[0].SchemaURL != "https://example.test/schema" {
		t.Fatalf("scope schema URL = %q, want %q", results[0].SchemaURL, "https://example.test/schema")
	}
	if len(results[0].Records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(results[0].Records))
	}

	attrs := recordAttrs(results[0].Records[0].Record)
	tests := map[string]string{
		"service.name": "helloworld",
		"request_id":   "req-1",
		"trace_id":     traceID.String(),
		"span_id":      spanID.String(),
		"trace_flags":  trace.FlagsSampled.String(),
		"user_id":      "42",
		"password":     "***",
	}
	for key, want := range tests {
		if attrs[key] != want {
			t.Fatalf("%s = %q, want %q", key, attrs[key], want)
		}
	}
}

func recordAttrs(record otellog.Record) map[string]string {
	attrs := map[string]string{}
	record.WalkAttributes(func(kv otellog.KeyValue) bool {
		if kv.Value.Kind() == otellog.KindString {
			attrs[kv.Key] = kv.Value.AsString()
		}
		return true
	})
	return attrs
}
