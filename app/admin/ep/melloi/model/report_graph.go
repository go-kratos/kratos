package model

import (
	"time"
)

//ReportGraph reportgraph
type ReportGraph struct {
	ID             int       `json:"id" gorm:"AUTO_INCREMENT;primary_key;" form:"id"`
	TestName       string    `json:"test_name" form:"test_name"`
	TestNameNick   string    `json:"test_name_nick" form:"test_name_nick"`
	Count          int       `json:"count"`
	QPS            int       `json:"qps"`
	AvgTime        int       `json:"avg_time"`
	Min            int       `json:"min"`
	Max            int       `json:"max"`
	Error          int       `json:"error"`
	FailPercent    string    `json:"fail_percent"`
	NinetyTime     int       `json:"ninety_time"`
	NinetyFiveTime int       `json:"ninety_five_time"`
	NinetyNineTime int       `json:"ninety_nine_time"`
	NetIo          int       `json:"net_io"`
	CodeEll        int       `json:"code_ell"`
	CodeWll        int       `json:"code_wll"`
	CodeWly        int       `json:"code_wly"`
	CodeWle        int       `json:"code_wle"`
	CodeWls        int       `json:"code_wls"`
	CodeSll        int       `json:"code_sll"`
	CodeSly        int       `json:"code_sly"`
	CodeSls        int       `json:"code_sls"`
	CodeKong       int       `json:"code_kong"`
	CodeNonHTTP    int       `json:"code_non_http"`
	CodeOthers     int       `json:"code_others"`
	ElapsdTime     int       `json:"elapsd_time"`
	PodName        string    `json:"pod_name" form:"pod_name"`
	ThreadsSum     int       `json:"threads_sum"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
	FiftyTime      int       `json:"fifty_time"`
	Code301        int       `json:"code301"`
	Code302        int       `json:"code302"`
	BeginTime      time.Time `json:"begin_time" gorm:"-"`
	QpsRecent      int       `json:"qps_recent" gorm:"-"`
}

//TableName table name
func (r ReportGraph) TableName() string {
	return "report_graph"
}
