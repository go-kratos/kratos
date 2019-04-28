package archive

import (
	"go-common/app/service/main/archive/api"
	a "go-common/app/service/main/archive/model/archive"
	"go-common/library/time"
)

var (
	// VjInfo 审核当前拥挤状态的数据映射
	VjInfo = map[int8]*VideoJam{
		0: {
			Level:   0,
			State:   "状态正在计算中",
			Comment: "状态正在计算中",
		},
		1: {
			Level:   1,
			State:   "畅通",
			Comment: "预计稿件过审时间小于20分钟（剧集、活动投稿除外）",
		},
		2: {
			Level:   2,
			State:   "繁忙",
			Comment: "预计稿件过审时间小于40分钟（剧集、活动投稿除外）",
		},
		3: {
			Level:   3,
			State:   "拥挤",
			Comment: "预计稿件过审时间小于60分钟（剧集、活动投稿除外）",
		},
		4: {
			Level:   4,
			State:   "爆满",
			Comment: "预计稿件过审时间小于120分钟（剧集、活动投稿除外）",
		},
		5: {
			Level:   5,
			State:   "阻塞",
			Comment: "预计稿件过审时间大于120分钟（剧集、活动投稿除外）",
		},
	}
)

const (
	// CopyrightOrigin 自制
	CopyrightOrigin = int64(1)
	// CopyrightReprint 转载
	CopyrightReprint = int64(2)
	// TagPredictFromWeb web tag推荐
	TagPredictFromWeb = int8(0)
	// TagPredictFromAPP app tag推荐
	TagPredictFromAPP = int8(1)
	// TagPredictFromWindows windows tag推荐
	TagPredictFromWindows = int8(2)
)

// Const State
const (
	// attribute yes and no
	AttrYes = int32(1)
	AttrNo  = int32(0)
	// attribute bit
	AttrBitNoRank       = uint(0)
	AttrBitNoDynamic    = uint(1)
	AttrBitNoWeb        = uint(2)
	AttrBitNoMobile     = uint(3)
	AttrBitNoSearch     = uint(4)
	AttrBitOverseaLock  = uint(5)
	AttrBitNoRecommend  = uint(6)
	AttrBitNoReprint    = uint(7)
	AttrBitHasHD5       = uint(8)
	AttrBitIsPGC        = uint(9)
	AttrBitAllowBp      = uint(10)
	AttrBitIsBangumi    = uint(11)
	AttrBitIsPorder     = uint(12)
	AttrBitLimitArea    = uint(13)
	AttrBitAllowTag     = uint(14)
	AttrBitIsFromArcAPI = uint(15) // TODO: delete
	AttrBitJumpURL      = uint(16)
	AttrBitIsMovie      = uint(17)
	AttrBitBadgepay     = uint(18)
	AttrBitIsJapan      = uint(19) //日文稿件
	AttrBitNoPushBplus  = uint(20) //是否动态禁止
	AttrBitParentMode   = uint(21) //家长模式
	AttrBitUGCPay       = uint(22) //UGC付费
	AttrBitHasBGM       = uint(23) //稿件带有BGM
	AttrBitIsCoop       = uint(24) //联合投稿
)

// OldArchiveVideoAudit archive with audit.
// NOTE: old struct, will delete!!!
type OldArchiveVideoAudit struct {
	*api.Arc
	RejectReson string           `json:"reject,omitempty"`
	Dtime       time.Time        `json:"dtime,omitempty"`
	VideoAudits []*OldVideoAudit `json:"video_audit,omitempty"`
	StateDesc   string           `json:"state_desc"`
	StatePanel  int              `json:"state_panel"`
	ParentTName string           `json:"parent_tname"`
	Attrs       *Attrs           `json:"attrs"`
	UgcPay      int8             `json:"ugcpay"`
}

// OldVideoAudit video audit.
// NOTE: old struct, will delete!!!
type OldVideoAudit struct {
	Reason     string `json:"reason,omitempty"`
	Eptitle    string `json:"eptitle,omitempty"`
	IndexOrder int    `json:"index_order"`
}

// ArcVideoAudit archive video audit.
type ArcVideoAudit struct {
	*ArcVideo
	Stat        *api.Stat `json:"stat"`
	StatePanel  int       `json:"state_panel"`
	ParentTName string    `json:"parent_tname"`
	TypeName    string    `json:"typename"`
	OpenAppeal  int64     `json:"open_appeal"`
}

// Flow type
type Flow struct {
	ID     uint   `json:"id"`
	Remark string `json:"remark"`
}

// Porder type
type Porder struct {
	ID         int64     `json:"id"`
	AID        int64     `json:"aid"`
	IndustryID int64     `json:"industry_id"`
	BrandID    int64     `json:"brand_id"`
	BrandName  string    `json:"brand_name"`
	Official   int8      `json:"is_official"`
	ShowType   string    `json:"show_type"`
	Advertiser string    `json:"advertiser"`
	Agent      string    `json:"agent"`
	Ctime      time.Time `json:"ctime,omitempty"`
	Mtime      time.Time `json:"mtime,omitempty"`
}

// Staff type
type Staff struct {
	ID         int64  `json:"id"`
	AID        int64  `json:"aid"`
	MID        int64  `json:"mid"`
	StaffMID   int64  `json:"staff_mid"`
	StaffTitle string `json:"staff_title"`
}

// StaffApply type
type StaffApply struct {
	ID            int64  `json:"id"`
	Type          int8   `json:"apply_type"`
	ASID          int64  `json:"apply_as_id"`
	ApplyAID      int64  `json:"apply_aid"`
	ApplyUpMID    int64  `json:"apply_up_mid"`
	ApplyStaffMID int64  `json:"apply_staff_mid"`
	ApplyTitle    string `json:"apply_title"`
	ApplyTitleID  int64  `json:"apply_title_id"`
	State         int8   `json:"apply_state"`
	StaffState    int8   `json:"staff_state"`
	StaffTitle    string `json:"staff_title"`
}

// Commercial type
type Commercial struct {
	AID      int64 `json:"aid"`
	POrderID int64 `json:"porder_id"` // 私
	OrderID  int64 `json:"order_id"`  // 商
	GameID   int64 `json:"game_id"`
	//IndustryID int64     `json:"industry_id"` // open after mall
	//BrandID    int64     `json:"brand_id"`
}

// InMovieType judge type for pubdate
func InMovieType(tid int16) bool {
	return tid == 83 || tid == 145 || tid == 146 || tid == 147
}

// StatePanel judge archive state for app panel.
func StatePanel(s int8) (st int) {
	if s == a.StateForbidWait ||
		s == a.StateForbidFixed ||
		s == a.StateForbidLater ||
		s == a.StateForbidAdminDelay ||
		s == a.StateForbidSubmit ||
		s == a.StateForbidUserDelay {
		st = 1 //处理中
	} else if s == a.StateForbidRecicle {
		st = 2 //退回可编辑
	} else if s == a.StateForbidPolice || s == a.StateForbidLock {
		st = 3 //退回全部不可编辑
	} else if s == a.StateForbidXcodeFail {
		st = 4 //退回分区不可编辑
	} else {
		st = 0 //正常开放
	}
	return
}

// IsCloseState judge arc state.
func IsCloseState(s int) bool {
	return s == -2 || s == -4 || s == -5 || s == -14
}

// ShortDesc cut down to short desc for adapter app and windows clients
func ShortDesc(desc string) string {
	rs := []rune(desc)
	length := len(rs)
	max := 250
	if length < max {
		max = length
	}
	return string(rs[:max])
}

// AttrVal get attribute.
func (a *ArcVideoAudit) AttrVal(bit uint) int32 {
	return (a.Archive.Attribute >> bit) & int32(1)
}

// IsOwner fn
func (a *ArcVideoAudit) IsOwner(mid int64) int8 {
	if a.Archive.Mid == mid {
		return 1
	}
	return 0
}
