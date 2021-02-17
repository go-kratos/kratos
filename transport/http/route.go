package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RouteGroup adds a matcher for the URL path and method. This matches if the given
// template is a prefix of the full URL path. See route.Path() for details on
// the tpl argument.
type RouteGroup struct {
	prefix string
	router *mux.Router
}

// ANY maps an HTTP Any request to the path and the specified handler.
func (r *RouteGroup) ANY(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler)
}

// GET maps an HTTP Get request to the path and the specified handler.
func (r *RouteGroup) GET(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("GET")
}

// HEAD maps an HTTP Head request to the path and the specified handler.
func (r *RouteGroup) HEAD(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("HEAD")
}

// POST maps an HTTP Post request to the path and the specified handler.
func (r *RouteGroup) POST(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("POST")
}

// PUT maps an HTTP Put request to the path and the specified handler.
func (r *RouteGroup) PUT(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("PUT")
}

// DELETE maps an HTTP Delete request to the path and the specified handler.
func (r *RouteGroup) DELETE(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("DELETE")
}

// PATCH maps an HTTP Patch request to the path and the specified handler.
func (r *RouteGroup) PATCH(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("PATCH")
}

// OPTIONS maps an HTTP Options request to the path and the specified handler.
func (r *RouteGroup) OPTIONS(path string, handler http.HandlerFunc) {
	r.router.PathPrefix(r.prefix).Path(path).HandlerFunc(handler).Methods("OPTIONS")
}
