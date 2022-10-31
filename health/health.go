package health

import (
	"context"
	"encoding/json"
	"net/http"
)

type Option func(*Health)

type Status string

const (
	StatusUp   = Status("UP")
	StatusDown = Status("DOWN")
)

type Result struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}

// Checker returns nil if the resource is healthy, or a non-nil
// error if the resource is not healthy.
type Checker interface {
	CheckHealth(ctx context.Context) error
}

// CheckerFunc is an adapter type to allow the use of ordinary functions as
// health checks.
type CheckerFunc func() error

// CheckHealth calls f().
func (f CheckerFunc) CheckHealth() error {
	return f()
}

// Health .
type Health struct {
	liveness map[string]Checker
	readness map[string]Checker
}

func New() *Health {
	return &Health{
		liveness: map[string]Checker{},
		readness: map[string]Checker{},
	}
}

func (h *Health) AddLiveness(name string, checker Checker) {
	h.liveness[name] = checker
}

func (h *Health) AddReadness(name string, checker Checker) {
	h.readness[name] = checker
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var res struct {
		Status  Status
		Details map[string]Result `json:"details"`
	}
	res.Status = StatusUp
	for n, c := range h.liveness {
		if err := c.CheckHealth(r.Context()); err != nil {
			res.Status = StatusDown
			res.Details[n] = Result{Status: StatusDown, Error: err.Error()}
		} else {
			res.Details[n] = Result{Status: StatusUp}
		}
	}
	for n, c := range h.readness {
		if err := c.CheckHealth(r.Context()); err != nil {
			res.Details[n] = Result{StatusUp, err.Error()}
		} else {
			res.Details[n] = Result{Status: StatusUp}
		}
	}
	if res.Status == StatusDown {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(res)
}
