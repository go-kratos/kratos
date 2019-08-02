package zipkin

import (
	"fmt"
	protogen "github.com/bilibili/kratos/pkg/net/trace/proto"
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
	ctx := raw.Context()
	traceID := model.TraceID{Low: ctx.TraceID}
	spanID := model.ID(ctx.SpanID)
	parentID := model.ID(ctx.ParentID)
	tags := raw.Tags()
	span := model.SpanModel{
		SpanContext: model.SpanContext{
			TraceID:  traceID,
			ID:       spanID,
			ParentID: &parentID,
		},
		Name:      raw.OperationName(),
		Timestamp: raw.StartTime(),
		Duration:  raw.Duration(),
		Tags:      make(map[string]string, len(tags)),
	}
	span.LocalEndpoint = &model.Endpoint{ServiceName: raw.ServiceName()}
	for _, tag := range tags {
		switch tag.Key {
		case trace.TagSpanKind:
			switch tag.Value.(string) {
			case "client":
				span.Kind = model.Client
			case "server":
				span.Kind = model.Server
			case "producer":
				span.Kind = model.Producer
			case "consumer":
				span.Kind = model.Consumer
			}
		default:
			v, ok := tag.Value.(string)
			if ok {
				span.Tags[tag.Key] = v
			} else {
				span.Tags[tag.Key] = fmt.Sprint(v)
			}
		}
	}
	//log save to zipkin annotation
	span.Annotations = r.converLogsToAnnotations(raw.Logs())
	r.rpt.Send(span)
	return
}

func (r *report) converLogsToAnnotations(logs []*protogen.Log) (annotations []model.Annotation) {
	annotations = make([]model.Annotation, 0, len(annotations))
	for _, lg := range logs {
		annotations = append(annotations, r.converLogToAnnotation(lg)...)
	}
	return annotations
}
func (r *report) converLogToAnnotation(log *protogen.Log) (annotations []model.Annotation) {
	annotations = make([]model.Annotation, 0, len(log.Fields))
	for _, field := range log.Fields {
		val := string(field.Value)
		annotation := model.Annotation{
			Timestamp: time.Unix(0, log.Timestamp),
			Value:     field.Key + " : " + val,
		}
		annotations = append(annotations, annotation)
	}
	return annotations
}

// Close close the report.
func (r *report) Close() error {
	return r.rpt.Close()
}
