package model

import (
	"encoding/json"

	xtime "go-common/library/time"
)

// coupon state.
const (
	NotUsed = iota
	InUse
	Used
	Expire
)

// coupon type.
const (
	BangumiVideo = iota + 1
	Cartoon
)

// call back status
const (
	Unpaid = iota
	PaidSuccess
)

// coupon state.
const (
	WaitPay = iota
	InPay
	PaySuccess
	PayFaild
)

// blance change type
const (
	VipSalary int8 = iota + 1
	SystemAdminSalary
	Consume
	ConsumeFaildBack
)

//allowance origin
const (
	AllowanceNone = iota
	AllowanceSystemAdmin
	AllowanceBusinessReceive
	AllowanceBusinessNewYear
)

// CouponInfo coupon info.
type CouponInfo struct {
	ID          int64      `json:"_"`
	CouponToken string     `json:"coupon_token"`
	Mid         int64      `json:"mid"`
	State       int64      `json:"state"`
	StartTime   int64      `json:"start_time"`
	ExpireTime  int64      `json:"expire_time"`
	Origin      int64      `json:"origin"`
	CouponType  int64      `json:"coupon_type"`
	OrderNO     string     `json:"order_no"`
	Ver         int64      `json:"ver"`
	Oid         int64      `json:"oid"`
	Remark      string     `json:"remark"`
	UseVer      int64      `json:"use_ver"`
	CTime       xtime.Time `json:"-"`
	MTime       xtime.Time `json:"-"`
}

// MsgCanal canal message struct.
type MsgCanal struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// CallBackRet .
type CallBackRet struct {
	Ver    int64 `json:"ver"`
	IsPaid int8  `json:"is_paid"`
}

// NotifyParam notify param.
type NotifyParam struct {
	CouponToken string `json:"coupon_token"`
	Mid         int64  `json:"mid"`
	NotifyURL   string `json:"notify_url"`
	NotifyCount int    `json:"count"`
	Type        int64  `json:"type"`
}

// CouponChangeLog coupon change log.
type CouponChangeLog struct {
	ID          int64      `json:"-"`
	CouponToken string     `json:"coupon_token"`
	Mid         int64      `json:"mid"`
	State       int8       `json:"state"`
	Ctime       xtime.Time `json:"ctime"`
	Mtime       xtime.Time `json:"mtime"`
}

// CouponOrder coupon order info.
type CouponOrder struct {
	ID           int64      `json:"id"`
	OrderNo      string     `json:"order_no"`
	Mid          int64      `json:"mid"`
	Count        int64      `json:"count"`
	State        int8       `json:"state"`
	CouponType   int8       `json:"coupon_type"`
	ThirdTradeNo string     `json:"third_trade_no"`
	Remark       string     `json:"remark"`
	Tips         string     `json:"tips"`
	UseVer       int64      `json:"use_ver"`
	Ver          int64      `json:"ver"`
	Ctime        xtime.Time `json:"-"`
	Mtime        xtime.Time `json:"-"`
}

// CouponOrderLog coupon order log.
type CouponOrderLog struct {
	ID      int64      `json:"id"`
	OrderNo string     `json:"order_no"`
	Mid     int64      `json:"mid"`
	State   int8       `json:"state"`
	Ctime   xtime.Time `json:"ctime"`
	Mtime   xtime.Time `json:"mtime"`
}

// CouponBalanceChangeLog coupon balance change log.
type CouponBalanceChangeLog struct {
	ID            int64      `json:"id"`
	OrderNo       string     `json:"order_no"`
	Mid           int64      `json:"mid"`
	BatchToken    string     `json:"batch_token"`
	Balance       int64      `json:"balance"`
	ChangeBalance int64      `json:"change_balance"`
	ChangeType    int8       `json:"change_type"`
	Ctime         xtime.Time `json:"ctime"`
	Mtime         xtime.Time `json:"mtime"`
}

// CouponBalanceInfo def.
type CouponBalanceInfo struct {
	ID         int64      `protobuf:"varint,1,opt,name=ID,proto3" json:"_"`
	BatchToken string     `protobuf:"bytes,2,opt,name=BatchToken,proto3" json:"batch_token"`
	Mid        int64      `protobuf:"varint,3,opt,name=Mid,proto3" json:"mid"`
	Balance    int64      `protobuf:"varint,4,opt,name=Balance,proto3" json:"balance"`
	StartTime  int64      `protobuf:"varint,5,opt,name=StartTime,proto3" json:"start_time"`
	ExpireTime int64      `protobuf:"varint,6,opt,name=ExpireTime,proto3" json:"expire_time"`
	Origin     int64      `protobuf:"varint,7,opt,name=Origin,proto3" json:"origin"`
	CouponType int64      `protobuf:"varint,8,opt,name=CouponType,proto3" json:"coupon_type"`
	Ver        int64      `protobuf:"varint,9,opt,name=Ver,proto3" json:"ver"`
	CTime      xtime.Time `protobuf:"varint,10,opt,name=CTime,proto3,casttype=go-common/library/time.Time" json:"-"`
	MTime      xtime.Time `protobuf:"varint,11,opt,name=MTime,proto3,casttype=go-common/library/time.Time" json:"-"`
}

// CouponAllowanceInfo struct .
type CouponAllowanceInfo struct {
	ID          int64      `gorm:"column:id" json:"id" form:"id"`
	CouponToken string     `gorm:"column:coupon_token" json:"coupon_token" form:"coupon_token"`
	MID         int64      `gorm:"column:mid" json:"mid" form:"mid"`
	State       int8       `gorm:"column:state" json:"state" form:"state"`
	StartTime   int64      `gorm:"column:start_time" json:"start_time" form:"start_time"`
	ExpireTime  int64      `gorm:"column:expire_time" json:"expire_time" form:"expire_time"`
	Origin      int8       `gorm:"column:origin" json:"origin" form:"origin"`
	Ver         int64      `gorm:"column:ver" json:"ver" form:"ver"`
	BatchToken  string     `gorm:"column:batch_token" json:"batch_token" form:"batch_token"`
	OrderNo     string     `gorm:"column:order_no" json:"order_no" form:"order_no"`
	Amount      float64    `gorm:"column:amount" json:"amount" form:"amount"`
	FullAmount  float64    `gorm:"column:full_amount" json:"full_amount" form:"full_amount"`
	Ctime       xtime.Time `gorm:"column:ctime" json:"-" form:"ctime"`
	Mtime       xtime.Time `gorm:"column:mtime" json:"-" form:"mtime"`
	Remark      string     `gorm:"column:remark" json:"remark" form:"remark"`
	AppID       int64      `gorm:"column:app_id" json:"app_id" form:"app_id"`
}
