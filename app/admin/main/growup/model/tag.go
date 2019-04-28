package model

import (
	"time"

	xtime "go-common/library/time"
)

// TagInfo tag_info
type TagInfo struct {
	ID              int64      `json:"id" form:"id"`
	Tag             string     `json:"tag" form:"tag" validate:"required"`
	Dimension       int        `json:"dimension" form:"dimension"`
	Category        int        `json:"category" form:"category"`
	Business        int        `json:"business" form:"business"`
	StartTime       xtime.Time `json:"start_at" form:"start_time" validate:"required"`
	EndTime         xtime.Time `json:"end_at" form:"end_time" validate:"required"`
	CreateTime      xtime.Time `json:"ctime"`
	Creator         string     `json:"creator"`
	AdjustType      int        `json:"adjust_type" form:"adjust_type"`
	Ratio           int        `json:"-" form:"ratio" validate:"required"`
	RetRatio        float32    `json:"ratio"`
	IsCommon        int        `json:"is_common"`
	TotalIncome     int        `json:"total_income"`
	UpCount         int        `json:"up_count"`
	IsDeleted       int        `json:"is_deleted"`
	ActivityID      int64      `json:"activity_id" form:"activity_id"`
	Icon            string     `json:"icon" form:"icon"`
	UploadStartTime xtime.Time `json:"upload_start_time" form:"upload_start_time" validate:"required"`
	UploadEndTime   xtime.Time `json:"upload_end_time" form:"upload_end_time" validate:"required"`
	MIDs            []int64    `json:"-" form:"mids,split"`
}

// UpTagIncome calculate tag income
type UpTagIncome struct {
	ID          int64 `json:"id"`
	AvID        int64 `gorm:"column:av_id"`
	MID         int64 `gorm:"column:mid"`
	Income      int   `gorm:"column:income"`
	BaseIncome  int   `gorm:"base_income"`
	TotalIncome int   `gorm:"column:total_income"`
	TaxMoney    int   `gorm:"column:tax_money"`
	IsDeleted   int   `gorm:"column:is_deleted"`
	Date        time.Time
}

// UpIncomeInfo up info
type UpIncomeInfo struct {
	MID          int64      `json:"mid"`
	Nickname     string     `json:"nickname"`
	CreateTime   xtime.Time `json:"ctime"`
	BaseIncome   int        `json:"base_income"`
	AdjustIncome int        `json:"adjust_income"`
	TotalIncome  int        `json:"total_income"`
	IsDeleted    int        `json:"is_deleted"`
}

// AvIncomeInfo av income info
type AvIncomeInfo struct {
	AVID         int64      `json:"av_id"`
	MID          int64      `json:"mid"`
	Nickname     string     `json:"nickname"`
	Category     int        `json:"category"`
	CreateTime   xtime.Time `json:"ctime"`
	BaseIncome   int        `json:"base_income"`
	AdjustIncome int        `json:"adjust_income"`
	TotalIncome  int        `json:"total_income"`
}

// Nickname get nickname from up_category_info
type Nickname struct {
	Nickname  string `json:"nickname" gorm:"column:nick_name"`
	IsDeleted int    `json:"is_deleted" gorm:"column:is_deleted"`
}

// AVs get avid
type AVs struct {
	AVID      int64 `json:"av_id" gorm:"column:av_id"`
	IsDeleted int   `json:"is_deleted" gorm:"column:is_deleted"`
}

// Activity activity
type Activity struct {
	TagID      int64      `json:"tag_id"`
	ActivityID int        `json:"activity_id"`
	Category   int        `json:"category"`
	MID        int64      `json:"mid"`
	ArchiveID  int64      `json:"archive_id"`
	CreateTime xtime.Time `json:"create_time"`
}

// Details tag details.
type Details struct {
	Date         string `json:"date"`
	UpCnt        int    `json:"up_cnt"`
	AvCnt        int    `json:"av_cnt"`
	Income       int    `json:"income"`
	BaseIncome   int    `json:"base_income"`
	AdjustIncome int    `json:"adjust_income"`
}
