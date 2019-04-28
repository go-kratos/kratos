package upcrmmodel

import (
	"encoding/json"
	"fmt"
	"go-common/library/time"
	systime "time"
)

const (
	//BusinessTypeArticleAudit 1
	BusinessTypeArticleAudit = 1
)
const (
	//DateStr date format
	DateStr = "2006-01-02"
	//CreditLogTableCount all credit log table count
	CreditLogTableCount = 100
)

//ArgCreditLogAdd arg
type ArgCreditLogAdd struct {
	Type         int             `form:"type" json:"type"`                    // 日志类型，具体与业务方确定
	OpType       int             `form:"op_type" json:"optype"`               // 操作类型，具体与业务方确定
	Reason       int             `form:"reason" json:"reason"`                // 原因类型，具体与业务方确定
	BusinessType int             `form:"bussiness_type" json:"business_type"` // 业务类型
	Mid          int64           `form:"mid" validate:"required" json:"mid"`  // 用户id
	Oid          int64           `form:"oid" json:"oid"`                      // 对象类型，如aid
	UID          int             `form:"uid" json:"uid"`                      // 管理员id
	Content      string          `form:"content" json:"content"`              // 日志内容描述
	CTime        time.Time       `form:"ctime" json:"ctime"`                  // 创建时间
	Extra        json.RawMessage `form:"extra" json:"extra,omitempty"`        // 额外字段，与业务方确定
}

//ArgMidDate arg
type ArgMidDate struct {
	Mid       int64  `form:"mid" validate:"required"`
	Days      int    `form:"days"`                   // 最近n天内的数据，1表示最近1天（今天），2表示最近2天，默认为0
	FromDate  string `form:"from_date"`              // 2006-01-02
	ToDate    string `form:"to_date"`                // 2006-01-02
	ScoreType int    `form:"score_type" default:"3"` // 分数类型, 1,2,3，参见ScoreTypeCredit
}

//GetScoreParam arg
type GetScoreParam struct {
	Mid       int64
	FromDate  systime.Time
	ToDate    systime.Time
	ScoreType int
}

//ArgGetLogHistory arg
type ArgGetLogHistory struct {
	Mid      int64        `form:"mid" validate:"required"`
	FromDate systime.Time `form:"from_date"`
	ToDate   systime.Time `form:"to_date"`
	Limit    int          `form:"limit" default:"20"`
}

//CreditLog db struct
type CreditLog struct {
	ID           uint `gorm:"primary_key" json:"-"`
	Type         int
	OpType       int
	Reason       int
	BusinessType int
	Mid          int64
	Oid          int64
	UID          int `gorm:"column:uid"`
	Content      string
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
	Extra        string    `sql:"type:text;" json:"-"`
}

//TableName table name
func (c *CreditLog) TableName() string {
	return getTableName(c.Mid)
}

func getTableName(mid int64) string {
	return fmt.Sprintf("credit_log_%02d", mid%CreditLogTableCount)
}

//CopyFrom copy
func (c *CreditLog) CopyFrom(arg *ArgCreditLogAdd) *CreditLog {
	c.Type = arg.Type
	c.OpType = arg.OpType
	c.BusinessType = arg.BusinessType
	c.Reason = arg.Reason
	c.Mid = arg.Mid
	c.Oid = arg.Oid
	c.UID = arg.UID
	c.Content = arg.Content
	c.CTime = arg.CTime
	c.MTime = arg.CTime
	c.Extra = string(arg.Extra)
	return c
}

//SimpleCreditLog db struct
type SimpleCreditLog struct {
	ID           uint      `gorm:"primary_key" json:"-"`
	Type         int       `json:"type"`
	OpType       int       `json:"op_type"`
	Reason       int       `json:"reason"`
	BusinessType int       `json:"business_type"`
	Mid          int64     `json:"mid"`
	Oid          int64     `json:"oid"`
	CTime        time.Time `gorm:"column:ctime" json:"ctime"`
}

//TableName table name
func (c *SimpleCreditLog) TableName() string {
	return getTableName(c.Mid)
}

//SimpleCreditLogWithContent log with content
type SimpleCreditLogWithContent struct {
	SimpleCreditLog
	Content string `form:"content" json:"content"` // 日志内容描述
}
