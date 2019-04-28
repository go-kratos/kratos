package model

import "go-common/library/time"

// Jointly def.
type Jointly struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	StartTime int64     `json:"start_time"`
	EndTime   int64     `json:"end_time"`
	Link      string    `json:"link"`
	IsHot     int8      `json:"is_hot"`
	State     int8      `json:"state"`
	Operator  string    `json:"operator"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// JointlyResp jointly resp.
type JointlyResp struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsHot   int8   `json:"is_hot"`
	Link    string `json:"link"`
	EndTime int64  `json:"end_time"`
}
