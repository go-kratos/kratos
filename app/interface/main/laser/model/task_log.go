package model

import (
	xtime "go-common/library/time"
)

// TaskLog record the uploaded task details
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
