package prometheus

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func intToFloatString(in int) string {
	return strconv.FormatFloat(float64(in), 'f', -1, 64)
}

func TestHistogram(t *testing.T) {
	expect := `# HELP test_request_test_metrics test
# TYPE test_request_test_metrics histogram
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="0.05"} %s
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="0.1"} %s
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="0.25"} %s
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="0.5"} %s
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="1"} %s
test_request_test_metrics_bucket{code="test",kind="test",operation="test",reason="test",le="+Inf"} %s
test_request_test_metrics_sum{code="test",kind="test",operation="test",reason="test"} %s
test_request_test_metrics_count{code="test",kind="test",operation="test",reason="test"} %s
`

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "test",
		Name:      "test_metrics",
		Subsystem: "request",
		Help:      "test",
		Buckets:   []float64{0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation", "code", "reason"})

	histogram := NewHistogram(histogramVec)
	histogram.With("test", "test", "test", "test").Observe(0.5)

	reg := prometheus.NewRegistry()
	reg.MustRegister(histogramVec)

	result, err := gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect,
		intToFloatString(0),
		intToFloatString(0),
		intToFloatString(0),
		intToFloatString(1),
		intToFloatString(1),
		intToFloatString(1),
		"0.5",
		intToFloatString(1)) {
		t.Fatal("metrics error")
	}

	histogram.With("test", "test", "test", "test").Observe(0.1)
	result, err = gatherLatest(reg)

	if err != nil {
		t.Fatal(err)
	}
	if result != fmt.Sprintf(expect,
		intToFloatString(0),
		intToFloatString(1),
		intToFloatString(1),
		intToFloatString(2),
		intToFloatString(2),
		intToFloatString(2),
		"0.6",
		intToFloatString(2),
	) {
		t.Fatal("metrics error")
	}
}
