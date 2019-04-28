package model

import "go-common/library/time"

// UserBase user base.
type UserBase struct {
	Mid     int64     `json:"mid"`
	UserID  string    `json:"userid"`
	Pwd     []byte    `json:"pwd"`
	Salt    string    `json:"salt"`
	Status  int8      `json:"status"`
	Deleted int8      `json:"deleted"`
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"mtime"`
}

// UserEmail user email.
type UserEmail struct {
	Mid           int64     `json:"mid"`
	Email         []byte    `json:"email"`
	Verified      int8      `json:"verified"`
	EmailBindTime int64     `json:"email_bind_time"`
	CTime         time.Time `json:"ctime"`
	MTime         time.Time `json:"mtime"`
}

// UserTel user tel.
type UserTel struct {
	Mid         int64     `json:"mid"`
	Tel         []byte    `json:"tel"`
	Cid         string    `json:"cid"`
	TelBindTime int64     `json:"tel_bind_time"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"mtime"`
}

// UserRegOrigin user reg origin.
type UserRegOrigin struct {
	Mid      int64     `json:"mid"`
	JoinIP   int64     `json:"join_ip"`
	JoinTime int64     `json:"join_time"`
	Origin   int8      `json:"origin"`
	RegType  int8      `json:"reg_type"`
	AppID    int64     `json:"appid"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// UserSafeQuestion user safe question.
type UserSafeQuestion struct {
	Mid          int64     `json:"mid"`
	SafeQuestion int8      `json:"safe_question"`
	SafeAnswer   []byte    `json:"safe_answer"`
	SafeBindTime int64     `json:"safe_bind_time"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"mtime"`
}

// UserThirdBind user third bind.
type UserThirdBind struct {
	ID       int64     `json:"id"`
	Mid      int64     `json:"mid"`
	OpenID   string    `json:"openid"`
	PlatForm int64     `json:"platform"`
	Token    string    `json:"token"`
	Expires  int64     `json:"expires"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// UserTelDuplicate user tel duplicate.
type UserTelDuplicate struct {
	ID          int64  `json:"id"`
	Mid         int64  `json:"mid"`
	Tel         []byte `json:"tel"`
	Cid         string `json:"cid"`
	TelBindTime int64  `json:"tel_bind_time"`
	Status      int8   `json:"status"`
	Timestamp   int64  `json:"ts"`
}

// UserEmailDuplicate user email duplicate.
type UserEmailDuplicate struct {
	ID            int64  `json:"id"`
	Mid           int64  `json:"mid"`
	Email         []byte `json:"email"`
	Verified      int32  `json:"verified"`
	EmailBindTime int64  `json:"email_bind_time"`
	Status        int8   `json:"status"`
	Timestamp     int64  `json:"ts"`
}
