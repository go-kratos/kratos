package model

// SnsLog .
type SnsLog struct {
	Mid         int64  `json:"mid"`
	OpenID      string `json:"openid"`
	UnionID     string `json:"unionid"`
	AppID       string `json:"appid"`
	Platform    int    `json:"platform"`
	Operator    string `json:"operator"`
	Operate     int    `json:"operate"`
	Description string `json:"description"`
}

// SnsUser .
type SnsUser struct {
	Mid      int64  `json:"mid"`
	UnionID  string `json:"unionid"`
	Platform int    `json:"platform"`
	Expires  int64  `json:"expires"`
}

// SnsOpenID .
type SnsOpenID struct {
	Mid      int64  `json:"mid"`
	OpenID   string `json:"openid"`
	UnionID  string `json:"unionid"`
	AppID    string `json:"appid"`
	Platform int    `json:"platform"`
}

// SnsToken .
type SnsToken struct {
	Mid      int64  `json:"mid"`
	OpenID   string `json:"openid"`
	UnionID  string `json:"unionid"`
	Platform int    `json:"platform"`
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
	AppID    string `json:"appid"`
}

// AsoAccountSns aso account sns.
type AsoAccountSns struct {
	Mid               int64  `json:"mid"`
	SinaUID           int64  `json:"sina_uid"`
	SinaAccessToken   string `json:"sina_access_token"`
	SinaAccessExpires int64  `json:"sina_access_expires"`
	QQOpenid          string `json:"qq_openid"`
	QQAccessToken     string `json:"qq_access_token"`
	QQAccessExpires   int64  `json:"qq_access_expires"`
}

// QQUnionIDResp qq unionid resp.
type QQUnionIDResp struct {
	Code        int    `json:"error"`
	Description string `json:"error_description"`
	UnionID     string `json:"unionid"`
}
