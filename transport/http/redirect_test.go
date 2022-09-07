package http

import (
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	var (
		redirectURL  = "/redirect"
		redirectCode = 302
	)
	r := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()
	_ = DefaultResponseEncoder(w, r, NewRedirect(redirectURL, redirectCode))

	if w.Code != redirectCode {
		t.Fatalf("want %d but got %d", redirectCode, w.Code)
	}
	if v := w.Header().Get("Location"); v != redirectURL {
		t.Fatalf("want %s but got %s", redirectURL, v)
	}
}
