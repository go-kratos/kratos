package model

import xtime "go-common/library/time"

// ADMINID 后台管理员id.
const ADMINID = 1227

// AiConfig struct.
type AiConfig struct {
	// AI阀值
	Threshold float64 `json:"threshold"`
	// AI真实分标准
	TrueScore float64 `json:"true_score"`
}

// AiWhite struct.
type AiWhite struct {
	ID    int64      `json:"id"`
	MID   int64      `json:"mid"`
	State int8       `json:"state"`
	Ctime xtime.Time `json:"ctime"`
	Mtime xtime.Time `json:"mtime"`
}

// AiScore struct.
type AiScore struct {
	Scores    []float64 `json:"scores"`
	Threshold float64   `json:"threshold"`
	Note      string    `json:"note"`
}
