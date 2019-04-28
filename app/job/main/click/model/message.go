package model

// plat(web:0,h5:1,outer:2,ios:3,android:4),avid,cid,part,mid,lv,ftime,stime,buvid(device),ip,agent(version)

// is
const (
	LogTypeForNotUse      = int8(0)
	LogTypeForTurly       = int8(1)
	LogTypeForInlineBegin = int8(2)
)

// ClickMsg is
type ClickMsg struct {
	Plat       int8
	AID        int64
	MID        int64
	Lv         int8
	Buvid      string
	Did        string
	CTime      int64
	STime      int64
	IP         string
	KafkaBs    []byte
	EpID       int64
	SeasonType int
	UserAgent  string
}

// StatMsg is
type StatMsg struct {
	AID   int64 `json:"aid"`
	Click int   `json:"click"`
}

// StatViewMsg is
type StatViewMsg struct {
	Type  string `json:"type"`
	ID    int64  `json:"id"`
	Count int    `json:"count"`
	Ts    int64  `json:"timestamp"`
}

// BigDataMsg is
type BigDataMsg struct {
	Info string `json:"info"`
	Tp   int8   `json:"type"`
}
