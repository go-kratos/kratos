package metrics

import "time"

// A Summary captures individual observations from an event or sample stream and
// summarizes them in a manner similar to traditional summary statistics: 1. sum
// of observations, 2. observation count, 3. rank estimations.
//
// A typical use-case is the observation of request latencies. By default, a
// Summary provides the median, the 90th and the 99th percentile of the latency
// as rank estimations. However, the default behavior will change in the
// upcoming v1.0.0 of the library. There will be no rank estimations at all by
// default. For a sane transition, it is recommended to set the desired rank
// estimations explicitly.
type Summary interface {
	Metric
	// Observe adds a single observation to the summary.
	Observe(float64)
}

// SummaryOpts bundles the options for creating a Summary metric. It is
// mandatory to set Name to a non-empty string. While all other fields are
// optional and can safely be left at their zero value, it is recommended to set
// a help string and to explicitly set the Objectives field to the desired value
// as the default value will change in the upcoming v1.0.0 of the library.
type SummaryOpts struct {
	// Namespace, Subsystem, and Name are components of the fully-qualified
	// name of the Summary (created by joining these components with
	// "_"). Only Name is mandatory, the others merely help structuring the
	// name. Note that the fully-qualified name of the Summary must be a
	// valid Prometheus metric name.
	Namespace string
	Subsystem string
	Name      string

	// Help provides information about this Summary.
	//
	// Metrics with the same fully-qualified name must have the same Help
	// string.
	Help string

	// ConstLabels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their ConstLabels.
	//
	// Due to the way a Summary is represented in the Prometheus text format
	// and how it is handled by the Prometheus server internally, “quantile”
	// is an illegal label name. Construction of a Summary or SummaryVec
	// will panic if this label name is used in ConstLabels.
	//
	// ConstLabels are only used rarely. In particular, do not use them to
	// attach the same labels to all your metrics. Those use cases are
	// better covered by target labels set by the scraping Prometheus
	// server, or by one specific metric (e.g. a build_info or a
	// machine_role metric). See also
	// https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels
	ConstLabels map[string]string

	// Objectives defines the quantile rank estimates with their respective
	// absolute error. If Objectives[q] = e, then the value reported for q
	// will be the φ-quantile value for some φ between q-e and q+e.  The
	// default value is an empty map, resulting in a summary without
	// quantiles.
	Objectives map[float64]float64

	// MaxAge defines the duration for which an observation stays relevant
	// for the summary. Must be positive. The default value is DefMaxAge.
	MaxAge time.Duration

	// AgeBuckets is the number of buckets used to exclude observations that
	// are older than MaxAge from the summary. A higher number has a
	// resource penalty, so only increase it if the higher resolution is
	// really required. For very high observation rates, you might want to
	// reduce the number of age buckets. With only one age bucket, you will
	// effectively see a complete reset of the summary each time MaxAge has
	// passed. The default value is DefAgeBuckets.
	AgeBuckets uint32

	// BufCap defines the default sample stream buffer size.  The default
	// value of DefBufCap should suffice for most uses. If there is a need
	// to increase the value, a multiple of 500 is recommended (because that
	// is the internal buffer size of the underlying package
	// "github.com/bmizerany/perks/quantile").
	BufCap uint32
}
