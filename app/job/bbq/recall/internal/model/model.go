package model

import (
	xtime "go-common/library/time"
)

// RecommendVideoState 所有进推荐池的新发视频
var RecommendVideoState = []string{"5", "4", "3", "1", "0"}

// HDFSResult .
type HDFSResult struct {
	Code   int16    `json:"code"`
	Msg    string   `json:"msg"`
	Result []string `json:"result"`
}

// Video .
type Video struct {
	SVID     int64
	Title    string
	Content  string
	MID      int64
	AVID     int64
	CID      int64
	PubTime  xtime.Time
	CTime    xtime.Time
	MTime    xtime.Time
	Duration int32
	State    int16
	TID      int32
	SubTID   int32
}
