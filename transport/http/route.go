package http

import (
	"net/http"
	"path"
	"sync"
)

// Route is an HTTP route.
type Route struct {
	prefix string
	pool   sync.Pool
	srv    *Server
}

func newRoute(prefix string, srv *Server) *Route {
	r := &Route{
		prefix: prefix,
		srv:    srv,
	}
	r.pool.New = func() interface{} {
		return new(wrapper)
	}
	return r
}

// Handle registers a new route with a matcher for the URL path and method.
func (r *Route) Handle(method, relativePath string, h HandlerFunc) {
	r.srv.router.HandleFunc(path.Join(r.prefix, relativePath), func(res http.ResponseWriter, req *http.Request) {
		ctx := r.pool.Get().(Context)
		ctx.Reset(r, res, req)
		if err := h(ctx); err != nil {
			r.srv.ene(res, req, err)
		}
		r.pool.Put(ctx)
	}).Methods(method)
}

// GET registers a new GET route for a path with matching handler in the router.
func (r *Route) GET(path string, h HandlerFunc) { r.Handle(http.MethodGet, path, h) }

// HEAD registers a new HEAD route for a path with matching handler in the router.
func (r *Route) HEAD(path string, h HandlerFunc) { r.Handle(http.MethodHead, path, h) }

// POST registers a new POST route for a path with matching handler in the router.
func (r *Route) POST(path string, h HandlerFunc) { r.Handle(http.MethodPost, path, h) }

// PUT registers a new PUT route for a path with matching handler in the router.
func (r *Route) PUT(path string, h HandlerFunc) { r.Handle(http.MethodPut, path, h) }

// PATCH registers a new PATCH route for a path with matching handler in the router.
func (r *Route) PATCH(path string, h HandlerFunc) { r.Handle(http.MethodPatch, path, h) }

// DELETE registers a new DELETE route for a path with matching handler in the router.
func (r *Route) DELETE(path string, h HandlerFunc) { r.Handle(http.MethodDelete, path, h) }

// CONNECT registers a new CONNECT route for a path with matching handler in the router.
func (r *Route) CONNECT(path string, h HandlerFunc) { r.Handle(http.MethodConnect, path, h) }

// OPTIONS registers a new OPTIONS route for a path with matching handler in the router.
func (r *Route) OPTIONS(path string, h HandlerFunc) { r.Handle(http.MethodOptions, path, h) }

// TRACE registers a new TRACE route for a path with matching handler in the router.
func (r *Route) TRACE(path string, h HandlerFunc) { r.Handle(http.MethodTrace, path, h) }
