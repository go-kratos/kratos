package offlineactivity

import (
	"go-common/library/time"
	"strconv"
)

const (
	//BonusTypeThing 奖品
	BonusTypeThing = 0
	//BonusTypeMoney 奖金
	BonusTypeMoney = 1
)

//ActivityState activity state
type ActivityState int8

const (
	//StateDelete 删除了
	StateDelete = 100

	// 用来转换 float * moneyCont -> int
	moneyConst = 1000
)

const (
	//ActivityStateInit 初始状态
	ActivityStateInit ActivityState = 0
	//ActivityStateSending 发送状态
	ActivityStateSending ActivityState = 1
	//ActivityStateWaitResult 等待审核结果
	ActivityStateWaitResult ActivityState = 2
	//ActivityStateSucess 成功
	ActivityStateSucess ActivityState = 10
	//ActivityStateFail 发送失败
	ActivityStateFail ActivityState = 11
	//ActivityStateCreateFail 创建失败
	ActivityStateCreateFail ActivityState = 12
)

const (
	//TableOfflineActivityInfo info name
	TableOfflineActivityInfo = "offline_activity_info"
	//TableOfflineActivityBonus bonus name
	TableOfflineActivityBonus = "offline_activity_bonus"
	//TableOfflineActivityResult result name
	TableOfflineActivityResult = "offline_activity_result"
	//TableOfflineActivityShellOrder shell order name
	TableOfflineActivityShellOrder = "offline_activity_shell_order"
)

//OfflineActivityInfo table info
type OfflineActivityInfo struct {
	ID        int64     `json:"id" gorm:"column:id"`
	Title     string    `json:"title" gorm:"column:title"`
	Link      string    `json:"link" gorm:"column:link"`
	BonusType int8      `json:"bonus_type" gorm:"column:bonus_type"`
	Memo      string    `json:"memo" gorm:"column:memo"`
	Creator   string    `json:"creator" gorm:"column:creator"`
	State     int8      `json:"state" gorm:"column:state"`
	CTime     time.Time `json:"ctime" gorm:"column:ctime"`
	MTime     time.Time `json:"-" gorm:"column:mtime"`
}

//TableName get table name
func (o *OfflineActivityInfo) TableName() string {
	return TableOfflineActivityInfo
}

//BonusState bonus state
type BonusState int8

const (
	//BonusStateInit init state
	BonusStateInit BonusState = 0
)

//OfflineActivityBonus table bonus
type OfflineActivityBonus struct {
	ID          int64     `json:"id" gorm:"column:id"`
	ActivityID  int64     `json:"activity_id" gorm:"column:activity_id"`
	TotalMoney  int64     `json:"total_money" gorm:"column:total_money"`
	MemberCount uint32    `json:"member_count" gorm:"column:member_count"`
	State       int8      `json:"state" gorm:"column:state"`
	CTime       time.Time `json:"ctime" gorm:"column:ctime"`
	MTime       time.Time `json:"mtime" gorm:"column:mtime"`
}

//TableName tablename
func (o *OfflineActivityBonus) TableName() string {
	return TableOfflineActivityBonus
}

//OfflineActivityResult table result
type OfflineActivityResult struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`
	ActivityID int64     `json:"activity_id" gorm:"column:activity_id"`
	BonusID    int64     `json:"bonus_id" gorm:"column:bonus_id"`
	BonusType  int8      `json:"bonus_type" gorm:"column:bonus_type"`
	Mid        int64     `json:"mid" gorm:"column:mid"`
	BonusMoney int64     `json:"bonus_money" gorm:"column:bonus_money"`
	OrderID    string    `json:"order_id" gorm:"column:order_id"`
	State      int8      `json:"state" gorm:"column:state"`
	CTime      time.Time `json:"ctime" gorm:"column:ctime"`
	MTime      time.Time `json:"mtime" gorm:"column:mtime"`
}

//TableName table name
func (o *OfflineActivityResult) TableName() string {
	return TableOfflineActivityResult
}

//OfflineActivityShellOrder table shell order, order for shell system
type OfflineActivityShellOrder struct {
	ID          int64     `json:"id" gorm:"column:id"`
	ResultID    int64     `json:"result_id" gorm:"column:result_id"`
	OrderID     string    `json:"order_id" gorm:"column:order_id"`
	OrderStatus string    `json:"order_status" gorm:"column:order_status"`
	CTime       time.Time `json:"ctime" gorm:"column:ctime"`
	MTime       time.Time `json:"mtime" gorm:"column:mtime"`
}

//TableName table name
func (o *OfflineActivityShellOrder) TableName() string {
	return TableOfflineActivityShellOrder
}

//GetMoneyFromDb get money from db
func GetMoneyFromDb(dbmoney int64) float64 {
	return float64(dbmoney) / moneyConst
}

//GetMoneyForDb set money to db
func GetMoneyForDb(realmoney float64) int64 {
	return int64(realmoney * moneyConst)
}

//StateToString State to string
func StateToString(state int) string {
	switch state {
	case int(ActivityStateInit):
		return "初始"
	case int(ActivityStateSending):
		return "发送贝壳中"
	case int(ActivityStateWaitResult):
		return "等待审核结果"
	case int(ActivityStateSucess):
		return "成功"
	case int(ActivityStateFail):
		return "失败"
	case int(ActivityStateCreateFail):
		return "创建失败"
	case StateDelete:
		return "已删除"
	default:
		return strconv.Itoa(state)
	}
}
