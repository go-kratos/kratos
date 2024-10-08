package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type dummyExporter struct {
	mu       sync.Mutex
	writeBuf bytes.Buffer
}

func (x *dummyExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return sdkmetric.DefaultTemporalitySelector(kind)
}

func (x *dummyExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(kind)
}

func (x *dummyExporter) Export(ctx context.Context, resourceMetrics *metricdata.ResourceMetrics) error {
	select {
	case <-ctx.Done():
		// Don't do anything if the context has already timed out.
		return ctx.Err()
	default:
		// Context is still valid, continue.
	}
	x.mu.Lock()
	defer x.mu.Unlock()
	return json.NewEncoder(&x.writeBuf).Encode(resourceMetrics)
}

func (x *dummyExporter) ForceFlush(ctx context.Context) error {
	return ctx.Err()
}

func (x *dummyExporter) Shutdown(ctx context.Context) error {
	return ctx.Err()
}

func (x *dummyExporter) String() string {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.writeBuf.String()
}

var exporter = &dummyExporter{}

func init() {
	err := EnableOTELExemplar()
	if err != nil {
		panic(err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Microsecond*400))),
		sdkmetric.WithView(DefaultSecondsHistogramView(DefaultServerSecondsHistogramName), func(instrument sdkmetric.Instrument) (sdkmetric.Stream, bool) {
			return sdkmetric.Stream{
				Name:        instrument.Name,
				Description: instrument.Description,
				Unit:        instrument.Unit,
				AttributeFilter: func(attribute.KeyValue) bool {
					return true
				},
			}, true
		}),
	)
	otel.SetMeterProvider(mp)

	tr := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tr)
}

func TestWithRequests(t *testing.T) {
	o := options{}

	if o.requests != nil {
		t.Errorf(`The type of the option requests property must be of "nil"`)
		return
	}

	meter := otel.Meter("test_meter")
	requests, err := meter.Int64Counter(DefaultServerRequestsCounterName)
	if err != nil {
		t.Errorf("[Int64Counter] something went wrong: %v", err)
		return
	}

	WithRequests(requests)(&o)

	if o.requests == nil {
		t.Errorf(`The type of the option requests property must be of "mockCounter", %T given.`, o.requests)
	}
}

func TestWithSeconds(t *testing.T) {
	o := options{}

	if o.seconds != nil {
		t.Errorf(`The type of the option seconds property must be of "nil"`)
		return
	}

	meter := otel.Meter("test_meter")
	seconds, err := meter.Float64Histogram(DefaultServerSecondsHistogramName)
	if err != nil {
		t.Errorf("[Float64Histogram] something went wrong: %v", err)
		return
	}

	WithSeconds(seconds)(&o)

	if o.seconds == nil {
		t.Errorf(`The type of the option seconds property must be of "mockObserver", %T given.`, o.requests)
	}
}

func TestServer(t *testing.T) {
	tracer := otel.Tracer("test_trace")
	ctx, span := tracer.Start(context.Background(), "TestServer")
	defer span.End()

	e := errors.New("got an error")
	nextError := func(context.Context, interface{}) (interface{}, error) {
		return nil, e
	}
	nextValid := func(context.Context, interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond * time.Duration(rand.Int31n(100)))
		return "Hello valid", nil
	}

	// init server handler
	meter := otel.Meter("test_meter")
	requests, err := DefaultRequestsCounter(meter, DefaultServerRequestsCounterName)
	if err != nil {
		t.Errorf("[DefaultRequestsCounter] something went wrong: %v", err)
		return
	}
	seconds, err := DefaultSecondsHistogram(meter, DefaultServerSecondsHistogramName)
	if err != nil {
		t.Errorf("[DefaultSecondsHistogram] something went wrong: %v", err)
		return
	}
	serverHandler := Server(
		WithRequests(requests),
		WithSeconds(seconds),
	)

	_, err = serverHandler(nextError)(ctx, "test:")
	if err != e {
		t.Error("The given error mismatch the expected.")
		return
	}

	for i := 0; i < 20; i++ {
		res, err := serverHandler(nextValid)(transport.NewServerContext(ctx, &http.Transport{}), "test:")
		if err != nil {
			t.Error("The server must not throw an error.")
		}
		if res != "Hello valid" {
			t.Error(`The server must return a "Hello valid" response.`)
		}
	}
	bufStr := exporter.String()

	for _, label := range []string{
		metricLabelKind,
		metricLabelOperation,
		metricLabelCode,
		metricLabelReason,
		"SpanID",
		"TraceID",
	} {
		if !strings.Contains(bufStr, fmt.Sprintf("\"%s\"", label)) {
			t.Errorf("The expected metric label %s is not found in the output: %s", label, bufStr)
		}
	}
}

func TestClient(t *testing.T) {
	tracer := otel.Tracer("test_trace")
	ctx, span := tracer.Start(context.Background(), "TestClient")
	defer span.End()

	e := errors.New("got an error")
	nextError := func(context.Context, interface{}) (interface{}, error) {
		return nil, e
	}
	nextValid := func(context.Context, interface{}) (interface{}, error) {
		return "Hello valid", nil
	}

	// init client handler
	meter := otel.Meter("test_meter")
	requests, err := DefaultRequestsCounter(meter, DefaultServerRequestsCounterName)
	if err != nil {
		t.Errorf("[DefaultRequestsCounter] something went wrong: %v", err)
		return
	}
	seconds, err := DefaultSecondsHistogram(meter, DefaultServerSecondsHistogramName)
	if err != nil {
		t.Errorf("[DefaultSecondsHistogram] something went wrong: %v", err)
		return
	}
	clientHandler := Client(
		WithRequests(requests),
		WithSeconds(seconds),
	)

	_, err = clientHandler(nextError)(ctx, "test:")
	if err != e {
		t.Error("The given error mismatch the expected.")
	}

	for i := 0; i < 20; i++ {
		res, err := clientHandler(nextValid)(transport.NewClientContext(ctx, &http.Transport{}), "test:")
		if err != nil {
			t.Error("The server must not throw an error.")
		}
		if res != "Hello valid" {
			t.Error(`The server must return a "Hello valid" response.`)
		}
	}
	bufStr := exporter.String()

	for _, label := range []string{
		metricLabelKind,
		metricLabelOperation,
		metricLabelCode,
		metricLabelReason,
		"SpanID",
		"TraceID",
	} {
		if !strings.Contains(bufStr, fmt.Sprintf("\"%s\"", label)) {
			t.Errorf("The expected metric label %s is not found in the output: %s", label, bufStr)
		}
	}
}
