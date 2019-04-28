package dm

import "encoding/json"

// Canal canal message struct
type Canal struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Subject for table dm_subject_[0-9]+
type Subject struct {
	ID    int64 `json:"id"`
	Type  int64 `json:"type"`
	AID   int64 `json:"pid"`
	CID   int64 `json:"oid"`
	Count int64 `json:"count"`
}

// Count dm count
type Count struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Count     int64  `json:"count"`
	Timestamp int64  `json:"timestamp"`
}
