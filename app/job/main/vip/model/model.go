package model

//Eunm vip enum value
const (
	//ChangeType
	ChangeTypePointExhchange  = 1 // 积分兑换
	ChangeTypeRechange        = 2 //充值开通
	ChangeTypeSystem          = 3 // 系统发放
	ChangeTypeActiveGive      = 4 //活动赠送
	ChangeTypeRepeatDeduction = 5 //重复领取扣除

	VipDaysMonth = 31
	VipDaysYear  = 366

	NotVip    = 0 //非大会员
	Vip       = 1 //月度大会员
	AnnualVip = 2 //年度会员

	VipStatusOverTime    = 0 //过期
	VipStatusNotOverTime = 1 //未过期
	VipStatusFrozen      = 2 //冻结
	VipStatusBan         = 3 //封禁

	VipAppUser  = 1 //大会员对接业务方user缓存
	VipAppPoint = 2 //大会员对接业务方积分缓存

	VipChangeFrozen   = -1 //冻结
	VipChangeUnFrozen = 0  //解冻
	VipChangeOpen     = 1  //开通
	VipChangeModify   = 2  //变更

	VipBusinessStatusOpen  = 0 //有效
	VipBusinessStatusClose = 1 //无效

	VipUserFirstDiscount = 1

	AnnualVipBcoinDay              = "annual_vip_bcoin_day"                //年费VIPB券发放每月第几天
	AnnualVipBcoinCouponMoney      = "annual_vip_bcoin_coupon_money"       //年费VIP返回B券金额
	AnnualVipBcoinCouponActivityID = "annual_vip_bcoin_coupon_activity_id" //年费VIP返B券活动ID

	HadSalaryState = 1 // 已发放

	NormalVipSalaryType = 1
	AnnualVipSalaryType = 2
	VipSupplyType       = 3
	TimingSalaryType    = 4

	SalaryVipOrigin = 1

	TimeFormatSec = "2006-01-02 15:04:05"

	DayOfHour = 24

	IsAutoRenew = 1

	IAPChannelID = 100

	MsgSystemNotify   = 4
	MsgCouponSalaryMc = "10_99_2"

	CouponSalaryTitle       = `观影劵到账通知`
	CouponSalaryMsg         = `大会员专享观影券已到账，#{点击查看>>}{"https://big.bilibili.com/mobile/userticket"}`
	CouponToAnnualSalaryMsg = `升级年度大会员赠送观影券%d张已到账，#{点击查看>>}{"https://big.bilibili.com/mobile/userticket"}`

	CouponCartoonSalaryTitle       = `漫画阅读劵到账通知`
	CouponCartoonSalaryMsg         = `大会员漫画阅读券已到账 #{点击查看>>}{"https://big.bilibili.com/mobile/userticket"}`
	CouponCartoonToAnnualSalaryMsg = `升级年度大会员赠送漫画阅读券%d张已到账，#{点击查看>>}{"https://big.bilibili.com/mobile/userticket"}`
)

// dicount type
const (
	DiscountNotUse = iota
	DiscountUsed
)

// coupon salary type
const (
	CouponSalaryTiming int8 = iota
	CouponSalaryAtonce
)

// coupon type
const (
	SalaryCouponType int8 = iota + 1
	SalaryCartoonCouponType
)

//pay order type
const (
	Normal = iota
	AutoRenew
	IAPAutoRenew
)

//pay order status
const (
	PAYING = iota + 1
	SUCCESS
	FAILED
	SIGN
	UNSIGN
)
