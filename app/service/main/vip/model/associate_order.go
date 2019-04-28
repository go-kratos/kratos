package model

import "go-common/library/time"

// order grant state
const (
	AssociateGrantStateNone int8 = iota
	AssociateGrantStateHadGrant
)

// ArgBilibiliVipGrant bilibili vip grant args.
type ArgBilibiliVipGrant struct {
	OpenID     string
	OutOpenID  string
	OutOrderNO string
	Duration   int32
	AppID      int64
}

// VipOrderAssociateGrant vip order associate grant.
type VipOrderAssociateGrant struct {
	ID         int64
	AppID      int64
	Mid        int64
	Months     int32
	OutOpenID  string
	OutTradeNO string
	State      int8
	Ctime      time.Time
	Mtime      time.Time
}

// ArgEleVipGrant args ele vip grant.
type ArgEleVipGrant struct {
	OrderNO string `form:"order_no" validate:"required"`
}

// VipAssociateGrantCount associate grant count.
type VipAssociateGrantCount struct {
	ID           int64     `json:"id"`
	AppID        int64     `json:"app_id"`
	Mid          int64     `json:"mid"`
	Months       int32     `json:"months"`
	CurrentCount int64     `json:"current_count"`
	Ctime        time.Time `json:"-"`
	Mtime        time.Time `json:"-"`
}
