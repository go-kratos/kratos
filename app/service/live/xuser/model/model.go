package model

import (
	"go-common/app/service/live/xuser/api/grpc"
	"time"
)

// VipBuy buy vip request struct
type VipBuy struct {
	Uid      int64
	OrderID  string
	GoodID   int
	GoodNum  int
	Platform grpc.Platform
	Source   string
}

// VipInfo vip info struct
type VipInfo struct {
	Vip      int    `json:"vip"`
	VipTime  string `json:"vip_time"`
	Svip     int    `json:"svip"`
	SvipTime string `json:"svip_time"`
}

// VipRecord ap_vip_record log and notify message
type VipRecord struct {
	Uid           int64  `json:"uid"`
	Opcode        string `json:"opcode"`
	BuyType       int    `json:"buy_type"`
	BuyNum        int    `json:"buy_num"`
	VipType       int    `json:"vip_type"`
	BeforeVipTime string `json:"begin"`
	AfterVipTime  string `json:"end"`
	Platform      string
}

// GuardBuy buy guard request struct
type GuardBuy struct {
	OrderId    string
	Uid        int64
	Ruid       int64
	GuardLevel int
	Num        int
	Platform   grpc.Platform
	Source     string
}

// GuardInfo guard info struct for ap_user_privilege
type GuardInfo struct {
	Id            int64
	Uid           int64
	TargetId      int64
	PrivilegeType int
	StartTime     time.Time
	ExpiredTime   time.Time
}

// GuardEntryEffects entry effect message
type GuardEntryEffects struct {
	Business int                `json:"business"`
	Data     []GuardEntryEffect `json:"data"`
}

// GuardEntryEffect entry effect message
type GuardEntryEffect struct {
	EffectId int    `json:"effect_id"`
	Uid      int64  `json:"uid"`
	TargetId int64  `json:"target_id"`
	EndTime  string `json:"end_time"`
}

// Vip constants
var (
	Vip  = 1 // 月费姥爷
	Svip = 2 // 年费姥爷

	BuyStatusSuccess = 1 // 购买成功
	BuyStatusRetry   = 2 // 需要重试

	TimeNano  = "2006-01-02 15:04:05"
	TimeEmpty = "0000-00-00 00:00:00"

	OpcodeAdd = "add"
)
