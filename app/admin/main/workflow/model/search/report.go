package search

import (
	"time"

	"go-common/app/admin/main/workflow/model"
)

// business id const
const (
	Archive           = 3
	Workflow          = 11
	LogAuditAction    = "log_audit"
	_auditLogSrhComID = "log_audit_group"

	//// IndexTypeYear index by year
	//IndexTypeYear indexType = "year"
	//// IndexTypeMonth index by month
	//IndexTypeMonth indexType = "month"
	//// IndexTypeWeek index by week
	//IndexTypeWeek indexType = "week"
	//// IndexTypeDay index by day
	//IndexTypeDay indexType = "day"
)

// AuditReportSearchCond .
type AuditReportSearchCond struct {
	AppID         string    `json:"app_id"`
	Fields        []string  `json:"fields"`
	IndexTimeType string    `json:"index_time_type"`
	IndexTimeFrom time.Time `json:"index_time_from"`
	IndexTimeEnd  time.Time `json:"index_time_end"`
	Business      int       `json:"business"`
	UName         string    `json:"uname"`
	UID           []int64   `json:"uid"`
	Oid           []int64   `json:"oid"`
	Type          []int     `json:"type"`
	Action        string    `json:"action"`
	CTime         string    `json:"ctime"`
	Order         string    `json:"order"`
	Sort          string    `json:"sort"`
	Int0          []int64   `json:"int_0"`
	Int1          []int64   `json:"int_1"`
	Int2          []int64   `json:"int_2"`
	Str0          string    `json:"str_0"`
	Str1          string    `json:"str_1"`
	Str2          string    `json:"str_2"`
	Group         string    `json:"group"`
	Distinct      string    `json:"distinct"`
}

// AuditLogSearchCommonResult .
type AuditLogSearchCommonResult struct {
	Page   *model.Page  `json:"page"`
	Result []*ReportLog `json:"result"`
}

// ReportLog .
type ReportLog struct {
	Action    string `json:"action"`
	Business  int64  `json:"business"`
	CTime     string `json:"ctime"`
	ExtraData string `json:"extra_data"`
	Str0      string `json:"str_0"`
	Str1      string `json:"str_1"`
	Str2      string `json:"str_2"`
	Int0      int64  `json:"int_0"`
	Int1      int64  `json:"int_1"`
	Int2      int64  `json:"int_2"`
	Int3      int64  `json:"int_3"`
	Oid       int64  `json:"oid"`
	Type      int64  `json:"type"`
	UID       int64  `json:"uid"`
	UName     string `json:"uname"`
}
