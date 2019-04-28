package model

import (
	"encoding/json"
)

// CanalBinLog is
type CanalBinLog struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	Old    json.RawMessage `json:"old"`
	New    json.RawMessage `json:"new"`
}

// MidBinLog is
type MidBinLog struct {
	Mid int64 `json:"mid"`
}
