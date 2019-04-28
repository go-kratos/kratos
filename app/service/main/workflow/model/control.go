package model

import (
	"time"
)

const (
	// ControlTypeInput 文本类型控件
	ControlTypeInput = "input"
	// ControlTypeTextarea 多行文本类型控件
	ControlTypeTextarea = "textarea"
	// ControlTypeLink 链接类型控件
	ControlTypeLink = "link"
	// ControlTypeSelector 选择类型控件
	ControlTypeSelector = "selector"
	// ControlTypeFile 文件类型控件
	ControlTypeFile = "file"
	// ControlPageSize .
	ControlPageSize = int(1000)
)

// Control will describe how the tag be acted
type Control struct {
	Cid         int32     `gorm:"-" json:"-"`
	Tid         int32     `gorm:"column:tid" json:"tid"`
	Weight      int32     `gorm:"-" json:"-"`
	Name        string    `gorm:"column:name" json:"name"`
	Title       string    `gorm:"column:title" json:"title"`
	Component   string    `gorm:"column:component" json:"component"`
	Placeholder string    `gorm:"column:placeholder" json:"placeholder"`
	Required    bool      `gorm:"column:required" json:"required"`
	CTime       time.Time `gorm:"-" json:"-"`
	MTime       time.Time `gorm:"-" json:"-"`
}

// TableName by control
func (*Control) TableName() string {
	return "workflow_tag_control"
}

// Control3 .
type Control3 struct {
	TID         int64  `json:"tid"`
	BID         int64  `json:"bid"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Component   string `json:"component"`
	Placeholder string `json:"placeholder"`
	Required    int64  `json:"required"`
}

// ResponseControl3 .
type ResponseControl3 struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int32       `json:"ttl"`
	Data    []*Control3 `json:"data"`
}
