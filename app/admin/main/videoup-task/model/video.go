package model

import (
	"time"
)

//video status & attr.
const (
	VideoStatusOpen    = int16(0)
	VideoStatusOrange  = int16(10000)
	VideoStatusRecycle = int16(-2)
	VideoStatusLock    = int16(-4)

	VideoStatusDelete = -100
	ArcStateDelete    = -100
	RLStateDelete     = -100

	CopyrightOriginal  = int8(1)
	VideoXcodeSDFinish = int8(2)
	VideoXcodeHDFinish = int8(4)

	AttrBitNoRank = uint(0) // NOTE: double write for archive_forbid
	// AttrBitNoDynamic 动态禁止
	AttrBitNoDynamic = uint(1) // NOTE: double write for archive_forbid
	// AttrBitNoWeb 禁止网页输出
	AttrBitNoWeb = uint(2)
	// AttrBitNoMobile 禁止客户端列表
	AttrBitNoMobile = uint(3)
	// AttrBitNoSearch 搜索禁止
	AttrBitNoSearch = uint(4)
	// AttrBitOverseaLock 海外禁止
	AttrBitOverseaLock = uint(5)
	// AttrBitNoRecommend 禁止推荐
	AttrBitNoRecommend = uint(6) // NOTE: double write for archive_forbid
	// AttrBitNoReprint 禁止转载
	AttrBitNoReprint = uint(7)
	// AttrBitHasHD5 是否高清
	AttrBitHasHD5 = uint(8)
	// AttrBitIsPGC 是否PGC稿件
	AttrBitIsPGC = uint(9)
	// AttrBitAllowBp 允许承包
	AttrBitAllowBp = uint(10)
	// AttrBitIsBangumi 是否番剧
	AttrBitIsBangumi = uint(11)
	// AttrBitIsPorder 是否私单
	AttrBitIsPorder = uint(12)
	// AttrBitLimitArea 是否限制地区
	AttrBitLimitArea = uint(13)
	// AttrBitAllowTag 允许其他人添加tag
	AttrBitAllowTag = uint(14)
	// AttrBitJumpURL 跳转
	AttrBitJumpURL = uint(16)
	// AttrBitIsMovie 是否影视
	AttrBitIsMovie = uint(17)
	// AttrBitBadgepay 付费
	AttrBitBadgepay = uint(18)
	AttrBitPushBlog = uint(20)
)

//qa audit status & attr.
var (
	QAAuditStatus = map[int16]string{
		VideoStatusOpen:    "开放浏览",
		VideoStatusOrange:  "会员可见",
		VideoStatusRecycle: "打回",
		VideoStatusLock:    "锁定",
	}

	VideoAttribute = map[uint]string{
		AttrBitNoRank:      "norank",
		AttrBitNoDynamic:   "nodynamic",
		AttrBitNoWeb:       "noweb",
		AttrBitNoMobile:    "nomobile",
		AttrBitNoSearch:    "nosearch",
		AttrBitOverseaLock: "oversea_block",
		AttrBitNoRecommend: "norecommend",
		AttrBitNoReprint:   "no_reprint",
		AttrBitHasHD5:      "hd",
		AttrBitIsPGC:       "is_pgc",
		AttrBitAllowBp:     "allow_bp",
		AttrBitIsBangumi:   "bangumi",
		AttrBitIsPorder:    "is_porder",
		AttrBitLimitArea:   "limit_area",
		AttrBitAllowTag:    "allow_tag",
		AttrBitJumpURL:     "j",
		AttrBitIsMovie:     "is_movie",
		AttrBitBadgepay:    "badgepay",
		AttrBitPushBlog:    "push_blog",
	}
)

//Video video info
type Video struct {
	ID            int64            `json:"id"`
	AID           int64            `json:"aid"`
	CID           int64            `json:"cid"`
	MID           int64            `json:"mid"`
	Copyright     int8             `json:"copyright"`
	TypeID        int64            `json:"type_id"`
	Status        int16            `json:"status"`
	Attribute     int32            `json:"attribute"`
	XcodeState    int8             `json:"xcode_state"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Filename      string           `json:"filename"`
	TagID         int64            `json:"tag_id"`
	Reason        string           `json:"reason"`
	Note          string           `json:"note"`
	AttributeList map[string]int32 `json:"attribute_list"`
	Encoding      int32            `json:"encoding"`
}

//AttributeList get attr as map
func AttributeList(attr int32) (list map[string]int32) {
	list = map[string]int32{}
	for bit, name := range VideoAttribute {
		list[name] = int32(((attr >> bit) & 1))
	}
	return
}

// AttrSet video Attr set
func (v *Video) AttrSet(attr int32, bit uint) {
	v.Attribute = v.Attribute&(^(1 << bit)) | (attr << bit)
}

// ArcVideo is archive_video model.
type ArcVideo struct {
	ID           int64     `json:"-"`
	Aid          int64     `json:"aid"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Filename     string    `json:"filename"`
	SrcType      string    `json:"-"`
	Cid          int64     `json:"cid"`
	Duration     int64     `json:"-"`
	Filesize     int64     `json:"-"`
	Resolutions  string    `json:"-"`
	Index        int       `json:"index"`
	Playurl      string    `json:"-"`
	Status       int16     `json:"status"`
	StatusDesc   string    `json:"status_desc"`
	FailCode     int8      `json:"fail_code"`
	FailDesc     string    `json:"fail_desc"`
	XcodeState   int8      `json:"xcode"`
	Attribute    int32     `json:"-"`
	RejectReason string    `json:"reject_reason"`
	WebLink      string    `json:"weblink"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
}
