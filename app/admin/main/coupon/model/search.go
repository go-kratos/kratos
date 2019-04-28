package model

// SearchData search result detail.
type SearchData struct {
	Code int `json:"code"`
	Data *struct {
		Order  string                 `json:"order"`
		Sort   string                 `json:"sort"`
		Page   *SearchPage            `json:"page"`
		Result []*CouponAllowanceInfo `json:"result"`
	}
}

// SearchPage struct.
type SearchPage struct {
	PN    int `json:"num"`
	PS    int `json:"size"`
	Total int `json:"total"`
}

// ArgAllowanceSearch allowance search struct
type ArgAllowanceSearch struct {
	PN          int    `form:"pn" default:"1"`
	PS          int    `form:"ps" default:"20"`
	AppID       int64  `form:"app_id"`
	Mid         int64  `form:"mid" validate:"required,min=1,gte=1"`
	OrderNO     string `form:"order_no"`
	CouponToken string `form:"coupon_token"`
	BatchToken  string `form:"batch_token"`
}

// PageCouponInfo common page info.
type PageCouponInfo struct {
	Count int         `json:"count"`
	Item  interface{} `json:"item"`
}
