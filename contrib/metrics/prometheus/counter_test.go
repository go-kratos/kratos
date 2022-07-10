package prometheus

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func gatherLatest(reg *prometheus.Registry) (result string, err error) {
	mfs, err := reg.Gather()
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}

	enc := expfmt.NewEncoder(buf, expfmt.FmtText)
	for _, mf := range mfs {
		if err = enc.Encode(mf); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func TestCounter(t *testing.T) {
	expect := `# HELP test_request_test_metric test
# TYPE test_request_test_metric counter
test_request_test_metric{code="test",kind="test",operation="test",reason="test"} %d
`

	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "test",
		Name:      "test_metric",
		Subsystem: "request",
		Help:      "test",
	}, []string{"kind", "operation", "code", "reason"})

	counter := NewCounter(counterVec)
	counter.With("test", "test", "test", "test").Inc()

	reg := prometheus.NewRegistry()
	reg.MustRegister(counterVec)

	result, err := gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect, 1) {
		t.Fatal("metrics error")
	}

	counter.With("test", "test", "test", "test").Add(10)
	result, err = gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect, 11) {
		t.Fatal("metrics error")
	}
}
