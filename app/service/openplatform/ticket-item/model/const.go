package model

// 票价

// 额外属性字段名
const (
	TkBuyNumLimitNormal    = "buy_limit_0"
	TkBuyNumLimitVip       = "buy_limit_1"
	TkBuyNumLimitAnnualVip = "buy_limit_2"

	// TimeNull 空时间：0000-00-00 00:00:00
	TimeNull = -62135596800

	// 版本类型
	//VerTypeProject  = 1
	VerTypeBulletin = 2
	//VerTypePjCancel = 3
	VerTypeBanner = 4

	//版本状态
	VerStatusOffShelfManual = -1
	VerStatusOffShelfForced = -2
	VerStatusNotReviewed    = 0
	VerStatusReadyForReview = 1
	VerStatusRejected       = 2
	VerStatusReadyForSale   = 3
	VerStatusOnShelf        = 4
	VerStatusNoSalesinfo    = 5

	//版本审核操作
	VerReviewReject = 0
	VerReviewPass   = 1
)
