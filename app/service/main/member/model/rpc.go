package model

import (
	"encoding/json"

	"go-common/library/log"
	"go-common/library/time"
)

// ArgMid arg mid.
type ArgMid struct {
	Mid    int64
	RealIP string
}

// ArgMid2 arg mid2.
type ArgMid2 struct {
	Mid    int64 `form:"mid" validate:"min=1,required"` // 用户mid
	RealIP string
}

// ArgMemberMid is.
type ArgMemberMid struct {
	Mid      int64  `json:"mid"`
	RemoteIP string `json:"remoteIP"`
}

// ArgMemberMids are.
type ArgMemberMids struct {
	Mids     []int64 `json:"mids"`
	RemoteIP string  `json:"remoteIP"`
}

// ArgOfficialDoc arg official doc
type ArgOfficialDoc struct {
	Mid   int64  `json:"mid"`
	Name  string `json:"name"`
	Role  int8   `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`

	Realname          int8   `json:"realname"`
	Operator          string `json:"operator"`
	Telephone         string `json:"telephone"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	Company           string `json:"company"`
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

	SubmitSource string `json:"submit_source"` // 提交来源
}

// Log define user login log.
type Log struct {
	Mid        int64     `json:"mid,omitempty"`
	IP         uint32    `json:"loginip"`
	Location   string    `json:"location"`
	LocationID int64     `json:"location_id,omitempty"`
	Time       time.Time `json:"timestamp,omitempty"`
	Type       int8      `json:"type,omitempty"`
}

// Msg is user login status msg.
type Msg struct {
	Notify bool `json:"notify"`
	Log    *Log `json:"log"`
}

// ArgUpdateSex is.
type ArgUpdateSex struct {
	Mid      int64  `json:"mid"`
	Sex      int64  `json:"sex"`
	RemoteIP string `json:"remoteIP"`
}

// ArgUpdateFace is.
type ArgUpdateFace struct {
	Mid      int64  `json:"mid"`
	Face     string `json:"face"`
	RemoteIP string `json:"remoteIP"`
}

// ArgUpdateRank is.
type ArgUpdateRank struct {
	Mid      int64  `json:"mid"`
	Rank     int64  `json:"rank"`
	RemoteIP string `json:"remoteIP"`
}

// ArgUpdateBirthday is.
type ArgUpdateBirthday struct {
	Mid      int64     `json:"mid"`
	Birthday time.Time `json:"birthday"`
	RemoteIP string    `json:"remoteIP"`
}

// ArgUpdateUname arg for update uname.
type ArgUpdateUname struct {
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	RemoteIP string `json:"remoteIP"`
}

// ArgUpdateSign arg for udpate sign.
type ArgUpdateSign struct {
	Mid      int64  `json:"mid"`
	Sign     string `json:"sign"`
	RemoteIP string `json:"remoteIP"`
}

// ArgAddExp addexp arg.
type ArgAddExp struct {
	Mid     int64   `json:"mid,omitempty" form:"mid" validate:"min=1,required"`   // 用户mid
	Count   float64 `json:"count,omitempty"  form:"count" validate:"required"`    // 修改数量
	Reason  string  `json:"reason,omitempty" form:"reason" validate:"required"`   // 修改原因
	Operate string  `json:"operate,omitempty" form:"operate" validate:"required"` // 操作类型
	IP      string  `json:"ip" form:"ip"`
}

// ExpStat user exp stat.
type ExpStat struct {
	Login bool  `json:"login"`
	Watch bool  `json:"watch_av"`
	Coin  int64 `json:"coins_av"`
	Share bool  `json:"share_av"`
}

// ArgRealnameApply realname apply
type ArgRealnameApply struct {
	MID           int64
	CaptureCode   int
	Realname      string
	CardType      int8
	CardCode      string
	Country       int16
	HandIMGToken  string
	FrontIMGToken string
	BackIMGToken  string
}

// ArgRealnameAlipayConfirm is
type ArgRealnameAlipayConfirm struct {
	MID    int64
	Pass   bool
	Reason string
}

// ArgRealnameAlipayApply is
type ArgRealnameAlipayApply struct {
	MID         int64
	CaptureCode int
	Realname    string
	CardCode    string
	IMGToken    string
	Bizno       string
}

// ArgAddUserMonitor is
type ArgAddUserMonitor struct {
	Mid      int64
	Operator string
	Remark   string
}

// ArgAddPropertyReview is.
type ArgAddPropertyReview struct {
	Mid      int64                  `form:"mid" validate:"min=1,required"` // 用户mid
	New      string                 `form:"new"`                           // 新的值
	State    int8                   `form:"state"`                         // 0 待审核，1 通过，2 驳回，10 自动审核中
	Property int8                   `form:"property"`                      // 0 无意义，1 头像，2 签名，3 昵称
	Extra    map[string]interface{} // 审核扩展字段 extra
}

// ExtraStr is.
func (arg *ArgAddPropertyReview) ExtraStr() string {
	if arg.Extra == nil {
		return "{}"
	}
	bs, err := json.Marshal(arg.Extra)
	if err != nil {
		log.Error("Failed to marshal extra: %+v, error: %+v", arg.Extra, err)
		return "{}"
	}
	return string(bs)
}
