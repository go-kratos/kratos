package httpcache

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darkweak/souin/configurationtypes"
	"github.com/darkweak/souin/plugins"
)

var (
	devDefaultConfiguration = plugins.BaseConfiguration{
		API: configurationtypes.API{
			BasePath: "/httpcache_api",
			Prometheus: configurationtypes.APIEndpoint{
				Enable: true,
			},
			Souin: configurationtypes.APIEndpoint{
				BasePath: "/httpcache",
				Enable:   true,
			},
		},
		DefaultCache: &configurationtypes.DefaultCache{
			Regex: configurationtypes.Regex{
				Exclude: "/excluded",
			},
			TTL: configurationtypes.Duration{
				Duration: time.Second,
			},
		},
		LogLevel: "debug",
	}
)

type next struct{}

func (n *next) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello Kratos!"))
}

var nextFilter = &next{}

func prepare(endpoint string) (req *http.Request, res1 *httptest.ResponseRecorder, res2 *httptest.ResponseRecorder) {
	req = httptest.NewRequest(http.MethodGet, endpoint, nil)
	res1 = httptest.NewRecorder()
	res2 = httptest.NewRecorder()

	return
}

func Test_HttpcacheKratosPlugin_NewHTTPCacheFilterHandler(t *testing.T) {
	if NewHTTPCacheFilter(devDefaultConfiguration) == nil {
		t.Error("The NewHTTPCacheFilter method must return an HTTP Handler.")
	}
}

func Test_HttpcacheKratosPlugin_NewHTTPCacheFilter(t *testing.T) {
	time.Sleep(time.Second)
	handler := NewHTTPCacheFilter(devDefaultConfiguration)(nextFilter)
	req, res, res2 := prepare("/handled")
	handler.ServeHTTP(res, req)
	if res.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss; stored" {
		t.Error("The response must contain a Cache-Status header with the stored directive.")
	}
	handler.ServeHTTP(res2, req)
	if res2.Result().Header.Get("Cache-Status") != "Souin; hit; ttl=0" {
		t.Error("The response must contain a Cache-Status header with the hit and ttl directives.")
	}
	if res2.Result().Header.Get("Age") != "1" {
		t.Error("The response must contain a Age header with the value 1.")
	}
}

func Test_HttpcacheKratosPlugin_NewHTTPCacheFilter_Excluded(t *testing.T) {
	time.Sleep(time.Second)
	handler := NewHTTPCacheFilter(devDefaultConfiguration)(nextFilter)
	req, res, res2 := prepare("/excluded")
	handler.ServeHTTP(res, req)
	if res.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss" {
		t.Error("The response must contain a Cache-Status header without the stored directive and with the uri-miss only.")
	}
	handler.ServeHTTP(res2, req)
	if res2.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss" {
		t.Error("The response must contain a Cache-Status header without the stored directive and with the uri-miss only.")
	}
	if res2.Result().Header.Get("Age") != "" {
		t.Error("The response must not contain a Age header.")
	}
}

func Test_HttpcacheKratosPlugin_NewHTTPCacheFilter_Mutation(t *testing.T) {
	config := devDefaultConfiguration
	config.DefaultCache.AllowedHTTPVerbs = []string{http.MethodGet, http.MethodPost}
	time.Sleep(time.Second)
	handler := NewHTTPCacheFilter(config)(nextFilter)
	req, res, res2 := prepare("/handled")
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"query":"mutation":{something mutated}}`)))
	handler.ServeHTTP(res, req)
	if res.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss" {
		t.Error("The response must contain a Cache-Status header without the stored directive and with the uri-miss only.")
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"query":"mutation":{something mutated}}`)))
	handler.ServeHTTP(res2, req)
	if res2.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss" {
		t.Error("The response must contain a Cache-Status header without the stored directive and with the uri-miss only.")
	}
	if res2.Result().Header.Get("Age") != "" {
		t.Error("The response must not contain a Age header.")
	}
}

func Test_HttpcacheKratosPlugin_NewHTTPCacheFilter_API(t *testing.T) {
	time.Sleep(time.Second)
	handler := NewHTTPCacheFilter(devDefaultConfiguration)(nextFilter)
	req, res, res2 := prepare("/httpcache_api/httpcache")
	handler.ServeHTTP(res, req)
	if res.Result().Header.Get("Content-Type") != "application/json" {
		t.Error("The response must contain be in JSON.")
	}
	b, _ := ioutil.ReadAll(res.Result().Body)
	res.Result().Body.Close()
	if string(b) != "[]" {
		t.Error("The response body must be an empty array because no request has been stored")
	}
	req2 := httptest.NewRequest(http.MethodGet, "/handled", nil)
	handler.ServeHTTP(res2, req2)
	if res2.Result().Header.Get("Cache-Status") != "Souin; fwd=uri-miss; stored" {
		t.Error("The response must contain a Cache-Status header with the stored directive.")
	}
	res3 := httptest.NewRecorder()
	handler.ServeHTTP(res3, req)
	if res3.Result().Header.Get("Content-Type") != "application/json" {
		t.Error("The response must contain be in JSON.")
	}
	b, _ = ioutil.ReadAll(res3.Result().Body)
	res3.Result().Body.Close()
	var payload []string
	_ = json.Unmarshal(b, &payload)
	if len(payload) != 2 {
		t.Error("The system must store 2 items, the fresh and the stale one")
	}
	if payload[0] != "GET-example.com-/handled" || payload[1] != "STALE_GET-example.com-/handled" {
		t.Error("The payload items mismatch from the expectations.")
	}
}
