package http

import (
	"io/ioutil"
	"net/http"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	xhttp "github.com/go-kratos/kratos/v2/internal/http"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// HandleOption is handle option.
type HandleOption func(*HandleOptions)

// HandleOptions is handle options.
type HandleOptions struct {
	Decode     DecodeRequestFunc
	Encode     EncodeResponseFunc
	Error      EncodeErrorFunc
	Middleware middleware.Middleware
}

// DefaultHandleOptions returns a default handle options.
func DefaultHandleOptions() HandleOptions {
	return HandleOptions{
		Decode: decodeRequest,
		Encode: encodeResponse,
		Error:  encodeError,
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Decode = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Encode = en
	}
}

// ErrorEncoder with error encoder.
func ErrorEncoder(en EncodeErrorFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Error = en
	}
}

// Middleware with middleware option.
func Middleware(m middleware.Middleware) HandleOption {
	return func(o *HandleOptions) {
		o.Middleware = m
	}
}

// decodeRequest decodes the request body to object.
func decodeRequest(req *http.Request, v interface{}) error {
	subtype := xhttp.ContentSubtype(req.Header.Get(xhttp.HeaderContentType))
	if codec := encoding.GetCodec(subtype); codec != nil {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		return codec.Unmarshal(data, v)
	}
	return binding.BindForm(req, v)
}

// encodeResponse encodes the object to the HTTP response.
func encodeResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	codec := codecForRequest(r)
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set(xhttp.HeaderContentType, xhttp.ContentType(codec.Name()))
	_, _ = w.Write(data)
	return nil
}

// encodeError encodes the error to the HTTP response.
func encodeError(w http.ResponseWriter, r *http.Request, err error) {
	st, _ := status.FromError(err)
	data, err := protojson.Marshal(st.Proto())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(xhttp.HeaderContentType, "application/json; charset=utf-8")
	w.WriteHeader(xhttp.StatusFromGRPCCode(st.Code()))
	w.Write(data)
}

// codecForRequest get encoding.Codec via http.Request
func codecForRequest(r *http.Request) encoding.Codec {
	var codec encoding.Codec
	for _, accept := range r.Header[xhttp.HeaderAccept] {
		if codec = encoding.GetCodec(xhttp.ContentSubtype(accept)); codec != nil {
			break
		}
	}
	if codec == nil {
		codec = encoding.GetCodec(json.Name)
	}
	return codec
}
