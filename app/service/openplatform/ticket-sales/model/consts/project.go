package consts

//项目相关状态
const (
	DefaultBuyNumLimit = 8 //默认单张订单购买限制

	BuyerInfoTel    = 1 //需要手机号
	BuyerInfoPerID  = 2 //需要身份证号
	BuyerInfoPerPic = 3 //需要身份证图片（未实现）

	DeliverTypeNone    = 1 //不配送
	DeliverTypeSelf    = 2 //自取
	DeliverTypeExpress = 3 //快递配送

	TicketTypePaper = 1  //纸质票
	TicketTypeElec  = 2  //电子票
	TicketTypeExt   = 3  //外部票（电子票）
	TicketTypeExch  = 12 //兑换票（未实现，实为纸质票+自取）

	PickSeatYes = 1 //选座项目
	PickSeatNo  = 0 //不选座项目
)
