package signmodel

import (
	"go-common/app/admin/main/up/util"
	"go-common/library/time"
)

const (
	//TableNameSignTask .
	TableNameSignTask = "sign_task"
	//TableNameSignTaskHistory .
	TableNameSignTaskHistory = "sign_task_history"
	//TableNameSignUp .
	TableNameSignUp = "sign_up"
	//TableNameSignContract .
	TableNameSignContract = "sign_contract"
	//TableNameSignPay .
	TableNameSignPay = "sign_pay"
	//TableNameSignTaskAbsence .
	TableNameSignTaskAbsence = "sign_task_absence"
)

const (
	//TaskTypeAccumulate 0
	TaskTypeAccumulate = 0
	//TaskTypeDay 1
	TaskTypeDay = 1
	//TaskTypeWeek 2
	TaskTypeWeek = 2
	//TaskTypeMonth 3
	TaskTypeMonth = 3
	//TaskTypeQuarter 4 季度
	TaskTypeQuarter = 4
)

//TaskTypeStr get task type str
func TaskTypeStr(taskType int) string {
	switch taskType {
	case TaskTypeAccumulate:
		return "累计"
	case TaskTypeDay:
		return "每日"
	case TaskTypeWeek:
		return "每周"
	case TaskTypeMonth:
		return "每月"
	case TaskTypeQuarter:
		return "每季度"
	}
	return "未知"
}

const (
	//EmailStateNotSend 0
	EmailStateNotSend = 0
	//EmailStateSendSucc 1
	EmailStateSendSucc = 1
)

//SignUpOnlyID struct
type SignUpOnlyID struct {
	ID uint32
}

//SignUpOnlySignID  struct
type SignUpOnlySignID struct {
	SignID uint32
}

// const sign_up中的state定义
const (
	SignStateOnSign = 0
	SignStateExpire = 1
)

// const sign_up中的due_warn定义
const (
	DueWarnNoWarn = 1
	DueWarnWarn   = 2
)

// const sign_up中的pay_expire_state定义
const (
	// PayExpireStateNormal 未到期
	PayExpireStateNormal = 1
	// PayExpireStateDue 即将到期
	PayExpireStateDue = 2
)

//SignUp  struct
type SignUp struct {
	ID             uint32    `gorm:"column:id"`
	Sex            int8      `gorm:"column:sex"`
	Mid            int64     `gorm:"column:mid"`
	BeginDate      time.Time `gorm:"column:begin_date"`
	EndDate        time.Time `gorm:"column:end_date"`
	State          int8      `gorm:"column:state"`
	DueWarn        int8      `gorm:"column:due_warn"`
	PayExpireState int8      `gorm:"column:pay_expire_state"`
	Country        string    `gorm:"column:country"`
	Province       string    `gorm:"column:province"`
	City           string    `gorm:"column:city"`
	Note           string    `gorm:"column:note"`
	AdminID        int       `gorm:"column:admin_id"`
	AdminName      string    `gorm:"column:admin_name"`
	EmailState     int8      `gorm:"column:email_state"`
	Ctime          time.Time `gorm:"column:ctime"`
	Mtime          time.Time `gorm:"column:mtime"`
}

//TableName .
func (s *SignUp) TableName() string {
	return TableNameSignUp
}

//SignPay  struct
type SignPay struct {
	ID         uint32    `gorm:"column:id"`
	Mid        int64     `gorm:"column:mid"`
	SignID     uint32    `gorm:"column:sign_id"`
	DueDate    time.Time `gorm:"column:due_date"`
	PayValue   int64     `gorm:"column:pay_value"`
	State      int8      `gorm:"column:state"`
	Note       string    `gorm:"column:note"`
	EmailState int8      `gorm:"column:email_state"`
	Ctime      time.Time `gorm:"column:ctime"`
	Mtime      time.Time `gorm:"column:mtime"`
}

//TableName .
func (s *SignPay) TableName() string {
	return TableNameSignPay
}

//SignTaskState sign task's state
type SignTaskState int8

const (
	//SignTaskStateInit init
	SignTaskStateInit SignTaskState = 0
	//SignTaskStateRunning task running
	SignTaskStateRunning SignTaskState = 1
	//SignTaskStateFinish finish
	SignTaskStateFinish SignTaskState = 2
	//SignTaskStateDelete delete
	SignTaskStateDelete SignTaskState = 100
)

const (
	// SignTaskAttrBitBusiness 商单标记
	SignTaskAttrBitBusiness = 0
)

//SignTask  struct
type SignTask struct {
	ID            uint32        `gorm:"column:id"`
	Mid           int64         `gorm:"column:mid"`
	SignID        uint32        `gorm:"column:sign_id"`
	TaskType      int8          `gorm:"column:task_type"`
	TaskCounter   int32         `gorm:"column:task_counter"`
	TaskCondition int32         `gorm:"column:task_condition"`
	TaskData      string        `gorm:"column:task_data"`
	Attribute     int64         `gorm:"column:attribute"`
	GenerateDate  time.Time     `gorm:"column:generate_date"`
	State         SignTaskState `gorm:"column:state"`
	Ctime         time.Time     `gorm:"column:ctime"`
	Mtime         time.Time     `gorm:"column:mtime"`
}

//IsAttrSet is attribute set, see SignTaskAttrBitXXX above
func (s *SignTask) IsAttrSet(bit int) bool {
	return util.IsBitSet64(s.Attribute, uint(bit))
}

//TableName .
func (s *SignTask) TableName() string {
	return TableNameSignTask
}

//SignContract  struct
type SignContract struct {
	ID       uint32    `gorm:"column:id"`
	Mid      int64     `gorm:"column:mid"`
	SignID   uint32    `gorm:"column:sign_id"`
	Filename string    `gorm:"column:filename"`
	Filelink string    `gorm:"column:filelink"`
	State    int8      `gorm:"column:state"`
	Ctime    time.Time `gorm:"column:ctime"`
	Mtime    time.Time `gorm:"column:mtime"`
}

//TableName .
func (s *SignContract) TableName() string {
	return TableNameSignContract
}

//SignTaskHistory  struct
type SignTaskHistory struct {
	ID             uint32        `gorm:"column:id"`
	Mid            int64         `gorm:"column:mid"`
	SignID         uint32        `gorm:"column:sign_id"`
	TaskTemplateID uint32        `gorm:"column:task_template_id"`
	TaskType       int8          `gorm:"column:task_type"`
	TaskCounter    int32         `gorm:"column:task_counter"`
	TaskCondition  int32         `gorm:"column:task_condition"`
	TaskData       string        `gorm:"column:task_data"`
	Attribute      int64         `gorm:"column:attribute"`
	GenerateDate   time.Time     `gorm:"column:generate_date"`
	State          SignTaskState `gorm:"column:state"`
	Ctime          time.Time     `gorm:"column:ctime"`
	Mtime          time.Time     `gorm:"column:mtime"`
}

//TableName .
func (s *SignTaskHistory) TableName() string {
	return TableNameSignTaskHistory
}

//SignTaskAbsenceState .
type SignTaskAbsenceState int8

const (
	//SignTaskAbsenceStateInit initial
	SignTaskAbsenceStateInit SignTaskAbsenceState = 0
	//SignTaskAbsenceStateDelete deleted
	SignTaskAbsenceStateDelete SignTaskAbsenceState = 100
)

//SignTaskAbsence table
type SignTaskAbsence struct {
	ID            uint32    `gorm:"column:id" json:"id"`
	SignId        uint32    `gorm:"column:sign_id" json:"sign_id"`
	Mid           int64     `gorm:"column:mid" json:"mid"`
	TaskHistoryId uint32    `gorm:"column:task_history_id" json:"task_history_id"`
	AbsenceCount  uint32    `gorm:"column:absence_count" json:"absence_count"`
	Reason        string    `gorm:"column:reason" json:"reason"`
	State         int8      `gorm:"column:state" json:"state"`
	AdminId       int64     `gorm:"column:admin_id" json:"admin_id"`
	AdminName     string    `gorm:"column:admin_name" json:"admin_name"`
	Ctime         time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime         time.Time `gorm:"column:mtime" json:"mtime"`
}

//TableName .
func (s *SignTaskAbsence) TableName() string {
	return TableNameSignTaskAbsence
}
