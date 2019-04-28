package task

import (
	"time"
)

// Role .
type Role struct {
	ID   int64  `json:"id"`
	BID  int64  `json:"bid"`
	RID  int64  `json:"rid"`
	Type int8   `json:"type"`
	Name string `json:"name"`
}

// WatchItem 审核员状态
type WatchItem struct {
	UID          int64     `json:"uid"`
	Uname        string    `json:"uname"`
	IsOnLine     bool      `json:"is_online"`
	LastOn       string    `json:"laston"`
	LastOff      string    `json:"lastoff"`
	CompleteRate string    `json:"complete_rate"`
	PassRate     string    `json:"pass_rate"`
	Count        int64     `json:"count"`
	AvgUT        string    `json:"avgut"`
	BizID        int64     `json:"business_id"`
	FlowID       int64     `json:"flow_id"`
	Role         int8      `json:"role"`
	Mtime        time.Time `json:"-"`
	State        int8      `json:"-"`
}
