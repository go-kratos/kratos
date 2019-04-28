package model

// QQAccessResp qq access response
type QQAccessResp struct {
	Token       string `json:"access_token"`
	Refresh     string `json:"refresh_token"`
	Expires     int64  `json:"expires_in"`
	Code        int    `json:"error"`
	Description string `json:"error_description"`
}

// QQOpenIDResp qq open id response
type QQOpenIDResp struct {
	UnionID     string `json:"unionid"`
	OpenID      string `json:"openid"`
	Code        int    `json:"error"`
	Description string `json:"error_description"`
}
