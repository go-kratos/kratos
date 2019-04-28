package model

import (
	"time"
)

// AccountUser .
type AccountUser struct {
	ID       int64
	Biz      string
	MID      int64
	Currency string
	Balance  int64
	Ver      int64
	State    string
	CTime    time.Time
	MTime    time.Time
}

// AccountBiz .
type AccountBiz struct {
	ID       int64
	Biz      string
	Currency string
	Balance  int64
	Ver      int64
	State    string
	CTime    time.Time
	MTime    time.Time
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
