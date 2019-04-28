package model

import (
	"encoding/json"
)

// Message define binlog databus message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Stat is
type Stat struct {
	Mid       int64 `json:"mid,omitempty"`
	Following int64 `json:"following"`
	Whisper   int64 `json:"whisper"`
	Black     int64 `json:"black"`
	Follower  int64 `json:"follower"`
}

// Relation is
type Relation struct {
	Mid       int64  `json:"mid,omitempty"`
	Fid       int64  `json:"fid,omitempty"`
	Attribute uint32 `json:"attribute"`
	Status    int    `json:"status"`
	MTime     string `json:"mtime"`
	CTime     string `json:"ctime"`
}
