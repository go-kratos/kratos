package model

// 临时逻辑  --- start ----

import "go-common/library/time"

// OrderInfo order info.
type OrderInfo struct {
	ID           int64
	OrderNo      string
	AppID        int64
	OrderType    int8
	Platform     int8
	Mid          int64
	ToMid        int64
	BuyMonths    int16
	Money        float64
	RefundAmount float64
	Status       int8
	PayType      string
	RechargeBP   float64
	ThirdTradeNo string
	PaymentTime  time.Time
	Ver          int64
	Ctime        time.Time
	Mtime        time.Time
	AppSubID     string
}

// 临时逻辑  --- end ----
