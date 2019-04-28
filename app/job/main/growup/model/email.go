package model

import "time"

// IncomeInfo income info.
type IncomeInfo struct {
	ID          int64     `json:"id"`
	AVID        int64     `json:"av_id"`
	MID         int64     `json:"mid"`
	TagID       int64     `json:"-"`
	Income      int64     `json:"income"`
	BaseIncome  int64     `json:"-"`
	TotalIncome int64     `json:"total_income"`
	TaxMoney    int64     `json:"tax_money"`
	UploadTime  time.Time `json:"-"`
	IsDeleted   int       `json:"-"`
	Date        time.Time `json:"date"`
	BType       int       `json:"-"`
}

// TotalInfo total info.
type TotalInfo struct {
	TotalIncome int64 `json:"total_income"`
	MIDCount    int   `json:"mid_count"`
	AVCount     int   `json:"av_count"`
}

// MIDInfo mid info.
type MIDInfo struct {
	ID          int64  `json:"id"`
	MID         int64  `json:"mid"`
	Income      int64  `json:"income"`
	TotalIncome int64  `json:"total_income"`
	IsDeleted   int    `json:"is_deleted"`
	NickName    string `json:"nickname"`
}

// AVIDInfo av info.
type AVIDInfo struct {
	AVID        int64  `json:"av_id"`
	MID         int64  `json:"mid"`
	NickName    string `json:"nickname"`
	Income      int64  `json:"income"`
	TotalIncome int64  `json:"total_income"`
}

// TagInfo email tag info
type TagInfo struct {
	ID          int64  `json:"-"`
	Tag         string `json:"tag"`
	Category    int    `json:"-"`
	AVID        int64  `json:"-"`
	AVCount     int    `json:"av_count"`
	Income      int64  `json:"income"`
	TotalIncome int64  `json:"total_income"`
	IsCommon    int    `json:"-"`
	IsDeleted   int    `json:"-"`
}
