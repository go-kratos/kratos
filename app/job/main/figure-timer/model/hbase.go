package model

import (
	"time"
)

type FigureRecord struct {
	Mid            int64
	XPosLawful     int64
	XNegLawful     int64
	XPosWide       int64
	XNegWide       int64
	XPosFriendly   int64
	XNegFriendly   int64
	XPosCreativity int64
	XNegCreativity int64
	XPosBounty     int64
	XNegBounty     int64
	Version        time.Time
}

type UserInfo struct {
	Mid                 int64
	Exp                 uint64 // 当周期最终经验值
	SpyScore            uint64 // spy得分
	ArchiveViews        uint64 // 观看视频累计天使
	VIPStatus           uint64 // VIP状态
	DisciplineCommittee uint16 // 风纪委得分
}

type ActionCounter struct {
	Mid                   int64
	CoinCount             uint64 // 投币行为数
	ReplyCount            int64  // 评论次数
	DanmakuCount          int64  // 弹幕计次
	CoinLowRisk           uint64 // 投币疑似异常
	CoinHighRisk          uint64 // 投币高危异常
	ReplyLowRisk          uint64 // 评论疑似异常
	ReplyHighRisk         uint64 // 评论高危异常
	ReplyLiked            int64  // 评论被赞数
	ReplyUnliked          int64  // 评论被踩数
	ReportReplyPassed     int64  // 举报评论通过
	ReportDanmakuPassed   int64  // 举报弹幕通过
	PublishReplyDeleted   int64  // 评论被删除
	PublishDanmakuDeleted int64  // 弹幕被删除
	PayMoney              int64  // 支付消费金额（单位分）
	PayLiveMoney          int64  // 支付消费金额（单位分）
	Version               time.Time
}
