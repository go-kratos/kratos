package archive

import (
	"go-common/library/time"
)

//Platform const
const (
	PlatformWeb     = "web"
	PlatformWindows = "windows"
	PlatformH5      = "h5"
	PlatformAndroid = "android"
	PlatformIOS     = "ios"
)

// Archive is archive model.
type Archive struct {
	Aid          int64  `json:"aid"`
	Mid          int64  `json:"mid"`
	TypeID       int16  `json:"tid"`
	HumanRank    int    `json:"-"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	Cover        string `json:"cover"`
	RejectReason string `json:"reject_reason"`
	Tag          string `json:"tag"`
	Duration     int64  `json:"duration"`
	Copyright    int8   `json:"copyright"`
	NoReprint    int8   `json:"no_reprint"`
	UgcPay       int8   `json:"ugcpay"`
	OrderID      int64  `json:"order_id"`
	OrderName    string `json:"order_name"`
	Desc         string `json:"desc"`
	MissionID    int64  `json:"mission_id"`
	MissionName  string `json:"mission_name"`
	Round        int8   `json:"-"`
	Forward      int64  `json:"-"`
	Attribute    int32  `json:"attribute"`
	Access       int16  `json:"-"`
	State        int8   `json:"state"`
	StateDesc    string `json:"state_desc"`
	StatePanel   int    `json:"state_panel"`
	Source       string `json:"source"`
	DescFormatID int64  `json:"desc_format_id"`
	Attrs        *Attrs `json:"attrs"`
	// feature: private orders
	Porder  *ArcPorder `json:"porder"`
	Dynamic string     `json:"dynamic"`
	PoiObj  *PoiObj    `json:"poi_object"`
	// time
	DTime      time.Time    `json:"dtime"`
	PTime      time.Time    `json:"ptime"`
	CTime      time.Time    `json:"ctime"`
	MTime      time.Time    `json:"-"`
	UgcPayInfo *UgcPayInfo  `json:"ugcpay_info"`
	Staffs     []*StaffView `json:"staffs"`
	Vote       *Vote        `json:"vote"`
}

// Attrs str
type Attrs struct {
	IsCoop  int8 `json:"is_coop"`
	IsOwner int8 `json:"is_owner"`
}

// Vote  str
type Vote struct {
	VoteID    int64  `json:"vote_id"`
	VoteTitle string `json:"vote_title"`
}

// PayAct  str
type PayAct struct {
	Reason string `json:"reason"`
	State  int8   `json:"state"`
}

// PayAsset  str
type PayAsset struct {
	Price         int            `json:"price"`
	PlatformPrice map[string]int `json:"platform_price"`
}

// UgcPayInfo str
type UgcPayInfo struct {
	Acts  map[string]*PayAct `json:"acts"`
	Asset *PayAsset          `json:"asset"`
}

// NilPoiObj fn 防止非APP端的地理位置信息泄露
func (arc *Archive) NilPoiObj(platform string) {
	if (platform != PlatformAndroid) &&
		(platform != PlatformIOS) &&
		(platform != PlatformH5) {
		arc.PoiObj = nil
	}
}

// NilVote fn
func (arc *Archive) NilVote() {
	if arc.Vote != nil && arc.Vote.VoteID == 0 {
		arc.Vote = nil
	}
}

// ArcPorder str
type ArcPorder struct {
	FlowID     int64  `json:"flow_id"`
	IndustryID int64  `json:"industry_id"`
	BrandID    int64  `json:"brand_id"`
	BrandName  string `json:"brand_name"`
	Official   int8   `json:"official"`
	ShowType   string `json:"show_type"`
	// for admin operation
	Advertiser string `json:"advertiser"`
	Agent      string `json:"agent"`
	//state 0 自首  1  审核添加
	State int8 `json:"state"`
}

// Video is videos model.
type Video struct {
	ID           int64     `json:"-"`
	Aid          int64     `json:"aid"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Filename     string    `json:"filename"`
	SrcType      string    `json:"-"`
	Cid          int64     `json:"cid"`
	Duration     int64     `json:"duration"`
	Filesize     int64     `json:"-"`
	Resolutions  string    `json:"-"`
	Index        int       `json:"index"`
	Playurl      string    `json:"-"`
	Status       int16     `json:"status"`
	StatusDesc   string    `json:"status_desc"`
	RejectReason string    `json:"reject_reason"`
	FailCode     int8      `json:"fail_code"`
	FailDesc     string    `json:"fail_desc"`
	XcodeState   int8      `json:"-"`
	Attribute    int32     `json:"-"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
}

// ArcVideo str
type ArcVideo struct {
	Archive *Archive
	Videos  []*Video
}

// StaffView Archive staff
type StaffView struct {
	ID         int64  `json:"id"`
	ApMID      int64  `json:"apply_staff_mid"`
	ApName     string `json:"apply_staff_name"`
	ApTitle    string `json:"apply_title"`
	ApAID      int64  `json:"apply_aid"`
	ApType     int    `json:"apply_type"`
	ApState    int    `json:"apply_state"`
	ApStaffID  int64  `json:"apply_asid"` //Staff表的主键ID
	StaffState int    `json:"staff_state"`
	StaffTitle string `json:"staff_title"`
}

// ViewBGM bgm view
type ViewBGM struct {
	SID     int64  `json:"sid"`
	MID     int64  `json:"mid"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	JumpURL string `json:"jump_url"`
}

// Vcover muti cover for video.
type Vcover struct {
	Filename string `json:"filename"`
	BFSPath  string `json:"bfs_path"`
}

// SimpleArchive is simple model for videos
type SimpleArchive struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
}

// SimpleVideo is simple model for videos
type SimpleVideo struct {
	Cid   int64  `json:"cid"`
	Index int    `json:"index"`
	Title string `json:"title"`
}

// RecoArch is simple archive information for recommend
type RecoArch struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
	Owner string `json:"owner"`
}

// SpVideo is a simple model with danmu status
type SpVideo struct {
	Cid        int64     `json:"cid"`
	Index      int       `json:"part_id"`
	Title      string    `json:"part_name"`
	Status     int16     `json:"status"`
	DmActive   int       `json:"dm_active"`
	DmModified time.Time `json:"dm_modified"`
}

// SpArchive str
type SpArchive struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
	Mid   int64  `json:"mid,omitempty"`
}

// SimpleArchiveVideos str
type SimpleArchiveVideos struct {
	Archive   *SpArchive `json:"archive"`
	SpVideos  []*SpVideo `json:"part_list"`
	AcceptAss bool       `json:"accept_ass"`
}

// VideoJam is video traffic jam info for frontend
type VideoJam struct {
	Level   int8   `json:"level"`
	State   string `json:"state"`
	Comment string `json:"comment"`
}

// Dpub str
type Dpub struct {
	Deftime    time.Time `json:"deftime"`
	DeftimeEnd time.Time `json:"deftime_end"`
	DeftimeMsg string    `json:"deftime_msg"`
}
