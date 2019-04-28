package model

import (
	"net"

	"go-common/library/time"
)

// ArgDevice .
type ArgDevice struct {
	Device string `form:"device" default:"pc"`
	Build  int64  `form:"build" default:"0"`
}

// ArgMid .
type ArgMid struct {
	Mid int64 `form:"mid" validate:"required"`
}

//ArgPannel .
type ArgPannel struct {
	Mid      int64  `form:"mid" validate:"required,min=1,gte=1"`
	Platform string `form:"platform" validate:"required"`
}

//ArgChangeHistory .
type ArgChangeHistory struct {
	Mid int64 `form:"mid" validate:"required"`
	Pn  int   `form:"pn"`
	Ps  int   `form:"ps"`
}

// ArgAddOrder add order form.
type ArgAddOrder struct {
	AppID    int64  `form:"appId" default:"0"`
	Months   int64  `form:"months" validate:"required"`
	AppsubID string `form:"appsubId"`
	BmID     int64  `form:"bmid"`
}

//ArgCreateOrder .
type ArgCreateOrder struct {
	Mid       int64   `form:"mid" validate:"required,min=1,gte=1"`
	AppID     int64   `form:"app_id" default:"0"`
	AppSubID  string  `form:"app_sub_id"`
	Months    int16   `form:"months" validate:"required,min=1,gte=1"`
	OrderType int8    `form:"order_type" `
	DType     int8    `form:"dtype"`
	Bmid      int64   `form:"bmid"`
	Platform  string  `form:"platform"`
	Price     float64 `form:"price"`
	IP        string  `form:"ip"`
}

//ArgCreateOrder2 .
type ArgCreateOrder2 struct {
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	Month       int32  `form:"months" validate:"required,min=1,gte=1"`
	Platform    string `form:"platform"`
	MobiApp     string `form:"mobi_app"`
	Device      string `form:"device"`
	AppID       int64  `form:"appId"`
	AppSubID    string `form:"appSubId"`
	OrderType   int8   `form:"orderType"`
	Dtype       int8   `form:"dtype"`
	ReturnURL   string `form:"returnUrl"`
	CouponToken string `form:"coupon_token"`
	Bmid        int64  `form:"bmid"`
	PanelType   string `form:"panel_type" default:"normal"`
	Build       int64  `form:"build"`
	IP          net.IP
}

// ArgPrice def.
type ArgPrice struct {
	Month          int16 `form:"month"`
	Platform       int   `form:"platform"`
	Mt             int8  `form:"mt"`
	DiscountStatus int8
}

// ArgPriceV2 arg price v2.
type ArgPriceV2 struct {
	Mid       int64
	Month     int16
	SubType   int8
	Token     string
	Platform  string
	PanelType string
	MobiApp   string
	Device    string
	Build     int64
}

// ArgCreateOrderPlatform def.
type ArgCreateOrderPlatform struct {
	Mid       int64   `form:"mid"`
	AppID     int64   `form:"appId"`
	Bmid      int64   `form:"bmid"`
	Month     int16   `form:"months"`
	Price     float64 `form:"price"`
	Platform  int     `form:"platform"`
	Dtype     int8    `form:"dtype"`
	OrderType int8    `form:"orderType"`
	AppSubID  string  `form:"appsubId"`
}

// ArgUseBatch def.
type ArgUseBatch struct {
	BatchID int64  `form:"batch_id" validate:"required" json:"batch_id"`
	Mid     int64  `form:"mid" validate:"required" json:"mid"`
	OrderNo string `form:"order_no" validate:"required" json:"order_no"`
	Remark  string `form:"remark" validate:"required" json:"remark"`
	Appkey  string `form:"appkey" validate:"required" json:"appkey"`
	Sign    string `form:"sign"`
	Ts      int64  `form:"ts"`
}

//ToMap .
func (arg *ArgUseBatch) ToMap() map[string]interface{} {
	mapVal := make(map[string]interface{})
	mapVal["batch_id"] = arg.BatchID
	mapVal["mid"] = arg.Mid
	mapVal["order_no"] = arg.OrderNo
	mapVal["remark"] = arg.Remark
	mapVal["appkey"] = arg.Appkey
	mapVal["sign"] = arg.Sign
	mapVal["ts"] = arg.Ts
	return mapVal
}

// ArgBuyVip def.
type ArgBuyVip struct {
	AppID     int64  `form:"appId" default:"0"`
	Months    int16  `form:"months" validate:"required"`
	PayWay    string `form:"payWay" validate:"required"`
	Bmid      int64  `form:"bmid"`
	BankCode  string `form:"bank_code"`
	ProductID string `form:"productId"`
	AppSubID  string `form:"appsubId"`
	AccessKey string `form:"access_key"`
	Platform  string `form:"platform"`
}

//ArgBuyPoint .
type ArgBuyPoint struct {
	Mid   int64 `form:"mid" validate:"required,min=1,gte=1"`
	Month int16 `form:"month" validate:"required" `
}

// ArgOldPayOrder def.
type ArgOldPayOrder struct {
	OrderNo      string  `form:"order_no" validate:"required"`
	AppID        int64   `form:"app_id"`
	Platform     int8    `form:"platform" `
	OrderType    int8    `form:"order_type"`
	AppSubID     string  `form:"app_sub_id"`
	Mid          int64   `form:"mid"`
	ToMid        int64   `form:"to_mid"`
	BuyMonths    int16   `form:"buy_months" validate:"required,min=1,gte=1"`
	Money        float64 `form:"money" validate:"required"`
	Status       int8    `form:"status"`
	PayType      int8    `form:"pay_type"`
	RechargeBp   float64 `form:"recharge_bp"`
	ThirdTradeNo string  `form:"third_trade_no"`
}

// ArgVipConfig .
type ArgVipConfig struct {
	Mid       int64  `form:"mid" validate:"required,min=1,gte=1"`
	Device    string `form:"device"`
	MobiApp   string `form:"mobi_app"`
	SortType  int8   `form:"sort_type"`
	PanelType string `form:"panel_type" default:"normal"`
	Build     int64  `form:"build"`
}

//ArgCodeOpened code opened.
type ArgCodeOpened struct {
	BisAppkey string    `form:"bis_appkey"`
	BisSign   string    `form:"bis_sign"`
	BisTs     int64     `form:"bis_ts"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
	Cursor    int64     `form:"cursor"`
}

//ToMap .
func (arg *ArgCodeOpened) ToMap() map[string]interface{} {
	mapval := make(map[string]interface{})
	mapval["appkey"] = arg.BisAppkey
	mapval["sign"] = arg.BisSign
	mapval["ts"] = arg.BisTs
	mapval["start_time"] = arg.StartTime
	mapval["end_time"] = arg.EndTime
	mapval["cursor"] = arg.Cursor
	return mapval
}

// ArgVipConfigMonth .
type ArgVipConfigMonth struct {
	Mid                   int64  `form:"mid" validate:"required,min=1,gte=1"`
	Device                string `form:"device"`
	MobiApp               string `form:"mobi_app"`
	Month                 int16  `form:"month" validate:"required,min=1,gte=1"`
	SubType               int8   `form:"sub_type" validate:"min=0,max=1"`
	CouponToken           string `form:"coupon_token"`
	Platform              string `form:"platform" default:"pc"`
	PanelType             string `form:"panel_type" default:"normal"`
	IgnoreAutoRenewStatus int8   `form:"ignore_autorenew_status"`
	Build                 int64  `form:"build"`
}

//ArgCancelUseCoupon cancel use coupon.
type ArgCancelUseCoupon struct {
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	CouponToken string `form:"coupon_token" validate:"required"`
}

// ArgAssociateVip associate vip arg.
type ArgAssociateVip struct {
	Platform string `form:"platform"`
	MobiApp  string `form:"mobi_app"`
	Device   string `form:"device"`
}

// ArgPriceByProduct arg price by product.
type ArgPriceByProduct struct {
	ProductID string `form:"product_id" validate:"required"`
}

// ArgVipPriceByID arg vip price.
type ArgVipPriceByID struct {
	ID int64 `form:"id" validate:"required"`
}
