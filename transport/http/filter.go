package http

import "net/http"

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
