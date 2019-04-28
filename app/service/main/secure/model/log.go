package model

import (
	"encoding/json"
	"go-common/library/time"
)

// ArgSecure arg secure.
type ArgSecure struct {
	Mid  int64  `json:"mid,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

// ArgFeedBack arg feedback.
type ArgFeedBack struct {
	Mid  int64  `json:"mid,omitempty"`
	UUID string `json:"uuid,omitempty"`
	IP   string `json:"ip,omitempty"`
	Type int8   `json:"type,omitempty"`
	Ts   int64  `json:"ts"`
}

// Log define user login log.
type Log struct {
	Mid        int64     `json:"mid,omitempty"`
	IP         uint32    `json:"loginip"`
	Location   string    `json:"location"`
	LocationID int64     `json:"location_id,omitempty"`
	Time       time.Time `json:"timestamp,omitempty"`
	Type       int8      `json:"type,omitempty"`
}

// Msg is user login status msg.
type Msg struct {
	Notify bool `json:"notify"`
	Log    *Log `json:"log"`
}

// Message is databus message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Expection is user expection record.
type Expection struct {
	IP       uint64    `json:"ip"`
	Time     time.Time `json:"time"`
	FeedBack int8      `json:"feedback"`
}

// PWDlog is user change password log.
type PWDlog struct {
	Mid int64 `json:"mid"`
}

// Record user login record.
type Record struct {
	LocID int64 `json:"locid"`
	Count int64 `json:"count"`
}

// Often is user often use ipaddr
type Often struct {
	Result bool `json:"result"`
}
