package income

import (
	"go-common/library/time"
)

// UpIncome up income
type UpIncome struct {
	ID          int64     `json:"id"`
	MID         int64     `json:"mid"`
	AvCount     int64     `json:"-"`
	ColumnCount int64     `json:"-"`
	BgmCount    int64     `json:"-"`
	Count       int64     `json:"count"`
	TaxMoney    int64     `json:"tax_money"`
	Income      int64     `json:"income"`
	TotalIncome int64     `json:"total_income"`
	Date        time.Time `json:"-"`
	Nickname    string    `json:"nickname"`
	DateFormat  string    `json:"date_format"`
	BaseIncome  int64     `json:"base_income"`
	Breach      int64     `json:"breach"`
	ExtraIncome int64     `json:"extra_income"`
}

// UpDailyStatis up daily statis
type UpDailyStatis struct {
	Ups    int
	Income int64
	Date   time.Time
}

// UpStatisRsp response
type UpStatisRsp struct {
	Income int64
	Ups    int
}
