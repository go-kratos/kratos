package jpush

const (
	// AndroidAlertTypeAll 全开
	AndroidAlertTypeAll = -1
	// AndroidAlertTypeNone 全关
	AndroidAlertTypeNone = 0
	// AndroidAlertTypeSound 开声音
	AndroidAlertTypeSound = 1
	// AndroidAlertTypeVibrate 开振动
	AndroidAlertTypeVibrate = 2
	// AndroidAlertTypeLight 开呼吸灯
	AndroidAlertTypeLight = 4

	// AndroidStyleDefault 默认通知栏样式
	AndroidStyleDefault = 0
	// AndroidStyleBigTxt big_text 字段大文本的形式展示
	AndroidStyleBigTxt = 1
	// AndroidStyleInbox inbox 字段 json 的每个 key 对应的 value 会被当作文本条目逐条展示
	AndroidStyleInbox = 2
	// AndroidStylePic big_pic_path 字段的图片URL展示成图片
	AndroidStylePic = 3
)

// Notice .
type Notice struct {
	Alert    string          `json:"alert,omitempty"`
	Android  *AndroidNotice  `json:"android,omitempty"`
	IOS      *IOSNotice      `json:"ios,omitempty"`
	WINPhone *WinPhoneNotice `json:"winphone,omitempty"`
}

// AndroidNotice .
type AndroidNotice struct {
	Alert      string                 `json:"alert"`
	Title      string                 `json:"title,omitempty"`
	AlertType  int                    `json:"alert_type"`
	BuilderID  int                    `json:"builder_id,omitempty"`
	Style      int                    `json:"style,omitempty"`
	BigPicPath string                 `json:"big_pic_path,omitempty"`
	Extras     map[string]interface{} `json:"extras,omitempty"`
}

// SetPic sets Android notice pic.
func (an *AndroidNotice) SetPic(pic string) {
	an.Style = AndroidStylePic
	an.BigPicPath = pic
}

// IOSNotice .
type IOSNotice struct {
	Alert            interface{}            `json:"alert"`
	Sound            string                 `json:"sound,omitempty"`
	Badge            string                 `json:"badge,omitempty"`
	ContentAvailable bool                   `json:"content-available,omitempty"`
	MutableContent   bool                   `json:"mutable-content,omitempty"`
	Category         string                 `json:"category,omitempty"`
	Extras           map[string]interface{} `json:"extras,omitempty"`
}

// WinPhoneNotice .
type WinPhoneNotice struct {
	Alert    string                 `json:"alert"`
	Title    string                 `json:"title,omitempty"`
	OpenPage string                 `json:"_open_page,omitempty"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}

// SetAlert .
func (n *Notice) SetAlert(alert string) {
	n.Alert = alert
}

// SetAndroidNotice .
func (n *Notice) SetAndroidNotice(an *AndroidNotice) {
	n.Android = an
}

// SetIOSNotice .
func (n *Notice) SetIOSNotice(in *IOSNotice) {
	n.IOS = in
}

// SetWinPhoneNotice .
func (n *Notice) SetWinPhoneNotice(wn *WinPhoneNotice) {
	n.WINPhone = wn
}
