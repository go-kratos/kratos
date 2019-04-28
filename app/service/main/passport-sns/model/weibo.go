package model

// WeiboAccessResp weibo access response
type WeiboAccessResp struct {
	Token       string `json:"access_token"`
	Refresh     string `json:"refresh_token"`
	Expires     int64  `json:"expires_in"` // access_token的生命周期，单位是秒数
	OpenID      string `json:"uid"`
	Code        int    `json:"error_code"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}
