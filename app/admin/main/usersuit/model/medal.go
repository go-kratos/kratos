package model

import (
	xtime "go-common/library/time"
)

const (
	// OwnerActivated medal_owner is_activated=1 .
	OwnerActivated = int8(1)
	// OwnerNotActivated medal_owner is_activated=0 .
	OwnerNotActivated = int8(0)
	// MaxCount medal batch add max.
	MaxCount = 2000
	// MedalSourceTypeAdmin medal source type admin.
	MedalSourceTypeAdmin = int8(1)
)

// Medal medal info .
type Medal struct {
	ID          int64      `form:"id" json:"id"`
	GID         int64      `form:"gid" validate:"required" json:"gid"`
	Name        string     `form:"name" validate:"required" json:"name"`
	Description string     `form:"description" validate:"required" json:"description"`
	Image       string     `form:"image" validate:"required" json:"image"`
	ImageSmall  string     `form:"image_small" validate:"required" json:"image_small"`
	Condition   string     `form:"condition" validate:"required" json:"condition"`
	Level       int8       `form:"level" validate:"min=1,max=3" json:"level"`
	LevelRank   string     `form:"level_rank" validate:"required" json:"level_rank"`
	Sort        int        `form:"sort" validate:"required" json:"sort"`
	IsOnline    int        `form:"is_online" json:"is_online"`
	CTime       xtime.Time `json:"ctime,omitempty"`
	MTime       xtime.Time `json:"mtime,omitempty"`
}

// MedalGroup nameplate group .
type MedalGroup struct {
	ID       int64      `form:"id" json:"id"`
	PID      int64      `form:"pid" json:"pid"`
	Rank     int8       `form:"rank" validate:"required" json:"rank"`
	IsOnline int8       `form:"is_online"  json:"is_online"`
	Name     string     `form:"name" validate:"required" json:"name"`
	PName    string     `form:"pname" json:"pname,omitempty"`
	CTime    xtime.Time `json:"ctime,omitempty"`
	MTime    xtime.Time `json:"mtime,omitempty"`
}

// MedalOwner nameplate owner .
type MedalOwner struct {
	ID          int64      `json:"id"`
	MID         int64      `json:"mid"`
	NID         int64      `json:"nid"`
	IsActivated int8       `json:"is_activated"`
	IsDel       int8       `json:"is_del"`
	CTime       xtime.Time `json:"ctime"`
	MTime       xtime.Time `json:"mtime"`
}

// MedalInfo struct.
type MedalInfo struct {
	*Medal
	GroupName       string `json:"group_name"`
	ParentGroupName string `json:"parent_group_name"`
}

// MedalMemberMID struct.
type MedalMemberMID struct {
	ID          int64  `json:"id"`
	NID         int64  `json:"nid"`
	MedalName   string `json:"medal_name"`
	IsActivated int8   `json:"is_activated"`
	IsDel       int8   `json:"is_del"`
}

// MedalMemberAddList struct.
type MedalMemberAddList struct {
	ID        int64  `json:"id"`
	MedalName string `json:"medal_name"`
}

// MedalOperLog struct.
type MedalOperLog struct {
	OID        int64      `json:"oper_id"`
	Action     string     `json:"action"`
	CTime      xtime.Time `json:"ctime"`
	MTime      xtime.Time `json:"mtime"`
	OperName   string     `json:"oper_name"`
	MID        int64      `json:"mid"`
	MedalID    int64      `json:"medal_id"`
	SourceType int8       `json:"source_type"`
}
