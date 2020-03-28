package trace

import (
	"fmt"
	"time"

	protogen "github.com/go-kratos/kratos/pkg/net/trace/proto"
)

const (
	_maxChilds = 1024
	_maxTags   = 128
	_maxLogs   = 256
)

var _ Trace = &Span{}

// Span is a trace span.
type Span struct {
	dapper        *dapper
	context       spanContext
	operationName string
	startTime     time.Time
	duration      time.Duration
	tags          []Tag
	logs          []*protogen.Log
	childs        int
}

func (s *Span) ServiceName() string {
	return s.dapper.serviceName
}

func (s *Span) OperationName() string {
	return s.operationName
}

func (s *Span) StartTime() time.Time {
	return s.startTime
}

func (s *Span) Duration() time.Duration {
	return s.duration
}

func (s *Span) TraceID() string {
	return s.context.String()
}

func (s *Span) Context() spanContext {
	return s.context
}

func (s *Span) Tags() []Tag {
	return s.tags
}

func (s *Span) Logs() []*protogen.Log {
	return s.logs
}

func (s *Span) Fork(serviceName, operationName string) Trace {
	if s.childs > _maxChilds {
		// if child span more than max childs set return noopspan
		return noopspan{}
	}
	s.childs++
	// 为了兼容临时为 New 的 Span 设置 span.kind
	return s.dapper.newSpanWithContext(operationName, s.context).SetTag(TagString(TagSpanKind, "client"))
}

func (s *Span) Follow(serviceName, operationName string) Trace {
	return s.Fork(serviceName, operationName).SetTag(TagString(TagSpanKind, "producer"))
}

func (s *Span) Finish(perr *error) {
	s.duration = time.Since(s.startTime)
	if perr != nil && *perr != nil {
		err := *perr
		s.SetTag(TagBool(TagError, true))
		s.SetLog(Log(LogMessage, err.Error()))
		if err, ok := err.(stackTracer); ok {
			s.SetLog(Log(LogStack, fmt.Sprintf("%+v", err.StackTrace())))
		}
	}
	s.dapper.report(s)
}

func (s *Span) SetTag(tags ...Tag) Trace {
	if !s.context.isSampled() && !s.context.isDebug() {
		return s
	}
	if len(s.tags) < _maxTags {
		s.tags = append(s.tags, tags...)
	}
	if len(s.tags) == _maxTags {
		s.tags = append(s.tags, Tag{Key: "trace.error", Value: "too many tags"})
	}
	return s
}

// LogFields is an efficient and type-checked way to record key:value
// NOTE current unsupport
func (s *Span) SetLog(logs ...LogField) Trace {
	if !s.context.isSampled() && !s.context.isDebug() {
		return s
	}
	if len(s.logs) < _maxLogs {
		s.setLog(logs...)
	}
	if len(s.logs) == _maxLogs {
		s.setLog(LogField{Key: "trace.error", Value: "too many logs"})
	}
	return s
}

func (s *Span) setLog(logs ...LogField) Trace {
	protoLog := &protogen.Log{
		Timestamp: time.Now().UnixNano(),
		Fields:    make([]*protogen.Field, len(logs)),
	}
	for i := range logs {
		protoLog.Fields[i] = &protogen.Field{Key: logs[i].Key, Value: []byte(logs[i].Value)}
	}
	s.logs = append(s.logs, protoLog)
	return s
}

// Visit visits the k-v pair in trace, calling fn for each.
func (s *Span) Visit(fn func(k, v string)) {
	fn(KratosTraceID, s.context.String())
}

// SetTitle reset trace title
func (s *Span) SetTitle(operationName string) {
	s.operationName = operationName
}

func (s *Span) String() string {
	return s.context.String()
}
