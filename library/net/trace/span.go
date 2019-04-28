package trace

import (
	"fmt"
	"strconv"
	"time"

	protogen "go-common/library/net/trace/proto"
)

const (
	_maxChilds = 1024
	_maxTags   = 128
	_maxLogs   = 256
)

var _ Trace = &span{}

type span struct {
	dapper        *dapper
	context       spanContext
	operationName string
	startTime     time.Time
	duration      time.Duration
	tags          []Tag
	logs          []*protogen.Log
	childs        int
}

func (s *span) Fork(serviceName, operationName string) Trace {
	if s.childs > _maxChilds {
		// if child span more than max childs set return noopspan
		return noopspan{}
	}
	s.childs++
	// 为了兼容临时为 New 的 Span 设置 span.kind
	return s.dapper.newSpanWithContext(operationName, s.context).SetTag(TagString(TagSpanKind, "client"))
}

func (s *span) Follow(serviceName, operationName string) Trace {
	return s.Fork(serviceName, operationName).SetTag(TagString(TagSpanKind, "producer"))
}

func (s *span) Finish(perr *error) {
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

func (s *span) SetTag(tags ...Tag) Trace {
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
func (s *span) SetLog(logs ...LogField) Trace {
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

func (s *span) setLog(logs ...LogField) Trace {
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
func (s *span) Visit(fn func(k, v string)) {
	// NOTE: Deprecated key: delete in future
	fn(KeyTraceID, strconv.FormatUint(s.context.traceID, 10))
	fn(KeyTraceSpanID, strconv.FormatUint(s.context.spanID, 10))
	fn(KeyTraceParentID, strconv.FormatUint(s.context.parentID, 10))
	fn(KeyTraceSampled, strconv.FormatBool(s.context.isSampled()))
	fn(KeyTraceCaller, s.dapper.serviceName)

	fn(BiliTraceID, s.context.String())
}

// SetTitle reset trace title
func (s *span) SetTitle(operationName string) {
	s.operationName = operationName
}

func (s *span) String() string {
	return s.context.String()
}
