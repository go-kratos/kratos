package model

// User user info
type User struct {
	Mid      int64  `json:"mid"`
	Uname    string `json:"uname"`
	Userid   string `json:"userid"`
	Sign     string `json:"sign"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex"`
	NickFree bool   `json:"nick_free"`
}

// Settings settings
type Settings struct {
	Uname    string `json:"uname"`
	Sign     string `json:"sign"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
}

// LogCoins log money
type LogCoins struct {
	List  []*LogCoin `json:"list"`
	Count int        `json:"count"`
}

// Coin coin.
type Coin struct {
	Money float64 `json:"money"`
}

// LogCoin money
type LogCoin struct {
	Time   string  `json:"time"`
	Delta  float64 `json:"delta"`
	Reason string  `json:"reason"`
}

// LogMorals log moral
type LogMorals struct {
	Moral int64       `json:"moral"`
	List  []*LogMoral `json:"list"`
	Count int         `json:"count"`
}

// LogMoral moral
type LogMoral struct {
	Origin string  `json:"origin"`
	Delta  float64 `json:"delta"`
	Reason string  `json:"reason"`
	Time   string  `json:"time"`
}

// LogExps log exp
type LogExps struct {
	List  []*LogExp `json:"list"`
	Count int       `json:"count"`
}

// LogExp exp
type LogExp struct {
	Delta  float64 `json:"delta"`
	Time   string  `json:"time"`
	Reason string  `json:"reason"`
}

// LogLogins log login
type LogLogins struct {
	Count int         `json:"count"`
	List  []*LogLogin `json:"list"`
}

// LogLogin logLogin
type LogLogin struct {
	IP     string `json:"ip"`
	Time   int64  `json:"time"`
	TimeAt string `json:"time_at"`
	Status bool   `json:"status"`
	Type   int64  `json:"type"`
	Geo    string `json:"geo"`
}

// Reward exp reward.
type Reward struct {
	Login bool  `json:"login"`
	Watch bool  `json:"watch"`
	Coin  int64 `json:"coins"`
	Share bool  `json:"share"`
}

// OfficialApply .
type OfficialApply struct {
	Role  int8   `form:"role" validate:"min=0,max=6" json:"role"`
	Name  string `form:"name" validate:"required" json:"name"`
	Title string `form:"title" validate:"required" json:"title"`
	Desc  string `form:"desc" json:"desc"`

	Realname          int8   `form:"realname" json:"realname"`
	Operator          string `form:"operator" json:"operator"`
	Telephone         string `form:"telephone" json:"telephone"`
	TelVerifyCode     int64  `form:"tel_verify_code" json:"tel_verify_code"`
	Email             string `form:"email" json:"email"`
	Address           string `form:"address" json:"address"`
	Company           string `form:"company" json:"company"`
	CreditCode        string `form:"credit_code" json:"credit_code"`               // 社会信用代码
	Organization      string `form:"organization" json:"organization"`             // 政府或组织名称
	OrganizationType  string `form:"organization_type" json:"organization_type"`   // 组织或机构类型
	BusinessLicense   string `form:"business_license" json:"business_license"`     // 企业营业执照
	BusinessScale     string `form:"business_scale" json:"business_scale"`         // 企业规模
	BusinessLevel     string `form:"business_level" json:"business_level"`         // 企业登记
	BusinessAuth      string `form:"business_auth" json:"business_auth"`           // 企业授权函
	Supplement        string `form:"supplement" json:"supplement"`                 // 其他补充材料
	Professional      string `form:"professional" json:"professional"`             // 专业资质
	Identification    string `form:"identification" json:"identification"`         // 身份认证
	OfficialSite      string `form:"official_site" json:"official_site"`           // 官方站点
	RegisteredCapital string `form:"registered_capital" json:"registered_capital"` // 注册资本
}

// OfficialSubmittedTimes is
type OfficialSubmittedTimes struct {
	Submitted int64 `json:"submitted"`
	Remain    int64 `json:"remain"`
}

// OfficialConditions is official conditions
type OfficialConditions struct {
	IsFormal      bool `json:"is_formal"`
	BindTel       bool `json:"bind_tel"`
	Realname      bool `json:"realname"`
	FollowerCount bool `json:"follower_count"`
	ArchiveCount  bool `json:"archive_count"`
	// ViewCount     bool `json:"view_count"`
}

// ArgMobileVerify is.
type ArgMobileVerify struct {
	Mobile  string `form:"mobile" validate:"required"`
	Country int64  `form:"country"`
}

// AllPass is
func (cons *OfficialConditions) AllPass() bool {
	return cons.IsFormal &&
		cons.BindTel &&
		cons.Realname &&
		cons.FollowerCount &&
		cons.ArchiveCount // &&
	// cons.ViewCount
}
