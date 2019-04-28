package model

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
	LastWithdrawTime      time.Time `json:"last_withdraw_time"`
	Version               int64     `json:"version"`
	IsDeleted             int       `json:"is_deleted"`
	CTime                 time.Time `json:"ctime"`
	MTime                 time.Time `json:"mtime"`
	AllowanceState        int       `json:"allowance_state"`
	Nickname              string    `json:"nick_name"`
	WithdrawDateVersion   string    `json:"withdraw_date_version"`
}

// UpIncomeWithdraw up income withdraw
type UpIncomeWithdraw struct {
	ID             int64     `json:"id"`
	MID            int64     `json:"mid"`
	WithdrawIncome int64     `json:"withdraw_income"`
	DateVersion    string    `json:"date_version"`
	State          int       `json:"state"`
	IsDeleted      int       `json:"is_deleted"`
	CTime          time.Time `json:"ctime"`
	MTime          time.Time `json:"mtime"`
}

// WithdrawVo withdraw
type WithdrawVo struct {
	MID          int64   `json:"mid"`
	ThirdOrderNo string  `json:"third_order_no"`
	ThirdCoin    float64 `json:"third_coin"`
	CTime        string  `json:"ctime"`
	NotifyURL    string  `json:"notify_url"`
}
