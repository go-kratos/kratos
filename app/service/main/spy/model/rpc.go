package model

import (
	"go-common/library/time"
)

// ArgReBuild rebuild args
type ArgReBuild struct {
	Mid    int64
	Reason string
}

// ArgReset is.
type ArgReset struct {
	Mid        int64
	ReLiveTime bool
	EventScore bool
	BaseScore  bool
	Operator   string
}

// ArgUserScore rpc arg for getting user score.
type ArgUserScore struct {
	Mid int64
	IP  string
}

// ArgHandleEvent rpc arg for handling spy event.
type ArgHandleEvent struct {
	Time      time.Time
	IP        string      `json:"ip"`
	Service   string      `json:"service"`
	Event     string      `json:"event"`
	ActiveMid int64       `json:"active_mid"`
	TargetMid int64       `json:"target_mid"`
	TargetID  int64       `json:"target_id"`
	Args      interface{} `json:"args"`
	Result    string      `json:"result"`
	Effect    string      `json:"effect"`
	RiskLevel int8        `json:"risk_level"`
}

// ArgUser rpc arg for getting user info.
type ArgUser struct {
	Mid int64
	IP  string
}

// ArgStat rpc arg for getting user stat.
type ArgStat struct {
	ID  int64
	Mid int64
}
