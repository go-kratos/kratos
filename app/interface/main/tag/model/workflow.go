package model

// const const value.
const (
	// WorkflowBusinessChannel channel business.
	WorkflowBusinessChannel = int32(9)
	WorkflowFIDChannel      = int32(2)
	WorkflowRIDChannel      = int32(1)

	AppealStatEffective = int32(1)
	AppealStateInvalid  = int32(2)
)

// WorkflowAppeal workflow appeal.
type WorkflowAppeal struct {
	Eid      int64 `json:"eid"` // channel id.
	Mid      int64 `json:"mid"`
	Oid      int64 `json:"oid"`
	Business int32 `json:"business"`
}

// WorkflowAppealInfo WorkflowAppealInfo.
type WorkflowAppealInfo struct {
	Business int32  `json:"business"`
	FID      int32  `json:"fid"`
	RID      int32  `json:"rid"`
	EID      int64  `json:"eid"`
	Score    int8   `json:"score"`
	ReasonID int8   `json:"tid"`
	Oid      int64  `json:"oid"`
	RptMid   int64  `json:"mid"`
	RegionID int32  `json:"business_typeid"`
	Mid      int64  `json:"business_mid"`
	TName    string `json:"business_title"`
	RealIP   string `json:"ip"`
}
