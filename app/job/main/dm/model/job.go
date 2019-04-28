package model

import (
	"encoding/json"
)

// All const variable used in job
const (
	// binlog type
	SyncInsert = "insert"
	SyncUpdate = "update"
	SyncDelete = "delete"
)

// BinlogMsg dm binlog msg
type BinlogMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}
