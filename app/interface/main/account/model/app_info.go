package model

import xtime "go-common/library/time"

// const .
const (
	UpNameLogID  = 14
	UpNameAction = "log_name_update"
)

// Account account.
type Account struct {
	Mid      int64  `json:"mid"`
	Uname    string `json:"uname"`
	Face     string `json:"face"`
	Sex      int8   `json:"sex"`
	Birthday string `json:"birthday"`
	Sign     string `json:"sign"`
	NickFree bool   `json:"nick_free"`
}

// NickFree nickFree.
type NickFree struct {
	NickFree bool `json:"nick_free"`
}

// PassportProfile is
/*
   "mid": 288840748,
   "uname": "小学生pasami",
   "userid": "bili_43796499903",
   "telphone": "9028005779",
   "nickLock": 0,
   "bind_qq": false,
   "bind_sina": true,
   "spacesta": 2,
   "login_time": 1529165291,
   "login_ip": "103.228.109.204",
   "join_ip": "119.135.161.112",
   "join_time": 1518286685,
   "safe_question": 0,
   "country_code": 1
*/
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
