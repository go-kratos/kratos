package net

import (
	"time"
)

//..
const (
	TableTransition = "net_transition"
	//TriggerManual 人工触发，一旦enable，执行由人审核提交
	TriggerManual = int8(1)
	//TriggerAuto 自动触发，一旦enable就可执行
	TriggerAuto = int8(2)
	//TriggerMsg 消息触发，一旦enable，执行为接收到指定消息
	TriggerMsg = int8(3)
)

//TriggerDesc 变迁触发类型描述
var TriggerDesc = map[int8]string{
	TriggerAuto:   "直序",
	TriggerManual: "人工",
	TriggerMsg:    "消息",
}

//Transition 变迁
type Transition struct {
	ID          int64     `gorm:"primary_key" json:"id" form:"id" validate:"omitempty,gt=0"`
	NetID       int64     `gorm:"column:net_id" json:"net_id" form:"net_id" validate:"omitempty,gt=0"`
	Trigger     int8      `gorm:"column:trigger" json:"trigger" default:"1" form:"trigger"`
	Limit       int64     `gorm:"column:limit" json:"limit" form:"limit"`
	Name        string    `gorm:"column:name" json:"name" form:"name" validate:"required,max=32"`
	ChName      string    `gorm:"column:ch_name" json:"ch_name" form:"ch_name" validate:"required,max=16"`
	Description string    `gorm:"column:description" json:"description" form:"description" validate:"max=60"`
	UID         int64     `gorm:"column:uid" json:"uid"`
	DisableTime time.Time `gorm:"column:disable_time" json:"disable_time"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
}

//TableName .
func (t *Transition) TableName() string {
	return TableTransition
}

//IsAvailable .
func (t *Transition) IsAvailable() bool {
	return t.DisableTime.IsZero()
}
