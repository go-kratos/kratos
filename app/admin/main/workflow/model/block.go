package model

import credit "go-common/app/job/main/credit/model"

// BlockReason 封禁理由
var BlockReason = map[int8]string{
	credit.ReasonOtherType:                   "其他",
	credit.ReasonBrushScreen:                 "刷屏",
	credit.ReasonGrabFloor:                   "抢楼",
	credit.ReasonGamblingFraud:               "发布赌博诈骗信息",
	credit.ReasonProhibited:                  "发布违禁相关信息",
	credit.ReasonGarbageAds:                  "发布垃圾广告信息",
	credit.ReasonPersonalAttacks:             "发布人身攻击言论",
	credit.ReasonViolatePrivacy:              "发布侵犯他人隐私信息",
	credit.ReasonLeadBattle:                  "发布引战言论",
	credit.ReasonSpoiler:                     "发布剧透信息",
	credit.ReasonAddUnrelatedTags:            "恶意添加无关标签",
	credit.ReasonDelOtherTags:                "恶意删除他人标签",
	credit.ReasonPornographic:                "发布色情信息",
	credit.ReasonVulgar:                      "发布低俗信息",
	credit.ReasonBloodyViolence:              "发布暴力血腥信息",
	credit.ReasonAnimusVideoUp:               "涉及恶意投稿行为",
	credit.ReasonIllegalWebsite:              "发布非法网站信息",
	credit.ReasonSpreadErrinfo:               "发布传播不实信息",
	credit.ReasonAbettingEncouragement:       "发布怂恿教唆信息",
	credit.ReasonAnimusBrushScreen:           "恶意刷屏",
	credit.ReasonAccountViolation:            "账号违规",
	credit.ReasonMaliciousPlagiarism:         "恶意抄袭",
	credit.ReasonPosingAsHomemade:            "冒充自制原创",
	credit.ReasonPostTeenBadContent:          "发布青少年不良内容",
	credit.ReasonDestroyCyberSecurity:        "破坏网络安全",
	credit.ReasonPostingMisleadingInfo:       "发布虚假误导信息",
	credit.ReasonCounterfeitOfficialAuth:     "仿冒官方认证账号",
	credit.ReasonPublishInappropriateContent: "发布不适宜内容",
	credit.ReasonViolationOperatingRules:     "违反运营规则",
	credit.ReasonIllegalCreateTopic:          "恶意创建话题",
	credit.ReasonIllegalDrawLottery:          "发布违规抽奖",
	credit.ReasonIllegalFakeMan:              "恶意冒充他人",
}
