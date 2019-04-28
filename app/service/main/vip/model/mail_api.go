package model

// ArgMailCouponCodeCreate mail coupon code create args.
type ArgMailCouponCodeCreate struct {
	Mid      int64  `json:"mid"`
	Uname    string `json:"uname"`
	CouponID string `json:"couponId"`
}

// MailCouponCodeCreateResp mail code create resp.
type MailCouponCodeCreateResp struct {
	CouponCodeID       string `json:"couponCodeId"`
	RemainReceiveTimes int64  `json:"remainReceiveTimes"`
}
