package model

import "errors"

var (
	// ErrSearchReport search report error
	ErrSearchReport = errors.New("search report error")
	// ErrSearchReply search reply error
	ErrSearchReply = errors.New("search reply error")
	// ErrSearchMonitor search monitor error
	ErrSearchMonitor = errors.New("search monitor error")
	// ErrMsgSend send message error
	ErrMsgSend = errors.New("send message error")

	// AttrNo attribute no
	AttrNo = uint32(0)
	// AttrYes attribute yes
	AttrYes = uint32(1)

	// EventReportAdd event add a report
	EventReportAdd = "report_add"
	// EventReportDel event del a report
	EventReportDel = "report_del"
	// EventReportIgnore event ignore a report
	EventReportIgnore = "report_ignore"
	// EventReportRecover event recover a report
	EventReportRecover = "report_recover"
)

const (
	// WeightLike like sort weight
	WeightLike = 2
	// WeightHate hate sort weight
	WeightHate = 4
)

// Pager page info.
type Pager struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pagesize"`
	Total    int64 `json:"total"`
}

// NewPager NewPager
type NewPager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}
