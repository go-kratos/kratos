package model

import "encoding/json"

// MsgCanal canal message struct
type MsgCanal struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}
