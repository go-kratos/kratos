package trace

import (
	"fmt"
	"strings"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"

	protogen "github.com/bilibili/kratos/pkg/net/trace/proto"
)

type zipkinHTTPReport struct {
	rpt zipkinreporter.Reporter
}

// NewZipKinHTTPReport report trace to zipkin.
func NewZipKinHTTPReport(endpoint string, batchSize int, timeout time.Duration) *zipkinHTTPReport {
	// TODO: support multi entrypoint and custom path.
	if !strings.HasPrefix(endpoint, "http://") {
		endpoint = fmt.Sprintf("http://%s/api/v2/spans", endpoint)
	}
	if batchSize == 0 {
		batchSize = 100
	}
	if timeout == 0 {
		timeout = 200 * time.Millisecond
	}
	return &zipkinHTTPReport{
		rpt: http.NewReporter(endpoint,
			http.Timeout(time.Duration(timeout)),
			http.BatchSize(batchSize),
		),
	}
}

// WriteSpan write a trace span to queue.
func (r *zipkinHTTPReport) WriteSpan(raw *Span) (err error) {
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
		case TagSpanKind:
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

func (r *zipkinHTTPReport) converLogsToAnnotations(logs []*protogen.Log) (annotations []model.Annotation) {
	annotations = make([]model.Annotation, 0, len(annotations))
	for _, lg := range logs {
		annotations = append(annotations, r.converLogToAnnotation(lg)...)
	}
	return annotations
}
func (r *zipkinHTTPReport) converLogToAnnotation(log *protogen.Log) (annotations []model.Annotation) {
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

// Close close the zipkinHTTPReport.
func (r *zipkinHTTPReport) Close() error {
	return r.rpt.Close()
}
