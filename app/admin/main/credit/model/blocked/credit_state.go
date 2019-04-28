package blocked

// const credit state
const (
	// SearchDefaultStatus delault status
	SearchDefaultStatus = 0
	// JuryDay jury day
	JuryDay = 30
	// AccMaxPageSize max ps.
	AccMaxPageSize = 50
	// NotNeedSendMsg no send msg.
	NotNeedSendMsg = int8(0)
	// NeedSendMsg send msg.
	NeedSendMsg = int8(1)
	// UnBlockedForever no on forever
	UnBlockedForever = int8(0)
	// OnBlockedForever  on forever
	OnBlockedForever = int8(1)
	// SearchDefaultNum defalut num
	SearchDefaultNum = -100
	// SearchDefaultString defalut string
	SearchDefaultString = "-"
	// reasonType
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

	// AddJuryRemark add jury word
	AddJuryRemark           = "后台添加风纪委员"
	DefaultTime             = "1980-01-01 00:00:00"
	TableBlockedCase        = "blocked_case"
	TableBlockedInfo        = "blocked_info"
	TableBlockedJury        = "blocked_jury"
	TableBlockedOpinion     = "blocked_opinion"
	TableBlockedPublish     = "blocked_publish"
	TableBlockedKpiPoint    = "blocked_kpi_point"
	BusinessBlockedCase     = "block_case"
	BusinessBlockedInfo     = "block_info"
	BusinessBlockedJury     = "block_jury"
	BusinessBlockedOpinion  = "block_opinion"
	BusinessBlockedPublish  = "block_publish"
	BusinessBlockedKpiPoint = "block_kpi_point"
)

var _reasonType = map[int8]string{
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

// ReasonTypeDesc get reasonType desc
func ReasonTypeDesc(reasonType int8) string {
	return _reasonType[reasonType]
}
