package model

// DecryptBindLogParam DecryptBindLogParam
type DecryptBindLogParam struct {
	EncryptText []string `form:"text,split" validate:"min=1"` // 密文字段，','分割
}

// UserBindLogReq UserBindLogReq
type UserBindLogReq struct {
	// Action value : telBindLog or emailBindLog
	Action string `form:"action"`
	Mid    int64  `form:"mid"`
	//Query search tel or email
	Query string `form:"query"`
	Page  int    `form:"page"`
	Size  int    `form:"size"`
	From  int64  `form:"from"`
	To    int64  `form:"to"`
}

// EsRes EsRes
type EsRes struct {
	Page   Page            `json:"page"`
	Result []*UserActLogEs `json:"result"`
}

// Page Page
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// UserActLogEs UserActLogEs
type UserActLogEs struct {
	Mid       int64  `json:"mid"`
	Str0      string `json:"str_0"`
	ExtraData string `json:"extra_data"`
	CTime     string `json:"ctime"`
}

// UserBindLogRes UserBindLogRes
type UserBindLogRes struct {
	Page   Page           `json:"page"`
	Result []*UserBindLog `json:"result"`
}

// UserBindLog UserBindLog
type UserBindLog struct {
	Mid   int64  `json:"mid"`
	Phone string `json:"phone"`
	EMail string `json:"email"`
	Time  string `json:"time"`
}
