package model

import (
	"strconv"
)

// SpanListRef .
type SpanListRef struct {
	TraceID  uint64
	SpanID   uint64
	IsError  bool
	Duration int64
}

// TraceIDStr hex format traceid
func (s SpanListRef) TraceIDStr() string {
	return strconv.FormatUint(s.TraceID, 16)
}

// SpanIDStr hex format traceid
func (s SpanListRef) SpanIDStr() string {
	return strconv.FormatUint(s.SpanID, 16)
}
