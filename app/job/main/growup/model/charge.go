package model

import (
	"go-common/library/time"
)

// AvDailyCharge av_daily_charge
type AvDailyCharge struct {
	AvID        int64
	MID         int64
	IncCharge   int
	TotalCharge int
	Date        time.Time
	IsDeleted   int
	CTime       time.Time
	MTime       time.Time
}

// UpCharge up_charge
type UpCharge struct {
	MID         int64
	AvCount     int
	IncCharge   int
	TotalCharge int
	Date        time.Time
	IsDeleted   int
	CTime       time.Time
	MTime       time.Time
}

// AvChargeRatio av_charge_ratio
type AvChargeRatio struct {
	ID         int64
	AvID       int64
	Ratio      int64
	AdjustType int
}
