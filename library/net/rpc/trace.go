package rpc

import (
	"strconv"

	"go-common/library/net/trace"
)

// TraceInfo propagate trace propagate gorpc call
type TraceInfo struct {
	ID       uint64
	SpanID   uint64
	ParentID uint64
	Level    int32
	Sampled  bool
	Caller   string
	Title    string
	Time     int64
}

// Set implement trace.Carrier
func (i *TraceInfo) Set(key string, val string) {
	switch key {
	case trace.KeyTraceID:
		i.ID, _ = strconv.ParseUint(val, 10, 64)
	case trace.KeyTraceSpanID:
		i.SpanID, _ = strconv.ParseUint(val, 10, 64)
	case trace.KeyTraceParentID:
		i.ParentID, _ = strconv.ParseUint(val, 10, 64)
	case trace.KeyTraceSampled:
		i.Sampled, _ = strconv.ParseBool(val)
	case trace.KeyTraceLevel:
		lv, _ := strconv.Atoi(val)
		i.Level = int32(lv)
	case trace.KeyTraceCaller:
		i.Caller = val
	}
}

// Get implement trace.Carrier
func (i *TraceInfo) Get(key string) string {
	switch key {
	case trace.KeyTraceID:
		return strconv.FormatUint(i.ID, 10)
	case trace.KeyTraceSpanID:
		return strconv.FormatUint(i.SpanID, 10)
	case trace.KeyTraceParentID:
		return strconv.FormatUint(i.ParentID, 10)
	case trace.KeyTraceSampled:
		return strconv.FormatBool(i.Sampled)
	case trace.KeyTraceLevel:
		return strconv.Itoa(int(i.Level))
	case trace.KeyTraceCaller:
		return i.Caller
	}
	return ""
}
