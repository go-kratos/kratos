package model

import (
	xtime "go-common/library/time"
)

const (
	// DefaultButton .
	DefaultButton = 6
)

// TagType .
type TagType struct {
	ID    int64           `json:"id" gorm:"primary_key" form:"id"`
	Bid   int64           `json:"bid" gorm:"column:bid" form:"bid" validate:"required"`
	Name  string          `json:"name" gorm:"column:name" form:"name" validate:"required"`
	State int64           `json:"state" gorm:"column:state" form:"state"`
	Ctime xtime.Time      `json:"ctime" gorm:"column:ctime"`
	Mtime xtime.Time      `json:"mtime" gorm:"column:mtime"`
	Rids  []int64         `json:"rids" gorm:"-" form:"rids,split" validate:"required"`
	Roles []*BusinessRole `json:"roles" gorm:"-"`
}

// TagTypeList .
type TagTypeList struct {
	BID int64 `json:"bid" form:"bid" validate:"required"`
}

// TagTypeDel .
type TagTypeDel struct {
	ID int64 `json:"id" form:"id" validate:"required"`
}

// TableName .
func (tt *TagType) TableName() string {
	return "manager_tag_type"
}

// TagTypeRole .
type TagTypeRole struct {
	ID    int64      `json:"id" gorm:"primary_key"`
	Tid   int64      `json:"tid" gorm:"column:tid"`
	Rid   int64      `json:"rid" gorm:"column:rid"`
	Ctime xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (ttr *TagTypeRole) TableName() string {
	return "manager_tag_type_role"
}

// Tag .
type Tag struct {
	ID          int64      `json:"id" gorm:"primary_key" form:"id"`
	Bid         int64      `json:"bid" gorm:"column:bid" form:"bid" validate:"required"`
	Tid         int64      `json:"tid" gorm:"column:tid" form:"tid" validate:"required"`
	TagID       int64      `json:"tag_id" gorm:"column:tag_id"`
	TName       string     `json:"tname" gorm:"-"`
	Rid         int64      `json:"rid" gorm:"column:rid" form:"rid" validate:"required"`
	RName       string     `json:"rname" gorm:"-"`
	Name        string     `json:"name" gorm:"column:name" form:"name" validate:"required"`
	Weight      int64      `json:"weight" gorm:"column:weight" form:"weight" validate:"required"`
	State       int64      `json:"state" gorm:"column:state" form:"state" default:"1"`
	UID         int64      `json:"uid" gorm:"column:uid" form:"uid"`
	UName       string     `json:"uname" gorm:"-"`
	Description string     `json:"description" gorm:"column:description" form:"description"`
	Ctime       xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime       xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (t *Tag) TableName() string {
	return "manager_tag"
}

// TagBusinessAttr .
type TagBusinessAttr struct {
	ID     int64      `json:"id" gorm:"primary_key"`
	Bid    int64      `json:"bid" gorm:"column:bid" form:"bid" validate:"required"`
	Button int64      `json:"button" form:"button" gorm:"column:button" default:"6"`
	Extra  int64      `json:"extra" gorm:"column:extra"`
	Ctime  xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime  xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName .
func (t *TagBusinessAttr) TableName() string {
	return "manager_tag_business_attr"
}

// TagControl .
type TagControl struct {
	ID          int64      `json:"id" gorm:"primary_key" form:"id"`
	Tid         int64      `json:"tid" gorm:"column:tid" form:"tid" validate:"required"`
	Name        string     `json:"name" gorm:"column:name" form:"name" validate:"required"`
	Title       string     `json:"title" gorm:"column:title" form:"title" validate:"required"`
	Weight      int64      `json:"weight" gorm:"column:weight" form:"weight" validate:"required"`
	Component   string     `json:"component" gorm:"column:component" form:"component" validate:"required"`
	Placeholder string     `json:"placeholder" gorm:"column:placeholder" form:"placeholder" validate:"required"`
	Required    int64      `json:"required" gorm:"column:required" form:"required" default:"0"`
	Ctime       xtime.Time `json:"ctime" gorm:"-"`
	Mtime       xtime.Time `json:"mtime" gorm:"-"`
	BID         int64      `json:"bid" gorm:"column:bid" form:"bid"`
}

// TagControlParam .
type TagControlParam struct {
	BID int64 `json:"bid" form:"bid"`
	TID int64 `json:"tid" form:"tid"`
}

// TableName .
func (tc *TagControl) TableName() string {
	return "manager_tag_control"
}

// BatchUpdateState .
type BatchUpdateState struct {
	IDs   []int64 `json:"ids" form:"ids,split" validate:"required"`
	State int64   `json:"state" form:"state"`
}

// SearchTagParams .
type SearchTagParams struct {
	Bid     int64  `form:"bid" default:"-1"`
	KeyWord string `form:"keyword"`
	Tid     int64  `form:"tid" default:"-1"`
	Rid     int64  `form:"rid" default:"-1"`
	State   int64  `form:"state" default:"-1"`
	UID     int64  `form:"uid" `
	UName   string `form:"uname"`
	Order   string `form:"order"`
	Sort    string `form:"sort" default:"desc"`
	PS      int64  `form:"ps" default:"20"`
	PN      int64  `form:"pn" default:"1"`
}
