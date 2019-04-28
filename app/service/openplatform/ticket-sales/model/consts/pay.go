package consts

// 支付状态
const (
	PayStatusPaying   = "PAYING"   // 待支付
	PayStatusOverdue  = "OVERDUE"  // 已过期
	PayStatusClose    = "CLOSED"   // 支付关闭
	PayStatusFail     = "FAIL"     // 支付失败
	PayStatusSuccess  = "SUCCESS"  // 支付成功
	PayStatusFinished = "FINISHED" // 交易成功
)
