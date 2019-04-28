package model

import (
	"encoding/json"
)

// Action actions sent to job
const (
	ActionIdx     = "idx"
	ActionFlush   = "flush"
	ActAddDM      = "dm_add"       // 新增弹幕
	ActFlushDMSeg = "dm_seg_flush" // 刷新分段弹幕缓存
)

// Action job msg.
type Action struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// JobParam job param.
type JobParam struct {
	Type     int32 `json:"type,omitempty"`
	Pid      int64 `json:"pid,omitempty"`
	Oid      int64 `json:"oid,omitempty"`
	Cnt      int64 `json:"cnt,omitempty"`
	Num      int64 `json:"num,omitempty"`
	Duration int64 `json:"duration,omitempty"`
}

// Flush flush msg
type Flush struct {
	Type  int32 `json:"type,omitempty"`
	Oid   int64 `json:"oid,omitempty"`
	Force bool  `json:"force,omitempty"`
}

// ActionFlushDMSeg flush segment dm cache
type ActionFlushDMSeg struct {
	Type  int32 `json:"type"`
	Oid   int64 `json:"oid"`
	Force bool  `json:"force"`
	Page  *Page `json:"page"`
}
