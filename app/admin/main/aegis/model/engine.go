package model

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/resource"
	"go-common/app/admin/main/aegis/model/task"
	member "go-common/app/service/main/account/api"
	"go-common/library/log"
)

const (
	// DefaultZeroInt form default zero
	DefaultZeroInt = int64(-12345)
)

// EngineOption .
type EngineOption struct {
	common.BaseOptions
	TaskID    int64                  `json:"task_id" form:"task_id"`
	Result    *resource.Result       `json:"resource_result" form:"resource_result" submit:"object"` // 资源的修改内容，要传就全传，否则传空
	Forbid    *Forbid                `json:"forbid_params" submit:"json"`
	ExtraData map[string]interface{} `json:"extra_data" submit:"map"`
}

// Forbid 封禁参数
type Forbid struct {
	Notify   int    `json:"notify"`
	Reason   string `json:"reason"`
	Duration int64  `json:"duration"`
	Comment  string `json:"comment"`
	Img      string `json:"img"`
}

// InfoResOption .
type InfoResOption struct {
	common.BaseOptions
}

// ListResOption .
type ListResOption struct {
	common.Pager
	common.BaseOptions
	BusinessID int64 `form:"business_id"`
	FlowID     int64 `form:"flow_id"`
}

// SubmitOptions 审核结果提交参数
type SubmitOptions struct {
	EngineOption
	NewFlowID int64   `json:"newflow_id"`
	BindStr   string  `json:"binds"`
	Binds     []int64 `json:"-"`
}

// BatchOption 批量提交
type BatchOption struct {
	BusinessID   int64   `form:"business_id" json:"business_id" validate:"required"`
	NetID        int64   `form:"net_id" json:"net_id"`
	UID          int64   `form:"uid" json:"uid"`
	Debug        int8    `form:"debug" json:"debug"`
	Uname        string  `form:"uname" json:"uname"`
	RIDs         []int64 `form:"rids,split" json:"rids" validate:"gt=0,dive,gt=0"`
	Binds        []int64 `form:"binds,split" json:"binds" validate:"gt=0,dive,gt=0"`
	RejectReason string  `form:"reject_reason"`
	ReasonID     int64   `form:"reason_id" default:"0"`
	Notify       int8    `form:"notify"`
}

// Tip 批量提交的错误提示
type Tip struct {
	Success []int64          `json:"success"`
	Async   []int64          `json:"async"`
	Fail    map[int64]string `json:"fail"`
}

// AddOption add option
type AddOption struct {
	resource.Resource
	State int   `form:"state" json:"state"`
	NetID int64 `form:"net_id" json:"net_id" validate:"required"`
}

// UpdateKeys 可更新的资源字段
var UpdateKeys = map[string]struct{}{
	"mid":        {},
	"content":    {},
	"extra1":     {},
	"extra2":     {},
	"extra3":     {},
	"extra4":     {},
	"extra5":     {},
	"extra6":     {},
	"extra1s":    {},
	"extra2s":    {},
	"extra3s":    {},
	"extra4s":    {},
	"extratime1": {},
	"octime":     {},
	"ptime":      {},
	"metadata":   {},
}

// UpdateOption update option
type UpdateOption struct {
	BusinessID int64                  `form:"business_id"  json:"business_id" validate:"required"`
	NetID      int64                  `form:"net_id" json:"net_id" validate:"required"`
	OID        string                 `form:"oid" json:"oid" validate:"required"`
	Update     map[string]interface{} `json:"update"`
}

// AuditInfo 审核详情页
type AuditInfo struct {
	UnDoStat *task.UnDOStat `json:"undo_stat"`

	Task *task.Task `json:"task"`

	Flow *net.TransitionInfo `json:"flow"`

	Resource *resource.Res `json:"resource"`

	OperHistorys []string `json:"historys"`

	IFrame string `json:"iframe_url,omitempty"`

	UserInfo  *UserInfo               `json:"user_info"`
	UserGroup map[int64]*common.Group `json:"user_group"`
	Hit       []string                `json:"hit"`

	Actions []*Action `json:"actions"`
}

// UserInfo .
type UserInfo struct {
	Mid      int64               `json:"mid"`
	Name     string              `json:"name"`
	Official member.OfficialInfo `json:"official"`
	Follower int64               `json:"follower"`
}

// SyncMessage .
type SyncMessage struct {
	URL          string
	Uname        string                 `json:"auditor"`
	MID          int64                  `json:"uid"`
	OID          string                 `json:"oid"`
	RejectReason string                 `json:"reason"`
	FlowData     map[string]interface{} `json:"flow_data"`
	ExtraData    map[string]interface{} `json:"extra_data"`
	ForbidParams map[string]interface{} `json:"forbid_params"`
	Attribute    map[string]int8        `json:"attribute"`
}

// GetEmptyInfo empty info for debug
func GetEmptyInfo() *AuditInfo {
	info := &AuditInfo{
		UnDoStat: &task.UnDOStat{},
		Task:     &task.Task{},
		Flow: &net.TransitionInfo{
			Operations: []*net.TranOperation{{}},
		},
		UserInfo: &UserInfo{},
		Resource: &resource.Res{},
		OperHistorys: []string{
			"wulalalal",
		},
		IFrame: "iframe",
	}
	return info
}

// SearchParams ..
type SearchParams struct {
	common.Pager
	Debug      int8     `form:"debug" reflect:"ignore"`
	FilterOff  bool     `form:"state_filter_off" reflect:"ignore"`
	BusinessID int64    `form:"business_id" validate:"required" reflect:"business_id"`
	OID        []string `form:"oid,split" reflect:"oid"`
	FlowID     int64    `form:"flow_id" default:"-12345" reflect:"flow_id,omitdefault"`
	State      int64    `form:"state" default:"-12345" reflect:"state,omitdefault"`
	Mid        int64    `form:"mid" default:"-12345" reflect:"mid,omitdefault"`
	CtimeFrom  string   `form:"ctime_from" reflect:"ctime,from"`
	CtimeTo    string   `form:"ctime_to" reflect:"ctime,to"`
	Extra1     string   `form:"extra1" reflect:"extra1" parse:"[]int"`
	Extra2     string   `form:"extra2" reflect:"extra2" parse:"int"`
	Extra3     string   `form:"extra3" reflect:"extra3" parse:"int"`
	Extra4     string   `form:"extra4" reflect:"extra4" parse:"int"`
	Extra5     string   `form:"extra4" reflect:"extra5" parse:"int"`
	Extra6     string   `form:"extra4" reflect:"extra6" parse:"int"`
	Extra1s    string   `form:"extra1s" reflect:"extra1s" `
	Extra2s    string   `form:"extra2s" reflect:"extra2s" `
	Extra3s    string   `form:"extra3s" reflect:"extra3s" `
	Extra4s    string   `form:"extra4s" reflect:"extra4s" `
	ExtraTime1 string   `form:"extratime1" reflect:"extratime1" ` //TODO 具体用到了再看怎么解析
	KeyWord    string   `form:"keyWords" reflect:"ignore"`
	CtimeOrder string   `form:"ctime_order" reflect:"ignore"` //TODO 等octime导入后要按照octime筛选
}

// SearchRes search list.
type SearchRes struct {
	Resources []*ListRscItem `json:"result"`
	Page      struct {
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"page"`
}

//Column .
type Column struct {
	Name   string           `json:"name"`
	ChName string           `json:"chname"`
	Enum   map[int64]string `json:"enum,omitempty"`
}

// ListRscItem .
type ListRscItem struct {
	resource.Res
	FlowID    int64                   `json:"flow_id"`
	UserInfo  *UserInfo               `json:"user_info"`
	UserGroup map[int64]*common.Group `json:"user_group"`
	Hit       []string                `json:"hit"`
}

// ListTaskItem .
type ListTaskItem struct {
	ListTask
	WaitTime  string                  `json:"wait_time"`
	OID       string                  `json:"oid"`
	Content   string                  `json:"content"`
	UserInfo  *UserInfo               `json:"user_info"`
	UserGroup map[int64]*common.Group `json:"user_group"`
	Metas     map[string]interface{}  `json:"metas"`
}

// ListTask 转化gtime
type ListTask struct {
	*task.Task
	GTstr    string `json:"gtime"`
	CTstr    string `json:"ctime"`
	MTstr    string `json:"mtime"`
	UserName string `json:"uid"`
	MidStr   string `json:"mid"`
}

// EmptyListItem .
func EmptyListItem() *ListRscItem {
	return &ListRscItem{
		UserInfo: &UserInfo{},
	}
}

// IsImgType file type  is image
func IsImgType(fileType string) bool {
	return fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp"
}

//Action 下发到前端的组件
type Action struct {
	Name     string                  `json:"name"`
	URL      string                  `json:"url"`
	Method   string                  `json:"method"`
	Response string                  `json:"response"`
	Extra    string                  `json:"extra,omitempty"`
	Params   map[string]*ActionParam `json:"params"`
}

//ActionParam 烦人啊，动态的打回理由是必传参数
type ActionParam struct {
	Value   string `json:"value"`
	Default string `json:"default"`
}

//SubReflect .
func SubReflect(ot reflect.Type, ov reflect.Value, key string, vals []string, defaultval string, params map[string]interface{}) {
	var result string
	boommsg := fmt.Sprintf("key:%s vals:%v ot:%+v vt:%+v", key, vals, ot, reflect.TypeOf(ov).Kind())

	if len(vals) == 0 {
		return
	}
	val := vals[0]
	tfield, ok := ot.FieldByName(val)
	if !ok {
		log.Error("SubReflect1:" + boommsg)
		return
	}
	vfield := ov.FieldByName(val)
	if !vfield.IsValid() {
		log.Error("SubReflect7:" + boommsg)
		return
	}

	switch tfield.Tag.Get("submit") {
	case "object":
		vals = vals[1:]
		tt := tfield.Type
		vv := reflect.ValueOf(vfield)
		if tfield.Type.Kind() == reflect.Ptr {
			tt = tfield.Type.Elem()
			vv = vfield.Elem()
		}
		SubReflect(tt, vv, key, vals, defaultval, params)
		return
	case "map":
		if len(vals) != 2 || tfield.Type.Kind() != reflect.Map {
			log.Error("SubReflect2:" + boommsg)
		}

		v := vfield.MapIndex(reflect.ValueOf(vals[1]))
		if !v.IsValid() { //批量操作，可能会没有map的参数
			log.Warn("SubReflect4:" + boommsg)
			return
		}
		result = fmt.Sprint(v.Interface())
	case "json":
		bs, err := json.Marshal(vfield.Interface())
		if err != nil {
			log.Error("SubReflect5:" + boommsg)
			return
		}
		result = string(bs)
	case "int", "string":
		result = fmt.Sprint(vfield.Interface())
	default:
		log.Error("SubReflect6:" + boommsg)
	}
	if result == "" {
		result = defaultval
	}
	params[key] = result
}

//Auth auth
type Auth struct {
	OK       bool            `json:"ok"`
	Admin    bool            `json:"admin"`
	Business map[int64]int64 `json:"business"`
}

//UpsertItem 更新搜索的items
type UpsertItem struct {
	ID     int64 `json:"id"`
	State  int   `json:"state"`
	Extra1 int64 `json:"extra1"`
	Extra2 int64 `json:"extra2"`
	Extra3 int64 `json:"extra3"`
	Extra4 int64 `json:"extra4"`
}

//CancelOption .
type CancelOption struct {
	BusinessID int64    `form:"business_id" validate:"required" json:"business_id"`
	Oids       []string `form:"oids,split" validate:"required" json:"oids"`
	Reason     string   `form:"reason" json:"reason"`
}
