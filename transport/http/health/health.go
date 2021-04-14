package health

import (
	"context"
	"encoding/json"
	"net/http"
)

// CheckerFunc wraps the CheckHealth method.
//
// CheckHealth returns nil if the resource is healthy, or a non-nil
// error if the resource is not healthy.  CheckHealth must be safe to
// call from multiple goroutines.
type CheckerFunc func(ctx context.Context) error

// Handler is an HTTP handler that reports on the success of an
// aggregate of Checkers.  The zero value is always healthy.
type Handler struct {
	checkers  map[string]CheckerFunc
	observers map[string]CheckerFunc
}

// NewHandler new a health handler.
func NewHandler() *Handler {
	return &Handler{
		checkers:  make(map[string]CheckerFunc),
		observers: make(map[string]CheckerFunc),
	}
}

// AddChecker adds a new check to the handler.
func (h *Handler) AddChecker(name string, c CheckerFunc) {
	h.checkers[name] = c
}

// AddObserver adds a new check to the handler but it does not fail the entire status.
func (h *Handler) AddObserver(name string, c CheckerFunc) {
	h.observers[name] = c
}

// ServeHTTP returns 200 if it is healthy, 500 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	results := make(map[string]string, len(h.checkers))

	for name, checker := range h.checkers {
		if err := checker(r.Context()); err != nil {
			code = http.StatusInternalServerError
			results[name] = err.Error()
		} else {
			results[name] = "OK"
		}
	}

	for name, checker := range h.observers {
		if err := checker(r.Context()); err != nil {
			results[name] = err.Error()
		} else {
			results[name] = "OK"
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(results)
}
