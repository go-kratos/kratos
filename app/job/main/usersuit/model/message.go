package model

import "encoding/json"

// const .
const (
	TimeFormatSec = "2006-01-02 15:04:05"
)

// Message is simple message struct info.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// VipInfoMessage .
type VipInfoMessage struct {
	Mid            int64  `json:"mid"`
	VipType        int8   `json:"vip_type"`
	VipPayType     int8   `json:"vip_pay_type"`
	VipStatus      int8   `json:"vip_status"`
	VipOverdueTime string `json:"vip_overdue_time"`
}
