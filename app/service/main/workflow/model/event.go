package model

import (
	xtime "go-common/library/time"
)

const (
	// EventTypeAdminReply 管理员回复
	EventTypeAdminReply = int8(1)
	// EventTypeAdminNote 管理员备注
	EventTypeAdminNote = int8(2)
	// EventTypeUserReply 用户回复
	EventTypeUserReply = int8(3)
	// EventTypeSystemReply 系统回复
	EventTypeSystemReply = int(4)
)

// Event struct
type Event struct {
	ID          int32      `gorm:"column:id" json:"id"`
	Cid         int32      `gorm:"column:cid" json:"cid"`
	Event       int8       `gorm:"column:event" json:"event"`
	Adminid     int32      `gorm:"column:adminid" json:"adminid"`
	Content     string     `gorm:"column:content" json:"content"`
	Attachments string     `gorm:"column:attachments" json:"attachments"`
	Ctime       xtime.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       xtime.Time `gorm:"column:mtime" json:"mtime"`
}

// TableName by event
func (*Event) TableName() string {
	return "workflow_event"
}
