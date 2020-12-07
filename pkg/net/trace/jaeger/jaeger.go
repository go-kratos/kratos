package jaeger

import (
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/trace"
)

type Config struct {
	Endpoint  string
	BatchSize int
}

type JaegerReporter struct {
	transport *HTTPTransport
}

func newReport(c *Config) *JaegerReporter {
	transport := NewHTTPTransport(c.Endpoint)
	transport.batchSize = c.BatchSize
	return &JaegerReporter{transport: transport}
}

func (r *JaegerReporter) WriteSpan(raw *trace.Span) (err error) {
	ctx := raw.Context()
	traceID := TraceID{Low: ctx.TraceID}
	spanID := SpanID(ctx.SpanID)
	parentID := SpanID(ctx.ParentID)
	tags := raw.Tags()
	log.Info("[info] write span")
	span := &Span{
		context:       NewSpanContext(traceID, spanID, parentID, true, nil),
		operationName: raw.OperationName(),
		startTime:     raw.StartTime(),
		duration:      raw.Duration(),
	}

	span.serviceName = raw.ServiceName()

	for _, t := range tags {
		span.SetTag(t.Key, t.Value)
	}

	r.transport.Append(span)
	return nil
}

func (rpt *JaegerReporter) Close() error {
	return rpt.transport.Close()
}
