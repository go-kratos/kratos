package model

import (
	"go-common/app/service/main/passport-sns/api"
)

// SnsApps .
type SnsApps struct {
	AppID     string
	AppSecret string
	Platform  int
	Business  int
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

// CheckBindResp .
type CheckBindResp struct {
	Bind bool  `json:"bind"`
	Mid  int64 `json:"mid"`
}

// Oauth2Info oauth2 info
type Oauth2Info struct {
	UnionID string `json:"unionid"`
	OpenID  string `json:"openid"`
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	Expires int64  `json:"expires"`
}

// ConvertToProto .
func (t *SnsUser) ConvertToProto() *SnsProto {
	return &SnsProto{
		Mid:      t.Mid,
		UnionID:  t.UnionID,
		Platform: int32(t.Platform),
		Expires:  t.Expires,
	}
}

// ConvertToInfo .
func (p *SnsProto) ConvertToInfo() *api.Info {
	return &api.Info{
		Mid:      p.Mid,
		Platform: parsePlatformStr(p.Platform),
		UnionId:  p.UnionID,
		Expires:  p.Expires,
	}
}

func parsePlatformStr(platform int32) string {
	switch platform {
	case PlatformQQ:
		return PlatformQQStr
	case PlatformWEIBO:
		return PlatformWEIBOStr
	}
	return ""
}
