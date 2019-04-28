package model

// all variable used in advance dm
const (
	// mode
	AdvSpeMode = "sp"      // mode 7
	AdvMode    = "advance" // mode8 mode9
	AdvModeAll = "all"
	// type
	AdvTypeRequest = "request"
	AdvTypeAccept  = "accept"
	AdvTypeDeny    = "deny"
	AdvTypeAll     = "all"
)

// Advance advance dm list
type Advance struct {
	ID        int64  `json:"id"`        //高级弹幕ID
	Type      string `json:"bType"`     //处理结果
	Mode      string `json:"mode"`      //"sp" or 'advance"
	Mid       int64  `json:"mid"`       //申请人ID
	Timestamp int64  `json:"timestamp"` //申请时间
	Name      string `json:"name"`      //申请人昵称
}

// AdvanceRes advance dm list result including page info
type AdvanceRes struct {
	Result []*Advance `json:"result"`
	Page   *PageInfo  `json:"page"`
}

// PageInfo page info
type PageInfo struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// ArgMids advance dm mids
type ArgMids struct {
	Mids []int64
}
