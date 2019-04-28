package model

// consts for state
const (
	// Group and Challenge State field
	Pending       = int8(0)  // 未处理
	Effective     = int8(1)  // 有效
	Invalid       = int8(2)  // 无效
	RoleShift     = int8(3)  // 流转
	Deleted       = int8(9)  // 已删除
	PublicReferee = int8(10) // 移交众裁

	// dispatch_state offset bit
	AuditorStateOffset = 0
	CSStateOffset      = 4

	QueueState         = 15 // 队列中审核状态
	QueueBusinessState = 15 // 队列中客服状态

	QueueStateBefore         = 0 // 默认审核状态
	QueueBusinessStateBefore = 1 // 默认客服状态

	// 反馈状态
	FeedbackReplyNotRead = 6 // 已回复未读
)

// platform state
const (
	PlatformStateHandling = iota + 1
	PlatformStateDone
	PlatformStateClosed
)

// business round
const (
	AuditRoundMin = 1
	AuditRoundMax = 10
	FeedbackRound = 11
)
