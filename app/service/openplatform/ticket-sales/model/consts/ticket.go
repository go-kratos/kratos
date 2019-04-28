package consts

//票状态常量
const (
	TkStatusUnchecked  = 0 //未检票
	TkStatusChecked    = 1 //已检票
	TkStatusExpired    = 2 //已过期
	TkStatusRefunded   = 3 //已退票
	TkStatusRefunding  = 4 //退票中
	TkStatusRefundFail = 5 //退票失败
	TkStatusSended     = 6 //已转赠
)

// 票类型: 0-订单出票, 1-系统赠票, 2-用户赠票, 3-票代分销'
const (
	TkTypeOrder      = 0 // 订单出票
	TkTypeSystemSend = 1 // 系统赠票
	TkTypeUserSend   = 2 // 用户赠票
	TkTypeDistrib    = 3 // 票代分销
)

// ticket_id 类型
const (
	TIDTypeSend = "send"
	TIDTypeRecv = "recv"
)
