package model

import (
	"encoding/json"
)

// Message is
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
}

// ArcMsg is
var ArcMsg struct {
	Aid int64 `json:"aid"`
}
