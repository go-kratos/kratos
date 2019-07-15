package zipkin

import (
	"fmt"
	"time"

	"github.com/bilibili/kratos/pkg/net/trace"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

type report struct {
	rpt reporter.Reporter
}

func newReport(c *Config) *report {
	return &report{
		rpt: http.NewReporter(c.Endpoint,
			http.Timeout(time.Duration(c.Timeout)),
			http.BatchSize(c.BatchSize),
		),
	}
}

// WriteSpan write a trace span to queue.
func (r *report) WriteSpan(raw *trace.Span) (err error) {
	traceID := model.TraceID{Low: raw.Tid()}
	spanID := model.ID(raw.SpanID())
	parentID := model.ID(raw.ParentID())
	span := model.SpanModel{
		SpanContext: model.SpanContext{
			TraceID:  traceID,
			ID:       spanID,
			ParentID: &parentID,
		},
		Name:      raw.Name(),
		Timestamp: raw.StartTime(),
		Duration:  raw.Duration(),
		Tags:      make(map[string]string, len(raw.Tags())),
	}
	for _, tag := range raw.Tags() {
		switch tag.Key {
		case trace.TagSpanKind:
			span.Kind = model.Kind(tag.Value.(string))
		case trace.TagPeerService:
			span.LocalEndpoint = &model.Endpoint{ServiceName: tag.Value.(string)}
		default:
			v, ok := tag.Value.(string)
			if ok {
				span.Tags[tag.Key] = v
			} else {
				span.Tags[tag.Key] = fmt.Sprint(v)
			}
		}
	}
	r.rpt.Send(span)
	return
}

// Close close the report.
func (r *report) Close() error {
	return r.rpt.Close()
}
