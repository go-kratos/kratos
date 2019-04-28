package model

import (
	"go-common/library/time"
)

// Content content def.
type Content struct {
	ID         int64
	Title      string
	Subtitle   string
	Desc       string
	Cover      string
	SeasonID   int
	CID        int
	EPID       int
	MenuID     int
	State      int
	Valid      int
	AuditTime  int
	PayStatus  int
	IsDeleted  int
	Ctime      time.Time
	Mtime      time.Time
	InjectTime time.Time
	Reason     string
}

// TVEpSeason represents the season table
type TVEpSeason struct {
	ID         int64
	OriginName string
	Title      string
	Alias      string
	Category   int8
	Desc       string
	Style      string
	Area       string
	PlayTime   time.Time
	Info       int8
	State      int8
	Check      int8
	TotalNum   int32
	Upinfo     string
	Staff      string
	Role       string
	Copyright  string
	Cover      string
	AuditTime  int
	IsDeleted  int8
	Ctime      time.Time
	Mtime      time.Time
	Valid      int8
	InjectTime time.Time
	Reason     string
}

// Audit def.
type Audit struct {
	IDList   []*IDList `json:"id_list"`
	OpType   string    `json:"optype"`
	Count    int       `json:"count"`
	AuditMsg string    `json:"audit_msg"`
}

// Cont : UGC content struct
type Cont struct {
	ID    int
	Title string
}
