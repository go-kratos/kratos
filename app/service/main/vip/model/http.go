package model

// MemberResq response params .
type MemberResq struct {
	CommonResq
	Data *Member `json:"data"`
}

// Member response params .
type Member struct {
	Mid      string `json:"mid"`
	Name     string `json:"name"`
	Face     string `json:"face"`
	Sign     string `json:"sign"`
	Sex      string `json:"sex"`
	Cert     string `json:"cert"`
	Rank     string `json:"rank"`
	Certdesc string `json:"certdesc"`
}

// PayResq response params.
type PayResq struct {
	Errno   int64  `json:"errno"`
	Message string `json:"msg"`
}

// CommonResq response params.
type CommonResq struct {
	Code    int64  `json:"code"`
	TS      int64  `json:"ts"`
	Message string `json:"message"`
}

//TokenResq get token resq.
type TokenResq struct {
	CommonResq
	Data *Token `json:"data"`
}

//Token get token .
type Token struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

//OpenCodeResp openCode resq.
type OpenCodeResp struct {
	CommonResq
	Data int64 `json:"data"`
}

//PassportDetail .
type PassportDetail struct {
	Mid      int64  `json:"mid"`
	Email    string `json:"email"`
	Phone    string `json:"telphone"`
	Spacesta int8   `json:"spacesta"`
	JoinTime int64  `json:"join_time"`
}
