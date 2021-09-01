package metrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

var defaultClient *statsd.Client

func init() {
	var err error
	// use the default env: DD_AGENT_HOST/DD_DOGSTATSD_PORT/8125
	if defaultClient, err = statsd.New(""); err != nil {
		panic(fmt.Sprintf("datadog client initialize error %v\n", err))
	}
}

func withValues(labels, values []string) []string {
	if len(labels) != len(values) {
		return nil
	}
	lvs := make([]string, 0, len(labels))
	for i, v := range values {
		lvs = append(lvs, labels[i]+":"+v)
	}
	return lvs
}
