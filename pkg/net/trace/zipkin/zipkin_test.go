package zipkin

import (
	"net/http"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/net/trace"
	xtime "github.com/bilibili/kratos/pkg/time"
)

func TestZipkin(t *testing.T) {
	c := &Config{
		Endpoint:  "http://127.0.0.1:9411/api/v2/spans",
		Timeout:   xtime.Duration(time.Second * 5),
		BatchSize: 100,
	}
	report := newReport(c)
	t1 := trace.NewTracer("service1", report, true)
	t2 := trace.NewTracer("service2", report, true)
	sp1 := t1.New("opt_1")
	sp2 := sp1.Fork("", "opt_client")
	header := make(http.Header)
	t1.Inject(sp2, trace.HTTPFormat, header)
	sp3, err := t2.Extract(trace.HTTPFormat, header)
	if err != nil {
		t.Fatal(err)
	}
	sp3.Finish(nil)
	sp2.Finish(nil)
	sp1.Finish(nil)
	report.Close()
}
