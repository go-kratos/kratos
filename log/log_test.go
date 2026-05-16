package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestSetDefault(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	logger := slog.New(handler)
	SetDefault(logger)

	if got := Default(); got != logger {
		t.Fatalf("Default() = %v, want %v", got, logger)
	}

	Info("hello", "kind", "test")

	if len(handler.records) != 1 {
		t.Fatalf("records len = %d, want 1", len(handler.records))
	}
	if handler.records[0].Message != "hello" {
		t.Fatalf("message = %q, want %q", handler.records[0].Message, "hello")
	}
	if got := handler.attrs[0]["kind"]; got != "test" {
		t.Fatalf("kind = %v, want %q", got, "test")
	}
}

func TestWith(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	SetDefault(slog.New(handler))

	logger := With("service", "kratos", slog.Int("version", 2))
	logger.Info("hello", "kind", "test")

	if len(handler.records) != 1 {
		t.Fatalf("records len = %d, want 1", len(handler.records))
	}
	if got := handler.attrs[0]["service"]; got != "kratos" {
		t.Fatalf("service = %v, want %q", got, "kratos")
	}
	if got := handler.attrs[0]["version"]; got != int64(2) {
		t.Fatalf("version = %v, want 2", got)
	}
	if got := handler.attrs[0]["kind"]; got != "test" {
		t.Fatalf("kind = %v, want %q", got, "test")
	}
}

func TestWithGroup(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	var out bytes.Buffer
	SetDefault(slog.New(slog.NewJSONHandler(&out, nil)))

	logger := WithGroup("request")
	logger.Info("hello", "id", "r1")

	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v (%q)", err, out.String())
	}
	group, ok := got["request"].(map[string]any)
	if !ok {
		t.Fatalf("request group = %v, want object", got["request"])
	}
	if group["id"] != "r1" {
		t.Fatalf("request.id = %v, want r1", group["id"])
	}
}

func TestHandler(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	SetDefault(slog.New(handler))

	if got := Handler(); got != handler {
		t.Fatalf("Handler() = %v, want %v", got, handler)
	}
}

func TestEnabled(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	var out bytes.Buffer
	SetDefault(slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{Level: LevelWarn})))

	if Enabled(context.Background(), LevelInfo) {
		t.Fatal("Enabled(_, LevelInfo) = true, want false")
	}
	if !Enabled(context.Background(), LevelWarn) {
		t.Fatal("Enabled(_, LevelWarn) = false, want true")
	}
}

func TestLog(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	SetDefault(slog.New(handler))

	Log(context.Background(), LevelInfo, "hello", "kind", "test")

	if len(handler.records) != 1 {
		t.Fatalf("records len = %d, want 1", len(handler.records))
	}
	if handler.records[0].Message != "hello" {
		t.Fatalf("message = %q, want %q", handler.records[0].Message, "hello")
	}
	if got := handler.attrs[0]["kind"]; got != "test" {
		t.Fatalf("kind = %v, want %q", got, "test")
	}
}

func TestLogAttrs(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	SetDefault(slog.New(handler))

	LogAttrs(context.Background(), LevelWarn, "msg", slog.Int("n", 42))

	if got := handler.attrs[0]["n"]; got != int64(42) {
		t.Fatalf("n = %v, want 42", got)
	}
	if handler.records[0].Level != LevelWarn {
		t.Fatalf("level = %v, want %v", handler.records[0].Level, LevelWarn)
	}
}

func TestInfoContextPropagatesAttrs(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	handler := &captureHandler{}
	SetDefault(NewLogger(handler))

	ctx := ContextWithAttrs(context.Background(), slog.String("trace_id", "abc"))
	InfoContext(ctx, "hello")

	if got := handler.attrs[0]["trace_id"]; got != "abc" {
		t.Fatalf("trace_id = %v, want abc", got)
	}
}

func TestPackageHelpersReportCaller(t *testing.T) {
	old := Default()
	defer SetDefault(old)

	var out bytes.Buffer
	SetDefault(slog.New(NewHandler(
		WithWriter(&out),
		WithFormat(FormatJSON),
		WithAddSource(true),
	)))

	Info("source")
	assertLogSource(t, &out)

	out.Reset()
	LogAttrs(context.Background(), LevelInfo, "source")
	assertLogSource(t, &out)
}

func assertLogSource(t *testing.T, out *bytes.Buffer) {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v (%q)", err, out.String())
	}
	source, ok := got["source"].(map[string]any)
	if !ok {
		t.Fatalf("source = %v, want object", got["source"])
	}
	file, _ := source["file"].(string)
	if !strings.HasSuffix(file, "log/log_test.go") {
		t.Fatalf("source.file = %q, want log/log_test.go suffix", file)
	}
}

type captureHandler struct {
	parent *captureHandler
	prefix []slog.Attr
	groups []string

	records []slog.Record
	attrs   []map[string]any
}

func (h *captureHandler) root() *captureHandler {
	if h.parent != nil {
		return h.parent.root()
	}
	return h
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	for _, a := range h.collectPrefix() {
		attrs[a.Key] = a.Value.Any()
	}
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})
	root := h.root()
	root.records = append(root.records, record.Clone())
	root.attrs = append(root.attrs, attrs)
	return nil
}

func (h *captureHandler) collectPrefix() []slog.Attr {
	if h.parent == nil {
		return h.prefix
	}
	return append(h.parent.collectPrefix(), h.prefix...)
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &captureHandler{parent: h, prefix: append([]slog.Attr{}, attrs...)}
}

func (h *captureHandler) WithGroup(name string) slog.Handler {
	return &captureHandler{parent: h, groups: append([]string{}, name)}
}
