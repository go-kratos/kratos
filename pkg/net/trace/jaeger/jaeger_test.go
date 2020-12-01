package jaeger

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kratos/kratos/pkg/net/trace"
)

func TestJaegerReporter(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", r.Method)
		}

		aSpanPayload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		t.Logf("%s\n", aSpanPayload)
	}
	ht := httptest.NewServer(http.HandlerFunc(handler))
	defer ht.Close()

	c := &Config{
		Endpoint:  ht.URL,
		BatchSize: 1,
	}

	//c.Endpoint = "http://127.0.0.1:14268/api/traces"

	report := newReport(c)
	t1 := trace.NewTracer("jaeger_test_1", report, true)
	t2 := trace.NewTracer("jaeger_test_2", report, true)
	sp1 := t1.New("option_1")
	sp2 := sp1.Fork("service3", "opt_client")
	sp2.SetLog(trace.Log("log_k", "log_v"))
	// inject
	header := make(http.Header)
	t1.Inject(sp2, trace.HTTPFormat, header)
	t.Log(header)
	sp3, err := t2.Extract(trace.HTTPFormat, header)
	if err != nil {
		t.Fatal(err)
	}
	sp3.Finish(nil)
	sp2.Finish(nil)
	sp1.Finish(nil)
	report.Close()
}
