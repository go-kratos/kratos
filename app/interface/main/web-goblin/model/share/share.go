package share

import (
	suitmdl "go-common/app/service/main/usersuit/model"
	xtime "go-common/library/time"
)

// Share share .
type Share struct {
	ID        int64      `json:"id"`
	Mid       int64      `json:"mid"`
	DayCount  int64      `json:"day_count"`
	Cycle     int64      `json:"cycle"`
	ShareDate int64      `json:"share_date"`
	Ctime     xtime.Time `json:"ctime"`
	Mtime     xtime.Time `json:"mtime"`
}

// Encourage encourage .
type Encourage struct {
	UserInfo   interface{} `json:"user_info"`
	TodayShare int64       `json:"today_share"`
	ShareDays  int64       `json:"share_days"`
	Pendants   interface{} `json:"pendants"`
}

// GroupPendant groupPendant .
type GroupPendant struct {
	NeedDays int64                     `json:"need_days"`
	Pendant  *suitmdl.GroupPendantList `json:"pendant"`
}
