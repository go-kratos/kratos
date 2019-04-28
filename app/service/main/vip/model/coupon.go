package model

// coupon state.
const (
	CouponNotUsed = iota
	CouponInUse
	CouponUsed
	CouponExpire
)

// coupon remark
const (
	CouponUseRemark = "大会员券消费"
)

// MapProdLlimRenewal .
var MapProdLlimRenewal = map[int8]int8{0: 2, 1: 1}
