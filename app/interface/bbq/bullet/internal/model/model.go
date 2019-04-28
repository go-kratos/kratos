package model

// const 常量
const (
	SecondMaxNum         = 15
	CommonDurationSecond = 10
	BulletMaxLen         = 16
)

// 接口Action定义
const (
	ActionRecommend = iota
	ActionPlay
	ActionLike
	ActionCancelLike
	ActionFollow
	ActionCancelFollow
	ActionCommentAdd
	ActionCommentLike
	ActionCommentReport
	ActionFeedList
	ActionShare
	ActionDanmaku
	ActionPlayPause
	ActionPushRegister
	ActionPushSucced
	ActionPushCallback
	ActionBlack
	ActionCancelBlack
	ActionVideoSearch
	ActionUserSearch
)
