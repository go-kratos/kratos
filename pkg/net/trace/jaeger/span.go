package jaeger

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// Span implements opentracing.Span
type Span struct {
	// referenceCounter used to increase the lifetime of
	// the object before return it into the pool.
	referenceCounter int32

	serviceName string

	sync.RWMutex

	// TODO: (breaking change) change to use a pointer
	context SpanContext

	// The name of the "operation" this span is an instance of.
	// Known as a "span name" in some implementations.
	operationName string

	// firstInProcess, if true, indicates that this span is the root of the (sub)tree
	// of spans in the current process. In other words it's true for the root spans,
	// and the ingress spans when the process joins another trace.
	firstInProcess bool

	// startTime is the timestamp indicating when the span began, with microseconds precision.
	startTime time.Time

	// duration returns duration of the span with microseconds precision.
	// Zero value means duration is unknown.
	duration time.Duration

	// tags attached to this span
	tags []Tag

	// The span's "micro-log"
	logs []opentracing.LogRecord

	// The number of logs dropped because of MaxLogsPerSpan.
	numDroppedLogs int

	// references for this span
	references []Reference
}

// Tag is a simple key value wrapper.
// TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.
type Tag struct {
	key   string
	value interface{}
}

// NewTag creates a new Tag.
// TODO (breaking change) deprecate in the next major release, use opentracing.Tag instead.
func NewTag(key string, value interface{}) Tag {
	return Tag{key: key, value: value}
}

// SetOperationName sets or changes the operation name.
func (s *Span) SetOperationName(operationName string) opentracing.Span {
	s.Lock()
	s.operationName = operationName
	s.Unlock()
	return s
}

// SetTag implements SetTag() of opentracing.Span
func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	return s.setTagInternal(key, value, true)
}

func (s *Span) setTagInternal(key string, value interface{}, lock bool) opentracing.Span {
	if lock {
		s.Lock()
		defer s.Unlock()
	}
	s.appendTagNoLocking(key, value)
	return s
}

// SpanContext returns span context
func (s *Span) SpanContext() SpanContext {
	s.Lock()
	defer s.Unlock()
	return s.context
}

// StartTime returns span start time
func (s *Span) StartTime() time.Time {
	s.Lock()
	defer s.Unlock()
	return s.startTime
}

// Duration returns span duration
func (s *Span) Duration() time.Duration {
	s.Lock()
	defer s.Unlock()
	return s.duration
}

// Tags returns tags for span
func (s *Span) Tags() opentracing.Tags {
	s.Lock()
	defer s.Unlock()
	var result = make(opentracing.Tags, len(s.tags))
	for _, tag := range s.tags {
		result[tag.key] = tag.value
	}
	return result
}

// Logs returns micro logs for span
func (s *Span) Logs() []opentracing.LogRecord {
	s.Lock()
	defer s.Unlock()

	logs := append([]opentracing.LogRecord(nil), s.logs...)
	if s.numDroppedLogs != 0 {
		fixLogs(logs, s.numDroppedLogs)
	}

	return logs
}

// References returns references for this span
func (s *Span) References() []opentracing.SpanReference {
	s.Lock()
	defer s.Unlock()

	if s.references == nil || len(s.references) == 0 {
		return nil
	}

	result := make([]opentracing.SpanReference, len(s.references))
	for i, r := range s.references {
		result[i] = opentracing.SpanReference{Type: r.Type, ReferencedContext: r.Context}
	}
	return result
}

func (s *Span) appendTagNoLocking(key string, value interface{}) {
	s.tags = append(s.tags, Tag{key: key, value: value})
}

// LogFields implements opentracing.Span API
func (s *Span) LogFields(fields ...log.Field) {
	s.Lock()
	defer s.Unlock()
	if !s.context.IsSampled() {
		return
	}
	s.logFieldsNoLocking(fields...)
}

// this function should only be called while holding a Write lock
func (s *Span) logFieldsNoLocking(fields ...log.Field) {
	lr := opentracing.LogRecord{
		Fields:    fields,
		Timestamp: time.Now(),
	}
	s.appendLogNoLocking(lr)
}

// LogKV implements opentracing.Span API
func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	s.RLock()
	sampled := s.context.IsSampled()
	s.RUnlock()
	if !sampled {
		return
	}
	fields, err := log.InterleavedKVToFields(alternatingKeyValues...)
	if err != nil {
		s.LogFields(log.Error(err), log.String("function", "LogKV"))
		return
	}
	s.LogFields(fields...)
}

// LogEvent implements opentracing.Span API
func (s *Span) LogEvent(event string) {
	s.Log(opentracing.LogData{Event: event})
}

// LogEventWithPayload implements opentracing.Span API
func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{Event: event, Payload: payload})
}

// Log implements opentracing.Span API
func (s *Span) Log(ld opentracing.LogData) {
	s.Lock()
	defer s.Unlock()
	if s.context.IsSampled() {
		s.appendLogNoLocking(ld.ToLogRecord())
	}
}

// this function should only be called while holding a Write lock
func (s *Span) appendLogNoLocking(lr opentracing.LogRecord) {
	maxLogs := 100
	// We have too many logs. We don't touch the first numOld logs; we treat the
	// rest as a circular buffer and overwrite the oldest log among those.
	numOld := (maxLogs - 1) / 2
	numNew := maxLogs - numOld
	s.logs[numOld+s.numDroppedLogs%numNew] = lr
	s.numDroppedLogs++
}

// rotateLogBuffer rotates the records in the buffer: records 0 to pos-1 move at
// the end (i.e. pos circular left shifts).
func rotateLogBuffer(buf []opentracing.LogRecord, pos int) {
	// This algorithm is described in:
	//    http://www.cplusplus.com/reference/algorithm/rotate
	for first, middle, next := 0, pos, pos; first != middle; {
		buf[first], buf[next] = buf[next], buf[first]
		first++
		next++
		if next == len(buf) {
			next = middle
		} else if first == middle {
			middle = next
		}
	}
}

func fixLogs(logs []opentracing.LogRecord, numDroppedLogs int) {
	// We dropped some log events, which means that we used part of Logs as a
	// circular buffer (see appendLog). De-circularize it.
	numOld := (len(logs) - 1) / 2
	numNew := len(logs) - numOld
	rotateLogBuffer(logs[numOld:], numDroppedLogs%numNew)

	// Replace the log in the middle (the oldest "new" log) with information
	// about the dropped logs. This means that we are effectively dropping one
	// more "new" log.
	numDropped := numDroppedLogs + 1
	logs[numOld] = opentracing.LogRecord{
		// Keep the timestamp of the last dropped event.
		Timestamp: logs[numOld].Timestamp,
		Fields: []log.Field{
			log.String("event", "dropped Span logs"),
			log.Int("dropped_log_count", numDropped),
			log.String("component", "jaeger-client"),
		},
	}
}

func (s *Span) fixLogsIfDropped() {
	if s.numDroppedLogs == 0 {
		return
	}
	fixLogs(s.logs, s.numDroppedLogs)
	s.numDroppedLogs = 0
}

// SetBaggageItem implements SetBaggageItem() of opentracing.SpanContext
func (s *Span) SetBaggageItem(key, value string) opentracing.Span {
	s.context.baggage[key] = value
	return s
}

// BaggageItem implements BaggageItem() of opentracing.SpanContext
func (s *Span) BaggageItem(key string) string {
	s.RLock()
	defer s.RUnlock()
	return s.context.baggage[key]
}

// Finish implements opentracing.Span API
// After finishing the Span object it returns back to the allocator unless the reporter retains it again,
// so after that, the Span object should no longer be used because it won't be valid anymore.
func (s *Span) Finish() {
	s.FinishWithOptions(opentracing.FinishOptions{})
}

// FinishWithOptions implements opentracing.Span API
func (s *Span) FinishWithOptions(options opentracing.FinishOptions) {
}

// Context implements opentracing.Span API
func (s *Span) Context() opentracing.SpanContext {
	s.Lock()
	defer s.Unlock()
	return s.context
}

// Tracer implements opentracing.Span API
func (s *Span) Tracer() opentracing.Tracer {
	return nil
}

func (s *Span) String() string {
	s.RLock()
	defer s.RUnlock()
	return s.context.String()
}

// OperationName allows retrieving current operation name.
func (s *Span) OperationName() string {
	s.RLock()
	defer s.RUnlock()
	return s.operationName
}

// Retain increases object counter to increase the lifetime of the object
func (s *Span) Retain() *Span {
	atomic.AddInt32(&s.referenceCounter, 1)
	return s
}

// Release decrements object counter and return to the
// allocator manager  when counter will below zero
func (s *Span) Release() {

}

// reset span state and release unused data
func (s *Span) reset() {
	s.firstInProcess = false
	s.context = emptyContext
	s.operationName = ""
	s.startTime = time.Time{}
	s.duration = 0
	atomic.StoreInt32(&s.referenceCounter, 0)

	// Note: To reuse memory we can save the pointers on the heap
	s.tags = s.tags[:0]
	s.logs = s.logs[:0]
	s.numDroppedLogs = 0
	s.references = s.references[:0]
}

func (s *Span) ServiceName() string {
	return s.serviceName
}
