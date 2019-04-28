package model

// all variable used in advance dm
const (
	// mode
	AdvSpeMode = "sp"      // mode 7
	AdvMode    = "advance" // mode8 mode9
	// type
	AdvTypeRequest = "request"
	AdvTypeAccept  = "accept"
	AdvTypeBuy     = "buy"
	AdvTypeDeny    = "deny"
	// coin
	AdvSPCoin = 2
	AdvCoin   = 5
	// reason
	AdvSPCoinReason       = "购买特殊弹幕"
	AdvCoinReason         = "购买高级弹幕"
	AdvSPCoinCancelReason = "购买特殊弹幕被取消"
	AdvCoinCancelReason   = "购买高级弹幕被取消"
	// confirm state
	AdvStatConfirmDefault = 0
	AdvStatConfirmAgree   = 1
	AdvStatConfirmRequest = 2
	AdvStatConfirmDeny    = 3
	// 高级弹幕申请权限控制
	AdvPermitAll       = int8(0) // 任何人
	AdvPermitFollower  = int8(1) // 仅限粉丝
	AdvPermitAttention = int8(2) // 仅限相互关注
	AdvPermitForbid    = int8(3) // 始终拒绝
)

// BuyAdv user buy adv
type BuyAdv struct {
	CID       int64
	Owner     int64
	Mid       int64
	Type      string
	Timestamp int64
	Mode      string
	Refund    int
}

// ArgAdvBuy buy adv data
type ArgAdvBuy struct {
	Mid       int64
	Owner     int64
	Type      string
	Reason    string
	Cid       int64
	Coin      float64
	Mode      string
	Cookie    string
	AccessKey string
	Refund    int
	IsCoin    bool
}

// AdvState state
type AdvState struct {
	Coins   int  `json:"coins"`
	Confirm int  `json:"confirm"`
	Accept  bool `json:"accept"`
	HasBuy  bool `json:"hasBuy"`
}

// Advance dm_advancecomment
type Advance struct {
	ID        int64  `json:"id"`
	Owner     int64  `json:"owner"`
	Cid       int64  `json:"cid"`
	Aid       int64  `json:"aid"`
	Type      string `json:"type"`
	Mode      string `json:"mode"`
	Mid       int64  `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	Refund    int8   `json:"refund"`
	Uname     string `json:"uname"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
}
