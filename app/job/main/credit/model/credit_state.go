package model

import (
	"fmt"
)

// blocked const
const (
	// item type
	BUSSINESS = "credit-job"
	// deal bussinss time type
	DealTimeTypeNone = int8(0)
	DealTimeTypeDay  = int8(1)
	DealTimeTypeYear = int8(2)
	// default time
	DefaultTime   = "1979-12-31 16:00:00"
	TimeFormatSec = "2006-01-02 15:04:05"
	// Case Status
	CaseStatusGranting  = 1 // 发放中
	CaseStatusGrantStop = 2 // 停止发放
	CaseStatusDealing   = 3 // 结案中
	CaseStatusDealed    = 4 // 已裁决
	CaseStatusRestart   = 5 // 待重启
	CaseStatusUndealed  = 6 // 未裁决
	CaseStatusFreeze    = 7 // 冻结中
	CaseStatusQueueing  = 8 // 队列中

	// Judge Status.
	JudgeTypeUndeal  = 0 // 未裁决
	JudgeTypeViolate = 1 // 违规
	JudgeTypeLegal   = 2 // 未违规

	// Vote Status.
	VoteTypeUndo    = 0 // 未投票
	VoteTypeViolate = 1 // 违规-封禁
	VoteTypeDelete  = 4 // 违规-删除
	VoteTypeLegal   = 2 // 不违规
	VoteTypeGiveUp  = 3 // 放弃投票

	// punish type.
	PunishTypeMoral   = int8(1)
	PunishTypeBlock   = int8(2)
	PunishTypeForever = int8(3)

	// blocked_info blocked_forever
	NotInBlockedForever = int8(0)
	InBlockedForever    = int8(1)

	// Block Time.
	Punish3Days   = 1
	Punish7Days   = 2
	PunishForever = 3
	PunishCustom  = 4
	Punish15Days  = 5

	PunishBlock = 0
	PunishJury  = 1

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

	// Jury Invalid
	JuryBlocked = 1
	JuryExpire  = 2
	JuryAdmin   = 3

	// Case Load Switch
	StateCaseLoadClose = int8(0)
	StateCaseLoadOpen  = int8(1)

	// Blocked Opinio State
	OpinionStateOpen        = int8(0)
	OpinionStateClose       = int8(1)
	OpinionStateCloseAndMsg = int8(2)

	// blocked_jury.case_type
	JudeCaseTypePrivate = int8(0) // 小众众裁
	JudeCaseTypePublic  = int8(1) // 大众众裁

	// Reply regist type.
	ReplyBlocked = int8(6)
	ReplyPublish = int8(7)
	ReplyCase    = int8(15)

	// blocked_publish.publish_status
	PublishClose = int8(0)
	PublishOpen  = int8(1)

	// status
	StatusClose = int8(1)
	StatusOpen  = int8(0)

	// blocked_case.punish_result
	BlockNone    = int8(0)
	Block3Days   = int8(1)
	Block7Days   = int8(2)
	BlockForever = int8(3)
	BlockCustom  = int8(4)
	Block15Days  = int8(5)
	BlockOnlyDel = int8(6)

	// block time
	BlockTimeForever = 0  // 永久封禁
	BlockTimeThree   = 3  // 3天封禁
	BlockTimeSeven   = 7  // 7天封禁
	BlockTimeFifteen = 15 // 15天封禁

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

	// moral originType
	MoralOriginDM    = int8(1)
	MoralOriginReply = int8(2)
	MoralOriginTag   = int8(3)

	// block status
	BlockStatusNone    = int8(0)
	BlockStatusForever = int8(1)
	BlockStatusOn      = int8(2)

	// defealt deduct moral val
	DefealtMoralVal = -10

	// dm notify status
	DMNotifyNotDel = 0
	DMNotifyDel    = 1

	// msg content
	_dealMsgTitle   = "%s违规处理通知"
	_dealMsgContent = `您好，根据用户举报与风纪委众裁，您在#{"%s"}{%s}下的%s 『%s』，已被移除。请自觉遵守国家相关法律法规及《社区规则》，bilibili良好的社区氛围需要大家一起维护！
其中，《社区规则》为可点击超链接，地址：https://www.bilibili.com/blackboard/blackroom.html`

	// moral remark
	MoralRemark = "违规惩罚"
)

var (
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

	_reasonToFreeze = map[int8]bool{
		ReasonGamblingFraud:   true,
		ReasonViolatePrivacy:  true,
		ReasonProhibited:      true,
		ReasonPornographic:    true,
		ReasonVulgar:          true,
		ReasonSpoiler:         false,
		ReasonGrabFloor:       false,
		ReasonGarbageAds:      false,
		ReasonLeadBattle:      false,
		ReasonBrushScreen:     false,
		ReasonPersonalAttacks: false,
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

	// _orginMoralType 对应节操来源类型
	_orginMoralType = map[int8]int8{
		OriginReply: MoralOriginReply,
		OriginDM:    MoralOriginDM,
		OriginTag:   MoralOriginTag,
	}

	_blockDay = map[int8]string{
		BlockTimeForever: "永久封禁",
		BlockTimeThree:   "封禁3天",
		BlockTimeSeven:   "封禁7天",
		BlockTimeFifteen: "封禁15天",
	}
)

// ReasonToFreeze get reason yes or no to freeze.
func ReasonToFreeze(reasonType int8) bool {
	return _reasonToFreeze[reasonType]
}

// OrginMoralType get moral bussiness Type by blocked orgin.
func OrginMoralType(blockOrginType int8) int8 {
	return _orginMoralType[blockOrginType]
}

// PunishTypeDesc get punishType desc
func PunishTypeDesc(punishType int8) string {
	return _punishType[punishType]
}

// ReasonTypeDesc get reasonType desc
func ReasonTypeDesc(reasonType int8) string {
	return _reasonType[reasonType]
}

// OriginTypeDesc get originType desc
func OriginTypeDesc(originType int8) string {
	return _originType[originType]
}

// BlockedDayDesc is blocked day desc
func BlockedDayDesc(day int8) string {
	return _blockDay[day]
}

// BlockedReasonTypeByReply get blocked reason type.
func BlockedReasonTypeByReply(replyReasonType int8) int8 {
	return _replyReasonType[replyReasonType]
}

// BlockedReasonTypeByTag get blocked reason type.
func BlockedReasonTypeByTag(tagReasonType int8) int8 {
	return _tagReasonType[tagReasonType]
}

// OriginMsgContent get msg content by oTitle, oURL , oContent and oType
func OriginMsgContent(oTitle, oURL, oContent string, oType int8) (msgTitle, msgCon string) {
	msgTitle = fmt.Sprintf(_dealMsgTitle, _originType[oType])
	msgCon = fmt.Sprintf(_dealMsgContent, oTitle, oURL, _originType[oType], oContent)
	return
}
