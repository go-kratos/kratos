package telecom

import (
	"strconv"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

type TelecomJSON struct {
	FlowpackageID      int            `json:"flowPackageId"`
	FlowPackageSize    int            `json:"flowPackageSize"`
	FlowPackageType    int            `json:"flowPackageType"`
	TrafficAttribution int            `json:"trafficAttribution"`
	BeginTime          string         `json:"beginTime"`
	EndTime            string         `json:"endTime"`
	IsMultiplyOrder    int            `json:"isMultiplyOrder"`
	SettlementType     int            `json:"settlementType"`
	Operator           int            `json:"operator"`
	OrderStatus        int            `json:"orderStatus"`
	RemainedRebindNum  int            `json:"remainedRebindNum"`
	MaxbindNum         int            `json:"maxBindNum"`
	OrderID            string         `json:"orderId"`
	SignNo             string         `json:"signNo"`
	AccessToken        string         `json:"accessToken"`
	PhoneID            string         `json:"phoneId"`
	IsRepeatOrder      int            `json:"isRepeatOrder"`
	PayStatus          int            `json:"payStatus"`
	PayTime            string         `json:"payTime"`
	PayChannel         int            `json:"payChannel"`
	SignStatus         string         `json:"signStatus "`
	RefundStatus       int            `json:"refundStatus"`
	PayResult          *PayResultJSON `json:"payResult,omitempty"`
}

type PayResultJSON struct {
	IsRepeatOrder int `json:"isRepeatOrder"`
	RefundStatus  int `json:"refundStatus"`
	PayStatus     int `json:"payStatus"`
	PayChannel    int `json:"payChannel"`
}

type TelecomOrderJson struct {
	RequestNo  string       `json:"requestNo"`
	ResultType int          `json:"resultType"`
	Detail     *TelecomJSON `json:"detail"`
}

type TelecomRechargeJson struct {
	RequestNo  string        `json:"requestNo"`
	ResultType int           `json:"resultType"`
	Detail     *RechargeJSON `json:"detail"`
}

type RechargeJSON struct {
	RequestNo      string `json:"requestNo"`
	FcRechargeNo   string `json:"fcRechargeNo"`
	RechargeStatus int    `json:"rechargeStatus"`
	OrderTotalSize int    `json:"orderTotalSize"`
	FlowBalance    int    `json:"flowBalance"`
}

type OrderInfo struct {
	PhoneID       int        `json:"phone"`
	OrderID       int64      `json:"orderid"`
	OrderState    int        `json:"order_status"`
	IsRepeatorder int        `json:"isrepeatorder"`
	SignNo        string     `json:"sign_no"`
	Begintime     xtime.Time `json:"begintime"`
	Endtime       xtime.Time `json:"endtime"`
}

type Pay struct {
	OrderID   int64  `json:"orderid"`
	RequestNo int64  `json:"requestno,omitempty"`
	PayURL    string `json:"pay_url,omitempty"`
}

type SucOrder struct {
	FlowPackageID        string `json:"flowPackageId,omitempty"`
	Domain               string `json:"domain"`
	Port                 string `json:"port,omitempty"`
	PortInt              int    `json:"portInt"`
	KeyEffectiveDuration int    `json:"keyEffectiveDuration"`
	OrderKey             string `json:"orderKey"`
	FlowBalance          int    `json:"flowBalance"`
	FlowPackageSize      int    `json:"flowPackageSize"`
	AccessToken          string `json:"accessToken"`
	OrderIDStr           string `json:"orderId,omitempty"`
	OrderID              int64  `json:"orderid"`
}

type OrderFlow struct {
	FlowBalance int `json:"flowBalance"`
}

type PhoneConsent struct {
	Consent int `json:"consent"`
}

type TelecomMessageJSON struct {
	PhoneID       string `json:"phoneId"`
	ResultType    int    `json:"resultType"`
	ResultMessage string `json:"resultMsg"`
}

type OrderState struct {
	FlowBalance   int        `json:"flowBalance,omitempty"`
	FlowSize      int        `json:"flow_size"`
	OrderState    int        `json:"order_state"`
	Endtime       xtime.Time `json:"endtime,omitempty"`
	IsRepeatorder int        `json:"is_repeatorder"`
}

type OrderPhoneState struct {
	FlowPackageID int    `json:"flowPackageId"`
	FlowSize      int    `json:"flowPackageSize"`
	OrderState    int    `json:"orderStatus"`
	PhoneStr      string `json:"phoneId"`
}

func (s *TelecomJSON) TelecomJSONChange() {
	if s.PayResult != nil {
		s.IsRepeatOrder = s.PayResult.IsRepeatOrder
		s.RefundStatus = s.PayResult.RefundStatus
		s.PayStatus = s.PayResult.PayStatus
		s.PayChannel = s.PayResult.PayChannel
	}
}

func (t *OrderInfo) OrderInfoJSONChange(tjson *TelecomJSON) {
	t.PhoneID, _ = strconv.Atoi(tjson.PhoneID)
	t.OrderID, _ = strconv.ParseInt(tjson.OrderID, 10, 64)
	t.OrderState = tjson.OrderStatus
	t.IsRepeatorder = tjson.IsRepeatOrder
	t.SignNo = tjson.SignNo
	t.Begintime = timeStrToInt(tjson.BeginTime)
	t.Endtime = timeStrToInt(tjson.EndTime)
	t.TelecomChange()
}

// timeStrToInt
func timeStrToInt(timeStr string) (timeInt xtime.Time) {
	var err error
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	if err = timeInt.Scan(theTime); err != nil {
		log.Error("timeInt.Scan error(%v)", err)
	}
	return
}

// TelecomChange
func (t *OrderInfo) TelecomChange() {
	if t.Begintime.Time().IsZero() {
		t.Begintime = 0
	}
	if t.Endtime.Time().IsZero() {
		t.Endtime = 0
	}
}
