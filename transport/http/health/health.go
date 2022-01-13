package health

import (
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/health"
)

type Handler struct{}

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
	status, ok := info.Health().GetStatus(r.URL.Query().Get("service"))
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	switch status {
	case health.Status_SERVING:
		w.WriteHeader(http.StatusOK)
	case health.Status_NOT_SERVING:
		w.WriteHeader(http.StatusServiceUnavailable)
	case health.Status_SERVICE_UNKNOWN:
		w.WriteHeader(http.StatusNotImplemented)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(healthResponse{
		Status: status.String(),
	})
}
