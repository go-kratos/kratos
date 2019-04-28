package model

import (
	xtime "go-common/library/time"
)

// Group model is the group view for several challanges
type Group struct {
	ID           int64  `json:"id" gorm:"column:id"`
	Oid          int64  `json:"oid" gorm:"column:oid"`
	OidStr       string `json:"oid_str" gorm:"-"`
	Business     int8   `json:"business" gorm:"column:business"`
	Fid          int64  `json:"fid" gorm:"column:fid"`
	Rid          int8   `json:"rid" gorm:"column:rid"`
	Eid          int64  `json:"eid" gorm:"eid"`
	EidStr       string `json:"eid_str" gorm:"-"`
	State        int8   `json:"state" gorm:"column:state"`
	Tid          int64  `json:"tid" gorm:"column:tid"`
	FirstUserTid int64  `json:"first_user_tid" gorm:"-"`
	Note         string `json:"note" gorm:"column:note"`
	Score        int64  `json:"score" gorm:"column:score"`

	// Stat fields
	// this is a workround solution for calcuating appeals
	Count    int32 `json:"count" gorm:"column:count"`
	Handling int32 `json:"handling" gorm:"column:handling"`

	CTime    xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime    xtime.Time `json:"mtime" gorm:"column:mtime"`
	LastTime xtime.Time `json:"last_time" gorm:"column:lasttime"`

	LastLog        string    `json:"last_log" gorm:"-"`
	BusinessObject *Business `json:"business_object,omitempty" gorm:"-"`

	// Tags related to Group
	Tag           string        `json:"tag" gorm:"-"`
	Round         int8          `json:"round" gorm:"-"`
	ChallengeTags ChallTagSlice `json:"challenge_tags" gorm:"-"`

	Meta     interface{} `json:"meta" gorm:"-"`
	MetaData interface{} `json:"meta_data" gorm:""`
	TypeID   int64       `json:"type_id" gorm:"-"`

	LastProducer *Account `json:"last_producer" gorm:"-"`
	Defendant    *Account `json:"defendant" gorm:"-"`
}

// GroupListPage is the model for group list result
type GroupListPage struct {
	Items []*Group `json:"items"`
	Page  *Page    `json:"page"`
}

// GroupPendingCount .
type GroupPendingCount struct {
	Total int `json:"total"`
}

// GroupMeta .
type GroupMeta struct {
	Archive  *Archive    `json:"archive"`
	Object   *Business   `json:"object"`
	External interface{} `json:"external"`
}

// TableName is used to identify group table name in gorm
func (Group) TableName() string {
	return "workflow_group"
}
