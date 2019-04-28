package model

import (
	"encoding/json"
)

const (
	BinlogInsert = "insert"
	BinlogUpdate = "update"
	BinlogDelete = "delete"
)

// BinLog databus binlog message.
type BinLog struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
	MTS    int64
}

//RscMsg databus resource message
type RscMsg struct {
	Action string          `json:"action"`
	BizID  int64           `json:"business_id"`
	Raw    json.RawMessage `json:"raw"`
}
