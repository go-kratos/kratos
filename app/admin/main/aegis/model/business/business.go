package business

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// StateEnable 正常
	StateEnable = int8(0)
	// StateDisable 删除
	StateDisable = int8(1)

	// TypeIframe .
	TypeIframe = int8(0)
	// TypeSubmit 提交接口配置
	TypeSubmit = int8(1)
	// TypeAction 预留给前端组件配置
	TypeAction = int8(2)
	// TypeManagerBID 角色配置
	TypeManagerBID = int8(3)
	// TypeReverse 保留字配置
	TypeReverse = int8(4)
	// TypeAttribute 属性位配置
	TypeAttribute = int8(5)
	// TypeDeleteState 删除状态配置
	TypeDeleteState = int8(6)
	// TypeFiler 敏感词过滤
	TypeFiler = int8(7)
	// TypeBizBID 业务级别用户角色
	TypeBizBID = int8(8)
	// TypeAdapter 适配器配置
	TypeAdapter = int8(9)
	// TypeRscListAdapter 资源列表适配器
	TypeRscListAdapter = int8(10)
	// 中间件配置，用于前后端交互的自定义数据逻辑处理
	TypeMiddleware = int8(11)
	// TypeTempCodes 临时错误码
	TypeTempCodes = int8(12)
	// TypeCallback TODO 提交接口v2版本，为了兼容线上的TypeSubmit。待v2上线后，再去掉该配置，改为TypeSubmit
	TypeCallback = int8(13)

	// MngBIDMID 任务用户角色: {"manager_id":1, "flow_id":12}
	MngBIDMID = "manager_id"
	// MngBIDFlow .
	MngBIDFlow = "flow_id"
	// BizBIDMngID 业务用户角色: {"manager_id":1, "admin":12}
	BizBIDMngID = "manager_id"
	// BizBIDAdmin .
	BizBIDAdmin = "admin"
	// AccessBiz .
	AccessBiz = "access_biz"
	// AccessFlow .
	AccessFlow = "access_flow"
)

// Business .
type Business struct {
	ID        int64      `json:"id" gorm:"primary_key" form:"id"`
	TP        int64      `json:"type" gorm:"column:type" form:"type"`
	Name      string     `json:"name" gorm:"column:name" form:"name" validate:"max=10"`
	Desc      string     `json:"desc" gorm:"column:desc" form:"desc" validate:"max=160"`
	Developer string     `json:"developer" gorm:"column:developer" form:"developer"`
	UID       int64      `json:"uid" gorm:"column:uid"`
	Attribute int64      `json:"attribute" gorm:"column:attribute" form:"attribute"`
	State     int64      `json:"state" gorm:"column:state" form:"state"`
	UserName  string     `json:"username" gorm:"-"`
	Ctime     xtime.Time `json:"ctime" gorm:"column:ctime"`
	Mtime     xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// IAdapter .
type IAdapter interface {
	AdaptString(val string) string
	AdaptInt64(val int64) int64
}

// Adapter 业务适配器
type Adapter struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Value   string `json:"value"`
	Compare string `json:"compare"`
}

// AdaptString .
func (a Adapter) AdaptString(val string) string {
	return val
}

// AdaptInt .
func (a Adapter) AdaptInt(val int64) int64 {
	switch a.Rule {
	case "not":
		if strconv.FormatInt(val, 10) != a.Compare {
			if v, e := strconv.ParseInt(a.Value, 10, 64); e != nil {
				log.Error("AdaptInt strconv.ParseInt(%s)", a.Value)
			} else {
				val = v
			}
		}
	default:
		log.Error("AdaptInt unsupported rule(%s)", a.Rule)
	}
	return val
}

// AdaptAddOpt ..
func AdaptAddOpt(opt *model.AddOption, adps []*Adapter) {
	if len(adps) == 0 || opt == nil {
		return
	}

	log.Info("AdaptAddOpt before opt(%+v)", opt)
	v := reflect.ValueOf(opt).Elem()
	for _, adp := range adps {
		val := v.FieldByName(adp.Name)
		if !val.CanSet() {
			log.Error("AdaptAddOpt field(%s) can't set", adp.Name)
			return
		}

		switch val.Kind() {
		case reflect.String:
			val.SetString(adp.AdaptString(val.String()))
		case reflect.Int8, reflect.Int64, reflect.Int:
			val.SetInt(adp.AdaptInt(val.Int()))
		default:
			log.Error("AdaptAddOpt unsupported adp(%+v) type(%v)", adp, val.Kind())
		}
	}
	log.Info("AdaptAddOpt after opt(%+v)", opt)
}

// AdaptUpdateOpt ..
func AdaptUpdateOpt(opt *model.UpdateOption, adps []*Adapter) {
	if len(adps) == 0 || opt == nil {
		return
	}

	log.Info("AdaptUpdateOpt before opt(%+v)", opt)
	for _, adp := range adps {
		key := strings.ToLower(adp.Name)
		if iv, ok := opt.Update[key]; ok {
			switch reflect.TypeOf(iv).Kind() {
			case reflect.String:
				opt.Update[key] = adp.AdaptString(iv.(string))
			case reflect.Int8, reflect.Int64, reflect.Int:
				intv, err := strconv.Atoi(fmt.Sprint(iv))
				if err != nil {
					log.Error("AdaptUpdateOpt unsupported adp(%+v)", adp)
					continue
				}
				opt.Update[key] = adp.AdaptInt(int64(intv))
			default:
				log.Error("AdaptUpdateOpt unsupported adp(%+v)", adp)
			}
		}
	}

	log.Info("AdaptUpdateOpt after opt(%+v)", opt)
}

// TableName .
func (t *Business) TableName() string {
	return "business"
}

// OptList .
type OptList struct {
	common.Pager
	TP int8 `form:"type"`
}

// ListBusiness .
type ListBusiness struct {
	common.Pager
	Business []*Business `json:"business"`
}

// BizItem .
type BizItem struct {
	BizID   int64            `json:"business_id"`
	BizName string           `json:"business_name"`
	BizType int64            `json:"biz_type"`
	Flows   map[int64]string `json:"flows"`
}

// FlowItem .
type FlowItem struct {
	FlowID int64  `json:"flow_id"`
	NetID  int64  `json:"net_id"`
	Name   string `json:"name"`
	ChName string `json:"ch_name"`
}

// BizItemArr .
type BizItemArr []*BizItem

func (a BizItemArr) Len() int {
	return len(a)
}

func (a BizItemArr) Less(i, j int) bool {
	return a[i].BizID < a[j].BizID
}

func (a BizItemArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
