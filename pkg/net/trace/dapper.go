package trace

import (
	"log"
	"os"
	"sync"
	"time"
)

const _maxLevel = 64

func NewTracer(serviceName string, report reporter, cfg *Config) Tracer {
	// hard code reset probability at 0.00025, 1/4000
	cfg.Probability = 0.00025
	sampler := newSampler(cfg.Probability)

	// default internal tags
	tags := extendTag()
	stdlog := log.New(os.Stderr, "trace", log.LstdFlags)
	return &dapper{
		cfg:         cfg,
		serviceName: serviceName,
		propagators: map[interface{}]propagator{
			HTTPFormat: httpPropagator{},
			GRPCFormat: grpcPropagator{},
		},
		reporter: report,
		sampler:  sampler,
		tags:     tags,
		pool:     &sync.Pool{New: func() interface{} { return new(Span) }},
		stdlog:   stdlog,
	}
}

type dapper struct {
	cfg         *Config
	serviceName string
	tags        []Tag
	reporter    reporter
	propagators map[interface{}]propagator
	pool        *sync.Pool
	stdlog      *log.Logger
	sampler     sampler
}

func (d *dapper) New(operationName string, opts ...Option) Trace {
	opt := defaultOption
	for _, fn := range opts {
		fn(&opt)
	}
	traceID := genID()
	var sampled bool
	var probability float32
	if d.cfg.DisableSample {
		sampled = true
		probability = 1
	} else {
		sampled, probability = d.sampler.IsSampled(traceID, operationName)
	}
	pctx := spanContext{TraceID: traceID}
	if sampled {
		pctx.Flags = flagSampled
		pctx.Probability = probability
	}
	if opt.Debug {
		pctx.Flags |= flagDebug
		return d.newSpanWithContext(operationName, pctx).SetTag(TagString(TagSpanKind, "server")).SetTag(TagBool("debug", true))
	}
	// 为了兼容临时为 New 的 Span 设置 Span.kind
	return d.newSpanWithContext(operationName, pctx).SetTag(TagString(TagSpanKind, "server"))
}

func (d *dapper) newSpanWithContext(operationName string, pctx spanContext) Trace {
	sp := d.getSpan()
	// is Span is not sampled just return a Span with this context, no need clear it
	//if !pctx.isSampled() {
	//	sp.context = pctx
	//	return sp
	//}
	if pctx.Level > _maxLevel {
		// if Span reach max limit level return noopspan
		return noopspan{}
	}
	level := pctx.Level + 1
	nctx := spanContext{
		TraceID:  pctx.TraceID,
		ParentID: pctx.SpanID,
		Flags:    pctx.Flags,
		Level:    level,
	}
	if pctx.SpanID == 0 {
		nctx.SpanID = pctx.TraceID
	} else {
		nctx.SpanID = genID()
	}
	sp.operationName = operationName
	sp.context = nctx
	sp.startTime = time.Now()
	return sp
}

func (d *dapper) Inject(t Trace, format interface{}, carrier interface{}) error {
	// if carrier implement Carrier use direct, ignore format
	carr, ok := carrier.(Carrier)
	if ok {
		t.Visit(carr.Set)
		return nil
	}
	// use Built-in propagators
	pp, ok := d.propagators[format]
	if !ok {
		return ErrUnsupportedFormat
	}
	carr, err := pp.Inject(carrier)
	if err != nil {
		return err
	}
	if t != nil {
		t.Visit(carr.Set)
	}
	return nil
}

func (d *dapper) Extract(format interface{}, carrier interface{}) (Trace, error) {
	sp, err := d.extract(format, carrier)
	if err != nil {
		return sp, err
	}
	// 为了兼容临时为 New 的 Span 设置 Span.kind
	return sp.SetTag(TagString(TagSpanKind, "server")), nil
}

func (d *dapper) extract(format interface{}, carrier interface{}) (Trace, error) {
	// if carrier implement Carrier use direct, ignore format
	carr, ok := carrier.(Carrier)
	if !ok {
		// use Built-in propagators
		pp, ok := d.propagators[format]
		if !ok {
			return nil, ErrUnsupportedFormat
		}
		var err error
		if carr, err = pp.Extract(carrier); err != nil {
			return nil, err
		}
	}
	contextStr := carr.Get(KratosTraceID)
	pctx, err := contextFromString(contextStr)
	if err != nil {
		return nil, err
	}
	// NOTE: call SetTitle after extract trace
	return d.newSpanWithContext("", pctx), nil
}

func (d *dapper) Close() error {
	return d.reporter.Close()
}

func (d *dapper) report(sp *Span) {
	if sp.context.isSampled() {
		if err := d.reporter.WriteSpan(sp); err != nil {
			d.stdlog.Printf("marshal trace Span error: %s", err)
		}
	}
	d.putSpan(sp)
}

func (d *dapper) putSpan(sp *Span) {
	if len(sp.tags) > 32 {
		sp.tags = nil
	}
	if len(sp.logs) > 32 {
		sp.logs = nil
	}
	d.pool.Put(sp)
}

func (d *dapper) getSpan() *Span {
	sp := d.pool.Get().(*Span)
	sp.dapper = d
	sp.childs = 0
	sp.tags = sp.tags[:0]
	sp.logs = sp.logs[:0]
	return sp
}
