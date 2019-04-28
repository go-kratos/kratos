package pay

// UserWaller http://info.bilibili.co/pages/viewpage.action?pageId=7559096
type UserWallet struct {
	Code        int    `json:"errno"`
	Message     string `json:"msg"`
	ShowMessage string `json:"showMsg"`
	Data        struct {
		Mid               int64   `json:"mid"`
		TotalBp           float64 `json:"totalBp"`
		DefaultBp         float64 `json:"defaultBp"`
		IosBp             float64 `json:"iosBp"`
		CouponBalance     float64 `json:"couponBalance"`
		AvailableBp       float64 `json:"availableBp"`
		UnavailableBp     float64 `json:"unavailableBp"`
		UnavailableReason string  `json:"unavailableReason"`
	} `json:"data"`
}
