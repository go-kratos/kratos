package model

import (
	"go-common/library/time"
)

// EventMessage event from bigdata def.
type EventMessage struct {
	Time      time.Time   `json:"time"`
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
