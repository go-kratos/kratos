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
	checkers  map[string]Checker
	observers map[string]Checker
}

// NewHandler new a health handler.
func NewHandler() *Handler {
	return &Handler{
		checkers:  make(map[string]Checker),
		observers: make(map[string]Checker),
	}
}

// AddChecker adds a new check to the handler.
func (h *Handler) AddChecker(name string, c Checker) {
	h.checkers[name] = c
}

// AddObserver adds a new check to the handler but it does not fail the entire status.
func (h *Handler) AddObserver(name string, c Checker) {
	h.observers[name] = c
}

// ServeHTTP returns 200 if it is healthy, 500 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	errors := make(map[string]string, len(h.checkers))

	for name, c := range h.checkers {
		if err := c.CheckHealth(r.Context()); err != nil {
			code = http.StatusInternalServerError
			errors[name] = err.Error()
		} else {
			errors[name] = "ok"
		}
	}

	for name, c := range h.observers {
		if err := c.CheckHealth(r.Context()); err != nil {
			errors[name] = err.Error()
		} else {
			errors[name] = "ok"
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"status": code,
			"errors": errors,
		},
	)
}
