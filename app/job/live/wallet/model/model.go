package model

import "encoding/json"

// Message canal standary message
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// User canal base message
type User struct {
	Uid     int64 `json:"uid"`
	Gold    int64 `json:"gold"`
	IapGold int64 `json:"iap_gold"`
	Silver  int64 `json:"silver"`
}
