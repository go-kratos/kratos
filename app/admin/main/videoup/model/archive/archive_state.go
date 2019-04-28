package archive

const (
	// StateOpen 开放浏览
	StateOpen = int8(0)
	// StateOrange 橙色通过
	StateOrange = int8(1)
	// StateForbidWait 待审
	StateForbidWait = int8(-1)
	// StateForbidRecycle 被打回
	StateForbidRecycle = int8(-2)
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
	// StateForbitUpLoad 创建未提交
	StateForbitUpLoad = int8(-20) // NOTE:spell body can judge to change state
	// StateForbidSubmit 创建已提交
	StateForbidSubmit = int8(-30)
	// StateForbidUserDelay 定时发布
	StateForbidUserDelay = int8(-40)
	// StateForbidUpDelete 用户删除
	StateForbidUpDelete = int8(-100)
	// AttrYes attribute yes
	AttrYes = int32(1)
	// AttrNo attribute no
	AttrNo = int32(0)
	// AttrBitNoRank 禁止排行
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
	// AttrBitIsFromArcAPI useless
	AttrBitIsFromArcAPI = uint(15) // TODO: delete
	// AttrBitJumpURL 跳转
	AttrBitJumpURL = uint(16)
	// AttrBitIsMovie 是否影视
	AttrBitIsMovie = uint(17)
	// AttrBitBadgepay 付费
	AttrBitBadgepay = uint(18)
	//AttrBitPushBlog 推送动态
	AttrBitPushBlog = uint(20)
	//AttrBitParentMode 家长模式
	AttrBitParentMode = uint(21)
	//AttrBitUGCPay UGC付费
	AttrBitUGCPay = uint(22)

	// CopyrightUnknow 未知版权类型
	CopyrightUnknow = int8(0)
	// CopyrightOriginal 原创
	CopyrightOriginal = int8(1)
	// CopyrightCopy 转载
	CopyrightCopy = int8(2)
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
	// DelayTypeForAdmin 管理员定时发布
	DelayTypeForAdmin = int8(1)
	// DelayTypeForUser 用户定时发布
	DelayTypeForUser = int8(2)
	// RoundBegin 一审阶段
	RoundBegin = int8(0)
	// RoundAuditSecond 二审：选定分区的多P稿件 及 PGC/活动的单P多P稿件
	RoundAuditSecond = int8(10)
	// RoundAuditThird 三审：选定分区/PGC/活动 的单P多P稿件
	RoundAuditThird = int8(20)
	// RoundReviewFlow 私单回查：私单ID大于0
	RoundReviewFlow = int8(21)
	//RoundReviewBadgepayFlow 付费审核
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

	// LogClientType 日志服务类型
	//for buiness

	//LogClientVideo 视频business id
	LogClientVideo = int(2)
	//LogClientArchive 稿件business id
	LogClientArchive = int(3)
	//LogClientUp up主business id
	LogClientUp = int(4)
	//LogClientPorder 私单business id
	LogClientPorder = int(5)
	//LogClientArchiveMusic 稿件bgm business id
	LogClientArchiveMusic = int(6)
	//LogClientPolicy 策略business id
	LogClientPolicy = int(7) //稿件策略组
	//LogClientConsumer 一审任务 business id
	LogClientConsumer = int(131)
	//LogClientTypePorderLog for business type

	//LogClientTypePorderLog 私单type id
	LogClientTypePorderLog = int(1)
	//LogClientTypeVideo 视频 type id
	LogClientTypeVideo = int(1)
	//LogClientTypeArchive 稿件 type id
	LogClientTypeArchive = int(1)
	//LogClientTypePorder 私单 id
	LogClientTypePorder = int(14)
	//LogClientTypePolicy 策略type id
	LogClientTypePolicy = int(1) //稿件策略组修改记录

	//LogClientArchiveMusicTypeMusic 稿件bgm type id
	LogClientArchiveMusicTypeMusic = int(1)
	//LogClientArchiveMusicTypeMaterial 稿件bgm素材 type id
	LogClientArchiveMusicTypeMaterial = int(2)
	//LogClientArchiveMusicTypeCategory 稿件bgm分类 type id
	LogClientArchiveMusicTypeCategory = int(3)
	//LogClientArchiveMusicTypeMaterialRelation 稿件bgm关联 type id
	LogClientArchiveMusicTypeMaterialRelation = int(4)
	//LogClientArchiveMusicTypeCategoryRelation 稿件bgm分区关联 type id
	LogClientArchiveMusicTypeCategoryRelation = int(5)

	//InnerAttrChannelReview 内部属性-频道回查--已删除
	InnerAttrChannelReview = uint(0)

	//LogClientTypeConsumer 一审任务type id
	LogClientTypeConsumer = int(1)
)

var (
	_attr = map[int32]int32{
		AttrNo:  AttrNo,
		AttrYes: AttrYes,
	}
	_access = map[int16]string{
		AccessDefault: "非会员可见",
		AccessMember:  "会员可见",
	}
	_copyright = map[int8]string{
		CopyrightUnknow:   "未知",
		CopyrightOriginal: "自制",
		CopyrightCopy:     "转载",
	}
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
		AttrBitJumpURL:    "跳转",
		AttrBitIsMovie:    "电影",
		AttrBitBadgepay:   "付费", //pgc付费
		AttrBitPushBlog:   "禁止粉丝动态",
		AttrBitParentMode: "家长模式",
		AttrBitUGCPay:     "UGC付费",
	}
	//  oversea forbidden typeid
	_overseaTypes = map[int16]int16{
		15:  15,  //'连载剧集'
		29:  29,  //'三次元音乐'
		32:  32,  //'完结动画'
		33:  33,  //'连载动画'
		34:  34,  //'完结剧集'
		37:  37,  //'纪录片'
		51:  51,  //'资讯'
		54:  54,  //'OP/ED/OST'
		71:  71,  //'综艺'
		86:  86,  //'特摄布袋戏'
		96:  96,  //'星海'
		130: 130, //'音乐选集'
		131: 131, //'Korea相关'
		137: 137, //'明星'
		145: 145, //'欧美电影'
		146: 146, //'日本电影'
		147: 147, //'国产电影'
		152: 152, //'官方延伸'
		153: 153, //'国产动画'
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

// UpFrom get upfrom desc
func UpFrom(ufID int8) string {
	return _upFromTypes[ufID]
}

// Attr attribute
type Attr int32

// InCopyrights in correct copyrights.
func InCopyrights(cp int8) (ok bool) {
	_, ok = _copyright[cp]
	return
}

// CopyrightsDesc return copyrights desc.
func CopyrightsDesc(cp int8) (desc string) {
	desc = _copyright[cp]
	return
}

// AccessDesc return acces desc.
func AccessDesc(acces int16) (desc string) {
	desc = _access[acces]
	return
}

// BitDesc return bit desc.
func BitDesc(bit uint) (desc string) {
	return _bits[bit]
}

// InAttr in correct attrs.
func InAttr(attr int32) (ok bool) {
	_, ok = _attr[attr]
	return
}

// InOverseaType check in oversea forbid type.
func InOverseaType(typeID int16) (ok bool) {
	_, ok = _overseaTypes[typeID]
	return
}

// NormalState check state.
func NormalState(state int8) bool {
	return state == StateOpen || state == StateOrange
}

// NotAllowDelay check need delete dtime of state.
func NotAllowDelay(state int8) bool {
	return state == StateForbidRecycle || state == StateForbidLock
}

// AttrSet set attribute.
func (arc *Archive) AttrSet(v int32, bit uint) {
	arc.Attribute = arc.Attribute&(^(1 << bit)) | (v << bit)
}

// AttrVal get attribute.
func (arc *Archive) AttrVal(bit uint) int32 {
	return (arc.Attribute >> bit) & int32(1)
}

// WithAttr set attribute value with a attr value.
func (arc *Archive) WithAttr(attr Attr) {
	arc.Attribute = arc.Attribute | int32(attr)
}

// NotAllowUp check archive is or not allow update state.
func (arc *Archive) NotAllowUp() bool {
	return arc.State == StateForbidUpDelete || arc.State == StateForbidLater || arc.State == StateForbidLock || arc.State == StateForbidPolice
}

//InnerAttrSet set inner_attr
func (addit *Addit) InnerAttrSet(v int64, bit uint) {
	addit.InnerAttr = addit.InnerAttr&(^(1 << bit)) | (v << bit)
}
