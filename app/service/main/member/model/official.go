package model

import (
	"encoding/json"

	xtime "go-common/library/time"
)

// official state const.
const (
	OfficialStateWait = iota
	OfficialStatePass
	OfficialStateNoPass
	OfficialStateReWait
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

// OfficialDoc official doc.
type OfficialDoc struct {
	Mid          int64      `json:"mid"`
	Name         string     `json:"name"`
	State        int8       `json:"state"`
	Role         int8       `json:"role"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	Extra        string     `json:"-"`
	RejectReason string     `json:"reject_reason"` // 被拒绝理由
	SubmitSource string     `json:"submit_source"` // 提交来源
	SubmitTime   xtime.Time `json:"submit_time"`   // 最后提交时间

	OfficialExtra
}

// OfficialExtra official extra.
type OfficialExtra struct {
	Realname          int8   `json:"realname"`
	Operator          string `json:"operator"`           // 经营人
	Telephone         string `json:"telephone"`          // 电话号码
	Email             string `json:"email"`              // 邮箱
	Address           string `json:"address"`            // 地址
	Company           string `json:"company"`            // 公司
	CreditCode        string `json:"credit_code"`        // 社会信用代码
	Organization      string `json:"organization"`       // 政府或组织名称
	OrganizationType  string `json:"organization_type"`  // 组织或机构类型
	BusinessLicense   string `json:"business_license"`   // 企业营业执照
	BusinessScale     string `json:"business_scale"`     // 企业规模
	BusinessLevel     string `json:"business_level"`     // 企业登记
	BusinessAuth      string `json:"business_auth"`      // 企业授权函
	Supplement        string `json:"supplement"`         // 其他补充材料
	Professional      string `json:"professional"`       // 专业资质
	Identification    string `json:"identification"`     // 身份证明
	OfficialSite      string `json:"official_site"`      // 官网地址
	RegisteredCapital string `json:"registered_capital"` // 注册资金
}

// ParseExtra parse extra.
func (oc *OfficialDoc) ParseExtra() {
	oe := OfficialExtra{}
	if len(oc.Extra) > 0 {
		json.Unmarshal([]byte(oc.Extra), &oe)
	}
	oc.OfficialExtra = oe
}

// String is
func (oe OfficialExtra) String() string {
	bs, _ := json.Marshal(oe)
	if len(bs) == 0 {
		bs = []byte("{}")
	}
	return string(bs)
}

// Validate is.
func (oc OfficialDoc) Validate() bool {
	if oc.Mid <= 0 ||
		oc.Name == "" ||
		oc.Role <= 0 ||
		oc.Title == "" {
		return false
	}
	return true
}
