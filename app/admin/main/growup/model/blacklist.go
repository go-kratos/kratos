package model

import (
	"go-common/library/time"
)

// Blacklist black list
type Blacklist struct {
	ID        int64     `json:"id" gorm:"column:id"`
	AvID      int64     `json:"av_id" gorm:"column:av_id"`
	MID       int64     `json:"mid" gorm:"column:mid"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	HasSigned int       `json:"has_signed" gorm:"column:has_signed"`
	Income    int64     `json:"income"`
	Reason    int       `json:"reason" gorm:"column:reason"`
	CType     int       `json:"ctype" gorm:"column:ctype"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted int       `json:"-"`
}

// AvIncomeStatis av income statis
type AvIncomeStatis struct {
	AvID        int64 `json:"av_id" gorm:"column:av_id"`
	TotalIncome int64 `json:"total_income" gorm:"total_income"`
}
