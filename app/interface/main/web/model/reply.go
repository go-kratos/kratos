package model

import (
	"encoding/json"

	rplmdl "go-common/app/interface/main/reply/model/reply"
)

// ReplyHot reply hot
type ReplyHot struct {
	Page    json.RawMessage `json:"page"`
	Replies []*rplmdl.Reply `json:"replies"`
}
