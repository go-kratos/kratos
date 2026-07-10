package log

import (
	"context"
	"log/slog"
	"testing"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
	"go.opentelemetry.io/otel/trace"

	klog "github.com/go-kratos/kratos/v3/log"
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
	for _, attr := range traceAttrs(ctx) {
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

func TestNewHandlerAppliesOTelOptionsAndTraceAttrs(t *testing.T) {
	recorder := logtest.NewRecorder()
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	logger := slog.New(NewHandler("helloworld",
		WithLoggerProvider(recorder),
		WithVersion("v1.2.3"),
		WithSchemaURL("https://example.test/schema"),
	))
	logger.InfoContext(ctx, "user created", "user_id", "42")

	scope, records := onlyRecording(t, recorder.Result())
	if scope.Name != "helloworld" {
		t.Fatalf("scope name = %q, want %q", scope.Name, "helloworld")
	}
	if scope.Version != "v1.2.3" {
		t.Fatalf("scope version = %q, want %q", scope.Version, "v1.2.3")
	}
	if scope.SchemaURL != "https://example.test/schema" {
		t.Fatalf("scope schema URL = %q, want %q", scope.SchemaURL, "https://example.test/schema")
	}
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}

	attrs := recordAttrs(records[0])
	tests := map[string]string{
		"trace_id":    traceID.String(),
		"span_id":     spanID.String(),
		"trace_flags": trace.FlagsSampled.String(),
		"user_id":     "42",
	}
	for key, want := range tests {
		if attrs[key] != want {
			t.Fatalf("%s = %q, want %q", key, attrs[key], want)
		}
	}
}

func TestNewHandlerComposesWithCoreLoggerOptions(t *testing.T) {
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

	logger := klog.NewLogger(
		NewHandler("helloworld", WithLoggerProvider(recorder)),
		klog.WithFilter(klog.FilterKey("password")),
	).With(slog.String("service.name", "helloworld"))
	logger.InfoContext(ctx, "user created", "user_id", "42", "password", "secret")

	_, records := onlyRecording(t, recorder.Result())
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}

	record := records[0]
	attrs := recordAttrs(record)
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
	if count := recordAttrCount(record, "request_id"); count != 1 {
		t.Fatalf("request_id count = %d, want 1", count)
	}
	if count := recordAttrCount(record, "trace_id"); count != 1 {
		t.Fatalf("trace_id count = %d, want 1", count)
	}
}

func onlyRecording(t *testing.T, results logtest.Recording) (logtest.Scope, []logtest.Record) {
	t.Helper()
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	for scope, records := range results {
		return scope, records
	}
	t.Fatal("missing recording")
	return logtest.Scope{}, nil
}

func recordAttrs(record logtest.Record) map[string]string {
	attrs := map[string]string{}
	for _, kv := range record.Attributes {
		if kv.Value.Kind() == otellog.KindString {
			attrs[kv.Key] = kv.Value.AsString()
		}
	}
	return attrs
}

func recordAttrCount(record logtest.Record, key string) int {
	var count int
	for _, kv := range record.Attributes {
		if kv.Key == key {
			count++
		}
	}
	return count
}
