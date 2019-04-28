package model

import (
	xtime "go-common/library/time"
)

// State code const .
const (
	InvalidateState int64 = 0
	ValidateState   int64 = 1
	AllState        int64 = -1
	AllType         int64 = -1
	CategoryCode    int64 = 0
	SecondeCode     int64 = 1
	ExtensionCode   int64 = 2
)

// Reason .
type Reason struct {
	ID           int64      `json:"id" gorm:"primary_key" form:"id"`
	BusinessID   int64      `json:"bid" gorm:"column:bid" form:"bid"`
	RoleID       int64      `json:"role_id" gorm:"column:rid" form:"role_id" validate:"required"`
	RoleName     string     `json:"role_name" gorm:"-"`
	CategoryID   int64      `json:"cate_id" gorm:"column:cid" form:"cate_id" validate:"required"`
	CategoryName string     `json:"category_name" gorm:"-"`
	SecondID     int64      `json:"sec_id" gorm:"column:sid" form:"sec_id" validate:"required"`
	SecondName   string     `json:"sec_name" gorm:"-"`
	State        int64      `json:"state" gorm:"column:state" form:"state" default:"1"`
	Common       int64      `json:"common" gorm:"column:common" form:"common" default:"0"`
	UID          int64      `json:"uid" gorm:"column:uid" form:"uid"`
	UName        string     `json:"uname" gorm:"-"`
	Description  string     `json:"description" gorm:"column:description" form:"description" validate:"required"`
	Weight       int64      `json:"weight" gorm:"column:weight" form:"weight" default:"1" validate:"required"`
	Flag         int64      `json:"flag" gorm:"column:flag" form:"flag" default:"0"`
	LinkID       int64      `json:"link_id" gorm:"column:lid" form:"link_id" default:"0"`
	TypeID       int64      `json:"type_id" gorm:"column:type_id" form:"type_id" default:"0"`
	TagID        int64      `json:"tag_id" gorm:"column:tid" form:"tag_id" default:"0"`
	Ctime        xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime        xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// DropList .
type DropList struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Child []*DropList `json:"child"`
}

// CateSecExt .
type CateSecExt struct {
	ID         int64      `json:"id" gorm:"primary_key" form:"id"`
	BusinessID int64      `json:"bid" gorm:"column:bid" form:"bid"`
	Name       string     `json:"name" gorm:"column:name" form:"name"`
	Type       int64      `json:"type" gorm:"column:type" form:"type"`
	State      int64      `json:"state" gorm:"column:state" form:"state" default:"-1"`
	Ctime      xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime      xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// Association .
type Association struct {
	ID           int64         `json:"id" gorm:"primary_key" form:"id"`
	BusinessID   int64         `json:"bid" gorm:"column:bid" form:"bid"`
	RoleID       int64         `json:"role_id" gorm:"column:rid" form:"role_id"`
	RoleName     string        `json:"role_name" gorm:"-"`
	CategoryID   int64         `json:"cate_id" gorm:"column:cid" form:"cate_id"`
	CategoryName string        `json:"cate_name" gorm:"-"`
	SecondIDs    string        `json:"second_ids" gorm:"column:sids" form:"sec_ids"`
	State        int64         `json:"state" gorn:"column:state" form:"state" default:"1"`
	Ctime        xtime.Time    `json:"ctime" gorm:"column:ctime"`
	Mtime        xtime.Time    `json:"mtime" gorm:"column:mtime"`
	Child        []*CateSecExt `json:"child" gorm:"-"`
}

// BatchUpdateReasonState .
type BatchUpdateReasonState struct {
	IDs   []int64 `json:"ids" form:"ids,split" validate:"required"`
	State int64   `json:"state" form:"state"`
}

// SearchReasonParams .
type SearchReasonParams struct {
	BusinessID int64  `form:"bid" default:"-1"`
	KeyWord    string `json:"keyword" form:"keyword"`
	RoleID     int64  `json:"role_id" form:"role_id" default:"0"`
	CategoryID int64  `json:"cate_id" form:"cate_id" default:"0"`
	SecondID   int64  `json:"sec_id" form:"sec_id" default:"0"`
	State      int64  `json:"state" form:"state" default:"-1"`
	UID        int64  `json:"uid"`
	UName      string `json:"uname" form:"uname"`
	Order      string `json:"order" form:"order" default:"ctime"`
	Sort       string `json:"sort" form:"sort" default:"desc"`
	PS         int64  `form:"ps" default:"20"`
	PN         int64  `form:"pn" default:"1"`
}

// BusinessAttr .
type BusinessAttr struct {
	BID int64 `json:"bid" form:"bid" validate:"required"`
}
