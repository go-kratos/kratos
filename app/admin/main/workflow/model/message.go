package model

import (
	"fmt"

	"go-common/app/admin/main/workflow/model/param"
)

// message mc
const (
	ArcComplainDealMC = "1_13_1"
	ArcComplainRevMC  = "1_15_1"
	WkfNotifyMC       = "1_5_3"
)

// DealArcComplainMsg generate archive complain deal message param
func DealArcComplainMsg(aid int64, mids []int64) *param.MessageParam {
	return &param.MessageParam{
		Type:     "json",
		Source:   1,
		DataType: 4,
		MC:       ArcComplainDealMC,
		Title:    "您的投诉已被受理",
		Context:  fmt.Sprintf("您对稿件（av%d）的投诉已被受理。感谢您对 bilibili 社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~ ", aid),
		MidList:  mids,
	}
}

// ReceivedArcComplainMsg generate archive complain received message param
func ReceivedArcComplainMsg(aid int64, mids []int64) *param.MessageParam {
	return &param.MessageParam{
		Type:     "json",
		Source:   1,
		DataType: 4,
		MC:       ArcComplainRevMC,
		Title:    "您的投诉已收到",
		Context:  fmt.Sprintf("您对稿件（av%d）的投诉我们已经收到。感谢您对 bilibili 社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~", aid),
		MidList:  mids,
	}
}
