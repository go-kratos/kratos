package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"google.golang.org/genproto/googleapis/api/httpbody"

	"github.com/go-kratos/kratos/v3/encoding"
	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/internal/httputil"
)

// SupportPackageIsVersion3 These constants should not be referenced from any other code.
const SupportPackageIsVersion3 = true

const defaultHTTPBodyContentType = "application/octet-stream"

// Redirector replies to the request with a redirect to url
// which may be a path relative to the request path.
type Redirector interface {
	error
	Redirect() (string, int)
}

// Request type net/http.
type Request = http.Request

// ResponseWriter type net/http.
type ResponseWriter = http.ResponseWriter

// Flusher type net/http
type Flusher = http.Flusher

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, any) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, any) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// DefaultRequestVars decodes the request vars to object.
func DefaultRequestVars(r *http.Request, v any) error {
	raws := mux.Vars(r)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return bindQuery(vars, v)
}

// DefaultRequestQuery decodes the request vars to object.
func DefaultRequestQuery(r *http.Request, v any) error {
	return bindQuery(r.URL.Query(), v)
}

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(r *http.Request, v any) error {
	if body, ok := httpBody(v); ok {
		data, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		if err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
		body.ContentType = r.Header.Get("Content-Type")
		body.Data = data
		return nil
	}
	codec, ok := CodecForRequest(r, "Content-Type")
	if !ok {
		return errors.BadRequest("CODEC", fmt.Sprintf("unregister Content-Type: %s", r.Header.Get("Content-Type")))
	}
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}
	if body, ok := httpBody(v); ok {
		contentType := body.GetContentType()
		if contentType == "" {
			contentType = defaultHTTPBodyContentType
		}
		w.Header().Set("Content-Type", contentType)
		_, err := w.Write(body.GetData())
		return err
	}
	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	var rd *redirect
	if errors.As(err, &rd) {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return
	}
	se := errors.FromError(err)
	codec, _ := CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(httputil.ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}

func httpBody(v any) (*httpbody.HttpBody, bool) {
	switch body := v.(type) {
	case *httpbody.HttpBody:
		return body, body != nil
	case **httpbody.HttpBody:
		if body == nil {
			return nil, false
		}
		if *body == nil {
			*body = new(httpbody.HttpBody)
		}
		return *body, true
	default:
		return nil, false
	}
}

// BodyContentType returns the content type carried by v or a binary default.
func BodyContentType(v any) string {
	if body, ok := httpBody(v); ok && body.GetContentType() != "" {
		return body.GetContentType()
	}
	return defaultHTTPBodyContentType
}
