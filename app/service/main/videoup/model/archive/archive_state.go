package archive

// Const State
const (
	// open state
	StateOpen   = int8(0)
	StateOrange = int8(1)
	// forbit state
	StateForbidWait     = int8(-1)
	StateForbidRecicle  = int8(-2)
	StateForbidPolice   = int8(-3)
	StateForbidLock     = int8(-4)
	StateForbidFackLock = int8(-5)
	StateForbidFixed    = int8(-6)
	StateForbidLater    = int8(-7)
	// StateForbidPatched   = int8(-8)
	StateForbidWaitXcode  = int8(-9)
	StateForbidAdminDelay = int8(-10)
	StateForbidFixing     = int8(-11)
	// StateForbidStorageFail = int8(-12)
	StateForbidOnlyComment = int8(-13)
	// StateForbidTmpRecicle  = int8(-14)
	StateForbidDispatch  = int8(-15)
	StateForbidXcodeFail = int8(-16)
	StateForbitUpLoad    = int8(-20) // NOTE:spell body can judge to change state
	StateForbidSubmit    = int8(-30)
	StateForbidUserDelay = int8(-40)
	StateForbidUpDelete  = int8(-100)
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
	AttrBitSTAFF        = uint(24) //联合投稿

	// copyright state
	CopyrightUnknow   = int8(0)
	CopyrightOriginal = int8(1)
	CopyrightCopy     = int8(2)
	// up_from
	UpFromWeb       = int8(0)
	UpFromPGC       = int8(1)
	UpFromWindows   = int8(2)
	UpFromAPP       = int8(3)
	UpFromMAC       = int8(4)
	UpFromSecretPGC = int8(5)
	UpFromCoopera   = int8(6)
	UpFromCreator   = int8(7) // 创作姬
	// delay
	DelayTypeForAdmin = int8(1)
	DelayTypeForUser  = int8(2)
	// flow type
	FlowNotLimit  = int8(1)
	FlowBudgeting = int8(2)
	FlowCapping   = int8(3)
	FlowForbid    = int8(4)
	// flow design type
	FlowDesignAppFeed = int8(0)
	FlowDesignUp      = int8(1)
	FlowDesignPrivate = int8(2)
	// oper uid
	AutoOperUID = int64(399)
	CMOperUID   = int64(518)
	// archive list type for up
	UpArcAllIn    = int8(0)
	UpArcOpenIn   = int8(1)
	UpArcUnOpenIn = int8(2)

	VideoFilenameTimeout = int64(48 * 60 * 60)
)

var (
	_attr = map[int32]int32{
		AttrNo:  AttrNo,
		AttrYes: AttrYes,
	}
	_copyright = map[int8]int8{
		CopyrightUnknow:   CopyrightUnknow,
		CopyrightOriginal: CopyrightOriginal,
		CopyrightCopy:     CopyrightCopy,
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
		AttrBitIsPGC:       "PGC稿件",
		AttrBitAllowBp:     "允许承包",
		AttrBitIsBangumi:   "番剧",
		// AttrBitAllowDownload: AttrBitAllowDownload,
		// AttrBitHideClick:    AttrBitHideClick,
		AttrBitAllowTag: "允许操作TAG",
		// AttrBitIsFromArcApi: AttrBitIsFromArcApi,
		AttrBitJumpURL:  "跳转",
		AttrBitIsMovie:  "电影",
		AttrBitBadgepay: "付费",
	}

	//  oversea forbidden typeid
	_overseaTypes = map[int16]int16{
		32: 32, //'完结动画'
		33: 33, //'连载动画'
	}
)

// InCopyrights in correct copyrights.
func InCopyrights(cp int8) (ok bool) {
	_, ok = _copyright[cp]
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
