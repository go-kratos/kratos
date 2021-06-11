package http

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/gorilla/mux"
)

var _ Context = (*wrapper)(nil)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(Context) error

// Context is an HTTP Context.
type Context interface {
	context.Context
	Vars() url.Values
	Form() url.Values
	Header() http.Header
	Request() *http.Request
	Response() http.ResponseWriter
	Middleware(middleware.Handler) middleware.Handler
	Bind(interface{}) error
	Result(int, interface{}) error
	Returns(interface{}, error) error
	Reset(http.ResponseWriter, *http.Request)
}

type wrapper struct {
	route *Route
	req   *http.Request
	res   http.ResponseWriter
}

func (c *wrapper) Header() http.Header {
	return c.req.Header
}

func (c *wrapper) Vars() url.Values {
	raws := mux.Vars(c.req)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return vars
}
func (c *wrapper) Form() url.Values {
	if err := c.req.ParseForm(); err != nil {
		return url.Values{}
	}
	return c.req.Form
}
func (c *wrapper) Request() *http.Request        { return c.req }
func (c *wrapper) Response() http.ResponseWriter { return c.res }
func (c *wrapper) Middleware(h middleware.Handler) middleware.Handler {
	return middleware.Chain(c.route.srv.serviceM...)(h)
}
func (c *wrapper) Bind(v interface{}) error { return c.route.srv.dec(c.req, v) }
func (c *wrapper) Result(code int, v interface{}) error {
	c.res.WriteHeader(code)
	if err := c.route.srv.enc(c.res, c.req, v); err != nil {
		return err
	}
	return nil
}
func (c *wrapper) Returns(v interface{}, err error) error {
	if err != nil {
		return err
	}
	if err := c.route.srv.enc(c.res, c.req, v); err != nil {
		return err
	}
	return nil
}
func (c *wrapper) Reset(res http.ResponseWriter, req *http.Request) {
	c.res = res
	c.req = req
}

func (c *wrapper) Deadline() (time.Time, bool) {
	return c.req.Context().Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	return c.req.Context().Done()
}

func (c *wrapper) Err() error {
	return c.req.Context().Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	return c.req.Context().Value(key)
}
