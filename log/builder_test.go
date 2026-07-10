package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestBuilderDefaultsText(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewHandler(WithWriter(&buf)))
	logger.Info("hello", "k", "v")
	if !strings.Contains(buf.String(), "hello") || !strings.Contains(buf.String(), "k=v") {
		t.Fatalf("output = %q", buf.String())
	}
}

func TestBuilderJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewHandler(WithWriter(&buf), WithFormat(FormatJSON)))
	logger.LogAttrs(context.Background(), LevelFatal, "boom")
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v (%q)", err, buf.String())
	}
	if got["level"] != LevelFatal.String() {
		t.Fatalf("level = %v, want %s", got["level"], LevelFatal.String())
	}
}

func TestBuilderWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewHandler(
		WithWriter(&buf),
		WithFormat(FormatJSON),
	)).With(
		slog.String("service.id", "i"),
		slog.String("service.name", "n"),
		slog.String("service.version", "v"),
	)
	logger.Info("hi")
	var got map[string]any
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if got["service.name"] != "n" {
		t.Fatalf("service.name = %v", got["service.name"])
	}
}

func TestBuilderWithHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := NewLogger(handler).With(slog.String("service.name", "n"))
	logger.Info("hi")
	var got map[string]any
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if got["service.name"] != "n" {
		t.Fatalf("service.name = %v", got["service.name"])
	}
}

func TestBuilderWithFilter(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: LevelWarn})
	logger := NewLogger(handler, WithFilter(FilterKey("password")))
	logger.Info("ignored")
	logger.Warn("login", "password", "secret")
	out := buf.String()
	if strings.Contains(out, "ignored") {
		t.Fatalf("info should be filtered: %q", out)
	}
	if !strings.Contains(out, "password=***") {
		t.Fatalf("expected password redacted: %q", out)
	}
}

func TestBuilderWithStoredContextAttrs(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewHandler(WithWriter(&buf), WithFormat(FormatJSON)))
	ctx := ContextWithAttrs(context.Background(), slog.String("trace_id", "abc"))
	logger.InfoContext(ctx, "hi")
	var got map[string]any
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if got["trace_id"] != "abc" {
		t.Fatalf("trace_id = %v", got["trace_id"])
	}
}

func TestBuilderWithExtractor(t *testing.T) {
	type userIDKey struct{}
	type requestIDKey struct{}

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := NewLogger(
		handler,
		WithExtractor(func(ctx context.Context) []slog.Attr {
			id, _ := ctx.Value(userIDKey{}).(string)
			if id == "" {
				return nil
			}
			return []slog.Attr{slog.String("user_id", id)}
		}, func(ctx context.Context) []slog.Attr {
			id, _ := ctx.Value(requestIDKey{}).(string)
			if id == "" {
				return nil
			}
			return []slog.Attr{slog.String("request_id", id)}
		}),
	)
	ctx := context.WithValue(context.Background(), userIDKey{}, "u1")
	ctx = context.WithValue(ctx, requestIDKey{}, "r1")
	logger.InfoContext(ctx, "hi")
	var got map[string]any
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got)
	if got["user_id"] != "u1" {
		t.Fatalf("user_id = %v", got["user_id"])
	}
	if got["request_id"] != "r1" {
		t.Fatalf("request_id = %v", got["request_id"])
	}
}
