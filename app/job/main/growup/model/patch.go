package model

import (
	"go-common/library/time"
)

// Patch patch
type Patch struct {
	Tax       int64
	Income    int64
	OldTax    int64
	OldIncome int64
	MID       int64
	TagID     int64
}

// Av av_income
type Av struct {
	AvID     int64
	MID      int64
	TagID    int64
	Income   int64
	TaxMoney int64
}

// AvBaseIncome av_base_income
type AvBaseIncome struct {
	MID          int64
	Date         time.Time
	AvBaseIncome int64
}
