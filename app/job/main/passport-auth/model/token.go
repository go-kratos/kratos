package model

import (
	"encoding/json"

	go_common_time "go-common/library/time"
)

// BMsg databus binlog message.
type BMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// OldToken old token
type OldToken struct {
	ID           int64  `json:"id"`
	Mid          int64  `json:"mid"`
	AppID        int64  `json:"appid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AppSubID     int64  `json:"app_subid"`
	CreateAt     int64  `json:"create_at"`
	Expires      int64  `json:"expires"`
	Type         int64  `json:"type"`
	CTime        string `json:"ctime"`
}

//OldCookie old cookie
type OldCookie struct {
	ID         int64  `json:"id"`
	Mid        int64  `json:"mid"`
	Session    string `json:"session_data"`
	CSRFToken  string `json:"csrf_token"`
	Type       int64  `json:"type"`
	Expires    int64  `json:"expire_time"`
	ModifyTime string `json:"modify_time"`
}

// Cookie for auth
type Cookie struct {
	ID      int64
	Mid     int64
	Session string
	CSRF    string
	Type    int64
	Expires int64
	Ctime   go_common_time.Time
	Mtime   go_common_time.Time
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

// Token for auth
type Token struct {
	ID      int64
	Mid     int64
	AppID   int64
	Token   string
	Expires int64
	Type    int64
	Ctime   go_common_time.Time
}

// Refresh for auth
type Refresh struct {
	ID      int64
	Mid     int64
	AppID   int64
	Refresh string
	Token   string
	Expires int64
	Ctime   go_common_time.Time
}
