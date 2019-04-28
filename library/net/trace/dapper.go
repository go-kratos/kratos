package trace

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const _maxLevel = 64

func newTracer(serviceName string, report reporter, cfg *Config) Tracer {
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
		pool:     &sync.Pool{New: func() interface{} { return new(span) }},
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
	pctx := spanContext{traceID: traceID}
	if sampled {
		pctx.flags = flagSampled
		pctx.probability = probability
	}
	if opt.Debug {
		pctx.flags |= flagDebug
		return d.newSpanWithContext(operationName, pctx).SetTag(TagString(TagSpanKind, "server")).SetTag(TagBool("debug", true))
	}
	// 为了兼容临时为 New 的 Span 设置 span.kind
	return d.newSpanWithContext(operationName, pctx).SetTag(TagString(TagSpanKind, "server"))
}

func (d *dapper) newSpanWithContext(operationName string, pctx spanContext) Trace {
	sp := d.getSpan()
	// is span is not sampled just return a span with this context, no need clear it
	//if !pctx.isSampled() {
	//	sp.context = pctx
	//	return sp
	//}
	if pctx.level > _maxLevel {
		// if span reach max limit level return noopspan
		return noopspan{}
	}
	level := pctx.level + 1
	nctx := spanContext{
		traceID:  pctx.traceID,
		parentID: pctx.spanID,
		flags:    pctx.flags,
		level:    level,
	}
	if pctx.spanID == 0 {
		nctx.spanID = pctx.traceID
	} else {
		nctx.spanID = genID()
	}
	sp.operationName = operationName
	sp.context = nctx
	sp.startTime = time.Now()
	sp.tags = append(sp.tags, d.tags...)
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
	// 为了兼容临时为 New 的 Span 设置 span.kind
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
	contextStr := carr.Get(BiliTraceID)
	if contextStr == "" {
		return d.legacyExtract(carr)
	}
	pctx, err := contextFromString(contextStr)
	if err != nil {
		return nil, err
	}
	// NOTE: call SetTitle after extract trace
	return d.newSpanWithContext("", pctx), nil
}

func (d *dapper) legacyExtract(carr Carrier) (Trace, error) {
	traceIDstr := carr.Get(KeyTraceID)
	if traceIDstr == "" {
		return nil, ErrTraceNotFound
	}
	traceID, err := strconv.ParseUint(traceIDstr, 10, 64)
	if err != nil {
		return nil, ErrTraceCorrupted
	}
	sampled, _ := strconv.ParseBool(carr.Get(KeyTraceSampled))
	spanID, _ := strconv.ParseUint(carr.Get(KeyTraceSpanID), 10, 64)
	parentID, _ := strconv.ParseUint(carr.Get(KeyTraceSpanID), 10, 64)
	pctx := spanContext{traceID: traceID, spanID: spanID, parentID: parentID}
	if sampled {
		pctx.flags = flagSampled
	}
	return d.newSpanWithContext("", pctx), nil
}

func (d *dapper) Close() error {
	return d.reporter.Close()
}

func (d *dapper) report(sp *span) {
	if sp.context.isSampled() {
		if err := d.reporter.WriteSpan(sp); err != nil {
			d.stdlog.Printf("marshal trace span error: %s", err)
		}
	}
	d.putSpan(sp)
}

func (d *dapper) putSpan(sp *span) {
	if len(sp.tags) > 32 {
		sp.tags = nil
	}
	if len(sp.logs) > 32 {
		sp.logs = nil
	}
	d.pool.Put(sp)
}

func (d *dapper) getSpan() *span {
	sp := d.pool.Get().(*span)
	sp.dapper = d
	sp.childs = 0
	sp.tags = sp.tags[:0]
	sp.logs = sp.logs[:0]
	return sp
}
