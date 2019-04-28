package model

import (
	"math"
	"time"
)

// MainVip represents bilibili vip info.
type MainVip struct {
	Mid        int64 `json:"mid"`
	VipType    int8  `json:"vip_type"` // 大会员类型 0.非大会员 1.月度大会员 2.年度会员
	PayType    int8  `json:"pay_type"`
	VipStatus  int8  `json:"vip_status"` //大会员状态: 0.过期 1.未过期 2.冻结 3.封禁
	VipDueDate int64 `json:"vip_due_date"`
}

// IsVip returns true if user is vip.
func (mv *MainVip) IsVip() bool {
	return mv.VipType != 0 && (mv.VipStatus == 1 || mv.VipStatus == 3)
}

// Months returns vip months.
func (mv *MainVip) Months() int32 {
	if !mv.IsVip() {
		return 0
	}
	nowInMs := time.Now().UnixNano() / int64(time.Millisecond)
	span := mv.VipDueDate - nowInMs
	if span <= 0 {
		return 0
	}
	return int32(math.Floor(float64(span) / 1000 / 60 / 60 / 24 / 31))
}
