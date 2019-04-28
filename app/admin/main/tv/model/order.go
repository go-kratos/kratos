package model

import "go-common/library/time"

// TvPayOrder is table struct
type TvPayOrder struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"order_no"`
	Platform     int8      `json:"platform"`
	OrderType    int8      `json:"order_type"`
	ActiveType   int8      `json:"active_type"`
	MID          int64     `json:"mid" gorm:"column:mid"`
	BuyMonths    int8      `json:"buy_months"`
	ProductID    string    `json:"product_id"`
	Money        int64     `json:"money"`
	Quantity     int64     `json:"quantity"`
	RefundAmount int64     `json:"refund_amount"`
	Status       int8      `json:"status"`
	ThirdTradeNO string    `json:"third_trade_no"`
	PaymentMoney int64     `json:"payment_money"`
	PaymentType  string    `json:"payment_type"`
	PaymentTime  time.Time `json:"payment_time"`
	Ver          int64     `json:"ver"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// TvPayOrderResp is used to list in TV pay order list
type TvPayOrderResp struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"order_no"`
	OrderType    int8      `json:"order_type"`
	ActiveType   int8      `json:"active_type"`
	MID          int64     `json:"mid" form:"mid" gorm:"column:mid"`
	BuyMonths    int8      `json:"buy_months"`
	ProductID    string    `json:"product_id"`
	Money        int64     `json:"money"`
	Quantity     int64     `json:"quantity"`
	RefundAmount int64     `json:"refund_amount"`
	Status       int8      `json:"status"`
	ThirdTradeNO string    `json:"third_trade_no"`
	PaymentMoney int64     `json:"payment_money"`
	PaymentType  string    `json:"payment_type"`
	PaymentTime  time.Time `json:"payment_time"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// OrderPageHelper is used to list in TV pay order list count
type OrderPageHelper struct {
	Items []*TvPayOrderResp `json:"items"`
	Total *int64            `json:"total"`
}

// TableName tv_pay_order
func (*TvPayOrderResp) TableName() string {
	return "tv_pay_order"
}
