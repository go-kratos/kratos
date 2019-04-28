package model

import "go-common/library/time"

// UserBase user base.
type UserBase struct {
	Mid     int64     `json:"mid"`
	UserID  string    `json:"userid"`
	Pwd     []byte    `json:"pwd"`
	Salt    string    `json:"salt"`
	Status  int32     `json:"status"`
	Deleted int8      `json:"deleted"`
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"mtime"`
}

// UserEmail user email.
type UserEmail struct {
	Mid           int64     `json:"mid"`
	Email         []byte    `json:"email"`
	Verified      int32     `json:"verified"`
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
	JoinIPV6 []byte    `json:"join_ip_v6"`
	Port     int32     `json:"port"`
	JoinTime int64     `json:"join_time"`
	Origin   int32     `json:"origin"`
	RegType  int32     `json:"reg_type"`
	AppID    int64     `json:"appid"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

// UserSafeQuestion user safe question.
type UserSafeQuestion struct {
	Mid          int64     `json:"mid"`
	SafeQuestion int32     `json:"safe_question"`
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

// ConvertToProto convert to proto
func (u *UserBase) ConvertToProto() *UserBaseProto {
	return &UserBaseProto{
		Mid:    u.Mid,
		UserID: u.UserID,
		Pwd:    u.Pwd,
		Salt:   u.Salt,
		Status: u.Status,
	}
}

// ConvertToProto convert to proto
func (u *UserTel) ConvertToProto() *UserTelProto {
	return &UserTelProto{
		Mid:         u.Mid,
		Tel:         u.Tel,
		Cid:         u.Cid,
		TelBindTime: u.TelBindTime,
	}
}

// ConvertToProto convert to proto
func (u *UserEmail) ConvertToProto() *UserEmailProto {
	return &UserEmailProto{
		Mid:           u.Mid,
		Email:         u.Email,
		Verified:      u.Verified,
		EmailBindTime: u.EmailBindTime,
	}
}

// ConvertToProto convert to proto
func (u *UserRegOrigin) ConvertToProto() *UserRegOriginProto {
	return &UserRegOriginProto{
		Mid:      u.Mid,
		JoinIP:   u.JoinIP,
		JoinIPV6: u.JoinIPV6,
		Port:     u.Port,
		JoinTime: u.JoinTime,
		Origin:   u.Origin,
		RegType:  u.RegType,
		AppID:    u.AppID,
	}
}

// ConvertToProto convert to proto
func (u *UserThirdBind) ConvertToProto() *UserThirdBindProto {
	return &UserThirdBindProto{
		ID:       u.ID,
		Mid:      u.Mid,
		OpenID:   u.OpenID,
		PlatForm: u.PlatForm,
		Token:    u.Token,
		Expires:  u.Expires,
	}
}

// ConvertToProto convert to proto
func (u *UserSafeQuestion) ConvertToProto() *UserSafeQuestionProto {
	return &UserSafeQuestionProto{
		Mid:          u.Mid,
		SafeQuestion: u.SafeQuestion,
		SafeAnswer:   u.SafeAnswer,
		SafeBindTime: u.SafeBindTime,
	}
}
