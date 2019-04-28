package consts

// 库存操作日志 opType 常量
const (
	OpTypeOrder    int16 = 0  // 0 下单减库存
	OpTypePaid     int16 = 1  // 1 订单支付解锁库存
	OpTypeRefund   int16 = 2  // 退票补充库存(3为退 voucher 减库存，已弃用)
	OpTypeReserve  int16 = 4  // 预留扣库存
	OpTypeConfirm  int16 = 5  // 确认预留解锁库存
	OpTypeActive   int16 = 7  // 场次激活重置库存
	OpTypeBaseDecr int16 = 10 // 基础库存减少
	OpTypeBaseIncr int16 = 11 // 基础库存增加
	OpTypePrivDecr int16 = 12 // 活动库存减少
	OpTypePrivIncr int16 = 13 // 活动库存增加
	OpTypePrivInit int16 = 14 // 活动库存初始化
)
