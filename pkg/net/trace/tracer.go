package trace

import (
	"io"
)

var (
	// global tracer
	_tracer Tracer = nooptracer{}
)

// SetGlobalTracer SetGlobalTracer
func SetGlobalTracer(tracer Tracer) {
	_tracer = tracer
}

// Tracer is a simple, thin interface for Trace creation and propagation.
type Tracer interface {
	// New trace instance with given title.
	New(operationName string, opts ...Option) Trace
	// Inject takes the Trace instance and injects it for
	// propagation within `carrier`. The actual type of `carrier` depends on
	// the value of `format`.
	Inject(t Trace, format interface{}, carrier interface{}) error
	// Extract returns a Trace instance given `format` and `carrier`.
	// return `ErrTraceNotFound` if trace not found.
	Extract(format interface{}, carrier interface{}) (Trace, error)
}

// New trace instance with given operationName.
func New(operationName string, opts ...Option) Trace {
	return _tracer.New(operationName, opts...)
}

// Inject takes the Trace instance and injects it for
// propagation within `carrier`. The actual type of `carrier` depends on
// the value of `format`.
func Inject(t Trace, format interface{}, carrier interface{}) error {
	return _tracer.Inject(t, format, carrier)
}

// Extract returns a Trace instance given `format` and `carrier`.
// return `ErrTraceNotFound` if trace not found.
func Extract(format interface{}, carrier interface{}) (Trace, error) {
	return _tracer.Extract(format, carrier)
}

// Close trace flush data.
func Close() error {
	if closer, ok := _tracer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Trace trace common interface.
type Trace interface {
	// Fork fork a trace with client trace.
	Fork(serviceName, operationName string) Trace

	// Follow
	Follow(serviceName, operationName string) Trace

	// Finish when trace finish call it.
	Finish(err *error)

	// Scan scan trace into info.
	// Deprecated: method Scan is deprecated, use Inject instead of Scan
	// Scan(ti *Info)

	// Adds a tag to the trace.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	//
	// Tag values can be numeric types, strings, or bools. The behavior of
	// other tag value types is undefined at the OpenTracing level. If a
	// tracing system does not know how to handle a particular value type, it
	// may ignore the tag, but shall not panic.
	// NOTE current only support legacy tag: TagAnnotation TagAddress TagComment
	// other will be ignore
	SetTag(tags ...Tag) Trace

	// LogFields is an efficient and type-checked way to record key:value
	// NOTE current unsupport
	SetLog(logs ...LogField) Trace

	// Visit visits the k-v pair in trace, calling fn for each.
	Visit(fn func(k, v string))

	// SetTitle reset trace title
	SetTitle(title string)
}
