package model

import (
	"go-common/library/time"
)

// UpBill up_bill
type UpBill struct {
	MID          int64     `json:"mid"`
	Nickname     string    `json:"nickname"`
	Face         string    `json:"face"`
	FirstIncome  int64     `json:"first_income"`
	MaxIncome    int64     `json:"max_income"`
	TotalIncome  int64     `json:"total_income"`
	AvCount      int64     `json:"av_count"`
	AvMaxIncome  int64     `json:"av_max_income"`
	AvID         int64     `json:"-"`
	AvTitle      string    `json:"av_title"`
	QualityValue int64     `json:"quality_value"`
	DefeatNum    int       `json:"defeat_num"`
	Title        string    `json:"title"`
	ShareItems   string    `json:"share_items"`
	FirstTime    time.Time `json:"first_time"`
	MaxTime      time.Time `json:"max_time"`
	SignedAt     time.Time `json:"signed_at"`
	EndAt        time.Time `json:"end_at"`
	Join         bool      `json:"join"`
}
