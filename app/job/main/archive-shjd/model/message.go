package model

import "encoding/json"

// Message databus
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Notify is
type Notify struct {
	Table  string   `json:"table"`
	Action string   `json:"action"`
	Nw     *Archive `json:"new"`
	Old    *Archive `json:"old"`
}

// AccountNotify is
type AccountNotify struct {
	Mid    int64  `json:"mid"`
	Action string `json:"action"`
}
