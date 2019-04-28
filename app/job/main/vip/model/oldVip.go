package model

import "go-common/library/time"

//VipUserInfoOld vip_user_info table
type VipUserInfoOld struct {
	ID                   int64     `json:"id"`
	Mid                  int64     `json:"mid"`
	Type                 int8      `json:"vipType"`
	Status               int8      `json:"vipStatus"`
	StartTime            time.Time `json:"vipStartTime"`
	OverdueTime          time.Time `json:"vipOverdueTime"`
	AnnualVipOverdueTime time.Time `json:"annualVipOverdueTime"`
	RecentTime           time.Time `json:"vip_recent_time"`
	Wander               int8      `json:"wander"`
	AccessStatus         int8      `json:"accessStatus"`
	AutoRenewed          int8      `json:"auto_renewed"`
	IsAutoRenew          int8      `json:"is_auto_renew"`
	IosOverdueTime       time.Time `json:"ios_overdue_time"`
	PayChannelID         int64     `json:"pay_channel_id"`
	Ver                  int64     `json:"ver"`
	Ctime                time.Time `json:"ctime"`
	Mtime                time.Time `json:"mtime"`
}

// ToNew convert old model to new.
func (v *VipUserInfoOld) ToNew() (n *VipUserInfo) {
	return &VipUserInfo{
		Mid:                  v.Mid,
		Type:                 v.Type,
		PayType:              v.IsAutoRenew,
		PayChannelID:         v.PayChannelID,
		Status:               v.Status,
		StartTime:            v.StartTime,
		RecentTime:           v.RecentTime,
		OverdueTime:          v.OverdueTime,
		AnnualVipOverdueTime: v.AnnualVipOverdueTime,
		Ctime:                v.Ctime,
		Mtime:                v.Mtime,
		IosOverdueTime:       v.IosOverdueTime,
		Ver:                  v.Ver,
	}
}

//VipPayOrderOld vip pay order table
type VipPayOrderOld struct {
	ID          int64     `json:"id"`
	OrderNo     string    `json:"orderNo"`
	AppID       int64     `json:"appId"`
	Platform    int8      `json:"platform"`
	OrderType   int8      `json:"orderType"`
	Mid         int64     `json:"mid"`
	Bmid        int64     `json:"bmid"`
	BuyMonths   int16     `json:"buyMonths"`
	Money       float64   `json:"money"`
	Status      int8      `json:"status"`
	PayType     int8      `json:"payType"`
	PaymentTime time.Time `json:"paymentTime"`
	Ver         int64     `json:"ver"`
	AppSubID    string    `json:"appSubId"`
	CouponMoney float64   `json:"coupon_money"`
}

//VipRechargeOrder vip recharge order table
type VipRechargeOrder struct {
	ID           int     `json:"id"`
	AppID        int     `json:"appId"`
	PayMid       int     `json:"payMid"`
	OrderNo      string  `json:"orderNo"`
	RechargeBp   float64 `json:"rechargeBp"`
	ThirdTradeNo string  `json:"thirdTradeNo"`
	PayOrderNo   string  `json:"payOrderNo"`
	Status       int     `json:"status"`
	Ver          int     `json:"ver"`
	Bmid         int     `json:"bmid"`
}

//VipChangeHistory .
type VipChangeHistory struct {
	ID          int64     `json:"id"`
	Mid         int64     `json:"mid"`
	ChangeType  int8      `json:"change_type"`
	ChangeTime  time.Time `json:"change_time"`
	Month       int16     `json:"month"`
	Days        int32     `json:"days"`
	OperatorID  string    `json:"operator_id"`
	RelationID  string    `json:"relation_id"`
	BatchCodeID int64     `json:"batch_code_id"`
	BatchID     int64     `json:"batch_id"`
	Remark      string    `json:"remark"`
}
