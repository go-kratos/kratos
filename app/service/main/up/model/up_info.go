package model

// UpInfoActiveReq req
type UpInfoActiveReq struct {
	Mid int64 `form:"mid" validate:"required"`
}

// UpsInfoActiveReq req
type UpsInfoActiveReq struct {
	Mids []int64 `form:"mids,split" validate:"required,dive,gt=0"`
}

// UpInfoActiveReply Reply
type UpInfoActiveReply struct {
	ID        int64 `json:"id"`
	MID       int64 `json:"mid"`
	ActiveTid int16 `json:"active_tid"`
}
