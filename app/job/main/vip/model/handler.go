package model

import "go-common/library/time"

//HandlerVip vip handler
type HandlerVip struct {
	Days   int32
	Months int16
	Mid    int64
	Type   int8
}

//AppCache appCache
type AppCache struct {
	AppID int64 `json:"appID"`
	Mid   int64 `json:"mid"`
}

//BcoinSendInfo .
type BcoinSendInfo struct {
	Amount     int32     `json:"amount"`
	DueDate    time.Time `json:"dueDate"`
	DayOfMonth int       `json:"dayOfMonth"`
}

//VipBuyResq .
type VipBuyResq struct {
	Mid      int64  `json:"mid"`
	CouponID string `json:"coupon_id"`
	OrderNo  string `json:"order_no"`
}

//VipPushResq .
type VipPushResq struct {
	Code int64 `json:"code"`
	Data int64 `json:"data"`
}
