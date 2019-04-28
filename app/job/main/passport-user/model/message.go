package model

import "encoding/json"

// BMsg databus binlog message.
type BMsg struct {
	Action    string          `json:"action"`
	Table     string          `json:"table"`
	New       json.RawMessage `json:"new"`
	Old       json.RawMessage `json:"old"`
	Timestamp int64           `json:"timestamp"`
}
