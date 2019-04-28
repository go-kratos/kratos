package model

import (
	"time"
)

// LogTask .
type LogTask struct {
	ID      int64
	Name    string
	Expect  int64
	Success int64
	Failure int64
	State   string
	CTime   time.Time
	MTime   time.Time
}

// Asset .
type Asset struct {
	ID       int64
	MID      int64
	OID      int64
	OType    string
	Currency string
	Price    int64
	State    string
	CTime    time.Time
	MTime    time.Time
}

// Order .
type Order struct {
	ID         int64
	OrderID    string
	MID        int64
	Biz        string
	Platform   string
	OID        int64
	OType      string
	Fee        int64
	RealFee    int64
	Currency   string
	PayID      string
	PayReason  string
	PayTime    time.Time
	RefundTime time.Time
	State      string // created 已创建, paying 支付中, paid 已支付, failed 支付失败, closed 已关闭, expired 已超时, finished 已完成(支付成功且对账成功)
	Version    int64  // 乐观锁，每次更新，值++
	CTime      time.Time
	MTime      time.Time
}

// LogOrder .
type LogOrder struct {
	ID        int64
	OrderID   string
	FromState string
	ToState   string
	Desc      string
	CTime     time.Time
	MTime     time.Time
}

// Bill .
type Bill struct {
	ID       int64
	BillID   string
	MID      int64
	Biz      string
	Currency string
	In       int64
	Out      int64
	Ver      int64 // 账单版本，比如 201812, 20181202
	Version  int64 // 乐观锁，每次更新++
	CTime    time.Time
	MTime    time.Time
}

// DailyBill .
type DailyBill struct {
	Bill
	MonthVer int64
}

// BizAccount .
type BizAccount struct {
	ID       int64
	Biz      string
	Currency string
	Balance  int64
	Ver      int64
	State    string
	CTime    time.Time
	MTime    time.Time
}

// UserAccount .
type UserAccount struct {
	BizAccount
	MID int64
}

// AccountLog .
type AccountLog struct {
	ID        int64
	Name      string
	AccountID int64
	From      int64
	To        int64
	Ver       int64
	State     string
	CTime     time.Time
	MTime     time.Time
}

// AggrIncomeUser .
type AggrIncomeUser struct {
	ID         int64
	MID        int64
	Currency   string
	PaySuccess int64
	PayError   int64
	TotalIn    int64
	TotalOut   int64
	CTime      time.Time
	MTime      time.Time
}

// AggrIncomeUserAsset .
type AggrIncomeUserAsset struct {
	ID         int64
	MID        int64
	Currency   string
	Ver        int64
	OID        int64
	OType      string
	PaySuccess int64
	PayError   int64
	TotalIn    int64
	TotalOut   int64
	CTime      time.Time
	MTime      time.Time
}

// OrderBadDebt .
type OrderBadDebt struct {
	ID      int64
	OrderID string
	Type    string
	State   string
	CTime   time.Time
	MTime   time.Time
}

// LogBillDaily .
type LogBillDaily struct {
	ID      int64
	BillID  string
	FromIn  int64
	ToIn    int64
	FromOut int64
	ToOut   int64
	OrderID string
	CTime   time.Time
	MTime   time.Time
}

// LogBillMonthly .
type LogBillMonthly struct {
	ID              int64
	BillID          string
	FromIn          int64
	ToIn            int64
	FromOut         int64
	ToOut           int64
	BillUserDailyID string
	CTime           time.Time
	MTime           time.Time
}

// OrderRechargeShell .
type OrderRechargeShell struct {
	ID      int64
	MID     int64
	OrderID string
	Biz     string
	Amount  int64
	PayMSG  string
	State   string
	Ver     int64
	CTime   time.Time
	MTime   time.Time
}

// OrderRechargeShellLog .
type OrderRechargeShellLog struct {
	ID                int64
	OrderID           string
	FromState         string
	ToState           string
	Desc              string
	BillUserMonthlyID string
	CTime             time.Time
	MTime             time.Time
}
