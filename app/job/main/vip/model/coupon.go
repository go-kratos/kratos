package model

// VideoCouponSalaryLog videl coupon salary log.
type VideoCouponSalaryLog struct {
	ID          int64 `json:"id"`
	Mid         int64 `json:"mid"`
	CouponCount int64 `json:"coupon_count"`
	CouponType  int8  `json:"coupon_type"`
	State       int8  `json:"state"`
	Type        int8  `json:"type"`
	ExpireTime  int64 `json:"expire_time"`
	StartTime   int64 `json:"start_time"`
	Ver         int64 `json:"ver"`
}

// OldSalaryLog def.
type OldSalaryLog struct {
	ID          int64 `json:"id"`
	Mid         int64 `json:"mid"`
	CouponCount int64 `json:"coupon_count"`
	State       int8  `json:"state"`
	Type        int8  `json:"type"`
}
