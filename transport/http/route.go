package http

import (
	"net/http"
	"path"
	"sync"
)

// FilterFunc is a function which receives an http.Handler and returns another http.Handler.
type FilterFunc func(http.Handler) http.Handler

// FilterChain returns a FilterFunc that specifies the chained handler for HTTP Router.
func FilterChain(filters ...FilterFunc) FilterFunc {
	return func(next http.Handler) http.Handler {
		for i := len(filters) - 1; i >= 0; i-- {
			next = filters[i](next)
		}
		return next
	}
}

// Route is an HTTP route.
type Route struct {
	prefix  string
	pool    sync.Pool
	srv     *Server
	filters []FilterFunc
}

func newRoute(prefix string, srv *Server, filters ...FilterFunc) *Route {
	r := &Route{
		prefix:  prefix,
		srv:     srv,
		filters: filters,
	}
	r.pool.New = func() interface{} {
		return &wrapper{route: r}
	}
	return r
}

// Handle registers a new route with a matcher for the URL path and method.
func (r *Route) Handle(method, relativePath string, h HandlerFunc, filters ...FilterFunc) {
	next := http.Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := r.pool.Get().(Context)
		ctx.Reset(res, req)
		if err := h(ctx); err != nil {
			r.srv.ene(res, req, err)
		}
		ctx.Reset(nil, nil)
		r.pool.Put(ctx)
	}))
	next = FilterChain(filters...)(next)
	next = FilterChain(r.filters...)(next)
	r.srv.router.Handle(path.Join(r.prefix, relativePath), next).Methods(method)
}

// GET registers a new GET route for a path with matching handler in the router.
func (r *Route) GET(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodGet, path, h, m...)
}

// HEAD registers a new HEAD route for a path with matching handler in the router.
func (r *Route) HEAD(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodHead, path, h, m...)
}

// POST registers a new POST route for a path with matching handler in the router.
func (r *Route) POST(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodPost, path, h, m...)
}

// PUT registers a new PUT route for a path with matching handler in the router.
func (r *Route) PUT(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodPut, path, h, m...)
}

// PATCH registers a new PATCH route for a path with matching handler in the router.
func (r *Route) PATCH(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodPatch, path, h, m...)
}

// DELETE registers a new DELETE route for a path with matching handler in the router.
func (r *Route) DELETE(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodDelete, path, h, m...)
}

// CONNECT registers a new CONNECT route for a path with matching handler in the router.
func (r *Route) CONNECT(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodConnect, path, h, m...)
}

// OPTIONS registers a new OPTIONS route for a path with matching handler in the router.
func (r *Route) OPTIONS(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodOptions, path, h, m...)
}

// TRACE registers a new TRACE route for a path with matching handler in the router.
func (r *Route) TRACE(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodTrace, path, h, m...)
}
