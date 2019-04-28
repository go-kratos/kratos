package model

import (
	xtime "go-common/library/time"
)

// Action 点赞或踩
const (
	ActionNormal int32 = 0 // 未踩赞
	ActionLike   int32 = 1 // 赞
	ActionHate   int32 = 2 // 踩
)

// Action reply action info
type Action struct {
	ID     int64      `json:"-"`
	RpID   int64      `json:"rpid"`
	Action int8       `json:"action"`
	Mid    int64      `json:"mid"`
	CTime  xtime.Time `json:"-"`
}
