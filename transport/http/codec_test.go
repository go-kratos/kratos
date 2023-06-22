package http

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

func TestDefaultRequestDecoder(t *testing.T) {
	var (
		bodyStr = `{"a":"1", "b": 2}`
		r, _    = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(bodyStr)))
	)
	r.Header.Set("Content-Type", "application/json")

	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err := DefaultRequestDecoder(r, &v1)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != "1" {
		t.Errorf("expected %v, got %v", "1", v1.A)
	}
	if v1.B != int64(2) {
		t.Errorf("expected %v, got %v", 2, v1.B)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if bodyStr != string(data) {
		t.Errorf("expected %v, got %v", bodyStr, string(data))
	}
}

type mockResponseWriter struct {
	StatusCode int
	Data       []byte
	header     http.Header
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.Data = b
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func TestDefaultResponseEncoder(t *testing.T) {
	var (
		w    = &mockResponseWriter{StatusCode: 200, header: make(http.Header)}
		r, _ = http.NewRequest(http.MethodPost, "", nil)
		v    = &struct {
			A string `json:"a"`
			B int64  `json:"b"`
		}{
			A: "1",
			B: 2,
		}
	)
	r.Header.Set("Content-Type", "application/json")

	err := DefaultResponseEncoder(w, r, v)
	if err != nil {
		t.Fatal(err)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", w.Header().Get("Content-Type"))
	}
	if w.StatusCode != 200 {
		t.Errorf("expected %v, got %v", 200, w.StatusCode)
	}
	if w.Data == nil {
		t.Errorf("expected not nil, got %v", w.Data)
	}
}

func TestDefaultErrorEncoder(t *testing.T) {
	var (
		w    = &mockResponseWriter{header: make(http.Header)}
		r, _ = http.NewRequest(http.MethodPost, "", nil)
		err  = errors.New(511, "", "")
	)
	r.Header.Set("Content-Type", "application/json")

	DefaultErrorEncoder(w, r, err)
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected %v, got %v", "application/json", w.Header().Get("Content-Type"))
	}
	if w.StatusCode != 511 {
		t.Errorf("expected %v, got %v", 511, w.StatusCode)
	}
	if w.Data == nil {
		t.Errorf("expected not nil, got %v", w.Data)
	}
}

func TestDefaultResponseEncoderEncodeNil(t *testing.T) {
	var (
		w    = &mockResponseWriter{StatusCode: 204, header: make(http.Header)}
		r, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString("<xml></xml>")))
	)
	r.Header.Set("Content-Type", "application/json")

	err := DefaultResponseEncoder(w, r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if w.Header().Get("Content-Type") != "" {
		t.Errorf("expected empty string, got %v", w.Header().Get("Content-Type"))
	}
	if w.StatusCode != 204 {
		t.Errorf("expected %v, got %v", 204, w.StatusCode)
	}
	if w.Data != nil {
		t.Errorf("expected nil, got %v", w.Data)
	}
}

func TestCodecForRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString("<xml></xml>")))
	r.Header.Set("Content-Type", "application/xml")
	c, ok := CodecForRequest(r, "Content-Type")
	if !ok {
		t.Fatalf("expected true, got %v", ok)
	}
	if c.Name() != "xml" {
		t.Errorf("expected %v, got %v", "xml", c.Name())
	}

	r, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(`{"a":"1", "b": 2}`)))
	r.Header.Set("Content-Type", "blablablabla")
	c, ok = CodecForRequest(r, "Content-Type")
	if ok {
		t.Fatalf("expected false, got %v", ok)
	}
	if c.Name() != "json" {
		t.Errorf("expected %v, got %v", "json", c.Name())
	}
}
