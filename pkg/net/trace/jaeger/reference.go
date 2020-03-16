package jaeger

import "github.com/opentracing/opentracing-go"

// Reference represents a causal reference to other Spans (via their SpanContext).
type Reference struct {
	Type    opentracing.SpanReferenceType
	Context SpanContext
}
