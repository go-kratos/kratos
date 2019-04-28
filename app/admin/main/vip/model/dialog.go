package model

import "go-common/library/time"

// ConfDialog .
type ConfDialog struct {
	ID          int64     `gorm:"column:id" json:"id" form:"id"`
	AppID       int64     `gorm:"column:app_id" json:"app_id" form:"app_id"`
	Platform    int64     `gorm:"column:platform" json:"platform" form:"platform"`
	StartTime   time.Time `gorm:"column:start_time" json:"start_time" form:"start_time"`
	EndTime     time.Time `gorm:"column:end_time" json:"end_time" form:"end_time" default:"32503651200"` //3000-01-01 00:00:00
	Title       string    `gorm:"column:title" json:"title" form:"title" validate:"required"`
	Content     string    `gorm:"column:content" json:"content" form:"content" validate:"required"`
	Follow      bool      `gorm:"column:follow" json:"follow" form:"follow"`
	LeftButton  string    `gorm:"column:left_button" json:"left_button" form:"left_button"`
	LeftLink    string    `gorm:"column:left_link" json:"left_link" form:"left_link"`
	RightButton string    `gorm:"column:right_button" json:"right_button" form:"right_button" validate:"required"`
	RightLink   string    `gorm:"column:right_link" json:"right_link" form:"right_link"`
	Operator    string    `gorm:"column:operator" json:"operator"`
	Stage       bool      `gorm:"column:stage" json:"stage" form:"stage" default:"true"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime" form:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime" form:"mtime"`
}

// TableName for grom.
func (c *ConfDialog) TableName() string {
	return "vip_conf_dialog"
}

// ConfDialogList admin list model.
type ConfDialogList struct {
	*ConfDialog
	Status string `json:"status"`
}
