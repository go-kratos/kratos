package http

import (
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/gorilla/mux"
)

var _ Context = (*wrapper)(nil)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(Context) error

// Context is an HTTP Context.
type Context interface {
	Vars() url.Values
	Form() url.Values
	Request() *http.Request
	Response() http.ResponseWriter
	Middleware() middleware.Middleware
	Bind(interface{}) error
	Result(int, interface{}) error
	Returns(int, interface{}, error) error
	Reset(*Route, http.ResponseWriter, *http.Request)
}

type wrapper struct {
	res   http.ResponseWriter
	req   *http.Request
	route *Route
}

func (c *wrapper) Vars() url.Values {
	raws := mux.Vars(c.req)
	vars := make(url.Values, len(raws))
	for k, v := range vars {
		vars[k] = v
	}
	return vars
}
func (c *wrapper) Form() url.Values {
	if err := c.req.ParseForm(); err != nil {
		return url.Values{}
	}
	return c.req.Form
}
func (c *wrapper) Request() *http.Request            { return c.req }
func (c *wrapper) Response() http.ResponseWriter     { return c.res }
func (c *wrapper) Middleware() middleware.Middleware { return c.route.srv.m }
func (c *wrapper) Bind(v interface{}) error          { return c.route.srv.dec(c.req, v) }
func (c *wrapper) Result(code int, v interface{}) error {
	if err := c.route.srv.enc(c.res, c.req, v); err != nil {
		return err
	}
	c.res.WriteHeader(code)
	return nil
}
func (c *wrapper) Returns(code int, v interface{}, err error) error {
	if err != nil {
		return err
	}
	if err := c.route.srv.enc(c.res, c.req, v); err != nil {
		return err
	}
	c.res.WriteHeader(code)
	return nil
}
func (c *wrapper) Reset(r *Route, res http.ResponseWriter, req *http.Request) {
	c.route = r
	c.res = res
	c.req = req
}
