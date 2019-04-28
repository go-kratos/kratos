package model

import "go-common/library/time"

// OriginAccount origin aso account.
type OriginAccount struct {
	Mid            int64     `json:"mid"`
	UserID         string    `json:"userid"`
	Uname          string    `json:"uname"`
	Pwd            string    `json:"pwd"`
	Salt           string    `json:"salt"`
	Email          string    `json:"email"`
	Tel            string    `json:"tel"`
	CountryID      int64     `json:"country_id"`
	MobileVerified int8      `json:"mobile_verified"`
	Isleak         int8      `json:"isleak"`
	MTime          time.Time `json:"modify_time"`
}

// OriginAccountInfo origin aso account info.
type OriginAccountInfo struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Spacesta     int8      `json:"spacesta"`
	SafeQuestion int8      `json:"safe_question"`
	SafeAnswer   string    `json:"safe_answer"`
	JoinTime     int64     `json:"join_time"`
	JoinIP       string    `json:"join_ip"`
	ActiveTime   int64     `json:"active_time"`
	MTime        time.Time `json:"modify_time"`
}

// OriginAccountReg origin aso account reg.
type OriginAccountReg struct {
	ID         int64     `json:"id"`
	Mid        int64     `json:"mid"`
	OriginType int8      `json:"origintype"`
	RegType    int8      `json:"regtype"`
	AppID      int64     `json:"appid"`
	CTime      time.Time `json:"active_time"`
	MTime      time.Time `json:"modify_time"`
}

// OriginAccountSns origin aso account sns.
type OriginAccountSns struct {
	Mid               int64  `json:"mid"`
	SinaUID           int64  `json:"sina_uid"`
	SinaAccessToken   string `json:"sina_access_token"`
	SinaAccessExpires int64  `json:"sina_access_expires"`
	QQOpenid          string `json:"qq_openid"`
	QQAccessToken     string `json:"qq_access_token"`
	QQAccessExpires   int64  `json:"qq_access_expires"`
}

// CountryCode aso country code.
type CountryCode struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Cname string `json:"cname"`
	Rank  int64  `json:"rank"`
	Type  int8   `json:"type"`
	Ename string `json:"ename"`
}
