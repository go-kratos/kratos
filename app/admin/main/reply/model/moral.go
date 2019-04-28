package model

import "fmt"

const (
	// ReportReasonOther 其他
	ReportReasonOther = int32(0)
	// ReportReasonAd 广告
	ReportReasonAd = int32(1)
	// ReportReasonPorn 色情
	ReportReasonPorn = int32(2)
	// ReportReasonMeaningless 刷屏
	ReportReasonMeaningless = int32(3)
	// ReportReasonProvoke 引站
	ReportReasonProvoke = int32(4)
	// ReportReasonSpoiler 剧透
	ReportReasonSpoiler = int32(5)
	// ReportReasonPolitic 政治
	ReportReasonPolitic = int32(6)
	// ReportReasonAttack 人身攻击
	ReportReasonAttack = int32(7)
	// ReportReasonUnrelated 视频不相关
	ReportReasonUnrelated = int32(8)
	// ReportReasonProhibited 违禁
	ReportReasonProhibited = int32(9)
	// ReportReasonVulgar 低俗
	ReportReasonVulgar = int32(10)
	// ReportReasonIllegalWebsite 非法网站
	ReportReasonIllegalWebsite = int32(11)
	// ReportReasonGamblingFraud 赌博诈骗
	ReportReasonGamblingFraud = int32(12)
	// ReportReasonRumor 传播不实信息
	ReportReasonRumor = int32(13)
	// ReportReasonAbetting 怂恿教唆信息
	ReportReasonAbetting = int32(14)
	// ReportReasonPrivacyInvasion 侵犯隐私
	ReportReasonPrivacyInvasion = int32(15)
	// ReportReasonUnlimitedSign 抢楼
	ReportReasonUnlimitedSign = int32(16)
	// ForbidReasonSpoiler 发布剧透信息
	ForbidReasonSpoiler = int32(10)
	// ForbidReasonAd 发布垃圾广告信息
	ForbidReasonAd = int32(6)
	// ForbidReasonUnlimitedSign 抢楼
	ForbidReasonUnlimitedSign = int32(2)
	// ForbidReasonMeaningless 刷屏
	ForbidReasonMeaningless = int32(1)
	// ForbidReasonProvoke 发布引战言论
	ForbidReasonProvoke = int32(9)
	// ForbidReasonVulgar 发布低俗信息
	ForbidReasonVulgar = int32(14)
	// ForbidReasonGamblingFraud 发布赌博诈骗信息
	ForbidReasonGamblingFraud = int32(4)
	// ForbidReasonPorn 发布色情信息
	ForbidReasonPorn = int32(13)
	// ForbidReasonRumor 发布传播不实信息
	ForbidReasonRumor = int32(18)
	// ForbidReasonIllegalWebsite 发布非法网站信息
	ForbidReasonIllegalWebsite = int32(17)
	// ForbidReasonAbetting 发布怂恿教唆信息
	ForbidReasonAbetting = int32(19)
	// ForbidReasonProhibited 发布违禁信息
	ForbidReasonProhibited = int32(5)
	// ForbidReasonPrivacyInvasion 涉及侵犯他人隐私
	ForbidReasonPrivacyInvasion = int32(8)
	// ForbidReasonAttack 发布人身攻击言论
	ForbidReasonAttack = int32(7)
	// ForbidReasonInaptitude 发布不适宜内容
	ForbidReasonInaptitude = int32(28)
)

var (
	// NotifyComRules 社区规则
	NotifyComRules = fmt.Sprintf(`评论区是公众场所，而非私人场所，具体规范烦请参阅#{《社区规则》}{"%s"}，良好的社区氛围需要大家一起维护！`, "http://www.bilibili.com/blackboard/blackroom.html")
	// NotifyComRulesReport 举报
	NotifyComRulesReport = "感谢您对bilibili社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~"
	// NotifyComUnrelated NotifyComUnrelated
	NotifyComUnrelated = "bilibili倡导发送与视频相关的评论，希望大家尊重作品，尊重UP主。良好的社区氛围需要大家一起维护！"
	// NotifyComProvoke NotifyComProvoke
	NotifyComProvoke = "bilibili倡导平等友善的交流。良好的社区氛围需要大家一起维护！"
	// NofityComProhibited NofityComProhibited
	NofityComProhibited = fmt.Sprintf(`请自觉遵守国家相关法律法规及#{《社区规则》}{"%s"}，bilibili良好的社区氛围需要大家一起维护！`, "http://www.bilibili.com/blackboard/blackroom.html")
	// ReportReason 举报理由类型
	ReportReason = map[int32]string{
		ReportReasonAd:              "内容涉及垃圾广告",
		ReportReasonPorn:            "内容涉及色情",
		ReportReasonMeaningless:     "刷屏",
		ReportReasonProvoke:         "内容涉及引战",
		ReportReasonSpoiler:         "内容涉及视频剧透",
		ReportReasonPolitic:         "内容涉及政治相关",
		ReportReasonAttack:          "内容涉及人身攻击",
		ReportReasonUnrelated:       "视频不相关",
		ReportReasonProhibited:      "内容涉及违禁相关",
		ReportReasonVulgar:          "内容涉及低俗信息",
		ReportReasonIllegalWebsite:  "内容涉及非法网站信息",
		ReportReasonGamblingFraud:   "内容涉及赌博诈骗信息",
		ReportReasonRumor:           "内容涉及传播不实信息",
		ReportReasonAbetting:        "内容不适宜",
		ReportReasonPrivacyInvasion: "内容涉及侵犯他人隐私",
		ReportReasonUnlimitedSign:   "抢楼",
	}
	// ForbidReason 封禁理由类型
	ForbidReason = map[int32]string{
		ForbidReasonSpoiler:         "发布剧透信息",
		ForbidReasonAd:              "发布垃圾广告信息",
		ForbidReasonUnlimitedSign:   "抢楼",
		ForbidReasonMeaningless:     "刷屏",
		ForbidReasonProvoke:         "发布引战言论",
		ForbidReasonVulgar:          "发布低俗信息",
		ForbidReasonGamblingFraud:   "发布赌博诈骗信息",
		ForbidReasonPorn:            "发布色情信息",
		ForbidReasonRumor:           "发布传播不实信息",
		ForbidReasonIllegalWebsite:  "发布非法网站信息",
		ForbidReasonAbetting:        "发布怂恿教唆信息",
		ForbidReasonProhibited:      "发布违禁信息",
		ForbidReasonPrivacyInvasion: "涉及侵犯他人隐私",
		ForbidReasonAttack:          "发布人身攻击言论",
		ForbidReasonInaptitude:      "发布不适宜内容",
	}
)
