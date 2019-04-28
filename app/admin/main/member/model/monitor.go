package model

import (
	xtime "go-common/library/time"
)

// Monitor is.
type Monitor struct {
	ID        int64      `json:"id" gorm:"column:id"`
	Mid       int64      `json:"mid" gorm:"column:mid"`
	Operator  string     `json:"operator" gorm:"column:operator"`
	Remark    string     `json:"remark" gorm:"column:remark"`
	CTime     xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime     xtime.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted bool       `json:"is_deleted" gorm:"column:is_deleted"`

	// 昵称，后期拼进来
	Name string `json:"name" gorm:"-"`
}
