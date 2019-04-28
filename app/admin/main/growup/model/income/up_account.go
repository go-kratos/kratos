package income

import (
	"go-common/library/time"
)

// UpAccount up account
type UpAccount struct {
	ID                    int64     `json:"id"`
	MID                   int64     `json:"mid"`
	HasSignContract       int       `json:"has_sign_contract"`
	State                 int       `json:"state"`
	TotalIncome           int64     `json:"total_income"`
	TotalUnwithdrawIncome int64     `json:"total_unwithdraw_income"`
	TotalWithdrawIncome   int64     `json:"total_withdraw_income"`
	IncIncome             int64     `json:"inc_income"`
	LastWithdrawTime      time.Time `json:"last_withdraw_time"`
	Version               int64     `json:"version"`
	AllowanceState        int       `json:"allowance_state"`
	Nickname              string    `json:"nickname"`
	WithdrawDateVersion   string    `json:"withdraw_date_version"`
	IsDeleted             int       `json:"is_deleted"`
	CTime                 time.Time `json:"ctime"`
	MTime                 time.Time `json:"mtime"`
}
