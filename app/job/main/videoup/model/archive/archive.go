package archive

import (
	xtime "go-common/library/time"
)

//# state 稿件状态
//# 0 开放浏览 , -1 待审 , -2 打回稿件回收站 , -3 网警锁定删除
//# -4 锁定稿件 , -6 修复待审 , -7 暂缓审核 , -9 等待转码
//# -10 延迟发布 , -11 视频源待修 , -13 允许评论待审 , -15 分发中
//# -16 转码失败, -30 创建提交, -40 用户定时发布, -100 UP主删除

//# attribute bit位置
//# 右1 - norank 禁止排名 , 右2 - noindex 首页禁止 , 右3 - noweb 禁止网页端输出 , 右4 - nomobile 禁止移动端输出
//# 右5 - nosearch 禁止移动端未登录搜索 , 右6 - overseas 海外禁止 , 右7 - nocount 不计算点击
//# 右8 - hidecoins 禁止显示硬币 , 右9 - is_hdflv2 1080p 是否有高清1080p , 右10 - dm 是否允许游客发弹幕
//# 右11 - allow_bp 是否允许投放bp , 右12 - 是否番剧 , 右13 - allow_download 是否允许下载
//# 右14 - hideclick 是否隐藏点击数, 右15 - allow_tag 允许添加tag, 右16 - 是否api投稿
//# 右17 - jump 是否跳转别的url, 右18 - 是否付费影视, 右19 - 付费标识

//# access 会员状态
//# 10000 普通会员 , 15000 新番搬运 , 20000 字幕君, 25000 VIP , 30000 真职人
//# 35000 橙色通过开放浏览 , 40000 橙色通过会员浏览

const (
	// StateOpen 开发浏览
	StateOpen = int8(0)
	// StateOrange 橙色通过
	StateOrange = int8(1)
	// StateForbidWait 待审
	StateForbidWait = int8(-1)
	// StateForbidRecicle 打回
	StateForbidRecicle = int8(-2)
	// StateForbidPolice 网警锁定
	StateForbidPolice = int8(-3)
	// StateForbidLock 被锁定
	StateForbidLock = int8(-4)
	// StateForbidFackLock 管理员锁定（可浏览）
	StateForbidFackLock = int8(-5)
	// StateForbidFixed 修复待审
	StateForbidFixed = int8(-6)
	// StateForbidLater 暂缓审核
	StateForbidLater = int8(-7)
	// StateForbidPatched 补档待审
	StateForbidPatched = int8(-8)
	// StateForbidWaitXcode 等待转码
	StateForbidWaitXcode = int8(-9)
	// StateForbidAdminDelay 延迟审核
	StateForbidAdminDelay = int8(-10)
	// StateForbidFixing 视频源待修
	StateForbidFixing = int8(-11)
	// StateForbidStorageFail 转储失败
	StateForbidStorageFail = int8(-12)
	// StateForbidOnlyComment 允许评论待审
	StateForbidOnlyComment = int8(-13)
	// StateForbidTmpRecicle 临时回收站
	StateForbidTmpRecicle = int8(-14)
	// StateForbidDispatch 分发中
	StateForbidDispatch = int8(-15)
	// StateForbidXcodeFail 转码失败
	StateForbidXcodeFail = int8(-16)
	// StateForbidSubmit 创建已提交
	StateForbidSubmit = int8(-30)
	// StateForbidUserDelay 定时发布
	StateForbidUserDelay = int8(-40)
	// StateForbidUpDelete 用户删除
	StateForbidUpDelete = int8(-100)

	// RoundBegin 一审阶段
	RoundBegin = int8(0)
	// RoundAuditSecond 二审：选定分区的多P稿件 及 PGC/活动的单P多P稿件
	RoundAuditSecond = int8(10)
	// RoundAuditThird 三审：选定分区/PGC/活动 的单P多P稿件
	RoundAuditThird = int8(20)
	// RoundReviewFlow 私单回查：私单ID大于0
	RoundReviewFlow = int8(21)
	//RoundAuditUGCPayFlow 付费待审
	RoundAuditUGCPayFlow = int8(24)
	// RoundReviewFirst 分区回查：粉丝小于配置阈值 如 5000 且 指定分区
	RoundReviewFirst = int8(30)
	// RoundReviewFirstWaitTrigger 点击/粉丝 等待触发中间状态，7天内达到阈值进列表，未达到自动变99
	RoundReviewFirstWaitTrigger = int8(31)
	// RoundReviewSecond 社区回查：粉丝大于配置阈值 如 5000 或 优质高危up
	RoundReviewSecond = int8(40)
	// RoundTriggerFans 粉丝回查：粉丝量达到配置阈值
	RoundTriggerFans = int8(80)
	// RoundTriggerClick 点击回查：点击量达到配置阈值
	RoundTriggerClick = int8(90)
	// RoundEnd 结束
	RoundEnd = int8(99)

	// AccessDefault 非会员可见
	AccessDefault = int16(0)
	// AccessMember 会员可见
	AccessMember = int16(10000)

	// CopyrightUnknow 未知版权类型
	CopyrightUnknow = int8(0)
	// CopyrightOriginal 原创
	CopyrightOriginal = int8(1)
	// CopyrightCopy 转载
	CopyrightCopy = int8(2)

	// AttrYes attribute yes
	AttrYes = int32(1)
	// AttrNo attribute no
	AttrNo = int32(0)

	// AttrBitNoRank 禁止排行
	AttrBitNoRank = uint(0)
	// AttrBitNoDynamic 动态禁止
	AttrBitNoDynamic = uint(1)
	// AttrBitNoWeb 禁止网页输出
	AttrBitNoWeb = uint(2)
	// AttrBitNoMobile 禁止客户端列表
	AttrBitNoMobile = uint(3)
	// AttrBitNoSearch 搜索禁止
	AttrBitNoSearch = uint(4)
	// AttrBitOverseaLock 海外禁止
	AttrBitOverseaLock = uint(5)
	// AttrBitNoRecommend 禁止推荐
	AttrBitNoRecommend = uint(6)
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
	// AttrBitIsFromArcAPI useless
	AttrBitIsFromArcAPI = uint(15)
	// AttrBitJumpURL 跳转
	AttrBitJumpURL = uint(16)
	// AttrBitIsMovie 是否影视
	AttrBitIsMovie = uint(17)
	// AttrBitBadgepay 付费
	AttrBitBadgepay = uint(18)
	//AttrNoPushBplus 禁止Bplus动态
	AttrNoPushBplus = uint(20)
	//AttrParentMode 家长模式
	AttrBitParentMode = uint(21)
	//AttrUGCPay   UGC付费
	AttrBitUGCPay = uint(22)
	//AttrBitSTAFF 联合投稿
	AttrBitSTAFF = uint(24)

	STATESTAFFON  = int8(1)
	STATESTAFFOFF = int8(2)
)

// Archive is archive model.
type Archive struct {
	Aid       int64      `json:"aid"`
	Mid       int64      `json:"mid"`
	TypeID    int16      `json:"tid"`
	Title     string     `json:"title"`
	Author    string     `json:"author"`
	Cover     string     `json:"cover"`
	Tag       string     `json:"tag"`
	Duration  int64      `json:"duration"`
	Copyright int8       `json:"copyright"`
	Desc      string     `json:"desc"`
	Round     int8       `json:"round"`
	Forward   int64      `json:"forward"`
	Attribute int32      `json:"attribute"`
	HumanRank int        `json:"humanrank"`
	Access    int16      `json:"access"`
	State     int8       `json:"state"`
	Reason    string     `json:"reject_reason"`
	PTime     xtime.Time `json:"ptime"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"mtime"`
}

// Attr archive attribute
type Attr int32

// Set set archive attribute
func (a *Attr) Set(v int32, bit uint) {
	*a = Attr(int32(*a)&(^(1 << bit)) | (v << bit))
}

// IsNormal check archive is open.
func (a *Archive) IsNormal() bool {
	return a.State >= StateOpen || a.State == StateForbidFixed
}

// NotAllowUp check archive is or not allow update state.
func (a *Archive) NotAllowUp() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidLater || a.State == StateForbidLock || a.State == StateForbidPolice
}

// IsForbid check archive state forbid by admin or delete.
func (a *Archive) IsForbid() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidRecicle || a.State == StateForbidPolice || a.State == StateForbidLock || a.State == StateForbidLater || a.State == StateForbidXcodeFail
}

// AttrVal get attribute value.
func (a *Archive) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// AttrSet set attribute value.
func (a *Archive) AttrSet(v int32, bit uint) {
	a.Attribute = a.Attribute&(^(1 << bit)) | (v << bit)
}

// WithAttr set attribute value with a attr value.
func (a *Archive) WithAttr(attr Attr) {
	a.Attribute = a.Attribute | int32(attr)
}

// NormalState check archive state is normal
func NormalState(state int8) bool {
	return state == StateOpen || state == StateOrange
}

// History archive history model
type History struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Desc  string `json:"desc"`
	State int8   `json:"state"`
}

// BlogCard 粉丝动态
type BlogCard struct {
	Type int64 `json:"type"`
	//Stype   int64  `json:"stype"`
	Rid     int64        `json:"rid"`
	OwnerID int64        `json:"owner_id"`
	Show    int64        `json:"show"`
	Comment string       `json:"comment"`
	Ts      int64        `json:"ts"`
	Dynamic string       `json:"dynamic"`
	Ext     string       `json:"extension"`
	Staffs  []*StaffItem `json:"staffs,omitempty"`
}

//StaffItem 联合投稿人信息  type=1
type StaffItem struct {
	Type int8  `json:"uid_type"`
	UID  int64 `json:"uid"`
}

//Staff . 正式staff
type Staff struct {
	ID           int64  `json:"id"`
	AID          int64  `json:"aid"`
	MID          int64  `json:"mid"`
	StaffMID     int64  `json:"staff_mid"`
	StaffTitle   string `json:"staff_title"`
	StaffTitleID int64  `json:"staff_title_id"`
	State        int8   `json:"state"`
}

//Ext 动态 ext 配置
type Ext struct {
	LBS  string `json:"lbs_cfg"`
	Vote string `json:"vote_cfg"`
}
