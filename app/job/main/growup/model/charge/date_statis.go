package charge

import (
	"go-common/library/time"
)

// DateStatis archive_charge_daily_statis archive_charge_weekly_statis archive_charge_monthly_statis
type DateStatis struct {
	ID           int64
	Count        int64
	MoneySection int64
	MoneyTips    string
	Charge       int64
	MinCharge    int64
	MaxCharge    int64
	CategoryID   int64
	CDate        time.Time
	CTime        time.Time
	MTime        time.Time
}
