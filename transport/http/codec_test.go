package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/genproto/googleapis/api/httpbody"

	"github.com/go-kratos/kratos/v3/encoding"
	_ "github.com/go-kratos/kratos/v3/encoding/protojson"
	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/internal/testdata/binding"
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

func TestDefaultRequestDecoderHTTPBody(t *testing.T) {
	const bodyStr = "raw file content"
	r, _ := http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(bodyStr)))
	r.Header.Set("Content-Type", "text/plain")

	var body *httpbody.HttpBody
	if err := DefaultRequestDecoder(r, &body); err != nil {
		t.Fatal(err)
	}
	if body.GetContentType() != "text/plain" {
		t.Errorf("expected %v, got %v", "text/plain", body.GetContentType())
	}
	if string(body.GetData()) != bodyStr {
		t.Errorf("expected %v, got %v", bodyStr, string(body.GetData()))
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != bodyStr {
		t.Errorf("expected request body reset to %q, got %q", bodyStr, string(data))
	}
}

func TestDefaultRequestDecoderProtoJSONMessageFieldPointer(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(`{"naming":"go"}`)))
	r.Header.Set("Content-Type", "application/protojson")

	var sub *binding.Sub
	if err := DefaultRequestDecoder(r, &sub); err != nil {
		t.Fatal(err)
	}
	if sub == nil {
		t.Fatal("expected message field to be allocated")
	}
	if sub.Name != "go" {
		t.Errorf("expected %v, got %v", "go", sub.Name)
	}
}

func TestDefaultRequestDecoderProtoJSONRejectsScalarField(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(`"kratos"`)))
	r.Header.Set("Content-Type", "application/protojson")

	var name string
	err := DefaultRequestDecoder(r, &name)
	if err == nil {
		t.Fatal("expected scalar protojson body to fail")
	}
	if !strings.Contains(err.Error(), "want proto.Message") {
		t.Errorf("expected proto message type error, got %v", err)
	}
}

func TestDefaultResponseEncoderProtoJSONRejectsScalarField(t *testing.T) {
	w := &mockResponseWriter{StatusCode: http.StatusOK, header: make(http.Header)}
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	r.Header.Set("Accept", "application/protojson")

	err := DefaultResponseEncoder(w, r, "kratos")
	if err == nil {
		t.Fatal("expected scalar protojson response to fail")
	}
	if !strings.Contains(err.Error(), "want proto.Message") {
		t.Errorf("expected proto message type error, got %v", err)
	}
}

func TestDefaultResponseDecoderProtoJSONMessageFieldPointer(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"application/protojson"}},
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"naming":"go"}`)),
	}

	var sub *binding.Sub
	if err := DefaultResponseDecoder(context.TODO(), resp, &sub); err != nil {
		t.Fatal(err)
	}
	if sub == nil {
		t.Fatal("expected message field to be allocated")
	}
	if sub.Name != "go" {
		t.Errorf("expected %v, got %v", "go", sub.Name)
	}
}

func TestDefaultResponseDecoderProtoJSONRejectsScalarField(t *testing.T) {
	resp := &http.Response{
		Header:     http.Header{"Content-Type": []string{"application/protojson"}},
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`"kratos"`)),
	}

	var name string
	err := DefaultResponseDecoder(context.TODO(), resp, &name)
	if err == nil {
		t.Fatal("expected scalar protojson response to fail")
	}
	if !strings.Contains(err.Error(), "want proto.Message") {
		t.Errorf("expected proto message type error, got %v", err)
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

type errorCodec struct{}

func (errorCodec) Marshal(any) ([]byte, error) {
	return nil, errors.New(500, "mock", "marshal error")
}

func (errorCodec) Unmarshal([]byte, any) error {
	return nil
}

func (errorCodec) Name() string {
	return "mock"
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

func TestDefaultResponseEncoderHTTPBody(t *testing.T) {
	w := &mockResponseWriter{StatusCode: 200, header: make(http.Header)}
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	body := &httpbody.HttpBody{
		ContentType: "application/octet-stream",
		Data:        []byte("raw response"),
	}

	if err := DefaultResponseEncoder(w, r, body); err != nil {
		t.Fatal(err)
	}
	if got := w.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Errorf("expected %v, got %v", "application/octet-stream", got)
	}
	if string(w.Data) != "raw response" {
		t.Errorf("expected %v, got %v", "raw response", string(w.Data))
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

func TestDefaultErrorEncoderRedirect(t *testing.T) {
	w := &mockResponseWriter{header: make(http.Header)}
	r, _ := http.NewRequest(http.MethodGet, "/test", nil)

	DefaultErrorEncoder(w, r, NewRedirect("/redirect", http.StatusTemporaryRedirect))

	if w.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected %v, got %v", http.StatusTemporaryRedirect, w.StatusCode)
	}
	if w.Header().Get("Location") != "/redirect" {
		t.Errorf("expected %v, got %v", "/redirect", w.Header().Get("Location"))
	}
}

func TestDefaultErrorEncoderMarshalError(t *testing.T) {
	encoding.RegisterCodec(errorCodec{})
	w := &mockResponseWriter{header: make(http.Header)}
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	r.Header.Set("Accept", "application/mock")

	DefaultErrorEncoder(w, r, errors.New(500, "mock", "marshal error"))

	if w.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v, got %v", http.StatusInternalServerError, w.StatusCode)
	}
	if w.Header().Get("Content-Type") != "" {
		t.Errorf("expected empty content type, got %v", w.Header().Get("Content-Type"))
	}
	if w.Data != nil {
		t.Errorf("expected nil, got %v", w.Data)
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
