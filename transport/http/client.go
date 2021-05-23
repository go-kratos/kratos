package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	xhttp "github.com/go-kratos/kratos/v2/internal/http"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// RequestEncoder is request encoder
type RequestEncoder func(in interface{}, bodyPattern string) (io.Reader, error)

var defaultEncoder RequestEncoder = func(in interface{}, bodyPattern string) (io.Reader, error) {
	if bodyPattern == "" {
		return nil, nil
	}
	content, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), err
}

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*callInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*callInfo, *csAttempt)
}

type callInfo struct{}

type csAttempt struct{}

// Client is http client
type Client struct {
	cc     *http.Client
	encode RequestEncoder
}

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, w *http.Response) error

// ClientOption is HTTP client option.
type ClientOption func(*clientOptions)

// WithTimeout with client request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
	}
}

// WithUserAgent with client user agent.
func WithUserAgent(ua string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = ua
	}
}

// WithTransport with client transport.
func WithTransport(trans http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = middleware.Chain(m...)
	}
}

// WithSchema with client schema.
func WithSchema(schema string) ClientOption {
	return func(o *clientOptions) {
		o.schema = schema
	}
}

// WithEndpoint with client addr.
func WithEndpoint(addr string) ClientOption {
	return func(o *clientOptions) {
		o.addr = addr
	}
}

// WithEncoder with client request encode.
func WithEncoder(encoder RequestEncoder) ClientOption {
	return func(o *clientOptions) {
		o.requestEncoder = encoder
	}
}

// Client is a HTTP transport client.
type clientOptions struct {
	ctx            context.Context
	timeout        time.Duration
	userAgent      string
	transport      http.RoundTripper
	errorDecoder   DecodeErrorFunc
	middleware     middleware.Middleware
	schema         string
	addr           string
	requestEncoder RequestEncoder
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := &clientOptions{
		requestEncoder: defaultEncoder,
	}
	for _, o := range opts {
		o(options)
	}
	trans, err := NewTransport(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		cc:     &http.Client{Transport: trans},
		encode: options.requestEncoder,
	}, nil
}

// Encode encode request
func (c *Client) Encode(in interface{}, bodyPattern string) (io.Reader, error) {
	return c.encode(in, bodyPattern)
}

// NewTransport creates an http.RoundTripper.
func NewTransport(ctx context.Context, opts ...ClientOption) (http.RoundTripper, error) {
	options := &clientOptions{
		ctx:          ctx,
		timeout:      500 * time.Millisecond,
		transport:    http.DefaultTransport,
		errorDecoder: checkResponse,
		schema:       "http",
	}
	for _, o := range opts {
		o(options)
	}
	return &baseTransport{
		errorDecoder: options.errorDecoder,
		middleware:   options.middleware,
		userAgent:    options.userAgent,
		timeout:      options.timeout,
		base:         options.transport,
		schema:       options.schema,
		addr:         options.addr,
	}, nil
}

type baseTransport struct {
	userAgent    string
	timeout      time.Duration
	base         http.RoundTripper
	errorDecoder DecodeErrorFunc
	middleware   middleware.Middleware
	addr         string
	schema       string
}

func (t *baseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	if t.schema != "" {
		req.URL.Scheme = t.schema
	}
	if t.addr != "" {
		req.URL.Host = t.addr
	}
	ctx := transport.NewContext(req.Context(), transport.Transport{Kind: transport.KindHTTP})
	info, ok := FromClientContext(ctx)
	if ok {
		info.Request = req
	} else {
		ctx = NewClientContext(ctx, &ClientInfo{Request: req})
	}
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		res, err := t.base.RoundTrip(in.(*http.Request))
		if err != nil {
			return nil, err
		}
		if err := t.errorDecoder(ctx, res); err != nil {
			return nil, err
		}
		return res, nil
	}
	if t.middleware != nil {
		h = t.middleware(h)
	}
	res, err := h(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*http.Response), nil
}

// Do send an HTTP request and decodes the body of response into target.
// returns an error (of type *Error) if the response status code is not 2xx.
func Do(client *Client, req *http.Request, target interface{}, opts ...CallOption) error {

	res, err := client.cc.Do(req)
	if err != nil {
		return err
	}

	subtype := xhttp.ContentSubtype(res.Header.Get(xhttp.HeaderContentType))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return codec.Unmarshal(data, target)
}

// checkResponse returns an error (of type *Error) if the response
// status code is not 2xx.
func checkResponse(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	if data, err := ioutil.ReadAll(res.Body); err == nil {
		st := new(spb.Status)
		if err = protojson.Unmarshal(data, st); err == nil {
			return status.ErrorProto(st)
		}
	}
	return status.Error(xhttp.GRPCCodeFromStatus(res.StatusCode), res.Status)
}
