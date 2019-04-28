package model

import (
	"go-common/library/time"
)

// order status
const (
	PAYING int8 = iota + 1
	SUCCESS
	FAILED
	Sign
	UnSign
	REFUNDING
	REFUNDED
	REFUNDFAIL
	CANCEL = 20
)

// order type
const (
	General int8 = iota
	AutoRenew
)

// quick pay status
const (
	QuickPaySuccess int8 = iota + 1
	QuickPayFailed
)

// pay way
const (
	ALIPAY int8 = iota + 1
	WECHAT
	BCION
	BANK
	PAYPAL
	IOSPAY
	QPAY
)

// pay origin
const (
	MobileOrigin int8 = iota + 1
	WebOrigin
	H5
)

// PayType pay type
var PayType = map[int8]string{
	ALIPAY: "alipay",
	WECHAT: "wechat",
	BCION:  "bp",
	BANK:   "bank",
	PAYPAL: "paypal",
	IOSPAY: "iospay",
	QPAY:   "qpay",
}

// PayTypeName pay name.
var PayTypeName = map[string]int8{
	"alipay": ALIPAY,
	"wechat": WECHAT,
	"bp":     BCION,
	"bank":   BANK,
	"paypal": PAYPAL,
	"iospay": IOSPAY,
}

// AddPayOrderResp add pay order
type AddPayOrderResp struct {
	OrderNo    string `json:"order_no"`
	CashierURL string `json:"cashier_url"`
	QrcodeURL  string `json:"qrcode_url"`
}

//PayBankResp pay bank resp
type PayBankResp struct {
	BanckCode string `json:"bankCode"`
	Name      string `json:"name"`
	Res       string `json:"res"`
}

//PayAccountResp .
type PayAccountResp struct {
	Mid       int64   `json:"mid"`
	Brokerage float64 `json:"brokerage"`
	DefaultBp float64 `json:"default_bp"`
}

//APIPayOrderResp api pay resp.
type APIPayOrderResp struct {
	SDK             string `json:"sdk"`
	QrcodeURL       string `json:"qrcode_url"`
	CashierURL      string `json:"cashier_url"`
	RechargeOrderNo string `json:"recharge_order_no"`
}

// QucikPayResp quick pay token.
type QucikPayResp struct {
	Token string `json:"token"`
}

// PayRetResp pay ret response.
type PayRetResp struct {
	Status int8 `json:"status"`
}

// BuyVipResp buy vip resp.
type BuyVipResp struct {
	Qrcode          string `json:"qrcode"`
	CashierURL      string `json:"cashier_url"`
	OrderNo         string `json:"orderNo"`
	RechargeOrderNo string `json:"rechargeOrderNo"`
	PayPayOrderNo   string `json:"payPayOrderNo"`
	PaySign         string `json:"paySign"`
	Status          int8   `json:"status"`
	ProductID       string `json:"productId"`
}

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
	CouponMoney  float64   `json:"coupon_money"`
	PaymentTime  time.Time `json:"payment_time"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
	PID          int64     `json:"p_id"`
	UserIP       []byte    `json:"-"`
}

//PayOrderResp pay order resp.
type PayOrderResp struct {
	OrderNo   string    `json:"order_no"`
	BuyMonths int16     `json:"buy_months"`
	Money     float64   `json:"money"`
	Status    int8      `json:"status"`
	Remark    string    `json:"remark"`
	Ctime     time.Time `json:"ctime"`
}

// Month def.
type Month struct {
	ID        int64     `json:"_"`
	Month     int16     `json:"month"`
	MonthType int8      `json:"month_type"`
	Operator  string    `json:"operator"`
	Status    int8      `json:"status"`
	Deleted   int8      `json:"deleted"`
	Mtime     time.Time `json:"mtime"`
}

// PriceMapping vip month map.
type PriceMapping struct {
	ID                 int64     `json:"_"`
	MonthID            int64     `json:"month_id"`
	MonthType          int8      `json:"month_type"`
	Money              float64   `json:"money"`
	Selected           int8      `json:"selected"`
	FirstDiscountMoney float64   `json:"first_discount_money"`
	DiscountMoney      float64   `json:"discount_money"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Remark             string    `json:"remark"`
	Operator           string    `json:"operator"`
	Mtime              time.Time `json:"mtime"`
}

//PayPlatformOrder .
type PayPlatformOrder struct {
	CustomerID      string `json:"customerId"`
	DeviceType      int8   `json:"deviceType"`
	OrderID         string `json:"orderId"`
	OrderCreateTime string `json:"orderCreateTime"`
	OrderExpire     int64  `json:"orderExpire"`
	NotifyURL       string `json:"notifyUrl"`
	SignURL         string `json:"signUrl"`
	ShowTitle       string `json:"showTitle"`
	TraceID         string `json:"traceId"`
	Timestamp       string `json:"timestamp"`
	Version         string `json:"version"`
	SignType        string `json:"signType"`
	Sign            string `json:"sign"`
	ProductID       string `json:"productId"`
	PayAmount       int32  `json:"payAmount"`
	PlanID          int32  `json:"planId"`
	UID             int64  `json:"uid"`
	DisplayAccount  string `json:"displayAccount"`
	ServiceType     int32  `json:"serviceType"`
	OriginalAmout   int32  `json:"originalAmount"`
}

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

//OrderMng .
type OrderMng struct {
	Username         string `json:"username"`
	Mid              int64  `json:"mid"`
	IsAutoRenew      int8   `json:"isAuto_renew"`
	ExpireDate       string `json:"expire_date"`
	NextDedutionDate string `json:"next_dedution_date"`
	AutoRenewLoop    string `json:"auto_renew_loop"`
	PriceTip         string `json:"price_tip"`
	ChannelID        int32  `json:"channel_id"`
	PayType          string `json:"pay_type"`
}

// OrderMessage .
type OrderMessage struct {
	LeftButton      string `json:"left_button"`
	RightButton     string `json:"right_button"`
	LeftButtonLink  string `json:"left_button_link"`
	RightButtonLink string `json:"right_button_link"`
	Title           string `json:"title"`
	Content         string `json:"content"`
}

//VipUserDiscountHistory .
type VipUserDiscountHistory struct {
	ID         int64  `json:"id"`
	Mid        int64  `json:"mid"`
	DiscountID int64  `json:"discount_id"`
	OrderNo    string `json:"order_no"`
	Status     int8   `json:"status"`
}

//PannelInfo .
type PannelInfo struct {
	VipMonths []*VipMonthsPriceBo `json:"vipMonths"`
	PayTypes  []*PayTypeBo        `json:"payTypes"`
	BcoinTips string              `json:"bcoinTips"`
}

//PayTypeBo .
type PayTypeBo struct {
	Name  string       `json:"name"`
	Code  string       `json:"code"`
	Banks []*PayBankBo `json:"banks"`
}

//VipMonthsPriceBo .
type VipMonthsPriceBo struct {
	ID                 int64   `json:"_"`
	Month              int16   `json:"month"`
	DiscountRate       string  `json:"discount_rate"`
	MonthStr           string  `json:"month_str"`
	MonthID            int64   `json:"month_id"`
	OrderType          int8    `json:"order_type"`
	MonthType          int8    `json:"month_type"`
	OriginalPrice      float64 `json:"original_price"`
	Selected           int8    `json:"selected"`
	FirstDiscountMoney float64 `json:"first_discount_money"`
	Price              float64 `json:"price"`
	Remark             string  `json:"remark"`
}

//PayBankBo .
type PayBankBo struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

//VipPayOrderLog vip pay order log table
type VipPayOrderLog struct {
	ID           int64   `json:"id"`
	OrderNo      string  `json:"order_no"`
	RefundID     string  `json:"refund_id"`
	RefundAmount float64 `json:"refund_amount"`
	Mid          int64   `json:"mid"`
	Status       int8    `json:"status"`
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
	PID         int64     `json:"pid"`
	UserIP      []byte    `json:"user_ip"`
}

//VipOldPayOrder vip pay order table
type VipOldPayOrder struct {
	ID          int64
	OrderNo     string
	AppID       int64
	Platform    int8
	OrderType   int8
	Mid         int64
	Bmid        int64
	BuyMonths   int16
	Money       float64
	Status      int8
	PaymentTime time.Time
	Ver         int64
	AppSubID    string
	CouponMoney float64
	PID         int64 //套餐ID
	UserIP      []byte
}

// VipOldRechargeOrder vip recharge order.
type VipOldRechargeOrder struct {
	AppID        int64
	PayMid       int64
	OrderNo      string
	RechargeBp   float64
	PayOrderNO   string
	Status       int8
	Remark       string
	Ver          int64
	Ctime        time.Time
	Bmid         int64
	ThirdTradeNO string
}
