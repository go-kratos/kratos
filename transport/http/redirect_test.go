package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	var (
		redirectURL  = "/redirect"
		redirectCode = 302
	)
	r := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	_ = DefaultResponseEncoder(w, r, NewRedirect(redirectURL, redirectCode))

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != redirectCode {
		t.Fatalf("want %d but got %d", redirectCode, resp.StatusCode)
	}
	if v := resp.Header.Get("Location"); v != redirectURL {
		t.Fatalf("want %s but got %s", redirectURL, v)
	}
}
