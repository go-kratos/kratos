package model

import (
	"encoding/json"
)

// BMsg databus binlog message
type BMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// PMsg databus passport message
type PMsg struct {
	Action string `json:"action"`
	Table  string `json:"type"`
	Data   *Token `json:"data"`
	CTime  int64  `json:"ctime"`
}
