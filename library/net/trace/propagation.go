package trace

import (
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrUnsupportedFormat occurs when the `format` passed to Tracer.Inject() or
	// Tracer.Extract() is not recognized by the Tracer implementation.
	ErrUnsupportedFormat = errors.New("trace: Unknown or unsupported Inject/Extract format")

	// ErrTraceNotFound occurs when the `carrier` passed to
	// Tracer.Extract() is valid and uncorrupted but has insufficient
	// information to extract a Trace.
	ErrTraceNotFound = errors.New("trace: Trace not found in Extract carrier")

	// ErrInvalidTrace errors occur when Tracer.Inject() is asked to
	// operate on a Trace which it is not prepared to handle (for
	// example, since it was created by a different tracer implementation).
	ErrInvalidTrace = errors.New("trace: Trace type incompatible with tracer")

	// ErrInvalidCarrier errors occur when Tracer.Inject() or Tracer.Extract()
	// implementations expect a different type of `carrier` than they are
	// given.
	ErrInvalidCarrier = errors.New("trace: Invalid Inject/Extract carrier")

	// ErrTraceCorrupted occurs when the `carrier` passed to
	// Tracer.Extract() is of the expected type but is corrupted.
	ErrTraceCorrupted = errors.New("trace: Trace data corrupted in Extract carrier")
)

// BuiltinFormat is used to demarcate the values within package `trace`
// that are intended for use with the Tracer.Inject() and Tracer.Extract()
// methods.
type BuiltinFormat byte

// support format list
const (
	// HTTPFormat represents Trace as HTTP header string pairs.
	//
	// the HTTPFormat format requires that the keys and values
	// be valid as HTTP headers as-is (i.e., character casing may be unstable
	// and special characters are disallowed in keys, values should be
	// URL-escaped, etc).
	//
	// the carrier must be a `http.Header`.
	HTTPFormat BuiltinFormat = iota
	// GRPCFormat represents Trace as gRPC metadata.
	//
	// the carrier must be a `google.golang.org/grpc/metadata.MD`.
	GRPCFormat
)

// Carrier propagator must convert generic interface{} to something this
// implement Carrier interface, Trace can use Carrier to represents itself.
type Carrier interface {
	Set(key, val string)
	Get(key string) string
}

// propagator is responsible for injecting and extracting `Trace` instances
// from a format-specific "carrier"
type propagator interface {
	Inject(carrier interface{}) (Carrier, error)
	Extract(carrier interface{}) (Carrier, error)
}

type httpPropagator struct{}

type httpCarrier http.Header

func (h httpCarrier) Set(key, val string) {
	http.Header(h).Set(key, val)
}

func (h httpCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}

func (httpPropagator) Inject(carrier interface{}) (Carrier, error) {
	header, ok := carrier.(http.Header)
	if !ok {
		return nil, ErrInvalidCarrier
	}
	if header == nil {
		return nil, ErrInvalidTrace
	}
	return httpCarrier(header), nil
}

func (httpPropagator) Extract(carrier interface{}) (Carrier, error) {
	header, ok := carrier.(http.Header)
	if !ok {
		return nil, ErrInvalidCarrier
	}
	if header == nil {
		return nil, ErrTraceNotFound
	}
	return httpCarrier(header), nil
}

const legacyGRPCKey = "trace"

type grpcPropagator struct{}

type grpcCarrier map[string][]string

func (g grpcCarrier) Get(key string) string {
	if v, ok := g[key]; ok && len(v) > 0 {
		return v[0]
	}
	ts := g[legacyGRPCKey]
	if len(ts) != 8 {
		return ""
	}
	switch key {
	case KeyTraceID:
		return ts[0]
	case KeyTraceSpanID:
		return ts[1]
	case KeyTraceParentID:
		return ts[2]
	case KeyTraceLevel:
		return ts[3]
	case KeyTraceSampled:
		return ts[4]
	case KeyTraceCaller:
		return ts[5]
	}
	return ""
}

func (g grpcCarrier) Set(key, val string) {
	ts := make([]string, 8)
	g[legacyGRPCKey] = ts
	switch key {
	case KeyTraceID:
		ts[0] = val
	case KeyTraceSpanID:
		ts[1] = val
	case KeyTraceParentID:
		ts[2] = val
	case KeyTraceLevel:
		ts[3] = val
	case KeyTraceSampled:
		ts[4] = val
	case KeyTraceCaller:
		ts[5] = val
	default:
		g[key] = append(g[key], val)
	}
}

func (grpcPropagator) Inject(carrier interface{}) (Carrier, error) {
	md, ok := carrier.(metadata.MD)
	if !ok {
		return nil, ErrInvalidCarrier
	}
	if md == nil {
		return nil, ErrInvalidTrace
	}
	return grpcCarrier(md), nil
}

func (grpcPropagator) Extract(carrier interface{}) (Carrier, error) {
	md, ok := carrier.(metadata.MD)
	if !ok {
		return nil, ErrInvalidCarrier
	}
	if md == nil {
		return nil, ErrTraceNotFound
	}
	return grpcCarrier(md), nil
}
