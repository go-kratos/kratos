package tracing

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	_ transport.Transporter = &Transport{}
)

type headerCarrier http.Header

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	header    headerCarrier
}

func (tr *Transport) Kind() transport.Kind            { return tr.kind }
func (tr *Transport) Endpoint() string                { return tr.endpoint }
func (tr *Transport) Operation() string               { return tr.operation }
func (tr *Transport) RequestHeader() transport.Header { return tr.header }
func (tr *Transport) ReplyHeader() transport.Header   { return tr.header }

func TestTracing(t *testing.T) {
	var carrier = headerCarrier{}
	tp := tracesdk.NewTracerProvider(tracesdk.WithSampler(tracesdk.TraceIDRatioBased(0)))

	// caller use Inject
	tracer := NewTracer(trace.SpanKindClient, WithTracerProvider(tp), WithPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})))
	ts := &Transport{kind: transport.KindHTTP, header: carrier}

	ctx, aboveSpan := tracer.Start(transport.NewClientContext(context.Background(), ts), ts.Kind().String(), ts.Operation(), ts.RequestHeader())
	defer tracer.End(ctx, aboveSpan, nil)

	// server use Extract fetch traceInfo from carrier
	tracer = NewTracer(trace.SpanKindServer, WithPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})))
	ts = &Transport{kind: transport.KindHTTP, header: carrier}

	ctx, span := tracer.Start(transport.NewServerContext(ctx, ts), ts.Kind().String(), ts.Operation(), ts.RequestHeader())
	defer tracer.End(ctx, span, nil)

	if aboveSpan.SpanContext().TraceID() != span.SpanContext().TraceID() {
		t.Fatalf("TraceID failed to deliver")
	}

	if v, ok := transport.FromClientContext(ctx); !ok || len(v.RequestHeader().Keys()) == 0 {
		t.Fatalf("traceHeader failed to deliver")
	}
}
