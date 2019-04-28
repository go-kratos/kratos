package archive

const (
	StateOpen   = 0
	StateOrange = 1

	StateForbidWait      = -1
	StateForbidRecicle   = -2
	StateForbidPolice    = -3
	StateForbidLock      = -4
	StateForbidFixed     = -6
	StateForbidLater     = -7
	StateForbidSubmit    = -30
	StateForbidUserDelay = -40
	StateForbidUpDelete  = -100
	// round
	RoundBegin                  = 0
	RoundAuditSecond            = 10 // 二审：选定分区的多P稿件 及 PGC/活动的单P多P稿件
	RoundAuditThird             = 20 // 三审：选定分区/PGC/活动 的单P多P稿件
	RoundReviewFirst            = 30 // 分区回查：粉丝小于配置阈值 如 5000 且 指定分区
	RoundReviewFirstWaitTrigger = 31 // 点击/粉丝 等待触发中间状态，7天内达到阈值进列表，未达到自动变99
	RoundReviewSecond           = 40 // 社区回查：粉丝大于配置阈值 如 5000 或 优质高危up
	RoundTriggerFans            = 80 // 粉丝回查：粉丝量达到配置阈值
	RoundTriggerClick           = 90 // 点击回查：点击量达到配置阈值
	RoundEnd                    = 99
	// access
	AccessDefault = int16(0)
	AccessMember  = int16(10000)
	// copyright
	CopyrightUnknow   = 0
	CopyrightOriginal = 1
	CopyrightCopy     = 2

	// attribute yes and no
	AttrYes = int32(1)
	AttrNo  = int32(0)
	// attribute bit
	AttrBitNoRank      = uint(0)
	AttrBitNoIndex     = uint(1)
	AttrBitNoWeb       = uint(2)
	AttrBitNoMobile    = uint(3)
	AttrBitNoSearch    = uint(4)
	AttrBitOverseaLock = uint(5)
	AttrBitNoRecommend = uint(6)
	// AttrBitHideCoins     = uint(7)
	AttrBitHasHD5 = uint(8)
	// AttrBitVisitorDm     = uint(9)
	AttrBitAllowBp   = uint(10)
	AttrBitIsBangumi = uint(11)
	// AttrBitAllowDownload = uint(12)
	AttrBitLimitArea = uint(13)
	AttrBitAllowTag  = uint(14)
	// AttrBitIsFromArcApi = uint(15)
	AttrBitJumpUrl       = uint(16)
	AttrBitIsMovie       = uint(17)
	AttrBitBadgepay      = uint(18)
	AttrBitIsCooperation = uint(24)
)

type UpInfo struct {
	Nw  *Archive
	Old *Archive
}

// archive
type Archive struct {
	ID        int64  `json:"id"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	AreaLimit int8   `json:"arealimit"`
	State     int    `json:"state"`
	Author    string `json:"author"`
	Access    int    `json:"access"`
	Forward   int    `json:"forward"`
	PubTime   string `json:"pubtime"`
	Reason    string `json:"reject_reason"`
	Round     int8   `json:"round"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}

func (a *Archive) IsSyncState() bool {
	if a.State >= 0 || a.State == StateForbidUserDelay || a.State == StateForbidUpDelete || a.State == StateForbidRecicle || a.State == StateForbidPolice ||
		a.State == StateForbidLock {
		return true
	}
	return false
}

type ArgStat struct {
	Aid    int64
	Field  int
	Value  int
	RealIP string
}

// AttrVal get attribute value.
func (a *Archive) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// Staff is
type Staff struct {
	Aid   int64  `json:"aid"`
	Mid   int64  `json:"mid"`
	Title string `json:"title"`
	Ctime string `json:"ctime"`
	Mtime string `json:"mtime"`
}
