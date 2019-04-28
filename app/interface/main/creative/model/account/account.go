package account

import "go-common/library/time"

// MyInfo get user info  for pre archive.
type MyInfo struct {
	Mid          int64           `json:"mid"`
	Name         string          `json:"uname"`
	Face         string          `json:"face"`
	Banned       bool            `json:"banned"`
	Level        int             `json:"level"`
	Activated    bool            `json:"activated"`
	Deftime      time.Time       `json:"deftime"`
	DeftimeEnd   time.Time       `json:"deftime_end"`
	DeftimeMsg   string          `json:"deftime_msg"`
	Commercial   int             `json:"commercial"`
	VideoRate    uint            `json:"video_rate,omitempty"`
	AudioRate    uint            `json:"audio_rate,omitempty"`
	IdentifyInfo *IdentifyInfo   `json:"identify_check"`
	DmSubtitle   bool            `json:"subtitle"`   //弹幕子业务之字幕协同创作
	DymcLottery  bool            `json:"lottery"`    //动态子业务之抽奖
	UploadSize   map[string]bool `json:"uploadsize"` // upload_size <= 8G
}

// IdentifyInfo str
type IdentifyInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	// IdentifyEnum define
	IdentifyEnum = map[int]string{
		0: "已实名认证",
		1: "根据国家实名制认证的相关要求，您需要换绑一个非170/171的手机号，才能继续进行操作。",
		2: "根据国家实名制认证的相关要求，您需要绑定手机号，才能继续进行操作。",
	}
)

const (
	// IsUp is up
	IsUp = 1
	// NotUp not up
	NotUp = 0
)

// UpInfo up type infos.
type UpInfo struct {
	Archive int `json:"archive"`
	Article int `json:"article"`
	Pic     int `json:"pic"`
	Blink   int `json:"blink"`
}

//IsUper judge up auth by archive/article/pic/blink.
func IsUper(up *UpInfo) (ok bool) {
	if up.Archive == 1 || up.Article == 1 || up.Blink == 1 || up.Pic == 1 {
		ok = true
	}
	return
}

// Friend str
type Friend struct {
	Mid          int64  `json:"mid"`
	Name         string `json:"name"`
	Face         string `json:"face"`
	Sign         string `json:"sign"`
	Comment      string `json:"comment"`
	ShouldFollow int8   `json:"should_follow"`
}

// SearchUp UP主搜索结果
type SearchUp struct {
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	Face     string `json:"face"`
	IsBlock  bool   `json:"is_block"`
	Relation int    `json:"relation"`
	Silence  int32  `json:"silence"`
}
