package model

// Wallet struct.
type Wallet struct {
	Mid           int64   `json:"mid"`
	BcoinBalance  float32 `json:"bcoin_balance"`
	CouponBalance float32 `json:"coupon_balance"`
	DefaultBp     float32 `json:"default_bp"`
	IosBp         float32 `json:"ios_bp"`
}
