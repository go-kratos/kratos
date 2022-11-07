package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Status string

var ErrServiceNotFind = errors.New("service not find")

const (
	Up   Status = "UP"
	Down Status = "DOWN"
)

type Checker interface {
	Check(ctx context.Context) error
}

type Result struct {
	Status  Status            `json:"status"`
	Details map[string]Detail `json:"details"`
}

type Detail struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}

type Health struct {
	Status
	checkers map[string]Checker
}

func New() *Health {
	return &Health{
		checkers: map[string]Checker{},
	}
}

func (h *Health) Register(name string, checker Checker) {
	h.checkers[name] = checker
}

func (h *Health) CheckAll(ctx context.Context) Result {
	res := Result{Status: h.Status}
	for n, c := range h.checkers {
		if err := c.Check(ctx); err != nil {
			res.Status = Down
			res.Details[n] = Detail{Status: Down, Error: err.Error()}
		} else {
			res.Details[n] = Detail{Status: Up}
		}
	}
	return res
}

func (h *Health) Check(ctx context.Context, service string) Detail {
	checker, ok := h.checkers[service]
	if !ok {
		return Detail{Down, ErrServiceNotFind.Error()}
	}
	d := Detail{}
	err := checker.Check(ctx)
	if err != nil {
		d.Status = Down
		d.Error = err.Error()
	}
	return d
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	if service == "" {
		res := h.CheckAll(r.Context())
		if res.Status == Down {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(res)
	} else {
		d := h.Check(r.Context(), service)
		if d.Status == Down {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(d)
	}
}
