package model

import (
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/aegis/model/net"
)

//..
const (
	LogBusinessAudit   = int(231)
	LogTypeAuditSubmit = int(1)
	LogTypeAuditJump   = int(2)
	LogTypeAuditAdd    = int(3)
	LogTypeAuditCancel = int(4)

	LogBusinessTask     = int(232)
	LogTypeTaskDispatch = int(1)
	LogTypeTaskConsumer = int(2)
	LogTYpeTaskWeight   = int(3)

	LogBusinessResource = int(233)
	LogTypeFormAuditor  = int(1) // 审核员提交变更
	LogTypeFromAdd      = int(2) // 业务添加资源
	LogTypeFromCancle   = int(3) // 业务注销资源

	LogBusinessNet    = int(234) //流程流转
	LogTypeNetTrigger = int(1)   //流程流转
	LogFromBatch      = "批量提交"
	LogFromSingle     = "单个提交"
	LogFromStart      = "启动"
	LogFromJump       = "跳流程"
	LogFromCancle     = "取消流程"

	LogBusinessNetConf   = int(235) //流程配置
	LogTypeNetConf       = int(1)   //流程网配置
	LogTypeTokenConf     = int(2)   //令牌配置
	LogTypeFlowConf      = int(3)   //节点配置
	LogTypeTranConf      = int(4)   //变迁配置
	LogTypeTokenBindConf = int(5)   //令牌绑定配置
	LogTypeDirConf       = int(6)   //有向线配置

	//流程配置的操作
	LogNetActionNew       = "创建"
	LogNetActionUpdate    = "更新"
	LogNetActionDisable   = "禁用"
	LogNetActionAvailable = "启用"
	//流程配置的更新字段
	LogFieldPID       = "pid"
	LogFieldChName    = "ch_name"
	LogFieldName      = "name"
	LogFieldTokenID   = "token_id"
	LogFieldStartFlow = "start_flow"
	LogFieldTrigger   = "trigger"
	LogFieldLimit     = "limit"
	LogFieldDirection = "direction"
	LogFieldOrder     = "order"
	LogFieldGuard     = "guard"
	LogFieldOutput    = "output"
)

//LogFieldDesc 日志中的变更字段描述
var LogFieldDesc = map[string]string{
	LogFieldPID:       "父网",
	LogFieldChName:    "中文名",
	LogFieldName:      "英文名",
	LogFieldTokenID:   "令牌",
	LogFieldStartFlow: "初始节点",
	LogFieldTrigger:   "触发类型",
	LogFieldLimit:     "限制",
	LogFieldDirection: "指向",
	LogFieldOrder:     "顺序",
	LogFieldGuard:     "通行条件",
	LogFieldOutput:    "输出规则",
}

//LogFieldTemp 指定模板和字段，组建变更日志字段
func LogFieldTemp(field string, nw interface{}, old interface{}, diffTemp bool) (s string) {
	var (
		desc  string
		exist bool
		temp  string
	)
	if desc, exist = LogFieldDesc[field]; !exist {
		desc = field
	}
	if diffTemp {
		temp = "[%s]从[%v]变成[%v]"
		s = fmt.Sprintf(temp, desc, old, nw)
	} else {
		temp = "[%s]为[%v]"
		s = fmt.Sprintf(temp, desc, nw)
	}

	return
}

//NetConfOper 流程配置日志
type NetConfOper struct {
	OID    int64    `json:"oid"`
	Action string   `json:"action"`
	UID    int64    `json:"uid"`
	Uname  string   `json:"uname"`
	NetID  int64    `json:"net_id"`
	ChName string   `json:"ch_name"`
	FlowID int64    `json:"flow_id"`
	TranID int64    `json:"tran_id"`
	Diff   []string `json:"diff"`
}

//SearchAuditLogParam 搜索行为日志参数
type SearchAuditLogParam struct {
	RID        []int64
	OID        []string `form:"oid,split" validate:"omitempty,max=30,dive,gt=0"`
	Username   []string `form:"username,split"`
	CtimeFrom  string   `form:"ctime_from"`
	CtimeTo    string   `form:"ctime_to"`
	Ps         int      `form:"ps" validate:"gt=0,max=50" default:"10"`
	Pn         int      `form:"pn" validate:"gt=0" default:"1"`
	TaskID     []int64  `form:"task_id,split" validate:"omitempty,max=30"`
	State      string   `form:"state"`
	BusinessID int64    `form:"business_id" validate:"required,gt=0"`
}

//SearchAuditLog 审核日志记录
type SearchAuditLog struct {
	RID        int64  `json:"rid"`
	OID        string `json:"oid"`
	TaskID     int64  `json:"task_id"`
	State      string `json:"state"`
	Stime      string `json:"stime"`
	UID        int64  `json:"uid"`
	Uname      string `json:"uname"`
	Department string `json:"department"`
	Extra      string `json:"extra"`
}

type SearchLogResult struct {
	Result []struct {
		UID    int64  `json:"uid"`
		Uname  string `json:"uname"`
		OID    int64  `json:"oid"`
		Type   int8   `json:"type"`
		Action string `json:"action"`
		Str0   string `json:"str_0"`
		Str1   string `json:"str_1"`
		Str2   string `json:"str_2"`
		Int0   int64  `json:"int_0"`
		Int1   int64  `json:"int_1"`
		Int2   int64  `json:"int_2"`
		Ctime  string `json:"ctime"`
		Extra  string `json:"extra_data"`
	} `json:"result"`
	Page struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
}

// Change for log
type Change struct {
	Opt  *SubmitOptions `json:"opt"`
	Flow *Flow          `json:"flow"`
}

//GetSubmitOper .
func (change *Change) GetSubmitOper() (flowaction, submitopt string) {
	if change.Flow != nil && change.Flow.ResultToken != nil {
		flowaction = change.Flow.ResultToken.ChName
	}
	if change.Opt != nil {
		if change.Opt.Result != nil {
			if change.Opt.Result.ReasonID != 0 {
				submitopt += fmt.Sprintf(" 理由ID:%d ", change.Opt.Result.ReasonID)
			}
			if change.Opt.Result.RejectReason != "" && change.Opt.Result.RejectReason != "null" {
				submitopt += fmt.Sprintf(" 打回理由:%s ", change.Opt.Result.RejectReason)
			}
			if change.Opt.Result.Note != "" && change.Opt.Result.Note != "null" {
				submitopt += fmt.Sprintf(" 备注:%s ", change.Opt.Result.Note)
			}
		}
		if change.Opt.Forbid != nil {
			if change.Opt.Forbid.Reason != "" && change.Opt.Forbid.Reason != "null" {
				submitopt += fmt.Sprintf(" 封禁理由:%s ", change.Opt.Forbid.Reason)
				if du := change.Opt.Forbid.Duration; du == -1 {
					submitopt += " 封禁时长:永久 "
				} else {
					submitopt += fmt.Sprintf(" 封禁时长:%d天 ", du/86400)
				}
			}
			if change.Opt.Forbid.Comment != "" && change.Opt.Forbid.Comment != "null" {
				submitopt += fmt.Sprintf(" 备注:%s ", change.Opt.Forbid.Comment)
			}
		}
	}
	return
}

// Flow for log
type Flow struct {
	SubmitToken *net.TokenPackage `json:"submit_token"`
	ResultToken *net.TokenPackage `json:"result_token"`
	NewFlowID   int64             `json:"new_flow_id"`
	OldFlowID   json.Number       `json:"old_flow_id"`
}

// WeightLog task log
type WeightLog struct {
	UPtime      string        `json:"uptime"`
	Mid         int64         `json:"mid"`
	MemName     string        `json:"member_name"` // up主名称
	Fans        int64         `json:"fans"`
	FansWeight  int64         `json:"fans_weight"`
	Group       string        `json:"group"`
	GroupWeight int64         `json:"group_weight"`
	WaitTime    string        `json:"wait_time"`
	WaitWeight  int64         `json:"wait_weight"`
	EqualWeight int64         `json:"config_weight"`
	ConfigItems []*ConfigItem `json:"config_items"`
	Weight      int64         `json:"weight"`
}

// ConfigItem .
type ConfigItem struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Username string `json:"uname"`
}

//TrackRsc 资源信息追踪
type TrackRsc struct {
	Ctime   string                 `json:"ctime"`
	Content string                 `json:"content"`
	Detail  map[string]interface{} `json:"detail"`
}

//TrackAudit 操作日志信息追踪
type TrackAudit struct {
	Ctime    string  `json:"ctime"`
	FlowID   []int64 `json:"flow_id"`
	FlowName string  `json:"flow_name"`
	State    string  `json:"state"`
	Uname    string  `json:"uname"`
}

//TrackInfo 信息追踪汇总
type TrackInfo struct {
	Add      []*TrackRsc   `json:"add"`
	Audit    []*TrackAudit `json:"audit"`
	Relation [][]int       `json:"relation"`
}

//TrackParam 信息追踪参数
type TrackParam struct {
	BusinessID   int64  `form:"business_id" validate:"required,gt=0"`
	OID          string `form:"oid" validate:"required,gt=0"`
	Pn           int    `form:"pn" validate:"gt=0" default:"1"`
	Ps           int    `form:"ps" validate:"gt=0,max=30" default:"10"`
	LastPageTime string `form:"last_page_time"`
}

//ParamsQueryLog .
type ParamsQueryLog struct {
	Business  int      `reflect:"business"`
	Type      int      `default:"0" reflect:"type,omitdefault"`
	Uname     []string `reflect:"uname"`
	UID       []int64  `reflect:"uid"`
	Oid       []int64  `reflect:"oid"`
	Action    []string `reflect:"action"`
	CtimeFrom string   `reflect:"ctime,from"`
	CtimeTo   string   `reflect:"ctime,to"`
	Int0      []int64  `reflect:"int_0"`
	Int0From  string   `reflect:"int_0,from"`
	Int0To    string   `reflect:"int_0,to"`
	Int1      []int64  `reflect:"int_1"`
	Int1From  string   `reflect:"int_1,from"`
	Int1To    string   `reflect:"int_1,to"`
	Int2      []int64  `reflect:"int_2"`
	Int2From  string   `reflect:"int_2,from"`
	Int3To    string   `reflect:"int_2,to"`
	Str0      []string `reflect:"str_0"`
	Str1      []string `reflect:"str_1"`
	Str2      []string `reflect:"str_2"`
}

//EsCommon .
type EsCommon struct {
	Pn    int
	Ps    int
	Order string
	Sort  string
	Group string
}
