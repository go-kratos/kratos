package model

import (
	"time"
)

const (
	// StateTypePending 处理中
	StateTypePending = int8(0)
	// StateTypeYes 有效
	StateTypeYes = int(1)
	// StateTypeNo 无效
	StateTypeNo = int(2)
	// StateDelete 删除
	StateDelete = int(9)
	// StatePublicReferee 移交众裁
	StatePublicReferee = int(10)
)

// Group appeal group
type Group struct {
	ID       int32     `gorm:"column:id" json:"id"`
	Oid      int64     `gorm:"column:oid" json:"oid"`
	State    int8      `gorm:"column:state" json:"state"`
	Business int8      `gorm:"column:business" json:"business"`
	Tid      int32     `gorm:"column:tid" json:"tid"`
	Count    int32     `gorm:"column:count" json:"count"`
	Handling int32     `gorm:"column:handling" json:"handling"`
	Note     string    `gorm:"column:note" json:"note"`
	CTime    time.Time `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time `gorm:"column:mtime" json:"mtime"`
	Lasttime time.Time `gorm:"column:lasttime" json:"lasttime"`
}

// TableName by Group
func (*Group) TableName() string {
	return "workflow_group"
}

// Group3 .
type Group3 struct {
	ID       int64     `gorm:"column:id" json:"id"`
	Oid      int64     `gorm:"column:oid" json:"oid"`
	State    int64     `gorm:"column:state" json:"state"`
	Business int64     `gorm:"column:business" json:"business"`
	Fid      int64     `gorm:"column:fid" json:"fid"`
	Rid      int64     `gorm:"column:rid" json:"rid"`
	Eid      int64     `gorm:"column:eid" json:"eid"`
	Score    int64     `gorm:"column:score" json:"score"`
	Tid      int64     `gorm:"column:tid" json:"tid"`
	Count    int64     `gorm:"column:count" json:"count"`
	Handling int64     `gorm:"column:handling" json:"handling"`
	Note     string    `gorm:"column:note" json:"note"`
	CTime    time.Time `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time `gorm:"column:mtime" json:"mtime"`
	Lasttime time.Time `gorm:"column:lasttime" json:"lasttime"`
}

// TableName .
func (g3 *Group3) TableName() string {
	return "workflow_group"
}

// DeleteGroupParams .
type DeleteGroupParams struct {
	Business int64 `json:"business" form:"business" validate:"required"`
	OID      int64 `json:"oid" form:"oid" validate:"required"`
	EID      int64 `json:"eid" form:"eid"`
}

// PublicRefereeGroupParams .
type PublicRefereeGroupParams struct {
	Business int8   `json:"business" form:"business" validate:"required"`
	Oid      string `json:"oid" form:"oid" validate:"required"`
	Eid      int64  `json:"eid" form:"eid"`
}
