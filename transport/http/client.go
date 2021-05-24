package http

import (
	"bytes"
	"context"
	"encoding/json"
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

// RequestEncoder is request encoder
type RequestEncoder func(in interface{}) ([]byte, error)

var defaultEncoder RequestEncoder = func(in interface{}) ([]byte, error) {
	content, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return content, err
}

// Client is http client
type Client struct {
	cc *http.Client

	encode       RequestEncoder
	userAgent    string
	timeout      time.Duration
	errorDecoder DecodeErrorFunc
	middleware   middleware.Middleware
	addr         string
	schema       string
	contentType  string
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
func WithEncoder(encoder RequestEncoder, contentType string) ClientOption {
	return func(o *clientOptions) {
		o.requestEncoder = encoder
		o.contentType = contentType
	}
}

// Client is a HTTP transport client.
type clientOptions struct {
	ctx            context.Context
	timeout        time.Duration
	userAgent      string
	errorDecoder   DecodeErrorFunc
	middleware     middleware.Middleware
	schema         string
	addr           string
	requestEncoder RequestEncoder
	contentType    string
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := &clientOptions{
		ctx:            ctx,
		timeout:        1000 * time.Millisecond,
		errorDecoder:   checkResponse,
		requestEncoder: defaultEncoder,
		contentType:    "application/json",
	}
	for _, o := range opts {
		o(options)
	}

	return &Client{
		cc:           &http.Client{Timeout: options.timeout},
		encode:       options.requestEncoder,
		errorDecoder: options.errorDecoder,
		middleware:   options.middleware,
		userAgent:    options.userAgent,
		timeout:      options.timeout,
		schema:       options.schema,
		addr:         options.addr,
		contentType:  options.contentType,
	}, nil
}

// Invoke makes an rpc call procedure for remote service
func (client *Client) Invoke(ctx context.Context, method string, pathPattern string, args interface{}, reply interface{}, opts ...CallOption) error {
	var (
		c       callInfo
		reqBody io.Reader
	)
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
	url := fmt.Sprintf("%s://%s%s", schema, client.addr, path)
	if args != nil {
		if c.bodyPattern == nil || (c.bodyPattern != nil && *c.bodyPattern != "") {
			// TODO: only encode the target field of args
			content, err := client.encode(args)
			if err != nil {
				return err
			}
			reqBody = bytes.NewReader(content)
		}
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	c.pathPattern = pathPattern

	resp, err := client.do(req, c)
	if err != nil {
		return err
	}

	subtype := xhttp.ContentSubtype(resp.Header.Get(xhttp.HeaderContentType))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return codec.Unmarshal(data, reply)
}

// Do send an HTTP request and decodes the body of response into target.
// returns an error (of type *Error) if the response status code is not 2xx.
func (client *Client) Do(req *http.Request, opts ...CallOption) (*http.Response, error) {
	var c callInfo
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return nil, err
		}
	}

	return client.do(req, c)
}

func (client *Client) do(req *http.Request, c callInfo) (*http.Response, error) {
	if client.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", client.userAgent)
	}
	if req.Body != nil && client.contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", client.contentType)
	}
	if client.schema != "" {
		req.URL.Scheme = client.schema
	}
	if client.addr != "" {
		req.URL.Host = client.addr
	}
	ctx := transport.NewContext(req.Context(), transport.Transport{Kind: transport.KindHTTP})
	ctx = NewClientContext(ctx, ClientInfo{
		PathPattern: c.pathPattern,
		Request:     req,
	})

	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		res, err := client.cc.Do(in.(*http.Request))
		if err != nil {
			return nil, err
		}
		if err := client.errorDecoder(ctx, res); err != nil {
			return nil, err
		}
		return res, nil
	}
	if client.middleware != nil {
		h = client.middleware(h)
	}
	resp, err := h(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*http.Response), nil
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
