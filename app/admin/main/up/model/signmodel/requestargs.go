package signmodel

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/dao"
	xtime "go-common/library/time"
)

// const .
const (
	SignUpList   = 0
	SignUpDetail = 1
)

// -------------------------------------------------

// CommonResponse  result
type CommonResponse struct {
}

// CommonArg  arg
type CommonArg struct {
}

// SignUpBaseInfo struct
type SignUpBaseInfo struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Sex             int8       `json:"sex"`
	Mid             int64      `json:"mid"`
	BeginDate       xtime.Time `json:"begin_date"`
	EndDate         xtime.Time `json:"end_date"`
	State           int8       `json:"state"`
	Country         string     `json:"country"`
	Province        string     `json:"province"`
	City            string     `json:"city"`
	Note            string     `json:"note"`
	TypeName        string     `json:"type_name"`
	ActiveTid       int16      `json:"active_tid"`
	AdminID         int        `json:"admin_id"`
	AdminName       string     `json:"admin_name"`
	CreateTime      xtime.Time `json:"create_time"`
	Organization    int8       `json:"organization"`
	SignType        int8       `json:"sign_type"`
	Age             int8       `json:"age"`
	Residence       string     `json:"residence"`
	IDCard          string     `json:"id_card"`
	Phone           string     `json:"phone"`
	QQ              int64      `json:"qq"`
	Wechat          string     `json:"wechat"`
	IsEconomic      int8       `json:"is_economic"`
	EconomicCompany string     `json:"economic_company"`
	EconomicBegin   xtime.Time `json:"economic_begin"`
	EconomicEnd     xtime.Time `json:"economic_end"`
	ViolationTimes  int        `json:"violation_times"`
	LeaveTimes      int        `json:"leave_times"`
}

// CopyTo copy
func (s *SignUpBaseInfo) CopyTo(dbstruct *SignUp) {
	dbstruct.ID = s.ID
	dbstruct.Mid = s.Mid
	dbstruct.Sex = s.Sex
	dbstruct.BeginDate = s.BeginDate
	dbstruct.EndDate = s.EndDate
	dbstruct.Country = s.Country
	dbstruct.Province = s.Province
	dbstruct.City = s.City
	dbstruct.Note = s.Note
	dbstruct.AdminID = s.AdminID
	dbstruct.AdminName = s.AdminName
	dbstruct.Organization = s.Organization
	dbstruct.SignType = s.SignType
	dbstruct.Age = s.Age
	dbstruct.Residence = s.Residence
	dbstruct.IDCard = s.IDCard
	dbstruct.Phone = s.Phone
	dbstruct.QQ = s.QQ
	dbstruct.Wechat = s.Wechat
	dbstruct.IsEconomic = s.IsEconomic
	if s.IsEconomic == ContainEconomic {
		dbstruct.EconomicCompany = s.EconomicCompany
		dbstruct.EconomicBegin = s.EconomicBegin
		dbstruct.EconomicEnd = s.EconomicEnd
	}
}

// CopyFrom copy
func (s *SignUpBaseInfo) CopyFrom(dbstruct *SignUp) {
	s.ID = dbstruct.ID
	s.Mid = dbstruct.Mid
	s.Sex = dbstruct.Sex
	s.BeginDate = dbstruct.BeginDate
	s.EndDate = dbstruct.EndDate
	s.Country = dbstruct.Country
	s.Province = dbstruct.Province
	s.City = dbstruct.City
	s.Note = dbstruct.Note
	s.AdminID = dbstruct.AdminID
	s.AdminName = dbstruct.AdminName
	s.Organization = dbstruct.Organization
	s.SignType = dbstruct.SignType
	s.Age = dbstruct.Age
	s.Residence = dbstruct.Residence
	s.IDCard = dbstruct.IDCard
	s.Phone = dbstruct.Phone
	s.QQ = dbstruct.QQ
	s.Wechat = dbstruct.Wechat
	s.IsEconomic = dbstruct.IsEconomic
	if dbstruct.IsEconomic == ContainEconomic {
		s.EconomicCompany = dbstruct.EconomicCompany
		s.EconomicBegin = dbstruct.EconomicBegin
		s.EconomicEnd = dbstruct.EconomicEnd
	}
	s.State = dbstruct.State
	s.ActiveTid = dbstruct.ActiveTid
	s.LeaveTimes = dbstruct.LeaveTimes
	s.ViolationTimes = dbstruct.ViolationTimes
	s.CreateTime = dbstruct.Ctime
}

// CopyFrom .
func (t *SignTaskHistoryArg) CopyFrom(st *SignTaskHistory, absenceCounter int) {
	var taskBegin, taskEnd time.Time
	taskBegin, taskEnd = GetTaskDuration(st.GenerateDate.Time(), st.TaskType)
	t.TaskBegin = xtime.Time(taskBegin.Unix())
	t.TaskEnd = xtime.Time(taskEnd.Unix())
	t.TaskType = st.TaskType
	t.TaskState = int8(st.State)
	t.TaskCounter = st.TaskCounter
	t.TaskCondition = st.TaskCondition
	t.AbsenceCounter = absenceCounter
	t.IsBusinessArchive = st.AttrVal(AttrBitIsBusinessArchive)
}

// SignUpArg struct
type SignUpArg struct {
	SignUpBaseInfo
	PayInfo      []*SignPayInfoArg      `json:"pay_info"`
	TaskInfo     []*SignTaskInfoArg     `json:"task_info"`
	ContractInfo []*SignContractInfoArg `json:"contract_info"`
}

// SignUpsArg struct
type SignUpsArg struct {
	SignUpBaseInfo
	TaskHistoryInfo []*SignTaskHistoryArg  `json:"task_history_info"`
	PayInfo         []*SignPayInfoArg      `json:"pay_info"`
	ContractInfo    []*SignContractInfoArg `json:"contract_info"`
}

// SignQueryResult  result
type SignQueryResult struct {
	SignBaseInfo *SignUpBaseInfo `json:"sign_base_info"`
	Result       []*SignUpsArg   `json:"result"`
	TotalCount   int             `json:"total_count"`
	Page         int             `json:"page"`
	Size         int             `json:"size"`
}

// SignTaskHistoryArg struct
type SignTaskHistoryArg struct {
	TaskBegin         xtime.Time `json:"task_begin"`
	TaskEnd           xtime.Time `json:"task_end"`
	TaskType          int8       `json:"task_type"`
	TaskState         int8       `json:"task_state"`
	TaskCounter       int        `json:"task_counter"`
	TaskCondition     int        `json:"task_condition"`
	AbsenceCounter    int        `json:"absence_counter"`
	IsBusinessArchive int64      `json:"is_business_archive"`
}

// ViolationArg struct
type ViolationArg struct {
	ID              int64      `json:"id"`
	SignID          int64      `json:"sign_id"`
	Mid             int64      `json:"mid"`
	ViolationReason string     `json:"violation_reason"`
	AdminID         int64      `json:"admin_id"`
	AdminName       string     `json:"admin_name"`
	OpTime          xtime.Time `json:"op_time"`
	State           int8       `json:"state"`
}

// ViolationResult  result
type ViolationResult struct {
	Result     []*ViolationArg `json:"result"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	Size       int             `json:"size"`
}

// CopyTo  violationArg.
func (a *ViolationArg) CopyTo(v *SignViolationHistory) {
	v.SignID = a.SignID
	v.Mid = a.Mid
	v.AdminID = a.AdminID
	v.AdminName = a.AdminName
	v.ViolationReason = a.ViolationReason
}

// CopyFrom violationArg.
func (a *ViolationArg) CopyFrom(v *SignViolationHistory) {
	a.ID = v.ID
	a.SignID = v.SignID
	a.Mid = v.Mid
	a.AdminID = v.AdminID
	a.AdminName = v.AdminName
	a.ViolationReason = v.ViolationReason
	a.OpTime = v.Mtime
	a.State = v.State
}

// AbsenceArg struct
type AbsenceArg struct {
	ID           int64      `json:"id"`
	SignID       int64      `json:"sign_id"`
	Mid          int64      `json:"mid"`
	AbsenceCount int        `json:"absence_count"`
	Reason       string     `json:"reason"`
	AdminID      int64      `json:"admin_id"`
	AdminName    string     `json:"admin_name"`
	OpTime       xtime.Time `json:"op_time"`
	TaskBegin    xtime.Time `json:"task_begin"`
	TaskEnd      xtime.Time `json:"task_end"`
	State        int8       `json:"state"`
}

// CopyTo  AbsenceArg.
func (a *AbsenceArg) CopyTo(v *SignTaskAbsence) {
	v.SignID = a.SignID
	v.Mid = a.Mid
	v.AdminID = a.AdminID
	v.AdminName = a.AdminName
	v.AbsenceCount = a.AbsenceCount
	v.Reason = a.Reason
}

// CopyFrom AbsenceArg.
func (a *AbsenceArg) CopyFrom(v *SignTaskAbsence) {
	a.ID = v.ID
	a.SignID = v.SignID
	a.Mid = v.Mid
	a.AdminID = v.AdminID
	a.AdminName = v.AdminName
	a.AbsenceCount = v.AbsenceCount
	a.Reason = v.Reason
	a.State = v.State
	a.OpTime = v.Mtime
}

// AbsenceResult  result
type AbsenceResult struct {
	Result     []*AbsenceArg `json:"result"`
	TotalCount int           `json:"total_count"`
	Page       int           `json:"page"`
	Size       int           `json:"size"`
}

// PageArg .
type PageArg struct {
	SignID int64 `form:"sign_id"`
	Page   int   `form:"page"`
	Size   int   `form:"size"`
}

// IDArg .
type IDArg struct {
	ID        int64 `json:"id"`
	SignID    int64 `json:"sign_id"`
	AdminID   int64
	AdminName string
}

// PowerCheckArg .
type PowerCheckArg struct {
	TIDs []int16 `form:"tids,split"`
	Mid  int64   `form:"mid"`
}

// PowerCheckReply .
type PowerCheckReply struct {
	IsPower bool `json:"is_power"`
	IsSign  bool `json:"is_sign"`
}

// SignPayInfoArg =============
type SignPayInfoArg struct {
	ID       int64      `json:"id"`
	SignID   int64      `json:"sign_id"`
	Mid      int64      `json:"mid"`
	DueDate  xtime.Time `json:"due_date"`
	PayValue int64      `json:"pay_value"`
	Note     string     `json:"note"`
	State    int8       `json:"state"`
	InTax    int8       `json:"in_tax"`
}

// CopyTo copy
func (s *SignPayInfoArg) CopyTo(dbstruct *SignPay) {
	dbstruct.ID = s.ID
	dbstruct.SignID = s.SignID
	dbstruct.Mid = s.Mid
	dbstruct.DueDate = s.DueDate
	dbstruct.PayValue = s.PayValue
	dbstruct.Note = s.Note
	dbstruct.InTax = s.InTax
}

// CopyFrom copy
func (s *SignPayInfoArg) CopyFrom(dbstruct *SignPay) {
	s.ID = dbstruct.ID
	s.SignID = dbstruct.SignID
	s.Mid = dbstruct.Mid
	s.DueDate = dbstruct.DueDate
	s.PayValue = dbstruct.PayValue
	s.Note = dbstruct.Note
	s.State = dbstruct.State
	s.InTax = dbstruct.InTax
}

// SignTaskInfoArg =============
type SignTaskInfoArg struct {
	ID                int64 `json:"id"`
	SignID            int64 `json:"sign_id"`
	Mid               int64 `json:"mid"`
	TaskType          int8  `json:"task_type"`
	TaskCondition     int   `json:"task_condition"`
	TaskCounter       int   `json:"task_counter"`
	TaskState         int8  `json:"task_state"`
	IsBusinessArchive int64 `json:"is_business_archive"`
}

// CopyTo copy
func (s *SignTaskInfoArg) CopyTo(dbstruct *SignTask) {
	dbstruct.ID = s.ID
	dbstruct.SignID = s.SignID
	dbstruct.Mid = s.Mid
	dbstruct.TaskType = s.TaskType
	dbstruct.TaskCondition = s.TaskCondition
	dbstruct.AttrSet(s.IsBusinessArchive, AttrBitIsBusinessArchive)
	dbstruct.TaskData = ""
}

// CopyFrom copy
func (s *SignTaskInfoArg) CopyFrom(dbstruct *SignTask) {
	s.ID = dbstruct.ID
	s.SignID = dbstruct.SignID
	s.Mid = dbstruct.Mid
	s.TaskType = dbstruct.TaskType
	s.TaskCondition = dbstruct.TaskCondition
	s.TaskCounter = dbstruct.TaskCounter
	s.IsBusinessArchive = dbstruct.AttrVal(AttrBitIsBusinessArchive)
	if s.TaskCounter >= s.TaskCondition {
		s.TaskState = 1
	}
}

// SignContractInfoArg =============
type SignContractInfoArg struct {
	ID       int64  `json:"id"`
	SignID   int64  `json:"sign_id"`
	Mid      int64  `json:"mid"`
	Filename string `json:"filename"`
	Filelink string `json:"filelink"`
}

// CopyTo copy
func (s *SignContractInfoArg) CopyTo(dbstruct *SignContract) {
	dbstruct.ID = s.ID
	dbstruct.SignID = s.SignID
	dbstruct.Mid = s.Mid
	dbstruct.Filelink = BuildOrcBfsURL(s.Filelink)
	dbstruct.Filename = s.Filename
}

// CopyFrom copy
func (s *SignContractInfoArg) CopyFrom(dbstruct *SignContract) {
	s.ID = dbstruct.ID
	s.SignID = dbstruct.SignID
	s.Mid = dbstruct.Mid
	s.Filelink = BuildDownloadURL(dbstruct.Filename, dbstruct.Filelink)
	s.Filename = dbstruct.Filename
}

// SignQueryArg =============
type SignQueryArg struct {
	Tids       []int64    `form:"tids,split"` // 权限tid
	Mid        int64      `form:"mid"`
	DueSign    int8       `form:"due_sign"`    // 签约即将过期
	DuePay     int8       `form:"due_pay"`     // 支付周期即将过期
	ExpireSign int8       `form:"expire_sign"` // 签约已过期
	Sex        int8       `form:"sex" default:"-1"`
	Country    []string   `form:"country,split"`
	ActiveTID  int16      `form:"active_tid"`
	SignType   int8       `form:"sign_type"`
	TaskState  int8       `form:"task_state"`
	SignBegin  xtime.Time `form:"sign_begin"`
	SignEnd    xtime.Time `form:"sign_end"`
	IsDetail   int8       `form:"is_detail"` // 是否详情
	Page       int        `form:"page"`
	Size       int        `form:"size"`
}

// SignIDArg .
type SignIDArg struct {
	ID int64 `form:"id" validate:"required"`
}

// SignPayCompleteArg ==============
type SignPayCompleteArg struct {
	IDs []int64 `json:"ids"`
}

// SignPayCompleteResult  result
type SignPayCompleteResult struct {
}

// SignCheckTaskArg ==============
type SignCheckTaskArg struct {
	Date string `form:"date"`
}

// SignCheckExsitArg ==============
type SignCheckExsitArg struct {
	Mid int64 `form:"mid"`
}

// SignCheckExsitResult  result
type SignCheckExsitResult struct {
	Exist bool `json:"exist"`
}

// SignOpSearchArg .
type SignOpSearchArg struct {
	Mid    int64  `form:"mid"`
	OpID   int64  `form:"oper_id"` // 操作人
	SignID int64  `form:"sign_id"`
	Tp     int8   `form:"type" default:"2"` // 操作类型 1:新增 2:修改
	Order  string `form:"order" default:"ctime"`
	Sort   string `form:"sort"  default:"desc"`
	PN     int    `form:"page" default:"1"`
	PS     int    `form:"size" default:"50"`
}

// BaseAuditReply .
type BaseAuditReply struct {
	CTime     string `json:"ctime"`
	IntOne    int64  `json:"int_0"`
	OID       int64  `json:"oid"`
	Tp        int8   `json:"type"`
	UID       int64  `json:"uid"`
	UName     string `json:"uname"`
	ExtraData string `json:"extra_data"`
}

// BaseAuditListReply .
type BaseAuditListReply struct {
	Order  string            `json:"order"`
	Sort   string            `json:"sort"`
	Pager  *pager            `json:"page"`
	Result []*BaseAuditReply `json:"result"`
}

type pager struct {
	Page       int `json:"num"`
	Size       int `json:"size"`
	TotalCount int `json:"total"`
}

// SignAuditReply .
type SignAuditReply struct {
	CTime    xtime.Time        `json:"ctime"`
	SignID   int64             `json:"sign_id"`
	Mid      int64             `json:"mid"`
	Tp       int8              `json:"type"`
	OperID   int64             `json:"oper_id"`
	OperName string            `json:"oper_name"`
	Content  *SignContentReply `json:"content"`
}

// SignAuditListReply .
type SignAuditListReply struct {
	Order      string            `json:"order"`
	Sort       string            `json:"sort"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	TotalCount int               `json:"total_count"`
	Result     []*SignAuditReply `json:"result"`
}

// SignContentReply .
type SignContentReply struct {
	New        *SignUpArg `json:"new"`
	Old        *SignUpArg `json:"old"`
	ChangeType []int8     `json:"change_type"`
}

// SignCountrysReply .
type SignCountrysReply struct {
	List []string `json:"list"`
}

// SignTidsReply .
type SignTidsReply struct {
	List map[int64]string `json:"list"`
}

// BuildOrcBfsURL orc bfs url.
func BuildOrcBfsURL(raw string) string {
	if raw == "" {
		return ""
	}
	ori, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	values := ori.Query()
	values.Del("token")
	ori.RawQuery = values.Encode()
	return ori.String()
}

// BuildDownloadURL .
func BuildDownloadURL(fileName string, bfsurl string) (finalurl string) {
	var bfsConf = conf.Conf.BfsConf
	var index = strings.LastIndex(bfsurl, "/")
	if index >= 0 && index+1 < len(bfsurl) {
		fileName = bfsurl[index+1:]
	}
	finalurl = fmt.Sprintf("%s?token=%s", bfsurl, url.QueryEscape(dao.Authorize(bfsConf.Key, bfsConf.Secret, "GET", bfsConf.Bucket, fileName, time.Now().Unix())))
	return
}
