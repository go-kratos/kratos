package model

import (
	xtime "go-common/library/time"
	"strings"
)

// PassportProfile is
type PassportProfile struct {
	Mid          int64      `json:"mid"`
	UName        string     `json:"uname"`
	UserID       string     `json:"user_id"`
	Telphone     string     `json:"telphone"`
	Email        string     `json:"email"`
	NickLock     int        `json:"nick_lock"`
	BindQQ       bool       `json:"bind_qq"`
	BindSina     bool       `json:"bind_sina"`
	SpaceSta     int        `json:"spacesta"`
	LoginTime    xtime.Time `json:"login_time"`
	LoginIP      string     `json:"login_ip"`
	JoinIP       string     `json:"join_ip"`
	JoinTime     xtime.Time `json:"join_time"`
	SafeQuestion int        `json:"safe_question"`
	CountryCode  int64      `json:"country_code"`
}

func bool2int(in bool) int64 {
	if in {
		return 1
	}
	return 0
}

// TelStatus is
func (p *PassportProfile) TelStatus() int64 {
	return bool2int(p.Telphone != "")
}

// EmailSuffix is
func (p *PassportProfile) EmailSuffix() string {
	if !strings.Contains(p.Email, "@") {
		return p.Email
	}
	parts := strings.SplitN(p.Email, "@", 2)
	return parts[1]
}

// AsoAccountRegOrigin is
type AsoAccountRegOrigin struct {
	Mid        int64 `json:"mid"`
	OriginType int64 `json:"origin_type"`
	RegType    int64 `json:"reg_type"`
}
