package http

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/gorilla/mux"
)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(Context) error

// Context is an HTTP Context.
type Context interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Middleware() middleware.Middleware
	Bind(interface{}) error
	Result(int, interface{}) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ctx struct {
	req *http.Request
	res http.ResponseWriter
	m   middleware.Middleware
	dec DecodeRequestFunc
	enc EncodeResponseFunc
	err EncodeErrorFunc
	h   HandlerFunc
}

func (c *ctx) Middleware() middleware.Middleware    { return c.m }
func (c *ctx) Bind(v interface{}) error             { return c.dec(c.req, v) }
func (c *ctx) Result(code int, v interface{}) error { return c.enc(c.res, c.req, v) }
func (c *ctx) Request() *http.Request               { return c.req }
func (c *ctx) Response() http.ResponseWriter        { return c.res }

func (c *ctx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.res = w
	c.req = r
	if err := c.h(c); err != nil {
		c.err(w, r, err)
	}
}

// Route is an HTTP route.
type Route struct {
	r   *mux.Router
	m   middleware.Middleware
	dec DecodeRequestFunc
	enc EncodeResponseFunc
	err EncodeErrorFunc
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(res, req)
}

func (r *Route) newHandler(h HandlerFunc) Context {
	return &ctx{h: h, m: r.m, dec: r.dec, enc: r.enc, err: r.err}
}

// Handle .
func (r *Route) Handle(method, path string, h HandlerFunc) {
	r.r.Handle(path, r.newHandler(h)).Methods(method)
}

// GET .
func (r *Route) GET(path string, h HandlerFunc) { r.Handle(http.MethodGet, path, h) }

// POST .
func (r *Route) POST(path string, h HandlerFunc) { r.Handle(http.MethodPost, path, h) }
