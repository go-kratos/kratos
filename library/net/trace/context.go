package trace

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	flagSampled = 0x01
	flagDebug   = 0x02
)

var (
	errEmptyTracerString   = errors.New("trace: cannot convert empty string to spancontext")
	errInvalidTracerString = errors.New("trace: string does not match spancontext string format")
)

// SpanContext implements opentracing.SpanContext
type spanContext struct {
	// traceID represents globally unique ID of the trace.
	// Usually generated as a random number.
	traceID uint64

	// spanID represents span ID that must be unique within its trace,
	// but does not have to be globally unique.
	spanID uint64

	// parentID refers to the ID of the parent span.
	// Should be 0 if the current span is a root span.
	parentID uint64

	// flags is a bitmap containing such bits as 'sampled' and 'debug'.
	flags byte

	// probability
	probability float32

	// current level
	level int
}

func (c spanContext) isSampled() bool {
	return (c.flags & flagSampled) == flagSampled
}

func (c spanContext) isDebug() bool {
	return (c.flags & flagDebug) == flagDebug
}

// IsValid check spanContext valid
func (c spanContext) IsValid() bool {
	return c.traceID != 0 && c.spanID != 0
}

// emptyContext emptyContext
var emptyContext = spanContext{}

// String convert spanContext to String
// {TraceID}:{SpanID}:{ParentID}:{flags}:[extend...]
// TraceID: uint64 base16
// SpanID: uint64 base16
// ParentID: uint64 base16
// flags:
// - :0 sampled flag
// - :1 debug flag
// extend:
// sample-rate: s-{base16(BigEndian(float32))}
func (c spanContext) String() string {
	base := make([]string, 4)
	base[0] = strconv.FormatUint(uint64(c.traceID), 16)
	base[1] = strconv.FormatUint(uint64(c.spanID), 16)
	base[2] = strconv.FormatUint(uint64(c.parentID), 16)
	base[3] = strconv.FormatUint(uint64(c.flags), 16)
	return strings.Join(base, ":")
}

// ContextFromString parse spanContext form string
func contextFromString(value string) (spanContext, error) {
	if value == "" {
		return emptyContext, errEmptyTracerString
	}
	items := strings.Split(value, ":")
	if len(items) < 4 {
		return emptyContext, errInvalidTracerString
	}
	parseHexUint64 := func(hexs []string) ([]uint64, error) {
		rets := make([]uint64, len(hexs))
		var err error
		for i, hex := range hexs {
			rets[i], err = strconv.ParseUint(hex, 16, 64)
			if err != nil {
				break
			}
		}
		return rets, err
	}
	rets, err := parseHexUint64(items[0:4])
	if err != nil {
		return emptyContext, errInvalidTracerString
	}
	sctx := spanContext{
		traceID:  rets[0],
		spanID:   rets[1],
		parentID: rets[2],
		flags:    byte(rets[3]),
	}
	return sctx, nil
}
