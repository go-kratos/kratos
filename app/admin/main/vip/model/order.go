package model

import "go-common/library/time"

//PayOrder pay order
type PayOrder struct {
	ID           int64     `json:"_"`
	OrderNo      string    `json:"order_no"`
	AppID        int64     `json:"app_id"`
	Platform     int8      `json:"platform"`
	OrderType    int8      `json:"order_type"`
	AppSubID     string    `json:"app_sub_id"`
	Mid          int64     `json:"mid"`
	ToMid        int64     `json:"to_mid"`
	BuyMonths    int16     `json:"buy_months"`
	Money        float64   `json:"money"`
	RefundAmount float64   `json:"refund_amount"`
	Status       int8      `json:"status"`
	PayType      int8      `json:"pay_type"`
	RechargeBp   float64   `json:"recharge_bp"`
	ThirdTradeNo string    `json:"third_trade_no"`
	Ver          int64     `json:"ver"`
	PaymentTime  time.Time `json:"payment_time"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

//PayOrderLog pay order log.
type PayOrderLog struct {
	ID           int64   `json:"id"`
	OrderNo      string  `json:"order_no"`
	RefundID     string  `json:"refund_id"`
	RefundAmount float64 `json:"refund_amount"`
	Mid          int64   `json:"mid"`
	Status       int8    `json:"status"`
	Operator     string  `json:"operator"`
}
