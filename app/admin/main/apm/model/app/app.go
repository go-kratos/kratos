package app

import (
	xtime "go-common/library/time"
)

// TableName case tablename
func (*App) TableName() string {
	return "app"
}

// App app
type App struct {
	ID        int64      `gorm:"column:id" json:"id"`
	AppTreeID int64      `gorm:"column:app_tree_id" json:"app_tree_id"`
	AppID     string     `gorm:"column:app_id" json:"app_id"`
	Limit     int64      `gorm:"column:limit" json:"limit"`
	CTime     xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime     xtime.Time `gorm:"column:mtime" json:"mtime"`
}
