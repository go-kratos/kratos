package upcrmmodel

import (
	"go-common/library/time"
)

const (
	//BusinessTypeArticleAudit 稿件的审核
	BusinessTypeArticleAudit = 1
)

//SimpleCreditLog simple credit log
type SimpleCreditLog struct {
	ID           uint      `json:"-"`
	Type         int       `json:"type"`
	OpType       int       `json:"op_type"`
	Reason       int       `json:"reason"`
	BusinessType int       `json:"business_type"`
	Mid          int64     `json:"mid"`
	Oid          int64     `json:"oid"`
	CTime        time.Time `json:"ctime"`
}

//SimpleCreditLogWithContent simple credit log with content
type SimpleCreditLogWithContent struct {
	SimpleCreditLog
	Content string `form:"content" json:"content"` // 日志内容描述
}
