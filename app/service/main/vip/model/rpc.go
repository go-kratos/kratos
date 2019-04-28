package model

// ArgPanel .
type ArgPanel struct {
	Mid       int64
	SortTp    int8
	IP        string
	MobiApp   string
	Device    string
	Platform  string
	Plat      int64
	PanelType string
	SubType   int32
	Month     int32
	Build     int64
	Lang      string
}

// ArgRPCMid .
type ArgRPCMid struct {
	Mid int64 `json:"mid"`
}

// ArgRPCMids .
type ArgRPCMids struct {
	Mids []int64 `json:"mids"`
}

// ArgRPCHistory .
type ArgRPCHistory struct {
	Mid       int64  `form:"mid" validate:"required"`
	StartDate string `form:"start_time"`
	EndDate   string `form:"end_time"`
	Pn        int    `form:"pn"`
	Ps        int    `form:"ps"`
}

//ArgRPCCreateOrder .
type ArgRPCCreateOrder struct {
	Mid       int64   `form:"mid" validate:"required,min=1,gte=1"`
	AppID     int64   `form:"appId" default:"0"`
	AppSubID  string  `form:"appSubId"`
	Months    int16   `form:"months" validate:"required,min=1,gte=1"`
	OrderType int8    `form:"orderType" `
	DType     int8    `form:"dtype"`
	Bmid      int64   `form:"bmid"`
	Platform  string  `form:"platform"`
	Price     float64 `form:"price"`
	IP        string  `form:"ip"`
}

// ArgRPCOrderNo .
type ArgRPCOrderNo struct {
	OrderNo string `json:"order_no"`
}

// ArgTips arg tips.
type ArgTips struct {
	Version  int64  `json:"version" form:"version"`
	Platform string `json:"platform" form:"platform" validate:"required"`
	Position int8   `json:"position" form:"position"`
}

// ArgCouponPanel coupon panel arg.
type ArgCouponPanel struct {
	Mid       int64 `json:"mid"`
	Sid       int64 `json:"sid"`
	Platform  int   `json:"platform"`
	MobiApp   string
	Device    string
	PanelType string
	Build     int64
}

// ArgCouponPanelV2 coupon panel arg.
type ArgCouponPanelV2 struct {
	Mid int64 `json:"mid"`
	Sid int64 `json:"sid"`
}

// ArgCouponCancel coupon cancel use.
type ArgCouponCancel struct {
	Mid         int64  `json:"mid"`
	CouponToken string `json:"coupon_token"`
	IP          string `json:"ip"`
}

// ArgPrivilegeDetail privilege by type.
type ArgPrivilegeDetail struct {
	Type     int8   `json:"type" form:"type"`
	Lang     string `json:"lang" form:"lang"`
	Platform string `json:"platform" form:"platform" default:"pc"`
}

// ArgPrivilegeBySid privilege by sid .
type ArgPrivilegeBySid struct {
	Sid      int64  `json:"sid" form:"sid" validate:"required"`
	Lang     string `form:"lang"`
	Platform string `json:"platform" form:"platform" default:"pc"`
}

// ArgPanelExplain arg explain .
type ArgPanelExplain struct {
	Mid int64 `json:"mid"`
}

// ArgDialog get order result withc appID and Platform.
type ArgDialog struct {
	OrderNo   string `form:"order_no" validate:"required"`
	Mid       int64
	AppID     int64  `form:"app_id"`
	Platform  string `form:"platform" default:"pc"`
	MobiApp   string `form:"mobi_app"`
	Device    string `form:"device"`
	Build     int64  `form:"build"`
	PanelType string `form:"panel_type" default:"normal"`
}
