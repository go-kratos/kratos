package model

// CouponCode coupon code.
type CouponCode struct {
	ID          int64
	BatchToken  string
	State       int32
	Code        string
	Mid         int64
	CouponType  int32
	CouponToken string
	Ver         int64
}

//Token get token .
type Token struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// ArgUseCouponCode arg use coupon code.
type ArgUseCouponCode struct {
	Token  string `form:"token" validate:"required"`
	Code   string `form:"code" validate:"required"`
	Verify string `form:"verify" validate:"required"`
	IP     string
	Mid    int64
}

// UseCouponCodeResp use coupon code resp.
type UseCouponCodeResp struct {
	CouponToken          string  `json:"coupon_token"`
	CouponAmount         float64 `json:"coupon_amount"`
	FullAmount           float64 `json:"full_amount"`
	PlatfromLimitExplain string  `json:"platfrom_limit_explain"`
	ProductLimitMonth    int32   `json:"product_limit_month"`
	ProductLimitRenewal  int32   `json:"product_limit_renewal"`
}

// coupon code state.
const (
	CodeStateNotUse = iota + 1
	CodeStateUsed
	CodeStateBlock
)
