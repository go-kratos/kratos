package model

// 大会员类型
const (
	NotVip    = 0 //非大会员
	Vip       = 1 //月度大会员
	AnnualVip = 2 //年度会员
)

// 大会员状态
const (
	VipStatusOverTime    = 0 //过期
	VipStatusNotOverTime = 1 //未过期
	VipStatusFrozen      = 2 //冻结
	VipStatusBan         = 3 //封禁
)

// vip pay type.
const (
	NormalPay int32 = iota
	AutoRenewPay
)

// pay channel id
const (
	IapPayChannelID = 100
)
