package model

// Audit member audit info
type Audit struct {
	Mid      int64 `json:"mid"`
	BindTel  bool  `json:"bind_tel"`
	BindMail bool  `json:"bind_mail"`
	Rank     int64 `json:"rank"`
	Blocked  bool  `json:"blocked"`
}

// PassportDetail passportDetail
type PassportDetail struct {
	Mid      int64  `json:"mid"`
	Email    string `json:"email"`
	Phone    string `json:"telphone"`
	Spacesta int8   `json:"spacesta"`
	JoinTime int64  `json:"join_time"`
}
