package databus

import (
	"go-common/library/time"
)

// TableName case tablename
func (*OldAuth) TableName() string {
	return "auth"
}

// OldAuth group model
type OldAuth struct {
	ID        int64     `gorm:"column:id" json:"id"`
	GroupName string    `gorm:"column:group_name" json:"group_name"`
	Topic     string    `gorm:"column:topic" json:"topic"`
	Business  string    `gorm:"business" json:"business"`
	AppID     int64     `gorm:"column:app_id" json:"app_id"`
	Operation int8      `gorm:"column:operation" json:"operation"`
	Ctime     time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime     time.Time `gorm:"column:mtime" json:"mtime"`
}
