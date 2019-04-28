package model

// ReportAction report message
type ReportAction struct {
	Cid      int64 `json:"cid"`       // 视频id
	Did      int64 `json:"dmid"`      // 弹幕id
	HideTime int64 `json:"hide_time"` // 弹幕隐截止j时间
}
