package model

import "go-common/library/time"

// OrderResult order result.
type OrderResult struct {
	OrderNo string      `json:"order_no"`
	Status  int8        `json:"status"`
	Dialog  *ConfDialog `json:"dialog,omitempty"`
}

// ConfDialog .
type ConfDialog struct {
	ID          int64     `json:"-"`
	AppID       int64     `json:"app_id"`
	Platform    int64     `json:"platform"`
	StartTime   time.Time `json:"-"`
	EndTime     time.Time `json:"-"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Follow      bool      `json:"follow"`
	LeftButton  string    `json:"left_button"`
	LeftLink    string    `json:"left_link"`
	RightButton string    `json:"right_button"`
	RightLink   string    `json:"right_link"`
	Operator    string    `json:"-"`
	Stage       bool      `json:"stage"`
	Ctime       time.Time `json:"-"`
	Mtime       time.Time `json:"-"`
}
