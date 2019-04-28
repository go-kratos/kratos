package model

import "go-common/library/time"

const (
	// EnvPro is pro.
	EnvPro = "pro"
	// EnvTest is env.
	EnvTest = "test"
	// EnvDev is env.
	EnvDev = "dev"
	// MaxUploadSize for h5 upload
	MaxUploadSize = 20 * 1024 * 1024
)

// feedback model const.
const (
	// StateNoReply 未回复
	StateNoReply = 0
	// StateReplied 已回复
	StateReplied = 1
	// StateRepeated 二次追问
	StateRepeated = 2
	// StateOther 其它
	StateOther = 4

	// TypeCustomer 客户
	TypeCustomer = 0
	// TypeServer 客服
	TypeServer = 1
	// player cast screen
	AndroidPlayerScreen        = int64(464)
	AndroidPlayerScreenNothing = int64(465)
	AndroidPlayerScreenDlna    = int64(466)
	AndroidPlayerScreenTV      = int64(470)
	IOSPlayerScreen            = int64(467)
	IOSPlayerScreenNothing     = int64(468)
	IOSPlayerScreenDlna        = int64(469)
	IOSPlayerScreenTV          = int64(471)
)

// Session is feedback session
type Session struct {
	ID          int64     `json:"id"`
	Buvid       string    `json:"-"`
	System      string    `json:"-"`
	Version     string    `json:"-"`
	Mid         int64     `json:"-"`
	Aid         string    `json:"-"`
	Content     string    `json:"content"`
	ImgURL      string    `json:"-"`
	LogURL      string    `json:"-"`
	Device      string    `json:"-"`
	Channel     string    `json:"-"`
	IP          uint32    `json:"-"`
	NetState    string    `json:"-"`
	NetOperator string    `json:"-"`
	AgencyArea  string    `json:"-"`
	Platform    string    `json:"-"`
	Browser     string    `json:"-"`
	Email       string    `json:"-"`
	QQ          string    `json:"-"`
	State       int8      `json:"state"`
	ReplyID     string    `json:"-"`
	ReplyTime   time.Time `json:"-"`
	LasterTime  time.Time `json:"-"`
	CTime       time.Time `json:"ctime"`
	MTime       time.Time `json:"-"`
}

// WebSession web session.
type WebSession struct {
	Session *Session `json:"session,omitempty"`
	Tag     *Tag     `json:"tag,omitempty"`
}

// Reply is feedback reply
type Reply struct {
	ID        int64     `json:"-"`
	SessionID int64     `json:"-"`
	ReplyID   string    `json:"reply_id"`
	Type      int8      `json:"type"`
	Content   string    `json:"content"`
	ImgURL    string    `json:"img_url"`
	LogURL    string    `json:"log_url"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"-"`
}

// SsnAndTagID ssn and tagid.
type SsnAndTagID struct {
	TagID     int64 `json:"tag_id"`
	SessionID int64 `json:"session_id"`
}

// Replys for sort by ctime asc
type Replys []Reply

func (t Replys) Len() int { return len(t) }
func (t Replys) Less(i, j int) bool {
	return t[i].CTime < t[j].CTime
}
func (t Replys) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// FormPlatForm fro playerCheck.
func FormPlatForm(platStr string) (plat int) {
	switch platStr {
	case "web":
		plat = 1
	case "ios":
		plat = 2
	case "android":
		plat = 3
	}
	return
}
