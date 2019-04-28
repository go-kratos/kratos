package model

import (
	"encoding/json"
)

// RecoveryInfo recoveryInfo.
type RecoveryInfo struct {
	Rid    int64 `json:"rid" params:"rid;Required"` // 用户rid列表
	Status int64 `json:"status"`
}

// Message is databus message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// CommonResq CommonResq
type CommonResq struct {
	Code    int64  `json:"code"`
	TS      int64  `json:"ts"`
	Message string `json:"message"`
}
