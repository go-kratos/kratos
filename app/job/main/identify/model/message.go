package model

import "encoding/json"

// BMsg databus binlog message.
type BMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Token info.
type Token struct {
	Mid          int64  `json:"mid"`
	APPID        int64  `json:"appid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	APPSubID     int64  `json:"app_subid"`
	Expires      int32  `json:"expires"`
	Permission   string `json:"permission"`
	TP           int8   `json:"type"`
	Version      string `json:"version"`
}

// Cookie info.
type Cookie struct {
	Mid         int64  `json:"mid"`
	SessionData string `json:"session_data"`
	CSRFToken   string `json:"csrf_token"`
	TP          uint8  `json:"type"`
	ExpireTime  int32  `json:"expire_time"`
}

// AuthCookie for auth
type AuthCookie struct {
	ID      int64  `json:"id"`
	Mid     int64  `json:"mid"`
	Session string `json:"session"`
	CSRF    string `json:"csrf"`
	Type    int64  `json:"type"`
	Expires int64  `json:"expires"`
}

// AuthToken for auth
type AuthToken struct {
	ID      int64  `json:"id"`
	Mid     int64  `json:"mid"`
	AppID   int64  `json:"appid"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	Type    int64  `json:"type"`
}
