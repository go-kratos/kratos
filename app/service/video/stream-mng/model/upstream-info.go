package model

import (
	"time"
)

// UpStreamInfo 上行调度信息
type UpStreamInfo struct {
	ID       int64     `json:"id,omitempty"`
	RoomID   int64     `json:"room_id,omitempty"`
	CDN      int64     `json:"cdn,omitempty"`
	PlatForm string    `json:"platform,omitempty"`
	IP       string    `json:"ip,omitempty"`
	Country  string    `json:"country,omitempty"`
	City     string    `json:"city,omitempty"`
	ISP      string    `json:"isp,omitempty"`
	Ctime    time.Time `json:"ctime,omitempty"`
}

// SummaryUpStreamInfo 上行调度统计信息
type SummaryUpStreamRtmp struct {
	CDN      int64  `json:"cdn,omitempty"`
	ISP      string `json:"isp,omitempty"`
	Count    int64  `json:"count,omitempty"`
	Country  string `json:"country,omitempty"`
	City     string `json:"city,omitempty"`
	PlatForm string `json:"platform,omitempty"`
}
