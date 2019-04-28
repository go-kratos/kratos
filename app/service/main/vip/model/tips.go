package model

import "go-common/library/time"

// Tips def.
type Tips struct {
	ID        int64     `json:"id"`
	Platform  int64     `json:"platform"`
	Version   int64     `json:"version"`
	Tip       string    `json:"tip"`
	Link      string    `json:"link"`
	StartTime int64     `json:"start_time"`
	EndTime   int64     `json:"end_time"`
	Level     int8      `json:"level"`
	JudgeType int8      `json:"judge_type"`
	Operator  string    `json:"operator"`
	Deleted   int8      `json:"deleted"`
	Position  int8      `json:"position"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TipsResp tips resp.
type TipsResp struct {
	ID         int64  `json:"id"`
	Version    int64  `json:"version"`
	Tip        string `json:"tip"`
	Link       string `json:"link"`
	ButtonName string `json:"button_name"`
	ButtonLink string `json:"button_link"`
}
