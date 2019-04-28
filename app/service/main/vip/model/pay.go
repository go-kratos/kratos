package model

import "github.com/satori/go.uuid"

// qr tip
const (
	QrTip          = "请使用微信、支付宝或QQ扫码支付"
	QrAutoRenewTip = "请使用微信扫码支付"
)

//APIPayCancelResp api pay cancel resp.
type APIPayCancelResp struct {
	TraceID        string `json:"traceId"`
	ServerTime     int64  `json:"serverTime"`
	CustomerID     int64  `json:"customerId"`
	OrderID        string `json:"orderId"`
	OrderCloseTime int64  `json:"orderCloseTime"`
}

// UUID4 is generate uuid
func UUID4() string {
	return uuid.NewV4().String()
}

// PayParam call pay param.
type PayParam struct {
	CustomerID      int64  `json:"customerId,omitempty"`
	DeviceType      int8   `json:"deviceType,omitempty"`
	DefaultChoose   string `json:"defaultChoose,omitempty"`
	OrderID         string `json:"orderId,omitempty"`
	OrderCreateTime int64  `json:"orderCreateTime,omitempty"`
	OrderExpire     int    `json:"orderExpire,omitempty"`
	NotifyURL       string `json:"notifyUrl,omitempty"`
	CreateIP        string `json:"createIp,omitempty"`
	SignURL         string `json:"signUrl,omitempty"`
	ShowTitle       string `json:"showTitle,omitempty"`
	ShowContent     string `json:"showContent,omitempty"`
	TraceID         string `json:"traceId,omitempty"`
	Timestamp       int64  `json:"timestamp,omitempty"`
	Version         string `json:"version,omitempty"`
	SignType        string `json:"signType,omitempty"`
	Sign            string `json:"sign,omitempty"`
	ProductID       string `json:"productId,omitempty"`
	PayAmount       int32  `json:"payAmount,omitempty"`
	PlanID          int32  `json:"planId,omitempty"`
	DisplayAccount  string `json:"displayAccount,omitempty"`
	UID             int64  `json:"uid,omitempty"`
	ServiceType     int    `json:"serviceType,omitempty"`
	OriginalAmount  int32  `json:"originalAmount,omitempty"`
	ReturnURL       string `json:"returnUrl,omitempty"`
	FeeType         string `json:"feeType,omitempty"`
	SubscribeType   int    `json:"subscribeType,omitempty"`
}

// CreateOrderRet create order ret.
type CreateOrderRet struct {
	PayParam    map[string]interface{} `json:"pay_param"`
	Dprice      float64                `json:"dprice"`
	Oprice      float64                `json:"oprice"`
	CouponMoney float64                `json:"coupon_money"`
	UserIP      string                 `json:"user_ip"`
	PID         int64                  `json:"p_id"`
}

// PayQrCode resp.
type PayQrCode struct {
	CodeURL     string `json:"codeUrl"`
	ExpiredTime int64  `json:"expiredTime"`
}

// PayQrCodeResp pay qr resp.
type PayQrCodeResp struct {
	CodeURL     string  `json:"codeUrl"`
	ExpiredTime int64   `json:"expiredTime"`
	Amount      float64 `json:"amount"`
	SaveAmount  float64 `json:"saveAmount"`
	Tip         string  `json:"tip"`
	OrderNo     string  `json:"orderNo"`
}

// PayQrCodeRet pay qr resp.
type PayQrCodeRet struct {
	PayQrCodeResp *PayQrCodeResp `json:"pay_qr_data"`
	Dprice        float64        `json:"dprice"`
	CouponMoney   float64        `json:"coupon_money"`
	UserIP        string         `json:"user_ip"`
	PID           int64          `json:"p_id"`
}
