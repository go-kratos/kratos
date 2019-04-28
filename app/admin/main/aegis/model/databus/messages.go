package databus

// CreateTaskMsg .
type CreateTaskMsg struct {
	BizID         int64 `json:"business_id"`
	RID           int64 `json:"rid"`
	FlowID        int64 `json:"flow_id"`
	DispatchLimit int64 `json:"dispatch_limit"`
}
