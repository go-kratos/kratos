package model

import (
	"encoding/json"
)

// BMsg databus binlog message.
type BMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
	MTS    int64
}

// PMsg Push msg
type PMsg struct {
	Action string      `json:"action"`
	Table  string      `json:"table"`
	Data   *AsoAccount `json:"data"`
	Flag   int         `json:"flag"`
	MTS    int64
}
