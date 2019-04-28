package income

import (
	"go-common/library/time"
)

// UpIncomeWithdraw up income withdraw
type UpIncomeWithdraw struct {
	ID             int64     `json:"id"`
	MID            int64     `json:"mid"`
	WithdrawIncome int64     `json:"-"`
	Income         string    `json:"withdraw_income"`
	DateVersion    string    `json:"date_version"`
	MTime          time.Time `json:"withdraw_time"`
	Nickname       string    `json:"nickname"`
}

// UpWithdrawRes response
type UpWithdrawRes struct {
	MID                int64     `json:"mid"`
	Nickname           string    `json:"nickname"`
	WithdrawIncome     string    `json:"withdraw_income"`
	UnWithdrawIncome   string    `json:"unwithdraw_income"`
	LastWithdrawIncome string    `json:"last_withdraw_income"`
	WithdrawDate       string    `json:"withdraw_date"`
	MTime              time.Time `json:"mtime"`
}
