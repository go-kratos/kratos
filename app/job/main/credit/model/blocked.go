package model

import xtime "go-common/library/time"

// BlockAndMoralStatus blocked status and moral.
type BlockAndMoralStatus struct {
	MID    int64      `json:"mid"`
	Status int8       `json:"status"`
	STime  xtime.Time `json:"start_time"`
	ETime  xtime.Time `json:"end_time"`
}

// BlockLabourAnswerLog .
type BlockLabourAnswerLog struct {
	ID    int64  `json:"id"`
	MID   int64  `json:"mid"`
	Score int8   `json:"score"`
	CTime string `json:"ctime"`
}
