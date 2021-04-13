package health

import (
	"context"
	"encoding/json"
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

// Handler is an HTTP handler that reports on the success of an
// aggregate of Checkers.  The zero value is always healthy.
type Handler struct {
	checkers map[string]Checker
}

// NewHandler new a health handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Add adds a new check to the handler.
func (h *Handler) Add(name string, c Checker) {
	h.checkers[name] = c
}

// ServeHTTP returns 200 if it is healthy, 500 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	res := make(map[string]string, len(h.checkers))
	for name, c := range h.checkers {
		if err := c.CheckHealth(r.Context()); err != nil {
			code = http.StatusInternalServerError
			res[name] = err.Error()
		} else {
			res[name] = "ok"
		}
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
