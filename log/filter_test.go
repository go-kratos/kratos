package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"
)

func TestFilterKeyRedacts(t *testing.T) {
	cap := &captureHandler{}
	h := newFilterHandler(cap, FilterKey("password"))
	record := slog.NewRecord(time.Now(), LevelInfo, "login", 0)
	record.AddAttrs(slog.String("user", "alice"), slog.String("password", "secret"))
	_ = h.Handle(context.Background(), record)

	if got := cap.attrs[0]["password"]; got != redactedValue {
		t.Fatalf("password = %v, want %q", got, redactedValue)
	}
	if got := cap.attrs[0]["user"]; got != "alice" {
		t.Fatalf("user = %v", got)
	}
}

func TestFilterFuncDrops(t *testing.T) {
	cap := &captureHandler{}
	h := newFilterHandler(cap, FilterFunc(func(_ context.Context, r slog.Record) bool {
		return r.Message == "drop"
	}))
	r1 := slog.NewRecord(time.Now(), LevelInfo, "drop", 0)
	r2 := slog.NewRecord(time.Now(), LevelInfo, "keep", 0)
	_ = h.Handle(context.Background(), r1)
	_ = h.Handle(context.Background(), r2)
	if len(cap.records) != 1 || cap.records[0].Message != "keep" {
		t.Fatalf("records = %+v", cap.records)
	}
}

func TestFilterKeyRedactsGroupPath(t *testing.T) {
	var buf bytes.Buffer
	h := newFilterHandler(slog.NewJSONHandler(&buf, nil), FilterKey("user.password"))
	logger := slog.New(h)

	logger.Info("login", slog.Group("user",
		slog.String("name", "alice"),
		slog.String("password", "secret"),
	))

	got := decodeJSONLog(t, &buf)
	user, ok := got["user"].(map[string]any)
	if !ok {
		t.Fatalf("user attr = %#v, want object", got["user"])
	}
	if user["password"] != redactedValue {
		t.Fatalf("password = %v, want %q", user["password"], redactedValue)
	}
	if user["name"] != "alice" {
		t.Fatalf("name = %v, want alice", user["name"])
	}
}

func TestFilterKeyRedactsHandlerGroupPath(t *testing.T) {
	var buf bytes.Buffer
	h := newFilterHandler(slog.NewJSONHandler(&buf, nil), FilterKey("request.password"))
	logger := slog.New(h).WithGroup("request")

	logger.Info("login", "password", "secret")

	got := decodeJSONLog(t, &buf)
	request, ok := got["request"].(map[string]any)
	if !ok {
		t.Fatalf("request attr = %#v, want object", got["request"])
	}
	if request["password"] != redactedValue {
		t.Fatalf("password = %v, want %q", request["password"], redactedValue)
	}
}

func TestFilterEmptyReturnsNext(t *testing.T) {
	cap := &captureHandler{}
	h := newFilterHandler(cap)
	if h != cap {
		t.Fatal("expected next returned as-is when no options")
	}
}

func decodeJSONLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal log: %v (%q)", err, buf.String())
	}
	return got
}
