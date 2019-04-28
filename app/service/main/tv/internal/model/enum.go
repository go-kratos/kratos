package model

// enums.
const (
	PlatformAndriod  = 1
	PlatformWechatMp = 2
	PlatformSystem   = 3

	SuitTypeAll  = 0
	SuitTypeMvip = 10

	SubTypeOther    = 0
	SubTypeContract = 1

	PayOrderTypeNormal = 0
	PayOrderTypeSub    = 1

	PayOrderStatusPending = 0
	PayOrderStatusPaying  = 1
	PayOrderStatusSuccess = 2
	PayOrderStatusFail    = 3

	PaymentTypeAliPay = "alipay"
	PaymentTypeWechat = "wechat"

	UserChangeTypeRecharge       = 1
	UserChangeTypeSystem         = 2
	UserChangeTypeGift           = 3
	UserChangeTypeDup            = 4
	UserChangeTypeRenewVIp       = 5
	UserChangeTypeSignContract   = 6
	UserChangeTypeCancelContract = 7

	VipTypeVip       = 1
	VipTypeAnnualVip = 2

	VipPayTypeNormal = 0
	VipPayTypeSub    = 1

	PayChannelAli    = "alipay"
	PayChannelWechat = "wechat"

	VipStatusActive  = 1
	VipStatusExpired = 0

	YstPayWayQr = "1"

	YstPayTypeAliPay = "1"
	YstPayTypeWechat = "2"

	YstPayStatusPaied   = "1"
	YstPayStatusPending = "0"

	YstTradeStateSuccess = "SUCCESS"
	YstTradeStateRefund  = "REFUND"
	YstTradeStateNotPay  = "NOTPAY"
	YstTradeStateClosed  = "CLOSED"
	YstTradeStateAccept  = "ACCEPT"
	YstTradeStatePayFail = "PAY_FAIL"

	YstResultSuccess = "0"
	YstResultFail    = "998"
	YstResultSysErr  = "999"

	YST_CONTRACT_TYPE_SIGN   = "ADD"
	YST_CONTRSCT_TYPE_CANCEL = "DELETE"
)
