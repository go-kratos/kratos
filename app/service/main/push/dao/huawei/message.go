package huawei

import (
	"encoding/json"
)

const (
	// MsgTypePassthrough 消息类型:透传
	MsgTypePassthrough = 1
	// MsgTypeNotification 消息类型:通知栏消息
	MsgTypeNotification = 3

	// ActionTypeCustom 动作类型:自定义
	ActionTypeCustom = 1
	// ActionTypeURL 动作类型:打开URL
	ActionTypeURL = 2
	// ActionTypeAPP 动作类型:打开APP
	ActionTypeAPP = 3

	// CallbackTokenUninstalled 应用被卸载了
	CallbackTokenUninstalled = 2
	// CallbackTokenNotApply 终端安装了该应用，但从未打开过，未申请token，所以不能展示
	CallbackTokenNotApply = 5
	// CallbackTokenInactive 非活跃设备，消息丢弃
	CallbackTokenInactive = 10
)

// Response push response.
type Response struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Err       string `json:"error"`
	RequestID string `json:"requestId"`
}

// InvalidTokenResponse invalid tokens info in the push response.
type InvalidTokenResponse struct {
	Success       int      `json:"success"`
	Failure       int      `json:"failure"`
	IllegalTokens []string `json:"illegal_tokens"`
}

// Message request message.
type Message struct {
	Hps Hps `json:"hps"`
}

// Hps .
type Hps struct {
	Msg Msg `json:"msg"`
	Ext Ext `json:"ext"`
}

// Msg .
type Msg struct {
	Type   int    `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
}

// Body .
type Body struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

// Action .
type Action struct {
	Type  int   `json:"type"`
	Param Param `json:"param"`
}

// Param .
type Param struct {
	Intent     string `json:"intent"`
	AppPkgName string `json:"appPkgName"`
}

// Ext .
type Ext struct {
	BiTag     string              `json:"biTag"`
	Icon      string              `json:"icon"`
	Customize []map[string]string `json:"customize"`
}

// Callback 华为推送回执（回调）
type Callback struct {
	Statuses []*CallbackItem `json:"statuses"`
}

// CallbackItem http://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_devguide_server_agent#3.3 消息回执
type CallbackItem struct {
	BiTag     string `json:"biTag"`
	AppID     string `json:"appid"`
	Token     string `json:"token"`
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// NewMessage get message.
func NewMessage() *Message {
	return &Message{
		Hps: Hps{
			Msg: Msg{
				Type: MsgTypeNotification, //1  透传异步消息, 3 系统通知栏异步消息  注意:2和4以后为保留后续扩展使用
				Body: Body{
					Content: "",
					Title:   "",
				},
				Action: Action{
					Type:  ActionTypeAPP, //1 自定义行为，  2 打开URL ,3 打开App
					Param: Param{},
				},
			},
			Ext: Ext{ //扩展信息，含BI消息统计，特定展示风格，消息折叠。
				BiTag: "Trump", // 设置消息标签，如果带了这个标签，会在回执中推送给CP用于检测某种类型消息的到达率和状态
			},
		},
	}
}

// SetContent sets content.
func (m *Message) SetContent(content string) *Message {
	m.Hps.Msg.Body.Content = content
	return m
}

// SetTitle sets title.
func (m *Message) SetTitle(title string) *Message {
	m.Hps.Msg.Body.Title = title
	return m
}

// SetMsgType sets title.
func (m *Message) SetMsgType(typ int) *Message {
	m.Hps.Msg.Type = typ
	return m
}

// SetIntent sets intent.
func (m *Message) SetIntent(intent string) *Message {
	m.Hps.Msg.Action.Param.Intent = intent
	return m
}

// SetPkg sets app package name.
func (m *Message) SetPkg(pkg string) *Message {
	m.Hps.Msg.Action.Param.AppPkgName = pkg
	return m
}

// SetCustomize set ext info.
func (m *Message) SetCustomize(key, val string) *Message {
	mp := map[string]string{key: val}
	m.Hps.Ext.Customize = append(m.Hps.Ext.Customize, mp)
	return m
}

// SetBiTag set biTag.
func (m *Message) SetBiTag(tag string) *Message {
	m.Hps.Ext.BiTag = tag
	return m
}

// SetIcon sets icon.
func (m *Message) SetIcon(url string) *Message {
	m.Hps.Ext.Icon = url
	return m
}

// JSON encode the message.
func (m *Message) JSON() (res string, err error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	res = string(bytes)
	return
}
