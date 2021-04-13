package health

import (
	"context"
	"fmt"
	"net/http"
)

// Checker wraps the CheckHealth method.
//
// CheckHealth returns nil if the resource is healthy, or a non-nil
// error if the resource is not healthy.  CheckHealth must be safe to
// call from multiple goroutines.
type Checker interface {
	CheckHealth(ctx context.Context) error
}

// Option is a health option.
type Option func(*Handler)

// WithChecker with health checker.
func WithChecker(c ...Checker) Option {
	return func(o *Handler) {
		o.checkers = c
	}
}

// Handler is an HTTP handler that reports on the success of an
// aggregate of Checkers.  The zero value is always healthy.
type Handler struct {
	checkers []Checker
}

// NewHandler new a health handler.
func NewHandler(opts ...Option) *Handler {
	h := &Handler{}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Add adds a new check to the handler.
func (h *Handler) Add(c Checker) {
	h.checkers = append(h.checkers, c)
}

// ServeHTTP returns 200 if it is healthy, 500 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, c := range h.checkers {
		if err := c.CheckHealth(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "unhealthy: %s", err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}
