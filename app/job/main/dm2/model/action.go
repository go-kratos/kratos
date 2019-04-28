package model

import (
	"encoding/json"
)

// action from DMAction-T
const (
	ActFlushDM    = "flush"        // 刷新弹幕缓存
	ActReportDel  = "report_del"   // 删除举报弹幕
	ActAddDM      = "dm_add"       // 新增弹幕
	ActFlushDMSeg = "dm_seg_flush" // 刷新分段弹幕缓存
)

// Action action message
type Action struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// Flush flush cache message
type Flush struct {
	Type  int32 `json:"type"`
	Oid   int64 `json:"oid"`
	Force bool  `json:"force"`
}

// FlushDMSeg flush segment dm cache
type FlushDMSeg struct {
	Type  int32 `json:"type"`
	Oid   int64 `json:"oid"`
	Force bool  `json:"force"`
	Page  *Page `json:"page"`
}
