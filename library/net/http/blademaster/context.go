package blademaster

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

const (
	_abortIndex int8 = math.MaxInt8 / 2
)

var (
	_openParen  = []byte("(")
	_closeParen = []byte(")")
)

// Context is the most important part. It allows us to pass variables between
// middleware, manage the flow, validate the JSON of a request and render a
// JSON response for example.
type Context struct {
	context.Context

	Request *http.Request
	Writer  http.ResponseWriter

	// flow control
	index    int8
	handlers []HandlerFunc

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	Error error

	method string
	engine *Engine
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in godoc.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		// only check method on last handler, otherwise middlewares
		// will never be effected if request method is not matched
		if c.index == s-1 && c.method != c.Request.Method {
			code := http.StatusMethodNotAllowed
			c.Error = ecode.MethodNotAllowed
			http.Error(c.Writer, http.StatusText(code), code)
			return
		}

		c.handlers[c.index](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = _abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= _abortIndex
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// Render http response with http code by a render instance.
func (c *Context) Render(code int, r render.Render) {
	r.WriteContentType(c.Writer)
	if code > 0 {
		c.Status(code)
	}

	if !bodyAllowedForStatus(code) {
		return
	}

	params := c.Request.Form

	cb := params.Get("callback")
	jsonp := cb != "" && params.Get("jsonp") == "jsonp"
	if jsonp {
		c.Writer.Write([]byte(cb))
		c.Writer.Write(_openParen)
	}

	if err := r.Render(c.Writer); err != nil {
		c.Error = err
		return
	}

	if jsonp {
		if _, err := c.Writer.Write(_closeParen); err != nil {
			c.Error = errors.WithStack(err)
		}
	}
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(data interface{}, err error) {
	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)
	// TODO app allow 5xx?
	/*
		if bcode.Code() == -500 {
			code = http.StatusServiceUnavailable
		}
	*/
	writeStatusCode(c.Writer, bcode.Code())
	c.Render(code, render.JSON{
		Code:    bcode.Code(),
		Message: bcode.Message(),
		Data:    data,
	})
}

// JSONMap serializes the given map as map JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSONMap(data map[string]interface{}, err error) {
	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)
	// TODO app allow 5xx?
	/*
		if bcode.Code() == -500 {
			code = http.StatusServiceUnavailable
		}
	*/
	writeStatusCode(c.Writer, bcode.Code())
	data["code"] = bcode.Code()
	if _, ok := data["message"]; !ok {
		data["message"] = bcode.Message()
	}
	c.Render(code, render.MapJSON(data))
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(data interface{}, err error) {
	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)
	// TODO app allow 5xx?
	/*
		if bcode.Code() == -500 {
			code = http.StatusServiceUnavailable
		}
	*/
	writeStatusCode(c.Writer, bcode.Code())
	c.Render(code, render.XML{
		Code:    bcode.Code(),
		Message: bcode.Message(),
		Data:    data,
	})
}

// Protobuf serializes the given struct as PB into the response body.
// It also sets the ContentType as "application/x-protobuf".
func (c *Context) Protobuf(data proto.Message, err error) {
	var (
		bytes []byte
	)

	code := http.StatusOK
	c.Error = err
	bcode := ecode.Cause(err)

	any := new(types.Any)
	if data != nil {
		if bytes, err = proto.Marshal(data); err != nil {
			c.Error = errors.WithStack(err)
			return
		}
		any.TypeUrl = "type.googleapis.com/" + proto.MessageName(data)
		any.Value = bytes
	}
	writeStatusCode(c.Writer, bcode.Code())
	c.Render(code, render.PB{
		Code:    int64(bcode.Code()),
		Message: bcode.Message(),
		Data:    any,
	})
}

// Bytes writes some data into the body stream and updates the HTTP code.
func (c *Context) Bytes(code int, contentType string, data ...[]byte) {
	c.Render(code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.String{Format: format, Data: values})
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

// BindWith bind req arg with parser.
func (c *Context) BindWith(obj interface{}, b binding.Binding) error {
	return c.mustBindWith(obj, b)
}

// Bind bind req arg with defult form binding.
func (c *Context) Bind(obj interface{}) error {
	return c.mustBindWith(obj, binding.Form)
}

// mustBindWith binds the passed struct pointer using the specified binding engine.
// It will abort the request with HTTP 400 if any error ocurrs.
// See the binding package.
func (c *Context) mustBindWith(obj interface{}, b binding.Binding) (err error) {
	if err = b.Bind(c.Request, obj); err != nil {
		c.Error = ecode.RequestErr
		c.Render(http.StatusOK, render.JSON{
			Code:    ecode.RequestErr.Code(),
			Message: err.Error(),
			Data:    nil,
		})
		c.Abort()
	}
	return
}

func writeStatusCode(w http.ResponseWriter, ecode int) {
	header := w.Header()
	header.Set("bili-status-code", strconv.FormatInt(int64(ecode), 10))
}
