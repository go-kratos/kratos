package prometheus

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestGuage(t *testing.T) {
	expect := `# HELP test_request_test_guage_metric test
# TYPE test_request_test_guage_metric gauge
test_request_test_guage_metric{code="test",kind="test",operation="test",reason="test"} %d
`

	guageVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      "test_guage_metric",
		Subsystem: "request",
		Help:      "test",
	}, []string{"kind", "operation", "code", "reason"})

	guage := NewGauge(guageVec)
	guage.With("test", "test", "test", "test").Set(1)

	reg := prometheus.NewRegistry()
	reg.MustRegister(guageVec)

	result, err := gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}

	if result != fmt.Sprintf(expect, 1) {
		t.Fatal("metrics error")
	}

	guage.With("test", "test", "test", "test").Add(1)
	result, err = gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}
	if result != fmt.Sprintf(expect, 2) {
		t.Fatal("metrics error")
	}

	guage.With("test", "test", "test", "test").Sub(1)
	result, err = gatherLatest(reg)
	if err != nil {
		t.Fatal(err)
	}
	if result != fmt.Sprintf(expect, 1) {
		t.Fatal("metrics error")
	}
}
