package model

// const for SpanPoint
const ()

// SamplePoint SamplePoint
type SamplePoint struct {
	TraceID uint64
	SpanID  uint64
	Value   int64
}

// SpanPoint contains time series
type SpanPoint struct {
	Timestamp     int64
	ServiceName   string
	OperationName string
	PeerService   string
	SpanKind      string
	AvgDuration   SamplePoint // random sample point
	MaxDuration   SamplePoint
	MinDuration   SamplePoint
	Errors        []SamplePoint
}
