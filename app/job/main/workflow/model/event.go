package model

import xtime "go-common/library/time"

// consts for workflow event
const (
	// EventAdminReply 管理员回复
	EventAdminReply = 1
	// EventAdminNote 管理员回复并记录
	EventAdminNote = 2
	// EventUserReply 用户回复
	EventUserReply = 3
	// EventSystem 系统回复
	EventSystem = 4
)

// Event model is the model for challenge changes
type Event struct {
	Eid         int64      `json:"eid" gorm:"column:id"`
	Cid         int64      `json:"cid" gorm:"column:cid"`
	AdminID     int64      `json:"adminid" gorm:"column:adminid"`
	Content     string     `json:"content" gorm:"column:content"`
	Attachments string     `json:"attachments" gorm:"column:attachments"`
	Event       int8       `json:"event" gorm:"column:event"`
	CTime       xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime       xtime.Time `json:"mtime" gorm:"column:mtime"`
	Admin       string     `json:"admin" gorm:"-"`
}

// TableName is used to identify table name in gorm
func (Event) TableName() string {
	return "workflow_event"
}
