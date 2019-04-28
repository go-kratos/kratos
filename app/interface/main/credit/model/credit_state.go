package model

// const credit state
const (
	// blocked_opinion
	BlockedOpinionAttrOn  = int8(1)
	BlockedOPinionAttrOff = int8(0)

	// blocked_info.block_type
	PunishBlock = int8(0) // 系统封禁
	PunishJury  = int8(1) // 风纪仲裁

	CreditStatusBlocked = -2

	// blocked_info.punish_type
	PunishTypeMoral   = int8(1)
	PunishTypeBlock   = int8(2)
	PunishTypeForever = int8(3)

	// Publish type
	PublishTypedef      = int8(0)
	PublishTypePunish   = int8(1)
	PublishTypeBan      = int8(2)
	PublishTypeOptimize = int8(3)

	// publish status
	PublishStatusClose = int8(0) // 案件关闭状态
	PublishStatusOpen  = int8(1) // 案件公开状态

	// origin_type.
	OriginReply    = int8(1)  // 评论
	OriginDM       = int8(2)  // 弹幕
	OriginMsg      = int8(3)  // 私信
	OriginTag      = int8(4)  // 标签
	OriginMember   = int8(5)  // 个人资料
	OriginArchive  = int8(6)  // 投稿
	OriginMusic    = int8(7)  // 音频
	OriginArticle  = int8(8)  // 专栏
	OriginSpaceTop = int8(9)  // 空间头图
	OriginDsynamic = int8(10) // 动态
	OriginPhoto    = int8(11) // 相册
	OriginMinVideo = int8(12) // 小视频

	// 	blocked_publish init lenth
	PublishInitLen = 4

	// blocked_jury.status
	JuryStatusEffect   = int8(1)
	JuryStatusNoEffect = int8(2)

	// blocked_jury.case_type
	JudeCaseTypePrivate = int8(0) // 小众众裁
	JudeCaseTypePublic  = int8(1) // 大众众裁

	// case obtain day by mid
	CaseObtainToday   = true
	CaseObtainNoToday = false

	// message
	ApplyJuryTitle   = "获得风纪委员资格"
	ApplyJuryContext = `恭喜您获得%d天风纪委员资格！风纪委员应遵守以下原则：
			"1. 在了解举报案件背景后，公正客观投票。对不了解或难以判断的案件，可以选择弃权。
			"2. 以身作则，不在举报案件相关视频、评论下讨论或发布不相关内容。相关违规举报被落实处罚后，将会失去风纪委员资格。`
	AppealTitle   = "申诉处理通知"
	MaxAddCaseNum = 100 //批量增加案件最大数量

	// list multi juryer info
	JuryMultiJuryerInfoMax = 50

	// jury expired
	JuryExpiredDays = 30

	// one day
	OneDaySecond = 86400

	// black or white
	JuryBlack = int8(1)
	JuryWhite = int8(2)

	// blocked_info blocked_forever
	NotInBlockedForever = int8(0)
	InBlockedForever    = int8(1)

	// blocked_info blocked_forever bool
	BlockedStateForever   = true
	BlockedStateNoForever = false

	// guard
	GuardMedalPointA = int64(5000)
	GuardMedalPointB = int64(1000)
	GuardMedalPointC = int64(200)
	GuardMedalNone   = int64(0)
	GuardMedalA      = int64(69)
	GuardMedalB      = int64(68)
	GuardMedalC      = int64(67)

	// reasonType
	ReasonOtherType                   = int8(0)
	ReasonBrushScreen                 = int8(1)
	ReasonGrabFloor                   = int8(2)
	ReasonGamblingFraud               = int8(4)
	ReasonProhibited                  = int8(5)
	ReasonGarbageAds                  = int8(6)
	ReasonPersonalAttacks             = int8(7)
	ReasonViolatePrivacy              = int8(8)
	ReasonLeadBattle                  = int8(9)
	ReasonSpoiler                     = int8(10)
	ReasonAddUnrelatedTags            = int8(11)
	ReasonDelOtherTags                = int8(12)
	ReasonPornographic                = int8(13)
	ReasonVulgar                      = int8(14)
	ReasonBloodyViolence              = int8(15)
	ReasonAnimusVideoUp               = int8(16)
	ReasonIllegalWebsite              = int8(17)
	ReasonSpreadErrinfo               = int8(18)
	ReasonAbettingEncouragement       = int8(19)
	ReasonAnimusBrushScreen           = int8(20)
	ReasonAccountViolation            = int8(21)
	ReasonMaliciousPlagiarism         = int8(22)
	ReasonPosingAsHomemade            = int8(23)
	ReasonPostTeenBadContent          = int8(24)
	ReasonDestroyCyberSecurity        = int8(25)
	ReasonPostingMisleadingInfo       = int8(26)
	ReasonCounterfeitOfficialAuth     = int8(27)
	ReasonPublishInappropriateContent = int8(28)
	ReasonViolationOperatingRules     = int8(29)
	ReasonIllegalCreateTopic          = int8(30)
	ReasonIllegalDrawLottery          = int8(31)
	ReasonIllegalFakeMan              = int8(32)
	// reply reasonType
	ReplyReasonOtherType             = int8(0)
	ReplyReasonGarbageAds            = int8(1)
	ReplyReasonPornographic          = int8(2)
	ReplyReasonAnimusBrushScreen     = int8(3)
	ReplyReasonLeadBattle            = int8(4)
	ReplyReasonSpoiler               = int8(5)
	ReplyReasonPolitical             = int8(6)
	ReplyReasonPersonalAttacks       = int8(7)
	ReplyReasonIrrelevantVideo       = int8(8)
	ReplyReasonProhibited            = int8(9)
	ReplyReasonVulgar                = int8(10)
	ReplyReasonIllegalWebsite        = int8(11)
	ReplyReasonGamblingFraud         = int8(12)
	ReplyReasonSpreadErrinfo         = int8(13)
	ReplyReasonAbettingEncouragement = int8(14)
	ReplyReasonViolatePrivacy        = int8(15)
	ReplyReasonGrabFloor             = int8(16)
	ReplyReasonPostTeenBadContent    = int8(17)
	// tag reasonType
	TagReasonAddUnrelatedTags = int8(1)
	TagReasonProhibited       = int8(2)
	TagReasonPersonalAttacks  = int8(3)
	TagReasonSpoiler          = int8(4)
	TagReasonDelOtherTags     = int8(5)
	// case status.
	CaseStatusGranting  = 1 // 发放中
	CaseStatusGrantStop = 2 // 停止发放
	CaseStatusDealing   = 3 // 结案中
	CaseStatusDealed    = 4 // 已裁决
	CaseStatusRestart   = 5 // 待重启
	CaseStatusUndealed  = 6 // 未裁决
	CaseStatusFreeze    = 7 // 冻结中
	CaseStatusQueueing  = 8 // 队列中

	// blocked_case.punish_result
	BlockNone    = int8(0)
	Block3Days   = int8(1)
	Block7Days   = int8(2)
	BlockForever = int8(3)
	BlockCustom  = int8(4)
	Block15Days  = int8(5)
	BlockOnlyDel = int8(6)

	// judge status.
	JudgeTypeUndeal  = 0 // 未裁决
	JudgeTypeViolate = 1 // 违规
	JudgeTypeLegal   = 2 // 未违规

	// vote type
	VoteBanned  = 1 // 违规封禁
	VoteRule    = 2
	VoteAbstain = 3
	VoteDel     = 4 // 违规删除

	// opinion type
	OpinonBreak = 1 // 违规观点
	OpinionRule = 2 // 不违规观点

	// labour ans
	LabourNoAnswer = int8(0)
	LabourOkAnswer = int8(1)

	// opinion state
	OpinionStateOK   = int8(0)
	OpinionStateNoOK = int8(1)

	// kpi rate
	KPILevelS = int8(1)
	KPILevelA = int8(2)
	KPILevelB = int8(3)
	KPILevelC = int8(4)
	KPILevelD = int8(5)

	// block status
	BlockStatusNone    = int8(0)
	BlockStatusForever = int8(1)
	BlockStatusOn      = int8(2)
)

// var credit state
var (
	_punishResult = map[int8]string{
		BlockNone:    "",
		Block3Days:   "封禁3天",
		Block7Days:   "封禁7天",
		BlockForever: "永久封禁",
		BlockCustom:  "封禁%d天",
		Block15Days:  "封禁15天",
		BlockOnlyDel: "扣节操",
	}
	_punishType = map[int8]string{
		PunishTypeMoral:   "节操",
		PunishTypeBlock:   "封禁",
		PunishTypeForever: "永久封禁",
	}
	_reasonType = map[int8]string{
		ReasonOtherType:                   "其他",
		ReasonBrushScreen:                 "刷屏",
		ReasonGrabFloor:                   "抢楼",
		ReasonGamblingFraud:               "发布赌博诈骗信息",
		ReasonProhibited:                  "发布违禁相关信息",
		ReasonGarbageAds:                  "发布垃圾广告信息",
		ReasonPersonalAttacks:             "发布人身攻击言论",
		ReasonViolatePrivacy:              "发布侵犯他人隐私信息",
		ReasonLeadBattle:                  "发布引战言论",
		ReasonSpoiler:                     "发布剧透信息",
		ReasonAddUnrelatedTags:            "恶意添加无关标签",
		ReasonDelOtherTags:                "恶意删除他人标签",
		ReasonPornographic:                "发布色情信息",
		ReasonVulgar:                      "发布低俗信息",
		ReasonBloodyViolence:              "发布暴力血腥信息",
		ReasonAnimusVideoUp:               "涉及恶意投稿行为",
		ReasonIllegalWebsite:              "发布非法网站信息",
		ReasonSpreadErrinfo:               "发布传播不实信息",
		ReasonAbettingEncouragement:       "发布怂恿教唆信息",
		ReasonAnimusBrushScreen:           "恶意刷屏",
		ReasonAccountViolation:            "账号违规",
		ReasonMaliciousPlagiarism:         "恶意抄袭",
		ReasonPosingAsHomemade:            "冒充自制原创",
		ReasonPostTeenBadContent:          "发布青少年不良内容",
		ReasonDestroyCyberSecurity:        "破坏网络安全",
		ReasonPostingMisleadingInfo:       "发布虚假误导信息",
		ReasonCounterfeitOfficialAuth:     "仿冒官方认证账号",
		ReasonPublishInappropriateContent: "发布不适宜内容",
		ReasonViolationOperatingRules:     "违反运营规则",
		ReasonIllegalCreateTopic:          "恶意创建话题",
		ReasonIllegalDrawLottery:          "发布违规抽奖",
		ReasonIllegalFakeMan:              "恶意冒充他人",
	}
	_originType = map[int8]string{
		OriginReply:    "评论",
		OriginDM:       "弹幕",
		OriginMsg:      "私信",
		OriginTag:      "标签",
		OriginMember:   "个人资料",
		OriginArchive:  "投稿",
		OriginMusic:    "音频",
		OriginArticle:  "专栏",
		OriginSpaceTop: "空间头图",
		OriginDsynamic: "动态",
		OriginPhoto:    "相册",
		OriginMinVideo: "小视频",
	}
	_replyReasonType = map[int8]int8{
		ReplyReasonOtherType:             ReasonOtherType,
		ReplyReasonGarbageAds:            ReasonGarbageAds,
		ReplyReasonPornographic:          ReasonPornographic,
		ReplyReasonAnimusBrushScreen:     ReasonAnimusBrushScreen,
		ReplyReasonLeadBattle:            ReasonLeadBattle,
		ReplyReasonSpoiler:               ReasonSpoiler,
		ReplyReasonPolitical:             ReasonOtherType,
		ReplyReasonPersonalAttacks:       ReasonPersonalAttacks,
		ReplyReasonIrrelevantVideo:       ReasonOtherType,
		ReplyReasonProhibited:            ReasonProhibited,
		ReplyReasonVulgar:                ReasonVulgar,
		ReplyReasonIllegalWebsite:        ReasonIllegalWebsite,
		ReplyReasonGamblingFraud:         ReasonGamblingFraud,
		ReplyReasonSpreadErrinfo:         ReasonSpreadErrinfo,
		ReplyReasonAbettingEncouragement: ReasonAbettingEncouragement,
		ReplyReasonViolatePrivacy:        ReasonViolatePrivacy,
		ReplyReasonGrabFloor:             ReasonGrabFloor,
		ReplyReasonPostTeenBadContent:    ReasonPostTeenBadContent,
	}
	_tagReasonType = map[int8]int8{
		TagReasonAddUnrelatedTags: ReasonAddUnrelatedTags,
		TagReasonProhibited:       ReasonProhibited,
		TagReasonPersonalAttacks:  ReasonPersonalAttacks,
		TagReasonSpoiler:          ReasonSpoiler,
		TagReasonDelOtherTags:     ReasonDelOtherTags,
	}
)

// PunishResultDesc get PunishResult desc
func PunishResultDesc(punishResult int8) (desc string) {
	desc = _punishResult[punishResult]
	return
}

// PunishTypeDesc get punishType desc
func PunishTypeDesc(punishType int8) (desc string) {
	desc = _punishType[punishType]
	return
}

// ReasonTypeDesc get reasonType desc
func ReasonTypeDesc(reasonType int8) (desc string) {
	desc = _reasonType[reasonType]
	return
}

// OriginTypeDesc get originType desc
func OriginTypeDesc(originType int8) (desc string) {
	desc = _originType[originType]
	return
}

// BlockedReasonTypeByReply get blocked reason type.
func BlockedReasonTypeByReply(replyReasonType int8) (reasonType int8) {
	reasonType = _replyReasonType[replyReasonType]
	return
}

// BlockedReasonTypeByTag get blocked reason type.
func BlockedReasonTypeByTag(tagReasonType int8) (reasonType int8) {
	reasonType = _tagReasonType[tagReasonType]
	return
}
