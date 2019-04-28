package archive

import (
	"go-common/library/time"
)

// ArchiveAddLogID
const (
	// ArchiveAddLogID 投稿日志
	ArchiveAddLogID = int(81)
	// UgcpayAddarcProtocol ugc新增稿件时候记录的协议版本
	UgcpayAddarcProtocol = int(83)
	// LogTypeSuccess 投稿成功
	LogTypeSuccess = int(0)
	// LogTypeFail 投稿失败
	LogTypeFail = int(1)
)

// AppRequest str
type AppRequest struct {
	MobiApp  string
	Platform string
	Build    string
	Device   string
}

// ArcParam str
type ArcParam struct {
	Aid          int64         `json:"aid"`
	Mid          int64         `json:"mid"`
	Author       string        `json:"author"`
	TypeID       int16         `json:"tid"`
	Title        string        `json:"title"`
	Cover        string        `json:"cover"`
	Tag          string        `json:"tag"`
	Copyright    int8          `json:"copyright"`
	NoReprint    int8          `json:"no_reprint"`
	OrderID      int64         `json:"order_id"`
	Desc         string        `json:"desc"`
	Source       string        `json:"source"`
	Attribute    int32         `json:"-"` // NOTE: not allow user
	OpenElec     int8          `json:"open_elec"`
	MissionID    int           `json:"mission_id"`
	FromIP       int64         `json:"from_ip"`
	IPv6         []byte        `json:"ipv6"`
	UpFrom       int8          `json:"up_from"`
	BizFrom      int8          `json:"biz_from"`
	DTime        time.Time     `json:"dtime"`
	Videos       []*VideoParam `json:"videos"`
	Body         string        `json:"body,omitempty"`
	CodeMode     bool          `json:"code_mode,omitempty"`
	DescFormatID int           `json:"desc_format_id,omitempty"`
	Dynamic      string        `json:"dynamic,omitempty"`
	Porder       *Porder       `json:"porder"`
	Lang         string        `json:"lang"`
	Watermark    *Watermark    `json:"watermark"`
	Geetest      *Geetest      `json:"geetest"`
	LotteryID    int64         `json:"lottery_id"`
	Subtitle     *Subtitle     `json:"subtitle"`
	Pay          *Pay          `json:"pay"`
	UgcPay       int8          `json:"ugcpay"` // videoup-service 需要按照这个字段转成attribute
	FollowMids   []int64       `json:"follow_mids"`
	PoiObj       *PoiObj       `json:"poi_object"`
	Staffs       []*Staff      `json:"staffs"`
	HandleStaff  bool          `json:"handle_staff"`
	Vote         *Vote         `json:"vote"`
}

// Vote str
type Vote struct {
	VoteID    int64  `json:"vote_id"`
	VoteTitle string `json:"vote_title"`
}

// Pay str
type Pay struct {
	Open           int8   `json:"open"`
	Price          int    `json:"price"`
	ProtocolID     string `json:"protocol_id"`
	ProtocolAccept int8   `json:"protocol_accept"`
	RefuseUpdate   bool   `json:"-"`
}

// Subtitle str only for web add and edit
type Subtitle struct {
	Open int8   `json:"open"`
	Lan  string `json:"lan"`
}

// Geetest str
type Geetest struct {
	Challenge string `json:"challenge"`
	Validate  string `json:"validate"`
	Seccode   string `json:"seccode"`
	Success   int    `json:"success"`
}

// Watermark str
type Watermark struct {
	State int8 `json:"state"`
	Ty    int8 `json:"type"`
	Pos   int8 `json:"position"`
}

// Porder str
// new porder, ads provoder
type Porder struct {
	FlowID     uint   `json:"flow_id"`     // 0/1 是否确实参加了广告平台
	IndustryID int64  `json:"industry_id"` // 2 (游戏)
	BrandName  string `json:"brand_name"`  // FGO游戏
	BrandID    int64  `json:"brand_id"`    // 2
	Official   int8   `json:"official"`    // 0/1
	ShowType   string `json:"show_type"`   // 2,3,4
}

// VideoParam str
type VideoParam struct {
	Title    string  `json:"title"`
	Desc     string  `json:"desc"`
	Filename string  `json:"filename"`
	Cid      int64   `json:"cid"`
	Sid      int64   `json:"sid"`
	Editor   *Editor `json:"editor"`
}

// Editor str
type Editor struct {
	CID    int64 `json:"cid"`
	UpFrom int8  `json:"upfrom"` // filled by backend
	// ids set
	Filters         interface{} `json:"filters"`          // 滤镜
	Fonts           interface{} `json:"fonts"`            //字体
	Subtitles       interface{} `json:"subtitles"`        //字幕
	Bgms            interface{} `json:"bgms"`             //bgm
	Stickers        interface{} `json:"stickers"`         //3d拍摄贴纸
	VideoupStickers interface{} `json:"videoup_stickers"` //2d投稿贴纸
	Transitions     interface{} `json:"trans"`            //视频转场特效
	// add from app535
	Themes     interface{} `json:"themes"`     //编辑器的主题使用相关
	Cooperates interface{} `json:"cooperates"` //拍摄之稿件合拍
	// switch env 0/1
	AudioRecord  int8 `json:"audio_record"`  //录音
	Camera       int8 `json:"camera"`        //拍摄
	Speed        int8 `json:"speed"`         //变速
	CameraRotate int8 `json:"camera_rotate"` //摄像头翻转
	// count from app536
	PicCount   uint16 `json:"pic_count"`   // 图片个数
	VideoCount uint16 `json:"video_count"` // 视频个数
}

// VideoExpire str
type VideoExpire struct {
	Filename string `json:"filename"`
	Expire   int64  `json:"expire"`
}

// CreatorParam struct
type CreatorParam struct {
	Aid      int64  `form:"aid" validate:"required"`
	Title    string `form:"title" validate:"required"`
	Desc     string `form:"desc" validate:"required"`
	Tag      string `form:"tag" validate:"required"`
	OpenElec int8   `form:"open_elec"`
	Build    string `form:"build" validate:"required"`
	Platform string `form:"platform" validate:"required"`
}

// Staff 稿件提交时的联合投稿人
type Staff struct {
	Title string `json:"title"`
	Mid   int64  `json:"mid"`
}

// StaffView Archive staff
type StaffView struct {
	ID         int64  `json:"id"`
	ApMID      int64  `json:"apply_staff_mid"`
	ApTitle    string `json:"apply_title"`
	ApAID      int64  `json:"apply_aid"`
	ApType     int    `json:"apply_type"`
	ApState    int    `json:"apply_state"`
	ApStaffID  int64  `json:"apply_asid"` //Staff表的主键ID
	StaffState int    `json:"staff_state"`
	StaffTitle string `json:"staff_title"`
}

// ForbidMultiVideoType fun
// 欧美电影,日本电影,国产电影,其他国家
func (ap *ArcParam) ForbidMultiVideoType() bool {
	return ap.TypeID == 145 || ap.TypeID == 146 || ap.TypeID == 147 || ap.TypeID == 83
}

// ForbidAddVideoType fun
// 连载剧集：15 完结剧集：34 电视剧相关：128 电影相关：82
func (ap *ArcParam) ForbidAddVideoType() bool {
	return ap.TypeID == 15 || ap.TypeID == 34 || ap.TypeID == 128 || ap.TypeID == 82
}

// ForbidCopyrightAndTypes fun
// // 32 完结动画; 33 连载动画
func (ap *ArcParam) ForbidCopyrightAndTypes() bool {
	return (ap.Copyright == CopyrightOriginal) && (ap.TypeID == 32 || ap.TypeID == 33)
}

// EmptyVideoEditInfo fn
func (ap *ArcParam) EmptyVideoEditInfo() {
	if (ap.UpFrom != UpFromAPPiOS) && (ap.UpFrom != UpFromAPPAndroid) {
		for _, v := range ap.Videos {
			v.Editor = nil
		}
	}
}

// NilPoiObj fn
func (ap *ArcParam) NilPoiObj() {
	if (ap.UpFrom != UpFromAPPiOS) && (ap.UpFrom != UpFromAPPAndroid) {
		ap.PoiObj = nil
	}
}

// DisableVideoDesc fn
func (ap *ArcParam) DisableVideoDesc(vs []*Video) {
	nvsMap := make(map[string]string)
	for _, v := range vs {
		nvsMap[v.Filename] = v.Desc
	}
	for _, pv := range ap.Videos {
		if nvFilename, ok := nvsMap[pv.Filename]; ok {
			pv.Desc = nvFilename
		}
	}
}
