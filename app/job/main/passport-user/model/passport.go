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
	MobileVerified int32     `json:"mobile_verified"`
	Isleak         int32     `json:"isleak"`
	MTime          time.Time `json:"-"`
}

// OriginAccountInfo origin aso account info.
type OriginAccountInfo struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Spacesta     int32     `json:"spacesta"`
	SafeQuestion int32     `json:"safe_question"`
	SafeAnswer   string    `json:"safe_answer"`
	JoinTime     int64     `json:"join_time"`
	JoinIP       string    `json:"join_ip"`
	JoinIPV6     []byte    `json:"join_ip_v6"`
	Port         int32     `json:"port"`
	ActiveTime   int64     `json:"active_time"`
	MTime        time.Time `json:"-"`
}

// OriginAccountReg origin aso account reg.
type OriginAccountReg struct {
	ID         int64     `json:"id"`
	Mid        int64     `json:"mid"`
	OriginType int32     `json:"origintype"`
	RegType    int32     `json:"regtype"`
	AppID      int64     `json:"appid"`
	CTime      time.Time `json:"-"`
	MTime      time.Time `json:"-"`
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
