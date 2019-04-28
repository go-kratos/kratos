package model

import (
	"go-common/library/time"
)

// AvIncome av income
type AvIncome struct {
	ID          int64
	AvID        int64
	MID         int64
	Income      int64
	TotalIncome int64
	TagID       int64
	Date        time.Time
	IsDeleted   int
}

// UpAccount up_account
type UpAccount struct {
	ID                    int64
	MID                   int64
	TotalIncome           int64
	TotalUnwithdrawIncome int64
	WithdrawDateVersion   string
	AvCount               int64
	MonthIncome           int64
	Nickname              string
}

// UpIncome up_income(weekly, monthly)
type UpIncome struct {
	ID           int64
	MID          int64
	AvCount      int64
	AvIncome     int64
	ColumnCount  int64
	ColumnIncome int64
	BgmCount     int64
	BgmIncome    int64
	Income       int64
	TaxMoney     int64
	BaseIncome   int64
	TotalIncome  int64
	Date         time.Time
}
