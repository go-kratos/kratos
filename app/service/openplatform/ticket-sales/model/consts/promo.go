package consts

//拼团常量
const (
	PromoWaitShelf   int16 = 1 //待上架
	PromoUpShelf     int16 = 2 //已上架
	PromoDelShelf    int16 = 3 //废弃
	PromoFinishShelf int16 = 4 //已结束

	DoUpShelf  int16 = 1 //上架
	DoDelShelf int16 = 2 //废弃

	GroupDoing   int16 = 0 //拼团中
	GroupSuccess int16 = 1 //拼团成功
	GroupFailed  int16 = 2 //拼团失败

	PromoOrderUnpaid int16 = 1 //待支付
	PromoOrderPaid   int16 = 2 //已支付
	PromoOrderRefund int16 = 3 //已退款
	PromoOrderCancel int16 = 4 //已取消
)
