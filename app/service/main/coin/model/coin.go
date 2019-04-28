package model

import (
	"go-common/app/service/main/archive/api"
	"go-common/library/time"
)

var (
	// TypeNone none type
	TypeNone = 0
	// TypeSend send type
	TypeSend = 1
	// TypeReceive receive type
	TypeReceive = 2
	// ReportType 上报business
	ReportType = 21
)

// Record coin added record.
type Record struct {
	Aid       int64
	Mid       int64
	Up        int64
	Timestamp int64
	Multiply  int64
	AvType    int64
	Business  string
	IP        uint32
	IPV6      string
}

// AddedArchive archive info.
type AddedArchive struct {
	*api.Arc
	IP    string `json:"ip"`
	Time  int64  `json:"time"`
	Coins int64  `json:"coins"`
}

// DataBus databus msg.
type DataBus struct {
	Mid      int64     `json:"mid"`      // user id
	Avid     int64     `json:"avid"`     // archive id
	AvType   int8      `json:"avtp"`     // archive type
	UpID     int64     `json:"upper_id"` // upper id
	Multiply int64     `json:"multiply"` // multiply
	Time     time.Time `json:"time"`     // archive pub date
	IP       string    `json:"ip"`       // userip
	TypeID   int16     `json:"rid"`      // zone id
	Tags     string    `json:"tags"`     // tag ids
	Ctime    int64     `json:"ctime"`    // add coin time
	MsgID    string    `json:"msg_id"`   // unique msg id
}

// CoinSettle .
type CoinSettle struct {
	ID        int64     `json:"id"`
	Mid       int64     `json:"mid"`
	Aid       int64     `json:"aid"`
	CoinCount int64     `json:"coin_count"`
	ExpTotal  int64     `json:"exp_total"`
	ExpSub    int64     `json:"exp_sub"`
	State     int       `json:"state"`
	Describe  string    `json:"describe"`
	ITime     time.Time `json:"itime"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// CoinSettlePeriod .
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

// AddCoins .
type AddCoins struct {
	Count int64 `json:"count"`
}

// Log coins log.
type Log struct {
	From      float64 `json:"from"`
	To        float64 `json:"to"`
	IP        string  `json:"ip"`
	Desc      string  `json:"desc"`
	TimeStamp int64   `json:"timestamp"`
}

// LogExp log exp
type LogExp struct {
	List  []*Exp `json:"list"`
	Count int    `json:"count"`
}

// Exp exp
type Exp struct {
	Delta  float64 `json:"delta"`
	Time   string  `json:"time"`
	Reason string  `json:"reason"`
}

// List define coin added list.
type List struct {
	Aid      int64  `json:"aid,omitempty"`
	Multiply int64  `json:"multiply,omitempty"`
	Ts       int64  `json:"ts,omitempty"`
	IP       uint32 `json:"ip,omitempty"`
}

// Business .
type Business struct {
	ID                 int64
	Name               string
	AddCoinReason      string
	AddCoinUpperReason string
	AddExpReason       string
}
