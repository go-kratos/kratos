package prometheus

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestSummary(t *testing.T) {
	expect := `# HELP test_request_test_metric test
# TYPE test_request_test_metric summary
test_request_test_metric_sum{code="test",kind="test",operation="test",reason="test"} %s
test_request_test_metric_count{code="test",kind="test",operation="test",reason="test"} %s
`

	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "test",
		Name:      "test_metric",
		Subsystem: "request",
		Help:      "test",
	}, []string{"kind", "operation", "code", "reason"})

	summary := NewSummary(summaryVec)
	summary.With("test", "test", "test", "test").Observe(1)

	reg := prometheus.NewRegistry()
	reg.MustRegister(summaryVec)
	result, err := gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect, intToFloatString(1), intToFloatString(1)) {
		t.Fatal("metrics error")
	}
	summary.With("test", "test", "test", "test").Observe(10)

	result, err = gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect, intToFloatString(11), intToFloatString(2)) {
		t.Fatal("metrics error")
	}
}
