package model

// ArgMid is rpc mid params.
type ArgMid struct {
	Mid    int64
	RealIP string
}

// ArgSend send sms
type ArgSend struct {
	Mid     int64
	RealIP  string
	Mobile  string
	Country string
	Tcode   string
	Tparam  string
}

// ArgSendBatch send batch
type ArgSendBatch struct {
	Mids    []int64
	RealIP  string
	Mobiles []string
	Tcode   string
	Tparam  string
}

// ArgUserActionLog add user action log
type ArgUserActionLog struct {
	MsgID    string // 发送短信时服务商返回的随机ID
	Mobile   string
	Content  string // 短信内容
	Status   string // 回执状态
	Desc     string // 回执状态描述
	Provider int    // 短信服务商ID
	Type     int    // 短信类型，验证码/国际/营销
	Action   int    // 操作类型，发送或回执
	Ts       int64  // 操作时间
}
