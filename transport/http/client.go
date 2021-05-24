package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	xhttp "github.com/go-kratos/kratos/v2/internal/http"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Client is http client
type Client struct {
	cc *http.Client

	encode       RequestEncodeFunc
	decode       RespDecodeFunc
	errorDecoder DecodeErrorFunc
	userAgent    string
	timeout      time.Duration
	middleware   middleware.Middleware
	endpoint     string
	schema       string
}

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, w *http.Response) error

// ClientOption is HTTP client option.
type ClientOption func(*clientOptions)

// WithTransport with client transport.
func WithTransport(trans http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

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
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// RequestEncodeFunc is request encoder
type RequestEncodeFunc func(ctx context.Context, in interface{}) (contentType string, body []byte, err error)

var defaultEncoder RequestEncodeFunc = func(ctx context.Context, in interface{}) (contentType string, body []byte, err error) {
	content, err := encoding.GetCodec("json").Marshal(in)
	if err != nil {
		return "", nil, err
	}
	return "application/json", content, err
}

// WithEncodeFunc with client request encode.
func WithEncodeFunc(encoder RequestEncodeFunc) ClientOption {
	return func(o *clientOptions) {
		o.encodeFunc = encoder
	}
}

// RespDecodeFunc is resp decoder
type RespDecodeFunc func(ctx context.Context, data []byte, v interface{}, contentType string) error

var defaultDecoder RespDecodeFunc = func(ctx context.Context, data []byte, v interface{}, contentType string) error {
	codec := encoding.GetCodec(contentType)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	return codec.Unmarshal(data, v)
}

// WithDecodeFunc with client response decode.
func WithDecodeFunc(decoder RespDecodeFunc) ClientOption {
	return func(o *clientOptions) {
		o.decodeFunc = decoder
	}
}

// Client is a HTTP transport client.
type clientOptions struct {
	ctx          context.Context
	transport    http.RoundTripper
	timeout      time.Duration
	userAgent    string
	middleware   middleware.Middleware
	schema       string
	endpoint     string
	encodeFunc   RequestEncodeFunc
	decodeFunc   RespDecodeFunc
	errorDecoder DecodeErrorFunc
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := &clientOptions{
		ctx:          ctx,
		timeout:      1000 * time.Millisecond,
		errorDecoder: checkResponse,
		encodeFunc:   defaultEncoder,
		decodeFunc:   defaultDecoder,
		transport:    http.DefaultTransport,
	}
	for _, o := range opts {
		o(options)
	}

	return &Client{
		cc:           &http.Client{Timeout: options.timeout, Transport: options.transport},
		encode:       options.encodeFunc,
		decode:       options.decodeFunc,
		errorDecoder: options.errorDecoder,
		middleware:   options.middleware,
		userAgent:    options.userAgent,
		timeout:      options.timeout,
		schema:       options.schema,
		endpoint:     options.endpoint,
	}, nil
}

// Invoke makes an rpc call procedure for remote service
func (client *Client) Invoke(ctx context.Context, pathPattern string, args interface{}, reply interface{}, opts ...CallOption) error {
	var (
		reqBody     io.Reader
		contentType string
	)

	c := defaultCallInfo()
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return err
		}
	}

	path := pathPattern
	if args != nil {
		path = binding.ProtoPath(path, args.(proto.Message))
	}
	schema := "http"
	if client.schema != "" {
		schema = client.schema
	}
	url := fmt.Sprintf("%s://%s%s", schema, client.endpoint, path)
	if args != nil && c.bodyPattern != "" {
		// TODO: only encode the target field of args
		var (
			content []byte
			err     error
		)
		contentType, content, err = client.encode(ctx, args)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(content)
	}
	req, err := http.NewRequest(c.method, url, reqBody)
	if err != nil {
		return err
	}
	if client.userAgent != "" {
		req.Header.Set("User-Agent", client.userAgent)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindHTTP})
	ctx = NewClientContext(ctx, ClientInfo{
		PathPattern: pathPattern,
		Request:     req,
	})

	return client.invoke(ctx, req, args, reply, c)
}

func (client *Client) invoke(ctx context.Context, req *http.Request, args interface{}, reply interface{}, c callInfo) error {
	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		info, _ := FromClientContext(ctx)

		resp, err := client.do(ctx, info.Request, c)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		subtype := xhttp.ContentSubtype(resp.Header.Get(xhttp.HeaderContentType))
		err = client.decode(ctx, data, reply, subtype)
		return reply, err
	}
	if client.middleware != nil {
		h = client.middleware(h)
	}
	_, err := h(ctx, args)
	return err
}

// Do send an HTTP request and decodes the body of response into target.
// returns an error (of type *Error) if the response status code is not 2xx.
func (client *Client) Do(req *http.Request, opts ...CallOption) (*http.Response, error) {
	c := defaultCallInfo()
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return nil, err
		}
	}

	return client.do(req.Context(), req, c)
}

func (client *Client) do(ctx context.Context, req *http.Request, c callInfo) (*http.Response, error) {
	resp, err := client.cc.Do(req)
	if err != nil {
		return nil, err
	}
	if err := client.errorDecoder(ctx, resp); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp, nil
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
