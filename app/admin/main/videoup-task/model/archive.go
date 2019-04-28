package model

import (
	"go-common/library/time"
)

const (
	// UpFromWeb 网页上传
	UpFromWeb = int8(0)
	// UpFromPGC PGC上传
	UpFromPGC = int8(1)
	// UpFromWindows Windows客户端上传
	UpFromWindows = int8(2)
	// UpFromAPP APP上传
	UpFromAPP = int8(3)
	// UpFromMAC Mac客户端上传
	UpFromMAC = int8(4)
	// UpFromSecretPGC 机密PGC上传
	UpFromSecretPGC = int8(5)
	// UpFromCoopera 合作方嵌套
	UpFromCoopera = int8(6)
	// UpFromCreator 创作姬上传
	UpFromCreator = int8(7)
	// UpFromAndroid 安卓上传
	UpFromAndroid = int8(8)
	// UpFromIOS IOS上传
	UpFromIOS = int8(9)

	// AttrYes attribute yes
	AttrYes = int32(1)
	// AttrNo attribute no
	AttrNo = int32(0)

	// StateForbidUpDelete 用户删除
	StateForbidUpDelete = int8(-100)
)

var (
	_bits = map[uint]string{
		AttrBitNoRank:      "排行禁止",
		AttrBitNoDynamic:   "动态禁止",
		AttrBitNoWeb:       "禁止web端输出",
		AttrBitNoMobile:    "禁止移动端输出",
		AttrBitNoSearch:    "禁止搜索",
		AttrBitOverseaLock: "海外禁止",
		AttrBitNoRecommend: "推荐禁止",
		AttrBitNoReprint:   "禁止转载",
		AttrBitHasHD5:      "高清1080P",
		// AttrBitVisitorDm:     AttrBitVisitorDm,
		AttrBitIsPGC:     "PGC",
		AttrBitAllowBp:   "允许承包",
		AttrBitIsBangumi: "番剧",
		AttrBitIsPorder:  "是否私单",
		AttrBitLimitArea: "是否地区限制",
		AttrBitAllowTag:  "允许操作TAG",
		// AttrBitIsFromArcAPI: AttrBitIsFromArcAPI,
		AttrBitJumpURL:  "跳转",
		AttrBitIsMovie:  "电影",
		AttrBitBadgepay: "付费",
		AttrBitPushBlog: "禁止粉丝动态",
	}

	_upFromTypes = map[int8]string{
		UpFromWeb:       "网页上传",
		UpFromPGC:       "PGC上传",
		UpFromWindows:   "Windows客户端上传",
		UpFromAPP:       "APP上传",
		UpFromMAC:       "Mac客户端上传",
		UpFromSecretPGC: "机密PGC上传",
		UpFromCoopera:   "合作方嵌套",
		UpFromCreator:   "创作姬上传",
		UpFromAndroid:   "安卓上传",
		UpFromIOS:       "IOS上传",
	}
)

// BitDesc return bit desc.
func BitDesc(bit uint) (desc string) {
	return _bits[bit]
}

// Archive is archive model.
type Archive struct {
	Aid          int64     `json:"aid"`
	Mid          int64     `json:"mid"`
	TypeID       int16     `json:"tid"`
	HumanRank    int       `json:"-"`
	Title        string    `json:"title"`
	Author       string    `json:"-"`
	Cover        string    `json:"cover"`
	RejectReason string    `json:"reject_reason"`
	Tag          string    `json:"tag"`
	Duration     int64     `json:"duration"`
	Copyright    int8      `json:"copyright"`
	Desc         string    `json:"desc"`
	MissionID    int64     `json:"mission_id"`
	Round        int8      `json:"-"`
	Forward      int64     `json:"-"`
	Attribute    int32     `json:"attribute"`
	Access       int16     `json:"-"`
	State        int8      `json:"state"`
	Source       string    `json:"source"`
	NoReprint    int32     `json:"no_reprint"`
	OrderID      int64     `json:"order_id"`
	Dynamic      string    `json:"dynamic"`
	DTime        time.Time `json:"dtime"`
	PTime        time.Time `json:"ptime"`
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
}

// Type is archive type info
type Type struct {
	ID   int16  `json:"id"`
	PID  int16  `json:"pid"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

// UpFrom get upfrom desc
func UpFrom(ufID int8) string {
	return _upFromTypes[ufID]
}
