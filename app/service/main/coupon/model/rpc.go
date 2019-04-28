package model

// ArgSalaryCoupon salary coupon.
type ArgSalaryCoupon struct {
	Mid        int64 `form:"mid" validate:"required,min=1,gte=1"`
	CouponType int64 `form:"type" validate:"required,min=1,gte=1"`
	Origin     int64
	Count      int    `form:"count" validate:"required,min=1,gte=1"`
	BatchToken string `form:"batch_no" validate:"required"`
	AppID      int64  `form:"app_id"`
	UniqueNo   string
}

// ArgRPCPage def .
type ArgRPCPage struct {
	Mid   int64
	State int8
	Pn    int
	Ps    int
}

// CouponPageRPCResp def.
type CouponPageRPCResp struct {
	Count int64             `json:"count"`
	Res   []*CouponPageResp `json:"list"`
}

// ArgAllowanceCoupon def .
type ArgAllowanceCoupon struct {
	Mid            int64
	Pirce          float64
	Platform       int
	ProdLimMonth   int8
	ProdLimRenewal int8
}

// ArgUsablePirces def .
type ArgUsablePirces struct {
	Mid            int64
	Pirce          []float64
	Platform       int
	ProdLimMonth   int8
	ProdLimRenewal int8
}

// ArgJuageUsable def .
type ArgJuageUsable struct {
	Mid            int64
	Pirce          float64
	CouponToken    string
	Platform       int
	ProdLimMonth   int8
	ProdLimRenewal int8
}

// ArgAllowance def .
type ArgAllowance struct {
	Mid         int64
	CouponToken string
}

// ArgNotify  .
type ArgNotify struct {
	Mid     int64
	OrderNo string
	State   int8
}

// ArgAllowanceList  .
type ArgAllowanceList struct {
	Mid   int64
	State int8
}

// ArgAllowanceMid  .
type ArgAllowanceMid struct {
	Mid int64
}
