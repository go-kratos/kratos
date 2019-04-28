package model

import (
	"encoding/csv"
	"fmt"
	"time"

	xtime "go-common/library/time"
)

// QueryActivityByIDReq arg
type QueryActivityByIDReq struct {
	FromDate string `form:"from_date"` // 20180101
	ToDate   string `form:"to_date"`   // 20180102 closed interval [20180101, 20180102]
	ID       int64  `form:"id"`        // activity id, if not 0, FromDate and toDate not used
	PageArg
	ExportArg
}

//UpSummaryBonusInfo bonus for one up info
type UpSummaryBonusInfo struct {
	Mid             int64      `json:"mid"`
	BilledMoney     float64    `json:"billed_money"`
	UnbilledMoney   float64    `json:"unbilled_money"`
	LastBillTime    string     `json:"last_bill_time"`    // 20180101， 最近结算时间
	TmpBillTime     xtime.Time `json:"-"`                 // 用来计算LastBillTime
	TotalBonusMoney float64    `json:"total_bonus_money"` // 所有的中奖金额
}

//QueryUpBonusByMidResult query result
type QueryUpBonusByMidResult struct {
	Result []*UpSummaryBonusInfo `json:"result"`
	PageResult
}

//GetFileName get file name
func (q *QueryUpBonusByMidResult) GetFileName() string {
	return fmt.Sprintf("%s_%s.csv", "结算记录", time.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (q *QueryUpBonusByMidResult) ToCsv(writer *csv.Writer) {
	var title = []string{
		"UID",
		"已结算",
		"未结算",
		"最近结算时间"}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.Result {
		var record []string
		record = append(record,
			intFormat(v.Mid),
			floatFormat(v.BilledMoney),
			floatFormat(v.UnbilledMoney),
			v.LastBillTime,
		)
		writer.Write(record)
	}
}
