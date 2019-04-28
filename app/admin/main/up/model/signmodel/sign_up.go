package signmodel

import (
	"strings"
	"time"

	"go-common/app/admin/main/up/util/now"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// TableSignPay table name
	TableSignPay = "sign_pay"
	// TableSignUp table name
	TableSignUp = "sign_up"
	// TableSignTask table name
	TableSignTask = "sign_task"
	// TableSignContract table name
	TableSignContract = "sign_contract"
	// TableSignTaskAbsence table name
	TableSignTaskAbsence = "sign_task_absence"
	// TableSignTaskHistory table name
	TableSignTaskHistory = "sign_task_history"
	// TableSignViolationHistory table name
	TableSignViolationHistory = "sign_violation_history"
)

const (
	// DateDefualtFromDB .
	DateDefualtFromDB = -28800
	// DateDefualt .
	DateDefualt = "0000-00-00"
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

const (
	//TaskStateRunning 0
	TaskStateRunning = 0
	//TaskStateExpire 1
	TaskStateExpire = 1
	//TaskStateFinish 2
	TaskStateFinish = 2
)

// const .
const (
	SignUpMidAdd    = 1
	SignUpMidUpdate = 2
)

// const 。
const (
	NotContainEconomic = 1
	ContainEconomic    = 2
)

const (
	// SignUpLogBizID 签约up信息修改日志
	SignUpLogBizID int = 261
)

const (
	//SignTaskStateInit init
	SignTaskStateInit = 0
	//SignTaskStateRunning task running
	SignTaskStateRunning = 1
	//SignTaskStateFinish finish
	SignTaskStateFinish = 2
	//SignTaskStateDelete delete
	SignTaskStateDelete = 100
)

// const 变更类型.
const (
	// 年龄
	ChangeSexHistory = iota + 1
	// 用户id
	ChangeMidHistory
	// 签约周期
	ChangeSignDateHistory
	// 地区
	ChangeAreaHistory
	// 组织
	ChangeOrganizationHistory
	// 签约类型
	ChangeSignTypeHistory
	// 年龄
	ChangeAgeHistory
	// 居住地
	ChangeResidenceHistory
	// 身份证
	ChangeIDCardHistory
	// 联系方式
	ChangePhoneHistory
	// QQ
	ChangeQQHistory
	// 微信
	ChangeWechatHistory
	// 经济公司
	ChangeEconomicHistory
	// 签约付款周期
	ChangeSignPayHistory
	// 签约任务
	ChangeSignTaskHistory
	// 签约合同
	ChangeSignContractHistory
	// 签约备注
	ChangeSignNoteHistory
)

const (
	//EmailStateNotSend 0
	EmailStateNotSend = 0
	//EmailStateSendSucc 1
	EmailStateSendSucc = 1
)

const (
	// AttrYes on
	AttrYes = int64(1)
	// AttrNo off
	AttrNo = int64(0)

	// AttrBitIsBusinessArchive bit
	AttrBitIsBusinessArchive = uint(0)
)

// SignUpOnlyID struct
type SignUpOnlyID struct {
	ID uint32
}

// SignUpOnlySignID  struct
type SignUpOnlySignID struct {
	SignID uint32
}

// SignUp  struct
type SignUp struct {
	ID              int64
	Sex             int8
	Mid             int64
	BeginDate       xtime.Time
	EndDate         xtime.Time
	State           int8
	Country         string
	Province        string
	City            string
	Note            string
	AdminID         int
	AdminName       string
	EmailState      int8
	Ctime           xtime.Time `gorm:"column:ctime"`
	Mtime           xtime.Time `gorm:"column:mtime"`
	Organization    int8
	SignType        int8
	Age             int8
	Residence       string
	IDCard          string `gorm:"column:id_card"`
	Phone           string
	QQ              int64  `gorm:"column:qq"`
	Wechat          string `gorm:"column:wechat"`
	IsEconomic      int8
	EconomicCompany string
	EconomicBegin   xtime.Time
	EconomicEnd     xtime.Time
	TaskState       int8
	LeaveTimes      int
	ViolationTimes  int
	ActiveTid       int16
}

// Diff .
func (su *SignUp) Diff(oriSu *SignUp, fields map[int8]struct{}) {
	if oriSu.Sex != su.Sex {
		fields[ChangeSexHistory] = struct{}{}
	}
	if oriSu.Mid != su.Mid {
		fields[ChangeMidHistory] = struct{}{}
	}
	if oriSu.BeginDate != su.BeginDate || oriSu.EndDate != su.EndDate {
		fields[ChangeSignDateHistory] = struct{}{}
	}
	if !strings.EqualFold(oriSu.Country, su.Country) || !strings.EqualFold(oriSu.Province, su.Province) || !strings.EqualFold(oriSu.City, su.City) {
		fields[ChangeAreaHistory] = struct{}{}
	}
	if oriSu.Organization != su.Organization {
		fields[ChangeOrganizationHistory] = struct{}{}
	}
	if oriSu.SignType != su.SignType {
		fields[ChangeSignTypeHistory] = struct{}{}
	}
	if oriSu.Age != su.Age {
		fields[ChangeAgeHistory] = struct{}{}
	}
	if oriSu.Residence != su.Residence {
		fields[ChangeResidenceHistory] = struct{}{}
	}
	if oriSu.IDCard != su.IDCard {
		fields[ChangeIDCardHistory] = struct{}{}
	}
	if oriSu.Phone != su.Phone {
		fields[ChangePhoneHistory] = struct{}{}
	}
	if oriSu.QQ != su.QQ {
		fields[ChangeQQHistory] = struct{}{}
	}
	if oriSu.Wechat != su.Wechat {
		fields[ChangeWechatHistory] = struct{}{}
	}
	if oriSu.EconomicBegin == DateDefualtFromDB {
		oriSu.EconomicBegin = 0
	}
	if oriSu.EconomicEnd == DateDefualtFromDB {
		oriSu.EconomicEnd = 0
	}
	if oriSu.IsEconomic != su.IsEconomic || !strings.EqualFold(oriSu.EconomicCompany, su.EconomicCompany) ||
		oriSu.EconomicBegin != su.EconomicBegin || oriSu.EconomicEnd != su.EconomicEnd {
		fields[ChangeEconomicHistory] = struct{}{}
	}
	if !strings.EqualFold(oriSu.Note, su.Note) {
		fields[ChangeSignNoteHistory] = struct{}{}
	}
	su.State = oriSu.State
	su.EmailState = oriSu.EmailState
	su.TaskState = oriSu.TaskState
	su.LeaveTimes = oriSu.LeaveTimes
	su.ViolationTimes = oriSu.ViolationTimes
	su.ActiveTid = oriSu.ActiveTid
	su.Ctime = oriSu.Ctime
}

// SignPay  struct
type SignPay struct {
	ID         int64
	Mid        int64
	SignID     int64
	DueDate    xtime.Time
	PayValue   int64
	State      int8
	Note       string
	EmailState int8
	Ctime      xtime.Time `gorm:"column:ctime"`
	Mtime      xtime.Time `gorm:"column:mtime"`
	InTax      int8
}

// Diff .
func (sp *SignPay) Diff(mOriSp map[int64]*SignPay, fields map[int8]struct{}) {
	var (
		ok    bool
		oriSp *SignPay
	)
	if oriSp, ok = mOriSp[sp.ID]; !ok {
		fields[ChangeSignPayHistory] = struct{}{}
		return
	}
	if sp.DueDate != oriSp.DueDate || sp.PayValue != oriSp.PayValue || sp.InTax != oriSp.InTax {
		fields[ChangeSignPayHistory] = struct{}{}
	}
	sp.Mid = oriSp.Mid
	sp.SignID = oriSp.SignID
	sp.State = oriSp.State
	sp.Note = oriSp.Note
	sp.EmailState = oriSp.EmailState
	sp.Ctime = oriSp.Ctime
}

// SignTask  struct
type SignTask struct {
	ID            int64      `gorm:"column:id"`
	Mid           int64      `gorm:"column:mid"`
	SignID        int64      `gorm:"column:sign_id"`
	TaskType      int8       `gorm:"column:task_type"`
	TaskCounter   int        `gorm:"column:task_counter"`
	TaskCondition int        `gorm:"column:task_condition"`
	TaskData      string     `gorm:"column:task_data"`
	State         int8       `gorm:"column:state"`
	Ctime         xtime.Time `gorm:"column:ctime"`
	Mtime         xtime.Time `gorm:"column:mtime"`
	Attribute     int64      `gorm:"column:attribute"`
	FinishNote    string     `gorm:"column:finish_note"`
}

// Diff .
func (st *SignTask) Diff(mOriSt map[int64]*SignTask, fields map[int8]struct{}) {
	var (
		ok    bool
		oriSt *SignTask
	)
	if oriSt, ok = mOriSt[st.ID]; !ok {
		fields[ChangeSignTaskHistory] = struct{}{}
		return
	}
	if st.TaskType != oriSt.TaskType || st.TaskCondition != oriSt.TaskCondition ||
		st.AttrVal(AttrBitIsBusinessArchive) != oriSt.AttrVal(AttrBitIsBusinessArchive) {
		fields[ChangeSignTaskHistory] = struct{}{}
	}
	st.Mid = oriSt.Mid
	st.SignID = oriSt.SignID
	st.TaskCounter = oriSt.TaskCounter
	st.TaskData = oriSt.TaskData
	st.State = oriSt.State
	st.Ctime = oriSt.Ctime
}

// AttrVal get attribute value.
func (st *SignTask) AttrVal(bit uint) int64 {
	return (st.Attribute >> bit) & int64(1)
}

// AttrSet set attribute value.
func (st *SignTask) AttrSet(v int64, bit uint) {
	st.Attribute = st.Attribute&(^(1 << bit)) | (v << bit)
}

// SignTaskHistory .
type SignTaskHistory struct {
	ID             int64      `gorm:"column:id"`
	Mid            int64      `gorm:"column:mid"`
	SignID         int64      `gorm:"column:sign_id"`
	TaskTemplateID int        `gorm:"column:task_template_id"`
	TaskType       int8       `gorm:"column:task_type"`
	TaskCounter    int        `gorm:"column:task_counter"`
	TaskCondition  int        `gorm:"column:task_condition"`
	TaskData       string     `gorm:"column:task_data"`
	Attribute      int64      `gorm:"column:attribute"`
	State          int        `gorm:"column:state"`
	GenerateDate   xtime.Time `gorm:"column:generate_date"`
	Ctime          xtime.Time `gorm:"column:ctime"`
	Mtime          xtime.Time `gorm:"column:mtime"`
}

// AttrVal get attribute value.
func (sth *SignTaskHistory) AttrVal(bit uint) int64 {
	return (sth.Attribute >> bit) & int64(1)
}

// AttrSet set attribute value.
func (sth *SignTaskHistory) AttrSet(v int64, bit uint) {
	sth.Attribute = sth.Attribute&(^(1 << bit)) | (v << bit)
}

//SignContract  struct
type SignContract struct {
	ID       int64 `gorm:"column:id"`
	Mid      int64
	SignID   int64
	Filename string
	Filelink string
	State    int8
	Ctime    xtime.Time `gorm:"column:ctime"`
	Mtime    xtime.Time `gorm:"column:mtime"`
}

// Diff .
func (sc *SignContract) Diff(mOriSc map[int64]*SignContract, fields map[int8]struct{}) {
	var (
		ok    bool
		oriSc *SignContract
	)
	if oriSc, ok = mOriSc[sc.ID]; !ok {
		log.Error("OriSc(%d) no exsits", sc.ID)
		fields[ChangeSignContractHistory] = struct{}{}
		return
	}
	if !strings.EqualFold(sc.Filelink, oriSc.Filelink) {
		log.Error("file(%s)----orc_file(%s) no exsits", sc.Filelink, oriSc.Filelink)
		fields[ChangeSignContractHistory] = struct{}{}
	}
	if !strings.EqualFold(sc.Filename, oriSc.Filename) {
		log.Error("filename(%s)----orc_filename(%s) no exsits", sc.Filename, oriSc.Filename)
		fields[ChangeSignContractHistory] = struct{}{}
	}
	sc.Mid = oriSc.Mid
	sc.SignID = oriSc.SignID
	sc.State = oriSc.State
	sc.Ctime = oriSc.Ctime
}

// SignTaskAbsence struct
type SignTaskAbsence struct {
	ID            int64 `gorm:"column:id"`
	SignID        int64
	Mid           int64
	TaskHistoryID int64
	AbsenceCount  int
	Reason        string
	State         int8
	AdminID       int64
	AdminName     string
	Ctime         xtime.Time `gorm:"column:ctime"`
	Mtime         xtime.Time `gorm:"column:mtime"`
}

// SignViolationHistory struct
type SignViolationHistory struct {
	ID              int64 `gorm:"column:id"`
	SignID          int64
	Mid             int64
	AdminID         int64
	AdminName       string
	ViolationReason string
	State           int8
	Ctime           xtime.Time `gorm:"column:ctime"`
	Mtime           xtime.Time `gorm:"column:mtime"`
}

// GetTaskDuration this will return task duration, [startDate, endDate)
func GetTaskDuration(date time.Time, taskType int8) (startDate, endDate time.Time) {
	var ndate = now.New(date)
	now.WeekStartDay = time.Monday
	switch taskType {
	case TaskTypeDay:
		var begin = ndate.BeginningOfDay()
		return begin, begin.AddDate(0, 0, 1)
	case TaskTypeWeek:
		var begin = ndate.BeginningOfWeek()
		return begin, begin.AddDate(0, 0, 7)
	case TaskTypeMonth:
		var begin = ndate.BeginningOfMonth()
		return begin, begin.AddDate(0, 1, 0)
	case TaskTypeQuarter:
		var begin = ndate.BeginningOfQuarter()
		return begin, begin.AddDate(0, 3, 0)
	}
	return
}
