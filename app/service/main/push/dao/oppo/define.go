package oppo

import "fmt"

const (
	_host = "https://api.push.oppomobile.com"

	_apiAuth          = _host + "/server/v1/auth"
	_apiMessage       = _host + "/server/v1/message/notification/save_message_content" // 保存通知栏消息内容体
	_apiPushUnicast   = _host + "/server/v1/message/notification/unicast"              // 单条推送
	_apiPushBroadcast = _host + "/server/v1/message/notification/broadcast"            // 批量推送
	// _apiStatistics    = _host + "/server/v1/message/statistics"                        // 推送统计

	_callbackURL = "https://api.bilibili.com/x/push/callback/oppo"

	_authExpire = 24 * 60 * 60 // auth token 过期秒数

	// _pushTypeAll   = "1" // 推送全部设备
	_pushTypeToken = "2" // 按token推

	// ResponseCodeServiceUnavalable service unavalable
	ResponseCodeServiceUnavalable = -1
	// ResponseCodeSuccess http normal response code
	ResponseCodeSuccess = 0
	// ResponseCodeInvalidToken invalid token response code
	ResponseCodeInvalidToken = 10000
	// ResponseCodeUnsubscribeToken unsubscribe token
	ResponseCodeUnsubscribeToken = 10001
	// ResponseCodeRepeatToken repeat token
	ResponseCodeRepeatToken = 10004

	// ActionTypeInner 打开应⽤内⻚（activity的intentaction）
	ActionTypeInner = 1
)

// Message message content.
type Message struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	ActionType     int    `json:"click_action_type"`     // 0:启动应⽤; 1:打开应⽤内⻚(activity的intentaction); 2:打开⽹⻚; 4:打开应⽤内⻚(activity); [⾮必填，默认值为0]
	ActionActivity string `json:"click_action_activity"` // 应⽤内⻚地址【click_action_type 为1或4时必填，⻓度500】
	ActionURL      string `json:"click_action_url"`      // ⽹⻚地址【click_action_type为2 必填，⻓度500】
	ActionParams   string `json:"action_parameters"`     // 传递给应⽤的参数,json格式
	OfflineTTL     int    `json:"off_line_ttl"`          // 离线消息的存活时间 (默认3600s) (单位：秒), 【off_line值为true时，必填，最 ⻓10天】
	CallbackURL    string `json:"call_back_url"`         // 应⽤接收消息到达回执的回调(仅⽀持registrationId或aliasName 两种推送⽅式)
}

// Response push response.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MsgID            string   `json:"message_id,omitempty"`
		TaskID           string   `json:"task_id,omitempty"`
		TokenInvalid     []string `json:"10000"`
		TokenUnsubscribe []string `json:"10001"`
		TokenRepeat      []string `json:"10004"`
	} `json:"data"`
}

// Callback oppo callback.
type Callback struct {
	MsgID     string `json:"messageId"`
	TaskID    string `json:"taskId"`
	Tokens    string `json:"registrationIds"` // regId1, regid2
	EventType string `json:"eventType"`       // push_arrive
}

// CallbackURL gets callback URL.
func CallbackURL(app int64, task string) string {
	return fmt.Sprintf("%s?app=%d&task=%s", _callbackURL, app, task)
}
