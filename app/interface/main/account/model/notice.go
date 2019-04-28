package model

import (
	"go-common/library/time"
)

const (
	// NoticeStatusNotify notify
	NoticeStatusNotify = 1
	// NoticeStatusNotNotify not notify
	NoticeStatusNotNotify = 0
	// NoticeTypeSecurity security
	NoticeTypeSecurity = "security"
)

// Notice2 v2.
type Notice2 struct {
	Status   int8      `json:"status"`
	Type     string    `json:"type,omitempty"`
	Realname *Realname `json:"realname,omitempty"`
	Security *Security `json:"security,omitempty"`
}

// Realname struct.
type Realname struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Security struct.
type Security struct {
	Location string    `json:"location,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	IP       string    `json:"ip,omitempty"`
	Mid      int64     `json:"mid,omitempty"`
}
