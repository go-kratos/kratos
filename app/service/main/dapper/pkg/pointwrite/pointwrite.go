package pointwrite

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go-common/app/service/main/dapper/model"
	"go-common/library/log"
)

// WriteFn .
type WriteFn func(ctx context.Context, points []*model.SpanPoint) error

// PointWriter writer span point
type PointWriter interface {
	WriteSpan(span *model.Span) error
	Close() error
}

// New PointWriter
func New(fn WriteFn, precision int64, timeout time.Duration) PointWriter {
	pw := &pointwriter{
		precision: precision,
		current:   make(map[string]*model.SpanPoint),
		timeout:   timeout,
		// TODO make it configurable
		tk: time.NewTicker(time.Second * 30),
		fn: fn,
	}
	go pw.start()
	return pw
}

type pointwriter struct {
	closed    bool
	rmx       sync.RWMutex
	precision int64
	timeout   time.Duration
	current   map[string]*model.SpanPoint
	fn        WriteFn
	tk        *time.Ticker
}

func (p *pointwriter) start() {
	for range p.tk.C {
		err := p.flush()
		if err != nil {
			log.Error("flush pointwriter error: %s", err)
		}
	}
}

func (p *pointwriter) flush() error {
	p.rmx.Lock()
	current := p.current
	p.current = make(map[string]*model.SpanPoint)
	p.rmx.Unlock()
	points := make([]*model.SpanPoint, 0, len(current))
	for _, point := range current {
		points = append(points, point)
	}
	if len(points) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	return p.fn(ctx, points)
}

// WriteSpan writespan
func (p *pointwriter) WriteSpan(span *model.Span) error {
	if p.closed {
		return fmt.Errorf("pointwriter already closed")
	}
	kind := "client"
	if span.IsServer() {
		kind = "server"
	}
	// NOTE: ingored sample ponit if is legacy span, DELETE it futrue
	if kind == "client" && !strings.Contains(span.ServiceName, ".") {
		return nil
	}
	peerService, ok := span.Tags["peer.service"].(string)
	if !ok {
		peerService = "unknown"
	}
	timestamp := span.StartTime.Unix() - (span.StartTime.Unix() % p.precision)
	key := fmt.Sprintf("%d_%s_%s_%s_%s",
		timestamp,
		span.ServiceName,
		span.OperationName,
		peerService,
		kind,
	)
	p.rmx.Lock()
	defer p.rmx.Unlock()
	point, ok := p.current[key]
	if !ok {
		point = &model.SpanPoint{
			Timestamp:     timestamp,
			ServiceName:   span.ServiceName,
			OperationName: span.OperationName,
			PeerService:   peerService,
			SpanKind:      kind,
			AvgDuration:   model.SamplePoint{TraceID: span.TraceID, SpanID: span.SpanID, Value: int64(span.Duration)},
		}
		p.current[key] = point
	}

	duration := int64(span.Duration)
	if duration > point.MaxDuration.Value {
		point.MaxDuration.TraceID = span.TraceID
		point.MaxDuration.SpanID = span.SpanID
		point.MaxDuration.Value = duration
	}
	if point.MinDuration.Value == 0 || duration < point.MinDuration.Value {
		point.MinDuration.TraceID = span.TraceID
		point.MinDuration.SpanID = span.SpanID
		point.MinDuration.Value = duration
	}
	if span.IsError() {
		point.Errors = append(point.Errors, model.SamplePoint{
			TraceID: span.TraceID,
			SpanID:  span.SpanID,
			Value:   duration,
		})
	}
	return nil
}

// Close pointwriter
func (p *pointwriter) Close() error {
	p.closed = true
	p.tk.Stop()
	return p.flush()
}
