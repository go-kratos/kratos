package health

import (
	"context"
	"encoding/json"
	"net/http"
)

type Status string

const (
	Up   Status = "UP"
	Down Status = "DOWN"
)

type Checker interface {
	Check(ctx context.Context) error
}

type CheckerFunc func(ctx context.Context) error

func (f CheckerFunc) Check(ctx context.Context) error {
	return f(ctx)
}

type Health struct {
	status   Status
	checkers map[string]Checker
}

func New() *Health {
	h := &Health{
		status:   Down,
		checkers: make(map[string]Checker),
	}
	return h
}

func (h *Health) Register(name string, checker Checker) {
	h.checkers[name] = checker
}

func (h *Health) Start(_ context.Context) error {
	h.status = Up
	return nil
}

func (h *Health) Stop(_ context.Context) error {
	h.status = Down
	return nil
}

func (h *Health) Check(ctx context.Context) Result {
	res := Result{Status: h.status}
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

func (h *Health) CheckService(ctx context.Context, svc string) Detail {
	c, ok := h.checkers[svc]
	if !ok {
		return Detail{Status: Down, Error: "service not find"}
	}
	err := c.Check(ctx)
	if err != nil {
		return Detail{Status: Down, Error: err.Error()}
	}
	return Detail{Status: Up}
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	if service == "" {
		res := h.Check(r.Context())
		if res.Status == Down {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(res)
	} else {
		detail := h.CheckService(r.Context(), service)
		if detail.Status == Down {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(detail)
	}
}

type Result struct {
	Status  Status            `json:"status"`
	Details map[string]Detail `json:"details"`
}

type Detail struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}
