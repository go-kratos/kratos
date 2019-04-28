package trace

// Trace key
const (
	// 历史遗留的 key 不要轻易修改
	KeyTraceID       = "x1-bilispy-id"
	KeyTraceSpanID   = "x1-bilispy-spanid"
	KeyTraceParentID = "x1-bilispy-parentid"
	KeyTraceSampled  = "x1-bilispy-sampled"
	KeyTraceLevel    = "x1-bilispy-lv"
	KeyTraceCaller   = "x1-bilispy-user"
	// trace sdk should use bili_trace_id to get trace info after this code be merged
	BiliTraceID    = "bili-trace-id"
	BiliTraceDebug = "bili-trace-debug"
)
