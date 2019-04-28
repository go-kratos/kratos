package model

// Pendant event const.
const (
	PendantPickOff int64 = iota + 1
	PendantPutOn
)

// VipInfo .
type VipInfo struct {
	Mid        int64 `json:"mid"`
	VipType    int64 `json:"vipType"`
	VipStatus  int64 `json:"vipStatus"`
	VipDueDate int64 `json:"vipDueDate"`
}

// ArgMIDNID struct.
type ArgMIDNID struct {
	MID int64 `form:"mid" validate:"gt=0,required"`
	NID int64 `form:"nid" validate:"gt=0,required"`
}

// ArgMID struct.
type ArgMID struct {
	MID int64 `form:"mid" validate:"gt=0,required"`
}

// AccountNotify .
type AccountNotify struct {
	UID    int64  `json:"mid"`
	Type   string `json:"type"`
	Action string `json:"action"`
}
