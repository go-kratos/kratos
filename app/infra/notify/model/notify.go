package model

import (
	"net/url"
	"time"
)

// Watcher define watcher object.
type Watcher struct {
	ID         int64
	Cluster    string
	Topic      string
	Group      string
	Offset     string
	Callback   string
	Callbacks  []*Callback
	Filter     bool
	Filters    []*Filter
	Concurrent int // concurrent goroutine for sub.
	Mtime      time.Time
}

// Pub define pub.
type Pub struct {
	ID        int64
	Cluster   string
	Topic     string
	Group     string
	Operation int8
	AppSecret string
}

// Callback define callback event
type Callback struct {
	URL      *NotifyURL
	Priority int8
	Finished bool
}

// NotifyURL callback url with parsed info
type NotifyURL struct {
	RawURL string
	Schema string
	Host   string
	Path   string
	Query  url.Values
}

// filter condition
const (
	ConditionEq  = 0
	ConditionPre = 1
)

// Filter define filter object.
type Filter struct {
	Field     string
	Condition int8 // 0 :eq 1:neq
	Value     string
}

// Message define canal message.
type Message struct {
	Table  string `json:"table,omitempty"`
	Action string `json:"action,omitempty"`
}

// ArgPub pub arg.
type ArgPub struct {
	AppKey    string `form:"appkey" validate:"min=1"`
	AppSecret string `form:"appsecret" validate:"min=1"`
	Group     string `form:"group" validate:"min=1"`
	Topic     string `form:"topic" validate:"min=1"`
	Key       string `form:"key" validate:"min=1"`
	Msg       string `form:"msg" validate:"min=1"`
}

// FailBackup fail backup msg.
type FailBackup struct {
	ID      int64
	Cluster string
	Topic   string
	Group   string
	Offset  int64
	Msg     string
	Index   int64
}

// Notify callback schema
const (
	LiverpcSchema = "liverpc"
	HTTPSchema    = "http"
)
