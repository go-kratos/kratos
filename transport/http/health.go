package http

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2"
	"net/http"
)

type Handler struct {
}

type healthResponse struct {
	Status string `json:"status"`
}

// NewHandler new a health handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP returns 200 if it is healthy, 500 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	info, _ := kratos.FromContext(r.Context())
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(healthResponse{
		Status: info.Health().GetStatus().String(),
	})
}