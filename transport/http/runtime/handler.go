package runtime

import (
	"net/http"
)

// Handler is like http.Handler except ServeHTTP may return an error.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, ParameterView) error
}

// Plugin is a middle layer which represents the traditional
// idea of plugin: it chains one Handler to the next by being
// passed the next Handler in the chain.
type Plugin func(Handler)

// HandlerFunc is a convenience type like http.HandlerFunc, except
// ServeHTTP returns an error. See Handler documentation for more information.
type HandlerFunc func(http.ResponseWriter, *http.Request, ParameterView) error

// ServeHTTP inplements the Handler interface.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, v ParameterView) error {
	return f(w, r, v)
}
