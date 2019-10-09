package trace

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZipkin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", r.Method)
		}

		aSpanPayload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		t.Logf("%s\n", aSpanPayload)
	}))
	defer ts.Close()
	//c.Endpoint = "http://127.0.0.1:9411/api/v2/spans"
	report := NewZipKinHTTPReport(ts.URL, 100, 0)
	t1 := NewTracer("service1", report, &testcfg)
	t2 := NewTracer("service2", report, &testcfg)
	sp1 := t1.New("option_1")
	sp2 := sp1.Fork("service3", "opt_client")
	sp2.SetLog(Log("log_k", "log_v"))
	// inject
	header := make(http.Header)
	t1.Inject(sp2, HTTPFormat, header)
	t.Log(header)
	sp3, err := t2.Extract(HTTPFormat, header)
	if err != nil {
		t.Fatal(err)
	}
	sp3.Finish(nil)
	sp2.Finish(nil)
	sp1.Finish(nil)
	report.Close()
}
