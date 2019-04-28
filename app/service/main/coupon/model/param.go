package model

// ArgMid .
type ArgMid struct {
	Mid  int64 `form:"mid" validate:"required,min=1,gte=1"`
	Type int8  `form:"type" validate:"required,min=1,gte=1"`
}

// ArgUseCoupon .
type ArgUseCoupon struct {
	Mid     int64  `form:"mid" validate:"required,min=1,gte=1"`
	Type    int8   `form:"type" validate:"required,min=1,gte=1"`
	Remark  string `form:"remark" validate:"required"`
	OrderNO string `form:"order_id" validate:"required"`
	Oid     int64  `form:"oid" validate:"required,min=1,gte=1"`
	Ver     int64  `form:"ver" validate:"required,min=1,gte=1"`
}

// ArgUseCartoonCoupon def .
type ArgUseCartoonCoupon struct {
	Mid     int64  `form:"mid" validate:"required,min=1,gte=1"`
	Type    int8   `form:"type" validate:"required,min=1,gte=1"`
	Tips    string `form:"tips" validate:"required"`
	Remark  string `form:"remark" validate:"required"`
	OrderNO string `form:"order_id" validate:"required"`
	Count   int64  `form:"count" validate:"required,min=1,gte=1"`
	Ver     int64  `form:"ver" validate:"required,min=1,gte=1"`
}

// ArgCoupon .
type ArgCoupon struct {
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	CouponToken string `form:"coupon_token" validate:"required"`
}

// ChangeCoupon .
type ChangeCoupon struct {
	Mid         int64  `form:"mid"`
	CouponToken string `form:"coupon_token"`
	Ver         int64  `form:"ver"`
	UseVer      int64  `form:"use_ver"`
}

// ArgAdd .
type ArgAdd struct {
	Mid        int64 `form:"mid" validate:"required,min=1,gte=1"`
	StartTime  int64 `form:"start_time"`
	ExpireTime int64 `form:"expire_time" validate:"required,min=1,gte=1"`
	Type       int64 `form:"type" validate:"required,min=1,gte=1"`
	Origin     int64 `form:"origin" validate:"required,min=1,gte=1"`
}

// ArgPage .
type ArgPage struct {
	State int8 `form:"state"`
	Pn    int  `form:"pn"`
	Ps    int  `form:"ps"`
}

// ArgSalary salary coupon.
type ArgSalary struct {
	Mid        int64  `form:"mid" validate:"required,min=1,gte=1"`
	CouponType int64  `form:"type" validate:"required,min=1,gte=1"`
	Count      int    `form:"count" validate:"required,min=1,gte=1"`
	BatchToken string `form:"batch_no" validate:"required"`
	AppID      int64  `form:"app_id"`
}

// ArgUseAllowance allowance coupon use.
type ArgUseAllowance struct {
	Mid            int64   `form:"mid" validate:"required,min=1,gte=1"`
	CouponToken    string  `form:"coupon_token" validate:"required"`
	Remark         string  `form:"remark" validate:"required"`
	OrderNO        string  `form:"order_id" validate:"required"`
	Price          float64 `form:"price" validate:"required"`
	Platform       string  `form:"platform" default:"pc"`
	MobiApp        string  `form:"mobi_app"`
	PanelType      string  `form:"panel_type" default:"normal"`
	Device         string  `form:"device"`
	Build          int64   `form:"build"`
	ProdLimMonth   int8    `form:"product_limit_month"`
	ProdLimRenewal int8    `form:"product_limit_renewal"`
}

// ArgCount allowance count.
type ArgCount struct {
	Mid int64 `form:"mid" validate:"required,min=1,gte=1"`
}

//ArgReceiveAllowance .
type ArgReceiveAllowance struct {
	Mid        int64  `form:"mid" validate:"required" json:"mid"`
	BatchToken string `form:"batch_token" validate:"required" json:"batch_token"`
	OrderNo    string `form:"order_no" validate:"required" json:"order_no"`
	Appkey     string `form:"appkey" validate:"required" json:"appkey"`
}

//ArgAllowanceCheck .
type ArgAllowanceCheck struct {
	Mid     int64  `form:"mid" validate:"required" json:"mid"`
	OrderNo string `form:"order_no" validate:"required" json:"order_no"`
}

// ArgPrizeDraw struct .
type ArgPrizeDraw struct {
	Mid      int64 `form:"mid" validate:"required,gte=1"`
	CardType int8  `form:"card_type" validate:"gte=0,lte=2" json:"card_type"`
}

// ArgAllowanceCoupons arg allowance coupon.
type ArgAllowanceCoupons struct {
	Mid   int64
	State int8
}
