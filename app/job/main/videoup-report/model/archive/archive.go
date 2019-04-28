package archive

import (
	"sync"
	"time"
)

const (
	//StateOpen state open
	StateOpen = 0
	//StateOrange 橙色通过
	StateOrange = 1

	//StateForbidWait 待审
	StateForbidWait = -1
	//StateForbidRecicle 打回
	StateForbidRecicle = -2
	//StateForbidPolice 网警锁定
	StateForbidPolice = -3
	//StateForbidLock 锁定
	StateForbidLock = -4
	//StateForbidFixed 修复待审
	StateForbidFixed = -6
	//StateForbidLater 暂缓待审
	StateForbidLater = -7
	//StateForbidXcodeFail 转码失败
	StateForbidXcodeFail = -16
	//StateForbidSubmit 创建提交
	StateForbidSubmit = -30
	//StateForbidUserDelay 定时
	StateForbidUserDelay = -40
	//StateForbidUpDelete 删除
	StateForbidUpDelete = -100
	//RoundBegin 开始流转
	RoundBegin = 0
	//RoundAuditSecond 二审：选定分区的多P稿件 及 PGC/活动的单P多P稿件
	RoundAuditSecond = 10
	//RoundAuditThird 三审：选定分区/PGC/活动 的单P多P稿件
	RoundAuditThird = 20
	//RoundReviewFirst 分区回查：粉丝小于配置阈值 如 5000 且 指定分区
	RoundReviewFirst = 30
	//RoundReviewFirstWaitTrigger 点击/粉丝 等待触发中间状态，7天内达到阈值进列表，未达到自动变99
	RoundReviewFirstWaitTrigger = 31
	//RoundReviewSecond 社区回查：粉丝大于配置阈值 如 5000 或 优质高危up
	RoundReviewSecond = 40
	//RoundTriggerFans  粉丝回查：粉丝量达到配置阈值
	RoundTriggerFans = 80
	//RoundTriggerClick 点击回查：点击量达到配置阈值
	RoundTriggerClick = 90
	//RoundEnd 流转结束
	RoundEnd = 99
	//AccessDefault access
	AccessDefault = int16(0)
	//AccessMember 会员可见
	AccessMember = int16(10000)
	//CopyrightUnknow copyright
	CopyrightUnknow = 0
	//CopyrightOriginal 原创
	CopyrightOriginal = 1
	//CopyrightCopy 转载
	CopyrightCopy = 2

	//AttrYes attribute yes
	AttrYes = int32(1)
	//AttrNo attribute no
	AttrNo = int32(0)
	//AttrBitNoRank 禁止排序
	AttrBitNoRank = uint(0)
	//AttrBitNoDynamic 禁止动态
	AttrBitNoDynamic = uint(1)
	//AttrBitNoWeb 禁止web
	AttrBitNoWeb = uint(2)
	//AttrBitNoMobile 禁止手机端
	AttrBitNoMobile = uint(3)
	//AttrBitNoSearch 禁止搜索
	AttrBitNoSearch = uint(4)
	//AttrBitOverseaLock 禁止海外
	AttrBitOverseaLock = uint(5)
	//AttrBitNoRecommend 禁止推荐
	AttrBitNoRecommend = uint(6)
	// AttrBitHideCoins     = uint(7)

	//AttrBitHasHD5 是否高清
	AttrBitHasHD5 = uint(8)
	// AttrBitVisitorDm     = uint(9)

	//AttrBitAllowBp 允许承包
	AttrBitAllowBp = uint(10)
	//AttrBitIsBangumi 番剧
	AttrBitIsBangumi = uint(11)
	//AttrBitIsPOrder 是否私单
	AttrBitIsPOrder = uint(12)

	//AttrBitHideClick 点击
	AttrBitHideClick = uint(13)
	//AttrBitAllowTag 允许操作tag
	AttrBitAllowTag = uint(14)
	// AttrBitIsFromArcApi = uint(15)

	//AttrBitJumpURL 跳转
	AttrBitJumpURL = uint(16)
	//AttrBitIsMovie is movie
	AttrBitIsMovie = uint(17)
	//AttrBitBadgepay 付费
	AttrBitBadgepay = uint(18)

	//ReplyDefault 默认评论状态
	ReplyDefault = int64(-1)
	//ReplyOn 开评论
	ReplyOn = int64(0)
	//ReplyOff 关评论
	ReplyOff = int64(1)

	//LogBusJob 稿件后台任务日志bus
	LogBusJob = 211
	//LogTypeReply 稿件后台任务type评论
	LogTypeReply = 1
)

//ReplyState 评论开关状态
var ReplyState = []int64{
	ReplyDefault,
	ReplyOn,
	ReplyOff,
}

//ReplyDesc 评论状态描述
var ReplyDesc = map[int64]string{
	ReplyDefault: "未知状态",
	ReplyOn:      "开",
	ReplyOff:     "关",
}

//UpInfo up info
type UpInfo struct {
	Nw  *Archive
	Old *Archive
}

// Oper is archive operate model.
type Oper struct {
	ID        int64     `json:"id"`
	AID       int64     `json:"aid"`
	UID       int64     `json:"uid"`
	TypeID    int16     `json:"typeid"`
	State     int       `json:"state"`
	Content   string    `json:"-"`
	Round     int8      `json:"round"`
	Attribute int32     `json:"attribute"`
	LastID    int64     `json:"last_id"`
	Remark    string    `json:"-"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// ArcMoveTypeCache archive move typeid count
type ArcMoveTypeCache struct {
	Data map[int8]map[int16]map[string]int
	sync.Mutex
}

// ArcRoundFlowCache archive round flow record
type ArcRoundFlowCache struct {
	Data map[int8]map[int64]map[string]int
	sync.Mutex
}

//Archive archive
type Archive struct {
	ID        int64  `json:"id"`
	AID       int64  `json:"aid"` //result库binlog={id:0,aid:xxx}
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Desc      string `json:"desc"`
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
	PTime     string `json:"ptime"`
}

//IsSyncState can archive sync
func (a *Archive) IsSyncState() bool {
	if a.State >= 0 || a.State == StateForbidUserDelay || a.State == StateForbidUpDelete || a.State == StateForbidRecicle || a.State == StateForbidPolice ||
		a.State == StateForbidLock {
		return true
	}
	return false
}

//ArgStat arg state
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

//NormalState normal state
func NormalState(state int) bool {
	return state == StateOpen || state == StateOrange
}

//Type archive_type
type Type struct {
	ID   int16  `json:"id"`
	PID  int16  `json:"pid"`
	Name string `json:"name"`
}

// StateMean the mean for archive state
var StateMean = map[int]string{
	StateOpen:   "开放浏览",
	StateOrange: "橙色通过",
	// forbid state
	StateForbidWait:    "待审",
	StateForbidRecicle: "打回",
	StateForbidPolice:  "网警锁定",
	StateForbidLock:    "锁定稿件",
	StateForbidFixed:   "修复待审",
	StateForbidLater:   "暂缓审核",
	//StateForbidAdminDelay: "延迟发布",
	StateForbidXcodeFail: "转码失败",
	StateForbidSubmit:    "创建提交",
	StateForbidUserDelay: "用户定时发布",
	StateForbidUpDelete:  "UP主删除",
}
