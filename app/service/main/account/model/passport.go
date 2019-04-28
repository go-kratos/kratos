package model

import (
	xtime "go-common/library/time"
)

//PassportDetail detail.
type PassportDetail struct {
	Mid       int64  `json:"mid"`
	Email     string `json:"email"`
	Phone     string `json:"telphone"`
	Spacesta  int8   `json:"spacesta"`
	JoinTime  int64  `json:"join_time"`
	IsTourist bool   `json:"is_tourist"`
}

//PassportProfile .
type PassportProfile struct {
	Mid          int64      `json:"mid"`
	UName        string     `json:"uname"`
	UserID       string     `json:"user_id"`
	Telphone     string     `json:"telphone"`
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
