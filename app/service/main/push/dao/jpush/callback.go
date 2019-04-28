package jpush

const (
	// CallbackTypeReceive 送达才回执
	CallbackTypeReceive = callbackType(1)
	// CallbackTypeClick 点击才回执
	CallbackTypeClick = callbackType(2)
	// CallbackTypeAll 送达和点击都回执
	CallbackTypeAll = callbackType(3)

	// StatusSwitchOn 回执时候通知栏开关状态：开
	StatusSwitchOn = int(1)
	// StatusSwitchOff 回执时候通知栏开关状态：关
	StatusSwitchOff = int(2)

	defaultCallbackURL = "https://api.bilibili.com/x/push/callback/jpush"
)

type callbackType int

// CallbackReq 消息回执请求体
type CallbackReq struct {
	// URL 接受回执数据的URL
	URL string `json:"url"`
	// Type 需要的回执类型
	Type callbackType `json:"type"`
	// Params 携带的自定义参数
	Params map[string]string `json:"params"`
}

// NewCallbackReq new Callback
func NewCallbackReq() *CallbackReq {
	return &CallbackReq{
		URL:    defaultCallbackURL,
		Type:   CallbackTypeReceive,
		Params: make(map[string]string),
	}
}

// SetURL 设置接收回执的URL
func (cb *CallbackReq) SetURL(url string) {
	if url == "" {
		return
	}
	cb.URL = url
}

// SetType 设置需要回执的类型
func (cb *CallbackReq) SetType(typ callbackType) {
	cb.Type = typ
}

// SetParam 设置自定义参数
func (cb *CallbackReq) SetParam(m map[string]string) {
	if m == nil {
		return
	}
	cb.Params = m
}

// CallbackReply 消息回执接收体
type CallbackReply struct {
	// Token device token
	Token string `json:"registration_id"`
	// Platform android or ios
	Platform string `json:"platform"`
	// Time 消息送达或点击的秒级时间戳
	Time int64 `json:"sent_time"`
	// Switch 通知栏消息开关
	Switch bool `json:"notification_state"`
	// Type 送达或点击
	Type callbackType `json:"callback_type"`
	// Channel 下发通道
	Channel int `json:"channel"`
	// Params 自定义参数
	Params map[string]string `json:"params"`
}
