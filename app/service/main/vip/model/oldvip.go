package model

import "go-common/library/time"

//VipUserInfo vip_user_info table for vip java
type VipUserInfo struct {
	ID                   int64     `json:"id"`
	Mid                  int64     `form:"mid" validate:"required" json:"mid"`
	VipType              int32     `form:"vipType" json:"vipType"`
	VipStatus            int32     `form:"vipStatus" json:"vipStatus"`
	VipStartTime         time.Time `form:"vipStartTime" validate:"required" json:"vipStartTime"`
	VipRecentTime        time.Time `form:"vipRecentTime" json:"vipRecentTime"`
	VipOverdueTime       time.Time `form:"vipOverdueTime" validate:"required" json:"vipOverdueTime"`
	AnnualVipOverdueTime time.Time `form:"annualVipOverdueTime" json:"annualVipOverdueTime"`
	Wander               int8      `json:"wander"`
	AccessStatus         int8      `json:"accessStatus"`
	Ctime                time.Time `form:"ctime" validate:"required" json:"ctime"`
	Mtime                time.Time `form:"mtime" validate:"required" json:"mtime"`
	Ver                  int64     `form:"ver" json:"ver"`
	AutoRenewed          int8      `form:"autoRenewed" json:"autoRenewed"`
	IsAutoRenew          int32     `form:"isAutoRenew" json:"isAutoRenew"`
	PayChannelID         int32     `form:"payChannelId" json:"payChannelId"`
	IosOverdueTime       time.Time `form:"iosOverdueTime" json:"iosOverdueTime"`
}

// ToNew convert old model to new.
func (v *VipUserInfo) ToNew() (res *VipInfoDB) {
	return &VipInfoDB{
		Mid:                  v.Mid,
		VipType:              v.VipType,
		VipPayType:           v.IsAutoRenew,
		PayChannelID:         v.PayChannelID,
		VipStatus:            v.VipStatus,
		VipStartTime:         v.VipStartTime,
		VipRecentTime:        v.VipRecentTime,
		VipOverdueTime:       v.VipOverdueTime,
		AnnualVipOverdueTime: v.AnnualVipOverdueTime,
		Ctime:                v.Ctime,
		Mtime:                v.Ctime,
		IosOverdueTime:       v.IosOverdueTime,
		Ver:                  v.Ver,
	}
}
