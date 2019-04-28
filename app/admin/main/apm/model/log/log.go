package log

import (
	xtime "go-common/library/time"
)

// TableName case tablename
func (*Log) TableName() string {
	return "log"
}

// Log log model
type Log struct {
	ID       int64      `gorm:"column:id" json:"id"`
	UserName string     `gorm:"column:username" json:"username"`
	Business string     `gorm:"column:business" json:"business"`
	Info     string     `gorm:"column:info" json:"info"`
	CTime    xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime    xtime.Time `gorm:"column:mtime" json:"mtime"`
}
