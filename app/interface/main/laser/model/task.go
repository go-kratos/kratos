package model

import (
	xtime "go-common/library/time"
)

// Task is laser task
type Task struct {
	ID           int64      `json:"id"`
	AdminID      int64      `json:"admin_id"`
	MID          int64      `json:"mid"`
	LogDate      xtime.Time `json:"log_date"`
	ContactEmail string     `json:"contact_email"`
	SourceType   int        `json:"source_type"`
	Platform     int        `json:"platform"`
	State        int        `json:"state"`
	IsDeleted    int        `json:"is_deleted"`
	CTime        xtime.Time `json:"ctime"`
	MTime        xtime.Time `json:"mtime"`
}

const (
	ALL_PLATFORM = iota
)

// TaskInfo is to set as value of memcache key(mid)
type TaskInfo struct {
	MID        int64
	LogDate    xtime.Time
	SourceType int
	Platform   int
	Empty      bool
}
