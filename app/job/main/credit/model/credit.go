package model

import (
	"encoding/json"

	xtime "go-common/library/time"
)

// Message is simple message struct info.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// CreditInfo is simple creditInfo struct info.
type CreditInfo struct {
	Mid        int64      `json:"mid"`
	Status     int64      `json:"status"`
	PunishType int64      `json:"punishType"`
	PunishEnd  xtime.Time `json:"punishEnd"`
	CTime      xtime.Time `json:"-"`
	MTime      xtime.Time `json:"-"`
}

// AutoCaseConf struct
type AutoCaseConf struct {
	Reasons     map[int8]struct{}
	ReasonStr   string
	ReportScore int
	Likes       int
}
