package model

import "encoding/json"

const (
	ActUpdate = "update"
	ActInsert = "insert"
	ActDelete = "delete"
)

// Message canal binlog message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// LikeMsg msg
type LikeMsg struct {
	BusinessID    int64 `json:"business_id"`
	MessageID     int64 `json:"message_id"`
	LikesCount    int64 `json:"likes_count"`
	DislikesCount int64 `json:"dislikes_count"`
}
