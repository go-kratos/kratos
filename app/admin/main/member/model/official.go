package model

import (
	"encoding/json"

	xtime "go-common/library/time"
)

// official role const.
const (
	OfficialRoleUnauth = iota
	OfficialRoleUp
	OfficialRoleIdentify
	OfficialRoleBusiness
	OfficialRoleGov
	OfficialRoleMedia
	OfficialRoleOther
)

// OfficialRoleName official role name.
func OfficialRoleName(role int8) string {
	switch role {
	case OfficialRoleUnauth:
		return "未认证"
	case OfficialRoleUp:
		return "UP主认证"
	case OfficialRoleIdentify:
		return "身份认证"
	case OfficialRoleBusiness:
		return "企业认证"
	case OfficialRoleGov:
		return "政府认证"
	case OfficialRoleMedia:
		return "媒体认证"
	case OfficialRoleOther:
		return "其他认证"
	}
	return ""
}

// official state const.
const (
	OfficialStateWait = iota
	OfficialStatePass
	OfficialStateNoPass
	OfficialStateReWait
)

// OfficialStateName official state name.
func OfficialStateName(state int8) string {
	switch state {
	case OfficialStateWait:
		return "待审核"
	case OfficialStatePass:
		return "审核通过"
	case OfficialStateNoPass:
		return "审核未通过"
	case OfficialStateReWait:
		return "待重新审核"
	}
	return ""
}

// all
var (
	AllRoles = []int8{
		OfficialRoleUnauth,
		OfficialRoleUp,
		OfficialRoleIdentify,
		OfficialRoleBusiness,
		OfficialRoleGov,
		OfficialRoleMedia,
		OfficialRoleOther,
	}
	AllStates = []int8{
		OfficialStateWait,
		OfficialStatePass,
		OfficialStateNoPass,
		OfficialStateReWait,
	}
)

// Official is.
type Official struct {
	ID    int64      `json:"id" gorm:"column:id"`
	Mid   int64      `json:"mid" gorm:"column:mid"`
	Role  int8       `json:"role" gorm:"column:role"`
	Title string     `json:"title" gorm:"column:title"`
	Desc  string     `json:"desc" gorm:"column:description"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime xtime.Time `json:"mtime" gorm:"column:mtime"`

	// 后台展示需求，需要查 official doc 对应的昵称
	Name string `json:"name" gorm:"-"`
}

// OfficialDoc is.
type OfficialDoc struct {
	ID             int64      `json:"id" gorm:"column:id"`
	Mid            int64      `json:"mid" gorm:"column:mid"`
	Name           string     `json:"name" gorm:"column:name"`
	State          int8       `json:"state" gorm:"column:state"`
	Role           int8       `json:"role" gorm:"column:role"`
	Title          string     `json:"title" gorm:"column:title"`
	Desc           string     `json:"desc" gorm:"column:description"`
	Uname          string     `json:"uname" gorm:"column:uname"`
	Extra          string     `json:"-" gorm:"column:extra"`
	IsInternal     bool       `json:"is_internal" gorm:"column:is_internal"`
	RejectReason   string     `json:"reject_reason" gorm:"column:reject_reason"`
	SubmitSource   string     `json:"submit_source" gorm:"column:submit_source"`
	SubmitTime     xtime.Time `json:"submit_time" gorm:"column:submit_time"`
	CTime          xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime          xtime.Time `json:"mtime" gorm:"column:mtime"`
	*OfficialExtra `gorm:"-"`
}

// OfficialExtra extra.
type OfficialExtra struct {
	Realname          int8   `json:"realname" gorm:"-"`
	Operator          string `json:"operator" gorm:"-"`
	Telephone         string `json:"telephone" gorm:"-"`
	Email             string `json:"email" gorm:"-"`
	Address           string `json:"address" gorm:"-"`
	Company           string `json:"company" gorm:"-"`
	CreditCode        string `json:"credit_code" gorm:"-"`        // 社会信用代码
	Organization      string `json:"organization" gorm:"-"`       // 政府或组织名称
	OrganizationType  string `json:"organization_type" gorm:"-"`  // 政府或机构类型
	BusinessLicense   string `json:"business_license" gorm:"-"`   // 企业营业执照
	BusinessScale     string `json:"business_scale" gorm:"-"`     // 企业规模
	BusinessLevel     string `json:"business_level" gorm:"-"`     // 行政级别
	BusinessAuth      string `json:"business_auth" gorm:"-"`      // 企业授权函
	Supplement        string `json:"supplement" gorm:"-"`         // 其他补充材料
	Professional      string `json:"professional" gorm:"-"`       // 专业资质
	Identification    string `json:"identification" gorm:"-"`     // 身份证明
	OfficalSite       string `json:"official_site" gorm:"-"`      // 官方站点
	RegisteredCapital string `json:"registered_capital" gorm:"-"` // 注册资本
}

// OfficialDocAddit .
type OfficialDocAddit struct {
	Mid      int64  `json:"mid" gorm:"mid"`
	Property string `json:"property" gorm:"property"`
	Vstring  string `json:"vstring" gorm:"vstring"`
}

// ParseExtra parse extra.
func (oc *OfficialDoc) ParseExtra() {
	oe := &OfficialExtra{}
	if len(oc.Extra) > 0 {
		json.Unmarshal([]byte(oc.Extra), oe)
	}
	oc.OfficialExtra = oe
}

// Validate is.
func (oc *OfficialDoc) Validate() bool {
	if oc.Mid <= 0 ||
		oc.Name == "" ||
		oc.Role <= 0 ||
		oc.Title == "" {
		return false
	}
	return true
}

// String is
func (oe *OfficialExtra) String() string {
	bs, _ := json.Marshal(oe)
	if len(bs) == 0 {
		bs = []byte("{}")
	}
	return string(bs)
}
