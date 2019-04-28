package model

import (
	xtime "go-common/library/time"
)

// TaskLog is record task
type TaskLog struct {
	ID        int64      `json:"id"`
	TaskID    int64      `json:"task_id"`
	MID       int64      `json:"mid"`
	Build     string     `json:"build"`
	Platform  int        `json:"platform"`
	TaskState int        `json:"task_state"`
	Reason    string     `json:"reason"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"mtime"`
}

// TaskLogPager TaskLog Pager
type TaskLogPager struct {
	Total    int64      `json:"total"`
	PageNo   int        `json:"page_no" default:"1"`
	PageSize int        `json:"page_size" default:"20"`
	Items    []*TaskLog `json:"items"`
}
