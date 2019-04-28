package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDapperProxy(t *testing.T) {
	dp, err := newDapperProxy("http://172.16.38.143:6193")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "http://example.com/x/admin/apm/dapper/familys", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	dp.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("proxy dapper error: %+v", resp)
	}
}
