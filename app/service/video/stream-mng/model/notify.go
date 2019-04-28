package model

import (
	"encoding/json"
)

// StreamingNotifyParam 开/关播回调请求参数
type StreamingNotifyParam struct {
	StreamName string      `json:"stream_name,omitempty"`
	Key        string      `json:"key,omitempty"`
	SRC        string      `json:"src,omitempty"`
	Type       json.Number `json:"type,omitempty"`
	TS         json.Number `json:"ts,omitempty"`
	SUID       string      `json:"suid,omitempty"`
	Sign       string      `json:"sign,omitempty"`
}
