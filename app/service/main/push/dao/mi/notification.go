package mi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// XMMessage define reference struct http://dev.xiaomi.com/doc/?p=533
type XMMessage struct {
	Payload               string     //	消息的内容。
	RestrictedPackageName string     // App的包名。备注：V2版本支持一个包名，V3版本支持多包名（中间用逗号分割）。
	PassThrough           int        //	pass_through的值可以为： 0 表示通知栏消息1 表示透传消息
	NotifyType            int        // 通知方式
	Title                 string     //	通知栏展示的通知的标题。
	Description           string     //	通知栏展示的通知的描述。
	TaskID                string     // 上报数据使用
	xmuv                  url.Values // 含有本条消息所有属性的数组
}

func (xm *XMMessage) buildXMPostParam() {
	xmuv := url.Values{}
	xmuv.Set("payload", xm.Payload)
	xmuv.Set("restricted_package_name", xm.RestrictedPackageName)
	xmuv.Set("pass_through", strconv.Itoa(xm.PassThrough))
	xmuv.Set("title", xm.Title)
	xmuv.Set("description", xm.Description)
	xmuv.Set("notify_type", strconv.Itoa(xm.NotifyType))
	xmuv.Set("extra.task_id", xm.TaskID)
	xmuv.Set("extra.jobkey", xm.TaskID)
	xmuv.Set("extra.callback", CallbackURL)
	xmuv.Set("extra.callback.type", "1") // 第三方所需要的回执类型。1:送达回执,2:点击回执,3:送达和点击回执,默认值为3。
	xm.xmuv = xmuv
}

// SetNotifyID 可选项
// 默认情况下，通知栏只显示一条推送消息。如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏消息会覆盖之前的）。
// notify_id 0-4 同一个notifyId在通知栏只会保留一条
func (xm *XMMessage) SetNotifyID(notifyID string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("notify_id", notifyID)
}

// SetNotifyType sound / vibration / led light
func (xm *XMMessage) SetNotifyType(typ int) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("notify_type", strconv.Itoa(typ))
}

// SetTimeToLive 可选项
// 如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
// time_to_live 可选项，当用户离线是，消息保留时间，默认两周，单位ms
func (xm *XMMessage) SetTimeToLive(expire int64) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	timeToLive := (expire - time.Now().Unix()) * 1000
	xm.xmuv.Set("time_to_live", fmt.Sprintf("%d", timeToLive))
}

// SetTimeToSend 可选项
// 定时发送消息。用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。注：仅支持七天内的定时消息。
func (xm *XMMessage) SetTimeToSend(timeToSend int64) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("time_to_send", fmt.Sprintf("%d", timeToSend))
}

// SetUserAccount 根据user_account，发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。参数仅适用于“/message/user_account”HTTP API。
func (xm *XMMessage) SetUserAccount(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("user_account", UserAccount)
}

// SetUserAccounts 针对不同的userAccount推送不同的消息
// 根据user_accounts，发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。
func (xm *XMMessage) SetUserAccounts(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("user_accounts", UserAccount)
}

// SetRegID 根据registration_id，发送消息到指定设备上。可以提供多个registration_id，发送给一组设备，不同的registration_id之间用“,”分割。
func (xm *XMMessage) SetRegID(deviceToken string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("registration_id", deviceToken)
}

// SetTopic 根据topic，发送消息给订阅了该topic的所有设备。参数仅适用于“/message/topic”HTTP API。
func (xm *XMMessage) SetTopic(UserAccount string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("topic", UserAccount)
}

// SetCallbackParam 把应用标识传过去，这样方便区分应用
func (xm *XMMessage) SetCallbackParam(p string) {
	if xm.xmuv == nil {
		xm.buildXMPostParam()
	}
	xm.xmuv.Set("extra.callback.param", p) // 可选字段。第三方自定义回执参数，最大长度64个字节（这里用来存应用ID）
}

// Response push result.
type Response struct {
	Result      string `json:"result,omitempty"`      //“result”: string，”ok” 表示成功, “error” 表示失败。
	Reason      string `json:"reason,omitempty"`      //reason: string，如果失败，reason失败原因详情。
	Code        int    `json:"code,omitempty"`        //“code”: integer，0表示成功，非0表示失败。
	Data        Data   `json:"data,omitempty"`        //“data”: string，本身就是一个json字符串（其中id字段的值就是消息的Id）。
	Description string `json:"description,omitempty"` //“description”: string， 对发送消息失败原因的解释。
	Info        string `json:"info,omitempty"`        //“info”: string，详细信息。
	TraceID     string `json:"trace_id,omitempty"`    // trace id for xiaomi
}

// Data response data.
type Data struct {
	ID   string          `json:"id,omitempty"`
	List []string        `json:"list,omitempty"` // for feedback
	Data json.RawMessage `json:"data,omitempty"` // for status
}

// UninstalledResponse .
type UninstalledResponse struct {
	Code   int      `json:"errorCode,omitempty"`
	Reason string   `json:"reason,omitempty"`
	Result []string `json:"result,omitempty"`
	Data   []string
}

// UninstalledData .
type UninstalledData struct {
	Token string `json:"regId"`
	Ts    int64  `json:"ts"`
	// Alias []string `json:"alias"` // 用不上
}

// RegidCallback regid callback
type RegidCallback struct {
	AppID     string `json:"app_id"`
	AppVer    string `json:"app_version"`
	AppPkg    string `json:"app_pkg"`
	AppSecret string `json:"app_secret"`
	Regid     string `json:"regid"`
}

// Callback 推送回执（回调）
type Callback struct {
	Param     string `json:"param"`     // 开发者上传的自定义参数值。
	BarStatus string `json:"barStatus"` // 消息送达时通知栏的状态。Enable:为用户允许此app展示通知栏消息, Disable:为通知栏消息已关闭, Unknown:通知栏状态未知。
	Type      int    `json:"type"`      // callback类型
	Targets   string `json:"targets"`   // 一批alias或者regId列表，之间是用逗号分割
	Jobkey    string `json:"jobkey"`
}
