package model

// HistoryPwdCheckParam history pwd check param
type HistoryPwdCheckParam struct {
	Mid int64  `form:"mid"`
	Pwd string `form:"pwd"`
}

// HistoryPwd history pwd
type HistoryPwd struct {
	OldPwd  string
	OldSalt string
}
