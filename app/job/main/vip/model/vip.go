package model

import (
	"go-common/library/time"
)

//VipAppInfo app info
type VipAppInfo struct {
	ID       int64  `json:"id"`
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	PurgeURL string `json:"purgeUrl"`
	AppKey   string `json:"appKey"`
}

// VipPoint vip_point table
type VipPoint struct {
	ID           int `json:"id"`
	Mid          int `json:"mid"`
	PointBalance int `json:"point_balance"`
	Ver          int `json:"ver"`
}

//VipPointChangeHistory vip_point_change_history table
type VipPointChangeHistory struct {
	ID           int       `json:"id"`
	Mid          int       `json:"mid"`
	Point        int       `json:"point"`
	OrderID      string    `json:"orderId"`
	ChangeType   int       `json:"changeType"`
	ChangeTime   time.Time `json:"changeTime"`
	RelationID   string    `json:"relationId"`
	PointBalance int       `json:"pointBalance"`
	Remark       string    `json:"remark"`
	Operator     string    `json:"operator"`
}

//VipPayOrderLog vip pay order log table
type VipPayOrderLog struct {
	ID      int64  `json:"id"`
	OrderNo string `json:"orderNo"`
	Mid     int64  `json:"mid"`
	Status  int8   `json:"status"`
}

//VipPayOrder VipPayOrder table
type VipPayOrder struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"orderNo"`
	AppID        int64     `json:"appId"`
	Platform     int8      `json:"platform"`
	OrderType    int8      `json:"orderType"`
	Mid          int64     `json:"mid"`
	ToMid        int64     `json:"toMid"`
	BuyMonths    int16     `json:"buyMonths"`
	Money        float64   `json:"money"`
	RechargeBp   float64   `json:"rechargeBp"`
	Status       int8      `json:"status"`
	PayType      int8      `json:"payType"`
	ThirdTradeNo string    `json:"thirdTradeNo"`
	PaymentTime  time.Time `json:"paymentTime"`
	Ver          int64     `json:"ver"`
	AppSubID     string    `json:"appSubId"`
	CouponMoney  float64   `json:"coupon_money"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

//VipUserInfo vip user info table
type VipUserInfo struct {
	ID                   int64     `json:"id"`
	Mid                  int64     `json:"mid"`
	Type                 int8      `json:"vipType"`
	PayType              int8      `json:"payType"`
	PayChannelID         int64     `json:"payChannelId"`
	Status               int8      `json:"vipStatus"`
	Ver                  int64     `json:"ver"`
	OldVer               int64     `json:"old_ver"`
	StartTime            time.Time `json:"vipStartTime"`
	RecentTime           time.Time `json:"vipRecentTime"`
	OverdueTime          time.Time `json:"vipOverdueTime"`
	AnnualVipOverdueTime time.Time `json:"annualVipOverdueTime"`
	AutoRenewed          int8      `json:"auto_renewed"`
	IosOverdueTime       time.Time `json:"ios_overdue_time"`
	Ctime                time.Time `json:"ctime"`
	Mtime                time.Time `json:"mtime"`
}

//VipPushData .
type VipPushData struct {
	ID              int64     `json:"id"`
	DisableType     int8      `json:"disable_type"`
	GroupName       string    `json:"group_name"`
	Title           string    `json:"title" `
	Content         string    `json:"content"`
	PushTotalCount  int32     `json:"-"`
	PushedCount     int32     `json:"-"`
	PushProgress    string    `json:"push_progress"`
	ProgressStatus  int8      `json:"progress_status"`
	Status          int8      `json:"status"`
	Platform        string    `json:"platform"`
	LinkType        int32     `json:"link_type"`
	ErrorCode       int32     `json:"error_code"`
	LinkURL         string    `json:"link_url"`
	ExpiredDayStart int32     `json:"expired_day_start" `
	ExpiredDayEnd   int64     `json:"expired_day_end" `
	EffectStartDate time.Time `json:"effect_start_date" `
	EffectEndDate   time.Time `json:"effect_end_date" `
	PushStartTime   string    `json:"push_start_time" `
	PushEndTime     string    `json:"push_end_time" `
}

//VipUserDiscountHistory vip user discount  history table
type VipUserDiscountHistory struct {
	ID         int64  `json:"id"`
	Mid        int64  `json:"mid"`
	DiscountID int32  `json:"discountId"`
	OrderNo    string `json:"orderNo"`
	Status     int8   `json:"status"`
}

//VipBcoinSalary .
type VipBcoinSalary struct {
	ID            int64     `json:"id"`
	Mid           int64     `json:"mid"`
	Status        int8      `json:"status"`
	GiveNowStatus int8      `json:"giveNowStatus"`
	Payday        time.Time `json:"month"`
	Amount        int32     `json:"amount"`
	Memo          string    `json:"memo"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

//VipInfoDB vip user info db
type VipInfoDB struct {
	ID                   int64     `json:"id"`
	Mid                  int64     `json:"mid"`
	Ver                  int64     `json:"ver"`
	Type                 int8      `json:"vip_type"`
	PayType              int8      `json:"vip_pay_type"`
	PayChannelID         int64     `json:"pay_channel_id"`
	Status               int8      `json:"vip_status"`
	StartTime            time.Time `json:"vip_start_time"`
	RecentTime           time.Time `json:"vip_recent_time"`
	OverdueTime          time.Time `json:"vip_overdue_time"`
	AnnualVipOverdueTime time.Time `json:"annual_vip_overdue_time"`
	IosOverdueTime       time.Time `json:"ios_overdue_time"`
	Ctime                time.Time `json:"ctime"`
	Mtime                time.Time `json:"mtime"`
}

//VipConfig .
type VipConfig struct {
	ID        int64  `json:"id"`
	ConfigKey string `json:"config_key"`
	Content   string `json:"content"`
}

//VipResourceBatchCode .
type VipResourceBatchCode struct {
	ID           int64     `json:"id"`
	BusinessID   int64     `json:"business_id"`
	PoolID       int64     `json:"pool_id"`
	Status       int8      `json:"status"`
	Type         int8      `json:"type"`
	BatchName    string    `json:"batch_name"`
	Reason       string    `json:"reason"`
	Unit         int32     `json:"unit"`
	Count        int64     `json:"count"`
	SurplusCount int64     `json:"surplus_count"`
	Price        float64   `json:"price"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

//VipResourceCode vip resource code
type VipResourceCode struct {
	ID          int64  `json:"id"`
	Bmid        int64  `json:"bmid"`
	RelationID  string `json:"relation_id"`
	Code        string `json:"code"`
	Status      int8   `json:"status"`
	BatchCodeID int64  `json:"batch_code_id"`
}
