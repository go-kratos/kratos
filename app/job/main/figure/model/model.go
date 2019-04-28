package model

import (
	"encoding/json"
	"time"
)

// Figure user figure model
type Figure struct {
	ID              int32     `json:"-"`
	Mid             int64     `json:"mid"`
	Score           int32     `json:"score"`
	LawfulScore     int32     `json:"lawful_score"`
	WideScore       int32     `json:"wide_score"`
	FriendlyScore   int32     `json:"friendly_score"`
	BountyScore     int32     `json:"bounty_score"`
	CreativityScore int32     `json:"creativity_score"`
	Ver             int32     `json:"-"`
	Ctime           time.Time `json:"-"`
	Mtime           time.Time `json:"-"`
}

// BinlogMsg dm binlog msg
type BinlogMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// ReplyInfo Reply info.
type ReplyInfo struct {
	Mid   int64 `json:"mid"`
	Oid   int64 `json:"oid"`
	Type  int8  `json:"type"`
	Ctime int64 `json:"ctime"`
}

// ReplyAction reply action.
type ReplyAction struct {
	Mid    int64 `json:"mid"`
	Action int8  `json:"action"`
	Oid    int64 `json:"oid"`
	Type   int8  `json:"type"`
	Ctime  int64 `json:"ctime"`
}

// DMAction danmaku report action.
type DMAction struct {
	Action string `json:"action"`
	Data   DMData `json:"data"`
}

// DMData danmaku report action.
type DMData struct {
	ID        int64 `json:"id"`
	Dmid      int64 `json:"dmid"`
	OwnerUID  int64 `json:"dm_owner_uid"`
	ReportUID int64 `json:"uid"`
}
