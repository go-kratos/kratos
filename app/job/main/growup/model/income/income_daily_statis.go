package income

import (
	"go-common/library/time"
)

// DateStatis up_income_daily_statis and av_income_daily_statis struct
type DateStatis struct {
	ID           int64
	Count        int64
	MoneySection int64
	MoneyTips    string
	Income       int64
	MinIncome    int64
	MaxIncome    int64
	CategoryID   int64
	CDate        time.Time
	CTime        time.Time
	MTime        time.Time
}
