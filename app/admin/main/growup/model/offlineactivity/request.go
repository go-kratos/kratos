package offlineactivity

import (
	"encoding/csv"
	"fmt"
	"go-common/library/time"
	"math"
	"strconv"
	time2 "time"
)

const (
	dateFmt     = "20060102"
	dateTimeFmt = "20060102_150405"
	filenamePre = "活动"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

//ToFixed fix float precision
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output // since go 1.9 doesn't have a math.Round function...
}

//ExportArg export arg
type ExportArg struct {
	Export string `form:"export"`
}

//ExportFormat export format
func (e *ExportArg) ExportFormat() string {
	return e.Export
}

//BonusInfo bonus struct
type BonusInfo struct {
	TotalMoney  float64 `json:"total_money"`
	MemberCount int64   `json:"member_count"`
	Mids        string  `json:"mids"`
	MidList     []int64 `json:"-"`        // 存储转换后的mids
	Filename    string  `json:"filename"` // 如果有文件名，优先使用文件
}

//AddActivityArg add activity arg
type AddActivityArg struct {
	Title     string       `json:"title"`
	Link      string       `json:"link"`
	BonusType int8         `json:"bonus_type"`
	Memo      string       `json:"memo"`
	BonusList []*BonusInfo `json:"bonus_list"`
	Creator   string       `json:"-"`
}

//PreAddActivityResult result, nothing to return
type PreAddActivityResult struct {
	BonusType   int8    `json:"bonus_type"`
	MemberCount int64   `json:"member_count"`
	TotalMoney  float64 `json:"total_money"`
}

//AddActivityResult result, nothing to return
type AddActivityResult struct {
}

//PageArg page arg
type PageArg struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

//CheckPageValidation check the page validte, return limit offset
func (arg *PageArg) CheckPageValidation() (limit, offset int) {
	if arg.Page < 1 {
		arg.Page = 1
	}
	if arg.Size > 100 || arg.Size <= 0 {
		arg.Size = 10
	}
	limit = arg.Size
	offset = (arg.Page - 1) * limit
	return
}

//ToPageResult cast to page result
func (arg *PageArg) ToPageResult(total int) (res PageResult) {
	res.TotalCount = total
	res.Page = arg.Page
	return
}

//PageResult page result
type PageResult struct {
	Page       int `json:"page"`
	TotalCount int `json:"total_count"`
}

//QueryActivityByIDArg arg
type QueryActivityByIDArg struct {
	FromDate string `form:"from_date"` // 20180101
	ToDate   string `form:"to_date"`   // 20180102 closed interval [20180101, 20180102]
	ID       int64  `form:"id"`        // activity id, if not 0, FromDate and toDate not used
	PageArg
	ExportArg
}

//QueryActivityInfo activity info
type QueryActivityInfo struct {
	OfflineActivityInfo
	CreateTime  string  `json:"create_time"`
	TotalMoney  float64 `json:"total_money"`
	MemberCount uint32  `json:"member_count"`
	PayDate     string  `json:"pay_date"`
}

//QueryActivityResult result
type QueryActivityResult struct {
	Result []*QueryActivityInfo `json:"result"`
	PageResult
}

//GetFileName get file name
func (q *QueryActivityResult) GetFileName() string {
	return fmt.Sprintf("%s-%s_%s.csv", filenamePre, "按活动", time2.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (q *QueryActivityResult) ToCsv(writer *csv.Writer) {
	var title = []string{
		"添加日期",
		"活动id",
		"活动标题",
		"活动链接",
		"支出金额",
		"中奖人数",
		"付款日期",
		"发奖状态"}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.Result {
		var record []string
		record = append(record,
			v.CreateTime,
			strconv.Itoa(int(v.ID)),
			v.Title,
			v.Link,
			fmt.Sprintf("%0.2f", v.TotalMoney),
			strconv.Itoa(int(v.MemberCount)),
			v.PayDate,
			StateToString(int(v.State)),
		)
		writer.Write(record)
	}

}

//CopyFromActivityDB copy from db
func (q *QueryActivityInfo) CopyFromActivityDB(v *OfflineActivityInfo) {
	q.OfflineActivityInfo = *v
	q.CreateTime = v.CTime.Time().Format(dateFmt)
	q.PayDate = q.CreateTime
}

//QueryUpBonusByMidArg arg
type QueryUpBonusByMidArg struct {
	Mid int64 `form:"mid"` // activity id, if not 0, FromDate and toDate not used
	PageArg
	ExportArg
}

//UpSummaryBonusInfo bonus for one up info
type UpSummaryBonusInfo struct {
	Mid             int64     `json:"mid"`
	BilledMoney     float64   `json:"billed_money"`
	UnbilledMoney   float64   `json:"unbilled_money"`
	LastBillTime    string    `json:"last_bill_time"`    // 20180101， 最近结算时间
	TmpBillTime     time.Time `json:"-"`                 // 用来计算LastBillTime
	TotalBonusMoney float64   `json:"total_bonus_money"` // 所有的中奖金额
	ActivityID      int64     `json:"activity_id"`       // 中奖活动ID
}

//Finish finish it before send back
func (s *UpSummaryBonusInfo) Finish() {
	if s.TmpBillTime > 0 {
		s.LastBillTime = s.TmpBillTime.Time().Format(dateFmt)
	}
	s.BilledMoney = ToFixed(s.BilledMoney, 2)
	s.UnbilledMoney = ToFixed(s.UnbilledMoney, 2)
	s.TotalBonusMoney = ToFixed(s.TotalBonusMoney, 2)
}

//QueryUpBonusByMidResult query result
type QueryUpBonusByMidResult struct {
	Result []*UpSummaryBonusInfo `json:"result"`
	PageResult
}

//GetFileName get file name
func (q *QueryUpBonusByMidResult) GetFileName() string {
	return fmt.Sprintf("%s-%s_%s.csv", filenamePre, "结算记录", time2.Now().Format(dateTimeFmt))
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

//QueryUpBonusByActivityResult use same result, we define another type, because we need to implement interface again
type QueryUpBonusByActivityResult QueryUpBonusByMidResult

//GetFileName get file name
func (q *QueryUpBonusByActivityResult) GetFileName() string {
	return fmt.Sprintf("%s-%s_%s.csv", filenamePre, "详细信息", time2.Now().Format(dateTimeFmt))
}

//ToCsv to buffer
func (q *QueryUpBonusByActivityResult) ToCsv(writer *csv.Writer) {
	var title = []string{
		"UID",
		"中奖活动id",
		"中奖金额",
		"结算日期",
		"结算收入",
		"未结算收入"}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.Result {
		var record []string
		record = append(record,
			intFormat(v.Mid),
			intFormat(v.ActivityID),
			floatFormat(v.TotalBonusMoney),
			v.LastBillTime,
			floatFormat(v.BilledMoney),
			floatFormat(v.UnbilledMoney),
		)
		writer.Write(record)
	}
}

//QueryActvityMonthArg arg
type QueryActvityMonthArg struct {
	ExportArg
	PageArg
	FromDate string `form:"from_date"`
	ToDate   string `form:"to_date"`
}

//ActivityMonthInfo info
type ActivityMonthInfo struct {
	CreateTime           string  `json:"create_time"`
	ActivityCount        int64   `json:"activity_count"`
	MemberCount          uint32  `json:"member_count"`
	AverageMoneyMember   float64 `json:"average_money_member"`
	AverageMoneyActivity float64 `json:"average_money_activity"`
	TotalBonusMoneyMonth float64 `json:"total_bonus_money"` // 所有的中奖金额
	TotalMoneyAccumulate float64 `json:"total_money_accumulate"`
	GenerateDay          string  `json:"generate_day"` // 生成日期

	activityMap map[int64]struct{} // 用来保存活动数
}

//AddBonus 增加bonus信息
func (s *ActivityMonthInfo) AddBonus(v *OfflineActivityBonus) {
	s.MemberCount += v.MemberCount
	s.TotalBonusMoneyMonth += GetMoneyFromDb(v.TotalMoney)

	if s.activityMap == nil {
		s.activityMap = make(map[int64]struct{})
	}

	s.activityMap[v.ActivityID] = struct{}{}
}

//Finish finish
func (s *ActivityMonthInfo) Finish() {
	if s.MemberCount > 0 {
		s.AverageMoneyMember = ToFixed(s.TotalBonusMoneyMonth/float64(s.MemberCount), 2)
	}
	s.ActivityCount = int64(len(s.activityMap))
	if s.ActivityCount > 0 {
		s.AverageMoneyActivity = ToFixed(s.TotalBonusMoneyMonth/float64(s.ActivityCount), 2)
	}
}

//QueryActivityMonthResult result
type QueryActivityMonthResult struct {
	PageResult
	Result []*ActivityMonthInfo `json:"result"`
}

//GetFileName get file name
func (q *QueryActivityMonthResult) GetFileName() string {
	return fmt.Sprintf("%s-%s_%s.csv", filenamePre, "按月", time2.Now().Format(dateTimeFmt))

}

func floatFormat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)

}
func intFormat(i int64) string {
	return strconv.Itoa(int(i))
}

//ToCsv to csv
func (q *QueryActivityMonthResult) ToCsv(writer *csv.Writer) {
	var title = []string{
		"月份",
		"最新数据截止日期",
		"月支出",
		"活动数",
		"中奖人数",
		"人均支出",
		"活动均支出",
		"累计支出"}
	writer.Write(title)
	if q == nil {
		return
	}
	for _, v := range q.Result {
		var record []string
		record = append(record,
			v.CreateTime,
			v.GenerateDay,
			floatFormat(v.TotalBonusMoneyMonth),
			intFormat(v.ActivityCount),
			intFormat(int64(v.MemberCount)),
			floatFormat(v.AverageMoneyMember),
			floatFormat(v.AverageMoneyActivity),
			floatFormat(v.TotalMoneyAccumulate),
		)
		writer.Write(record)
	}
}
