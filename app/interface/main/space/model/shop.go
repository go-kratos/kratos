package model

// ShopInfo shop info.
type ShopInfo struct {
	ID         int64  `json:"id"`
	Mid        int64  `json:"mid"`
	Name       string `json:"name"`
	Logo       string `json:"logo"`
	URL        string `json:"url"`
	Status     int    `json:"status"`
	GoodsNum   int64  `json:"goods_num"`
	MonthSales int64  `json:"month_sales"`
	StatusV    string `json:"status_v"`
}

// ShopLinkInfo shop link info.
type ShopLinkInfo struct {
	ShopID       int64  `json:"shopId"`
	VAppID       string `json:"vAppId"`
	AppID        string `json:"appId"`
	Name         string `json:"name"`
	JumpURL      string `json:"jumpUrl"`
	ShowItemsTab int    `json:"showItemsTab"`
}
