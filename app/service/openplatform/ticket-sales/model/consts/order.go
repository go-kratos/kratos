package consts

//订单常量
const (
	OrderTypeNormal  = 1 //普通订单
	OrderTypeGroup   = 2 //拼团订单
	OrderTypeDistrib = 4 //分销订单（已废弃）

	OrderStatusUnpaid   = 1 //未付款
	OrderStatusPaid     = 2 //已付款
	OrderStatusRefunded = 3 //已退款
	OrderStatusCancel   = 4 //已取消

	SubStatusUnpaid = 1 //未付款
	SubStatusClose  = 2 //已取消-系统关闭
	SubStatusPaid   = 3 //已付款-待出票
	//4-6为退款状态，已拆出为refundStatus，7-8为已废弃付款子状态
	SubStatusCompleted = 9  //已完成
	SubStatusUnshipped = 10 //已待发货
	SubStatusShipped   = 11 //已付款-待收货
	SubStatusCancel    = 12 //已取消-用户取消（待启用）

	RefundStatusNone        = 0 //无退款
	RefundStatusPtRefunding = 1 //部分退款中
	RefundStatusPtRefunded  = 2 //部分已退款
	RefundStatusRefunding   = 3 //退款中
	RefundStatusRefunded    = 4 //已退款

	RefundTxStatusCreated = 1 //退款流水创建
	RefundTxStatusSucc    = 2 //退款流水处理成功
)

var (
	//OrderTypes 订单类型
	OrderTypes = map[int16]string{
		OrderTypeNormal: "普通",
		OrderTypeGroup:  "拼团",
		//Deprecated
		OrderTypeDistrib: "分销",
	}
)
