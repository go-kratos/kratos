package model

import (
	"encoding/json"

	"go-common/library/time"
)

// Message binlog databus msg.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// UserCoin dede_member user coin.
type UserCoin struct {
	Mid   int64     `json:"mid"`
	Money float32   `money:"money"`
	Mtime time.Time `json:"mtime"`
}

// DatabusCoin databus coin msg.
type DatabusCoin struct {
	Mid   int64   `json:"mid"`
	Money float32 `money:"money"`
	Mtime string  `json:"modify_time"`
}

// CoinSettle coin settle.
type CoinSettle struct {
	ITime     time.Time `json:"itime"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	Aid       int64     `json:"aid"`
	AvType    int64     `json:"avtype"`
	CoinCount int64     `json:"coin_count"`
	ExpTotal  int64     `json:"exp_total"`
	ExpSub    int64     `json:"exp_sub"`
	State     int       `json:"state"`
	Describe  string    `json:"describe"`
}

// CoinSettlePeriod coin settle conf.
type CoinSettlePeriod struct {
	ID        int64     `json:"id"`
	FromYear  int       `json:"from_year"`
	FromMonth int       `json:"from_month"`
	FromDay   int       `json:"from_day"`
	ToYear    int       `json:"to_year"`
	ToMonth   int       `json:"to_month"`
	ToDay     int       `json:"to_day"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// LoginLog user login log.
type LoginLog struct {
	Mid       int64  `json:"mid,omitempty"`
	IP        string `json:"ip,omitempty"`
	CTime     string `json:"ctime"`
	Action    string `json:"action"`
	Business  int    `json:"business"`
	Type      int    `json:"type"`
	RawData   string
	Timestamp int64
}

// AddExp databus add exp arg.
type AddExp struct {
	Event string `json:"event,omitempty"`
	Mid   int64  `json:"mid,omitempty"`
	IP    string `json:"ip,omitempty"`
	Ts    int64  `json:"ts,omitempty"`
}
