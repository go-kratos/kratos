package model

import (
	"encoding/json"

	"go-common/library/time"
)

//VipUserInfoMsg get databus  vip_user_info msg
type VipUserInfoMsg struct {
	ID                   int64  `json:"id"`
	Mid                  int64  `json:"mid"`
	Type                 int8   `json:"vip_type"`
	Status               int8   `json:"vip_status"`
	StartTime            string `json:"vip_start_time"`
	OverdueTime          string `json:"vip_overdue_time"`
	AnnualVipOverdueTime string `json:"annual_vip_overdue_time"`
	RecentTime           string `json:"vip_recent_time"`
	Wander               int8   `json:"wander"`
	AutoRenewed          int8   `json:"auto_renewed"`
	IsAutoRenew          int8   `json:"is_auto_renew"`
	Ver                  int64  `json:"ver"`
	PayChannelID         int64  `json:"pay_channel_id"`
	IosOverdueTime       string `json:"ios_overdue_time"`
}

//VipUserInfoNewMsg .
type VipUserInfoNewMsg struct {
	ID                   int64  `json:"id"`
	Mid                  int64  `json:"mid"`
	Ver                  int64  `json:"ver"`
	VipType              int8   `json:"vip_type"`
	VipPayType           int8   `json:"vip_pay_type"`
	PayChannelID         int64  `json:"pay_channel_id"`
	VipStatus            int8   `json:"vip_status"`
	VipStartTime         string `json:"vip_start_time"`
	VipRecentTime        string `json:"vip_recent_time"`
	VipOverdueTime       string `json:"vip_overdue_time"`
	AnnualVipOverdueTime string `json:"annual_vip_overdue_time"`
	IosOverdueTime       string `json:"ios_overdue_time"`
}

//VipPointChangeHistoryMsg get databus json data
type VipPointChangeHistoryMsg struct {
	ID           int    `json:"id"`
	Mid          int    `json:"mid"`
	Point        int    `json:"point"`
	OrderID      string `json:"order_id"`
	ChangeType   int    `json:"change_type"`
	ChangeTime   string `json:"change_time"`
	RelationID   string `json:"relation_id"`
	PointBalance int    `json:"point_balance"`
	Remark       string `json:"remark"`
	Operator     string `json:"operator"`
}

//VipPayOrderOldMsg get databus json data
type VipPayOrderOldMsg struct {
	ID          int64   `json:"id"`
	OrderNo     string  `json:"order_no"`
	AppID       int64   `json:"app_id"`
	Platform    int8    `json:"platform"`
	OrderType   int8    `json:"order_type"`
	Mid         int64   `json:"mid"`
	Bmid        int64   `json:"bmid"`
	BuyMonths   int16   `json:"buy_months"`
	Money       float64 `json:"money"`
	RechargeBp  float64 `json:"recharge_bp"`
	Status      int8    `json:"status"`
	PayType     int8    `json:"pay_type"`
	PaymentTime string  `json:"payment_time"`
	Ver         int64   `json:"ver"`
	AppSubID    string  `json:"app_sub_id"`
	CouponMoney float64 `json:"coupon_money"`
	Ctime       string  `json:"ctime"`
	Mtime       string  `json:"mtime"`
}

//VipPayOrderNewMsg .
type VipPayOrderNewMsg struct {
	ID           int64   `json:"id"`
	OrderNo      string  `json:"order_no"`
	AppID        int64   `json:"app_id"`
	Platform     int8    `json:"platform"`
	OrderType    int8    `json:"order_type"`
	Mid          int64   `json:"mid"`
	ToMid        int64   `json:"to_mid"`
	BuyMonths    int16   `json:"buy_months"`
	Money        float64 `json:"money"`
	RechargeBp   float64 `json:"recharge_bp"`
	ThirdTradeNo string  `json:"third_trade_no"`
	Status       int8    `json:"status"`
	PayType      string  `json:"pay_type"`
	PaymentTime  string  `json:"payment_time"`
	Ver          int64   `json:"ver"`
	AppSubID     string  `json:"app_sub_id"`
	CouponMoney  float64 `json:"coupon_money"`
}

//VipRechargeOrderMsg get databus json data
type VipRechargeOrderMsg struct {
	ID           int64   `json:"id"`
	AppID        int64   `json:"app_id"`
	PayMid       int64   `json:"pay_mid"`
	OrderNo      string  `json:"order_no"`
	RechargeBp   float64 `json:"recharge_bp"`
	ThirdTradeNo string  `json:"third_trade_no"`
	PayOrderNo   string  `json:"pay_order_no"`
	Status       int8    `json:"status"`
	Ver          int64   `json:"ver"`
	Bmid         int64   `json:"bmid"`
}

//VipChangeHistoryMsg vip change history msg
type VipChangeHistoryMsg struct {
	Mid         int64  `json:"mid"`
	ChangeType  int8   `json:"change_type"`
	ChangeTime  string `json:"change_time"`
	Days        int32  `json:"days"`
	Month       int16  `json:"month"`
	OperatorID  string `json:"operator_id"`
	RelationID  string `json:"relation_id"`
	BatchID     int64  `json:"batch_id"`
	Remark      string `json:"remark"`
	BatchCodeID int64  `json:"batch_code_id"`
}

//VipBcoinSalaryMsg .
type VipBcoinSalaryMsg struct {
	ID            int64  `json:"id"`
	Mid           int64  `json:"mid"`
	Status        int8   `json:"status"`
	GiveNowStatus int8   `json:"give_now_status"`
	Payday        string `json:"month"`
	Amount        int32  `json:"amount"`
	Memo          string `json:"memo"`
}

//Message databus message
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// LoginLog login log.
type LoginLog struct {
	Mid        int64     `json:"mid,omitempty" form:"mid"`
	IP         uint32    `json:"loginip" form:"ip"`
	Location   string    `json:"location"`
	LocationID int64     `json:"location_id,omitempty"`
	Time       time.Time `json:"timestamp,omitempty"`
	Type       int8      `json:"type,omitempty"`
}
