package music

import (
	xtime "go-common/library/time"
)

// consts for workflow event
const (
	MusicDelete = -100
	MusicOpen   = 0
)

// Music model is the model for music
type Music struct {
	ID           int64      `json:"id" gorm:"column:id"`
	Sid          int64      `json:"sid" gorm:"column:sid"`
	Name         string     `json:"name" gorm:"column:name"`
	Musicians    string     `json:"musicians" gorm:"column:musicians"`
	Mid          int64      `json:"mid" gorm:"column:mid"`
	Tid          int64      `json:"tid" gorm:"-"`
	Rid          int64      `json:"rid" gorm:"-"`
	Pid          int64      `json:"pid" gorm:"-"`
	Cover        string     `json:"cover" gorm:"column:cover"`
	MaterialName string     `json:"material_name" gorm:"-"`
	CategoryName string     `json:"category_name" gorm:"-"`
	Stat         string     `json:"stat" gorm:"column:stat"`
	Categorys    string     `json:"categorys" gorm:"column:categorys"`
	Playurl      string     `json:"playurl" gorm:"column:playurl"`
	State        int8       `json:"state" gorm:"column:state"`
	Duration     int32      `json:"duration" gorm:"column:duration"`
	Filesize     int32      `json:"filesize" gorm:"column:filesize"`
	PubTime      xtime.Time `json:"pubtime" gorm:"column:pubtime"`
	SyncTime     xtime.Time `json:"synctime" gorm:"column:synctime"`
	CTime        xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime        xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Music) TableName() string {
	return "music"
}

// Param is used to parse user request
type Param struct {
	ID        int64      `form:"id" gorm:"column:id"`
	Sid       int64      `form:"sid" validate:"required"`
	Name      string     `form:"name" validate:"required"`
	Musicians string     `form:"musicians"`
	Mid       int64      `form:"mid" validate:"required"`
	Cover     string     `form:"cover" validate:"required"`
	Stat      string     `form:"stat" `
	Categorys string     `form:"categorys" `
	Playurl   string     `form:"playurl" `
	State     int8       `form:"state"`
	Duration  int32      `form:"duration" `
	Filesize  int32      `form:"filesize" `
	UID       int64      `form:"uid" `
	PubTime   xtime.Time `form:"pubtime"`
	SyncTime  xtime.Time `form:"synctime"`
}

// TableName is used to identify table name in gorm
func (Param) TableName() string {
	return "music"
}

// LogParam is used to parse user request
type LogParam struct {
	ID     int64  `json:"id"`
	UID    int64  `json:"uid"`
	UName  string `json:"uname"`
	Action string `json:"action"`
	Name   string `json:"name"`
}
