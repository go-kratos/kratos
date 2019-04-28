package model

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	protogen "go-common/library/net/trace/proto"
)

// ProtoSpan alias to tgo-common/library/net/trace/proto.Span
type ProtoSpan protogen.Span

// RefType Kind
const (
	RefTypeChildOf int8 = iota
	RefTypeFollowsFrom
)

// TagKind
const (
	TagString int8 = iota
	TagInt
	TagBool
	TagFloat
)

// SpanRef describes causal relationship of the current span to another span (e.g. 'child-of')
type SpanRef struct {
	RefType int8
	TraceID uint64
	SpanID  uint64
}

// Tag span tag
type Tag struct {
	Kind  int8
	Key   string
	Value interface{}
}

// Field log field
type Field struct {
	Key   string
	Value []byte
}

// Log span log
type Log struct {
	Timestamp int64
	Fields    []Field
}

// Span represents a named unit of work performed by a service.
type Span struct {
	ServiceName   string
	OperationName string
	TraceID       uint64
	SpanID        uint64
	ParentID      uint64
	Env           string
	StartTime     time.Time
	Duration      time.Duration
	References    []SpanRef
	Tags          map[string]interface{}
	Logs          []Log
	ProtoSpan     *ProtoSpan
}

// SetTag attach tag
func (s *Span) SetTag(key string, value interface{}) error {
	ptag, err := toProtoTag(key, value)
	if err != nil {
		return err
	}
	s.Tags[key] = value
	s.ProtoSpan.Tags = append(s.ProtoSpan.Tags, ptag)
	return nil
}

// SetOperationName .
func (s *Span) SetOperationName(operationName string) {
	s.OperationName = operationName
	s.ProtoSpan.OperationName = operationName
}

// TraceIDStr return hex format trace_id
func (s *Span) TraceIDStr() string {
	return strconv.FormatUint(s.TraceID, 16)
}

// SpanIDStr return hex format span_id
func (s *Span) SpanIDStr() string {
	return strconv.FormatUint(s.SpanID, 16)
}

// ParentIDStr return hex format parent_id
func (s *Span) ParentIDStr() string {
	return strconv.FormatUint(s.ParentID, 16)
}

// IsServer span kind is server
func (s *Span) IsServer() bool {
	kind, ok := s.Tags["span.kind"].(string)
	if !ok {
		return false
	}
	return kind == "server"
}

// IsError is error happend
func (s *Span) IsError() bool {
	isErr, _ := s.Tags["error"].(bool)
	return isErr
}

// StringTag get string type tag
func (s *Span) StringTag(key string) string {
	val, _ := s.Tags[key].(string)
	return val
}

// BoolTag get string type tag
func (s *Span) BoolTag(key string) bool {
	val, _ := s.Tags[key].(bool)
	return val
}

// GetTagString .
func (s *Span) GetTagString(key string) string {
	val, _ := s.Tags[key].(string)
	return val
}

// Marshal return
func (s *Span) Marshal() ([]byte, error) {
	if s.ProtoSpan == nil {
		return nil, nil
	}
	return proto.Marshal((*protogen.Span)(s.ProtoSpan))
}
