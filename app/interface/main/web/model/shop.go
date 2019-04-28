package model

// ShopInfo shop info.
type ShopInfo struct {
	ShopID       int64  `json:"shopId"`
	VAppID       string `json:"vAppId"`
	AppID        string `json:"appId"`
	Name         string `json:"name"`
	JumpURL      string `json:"jumpUrl"`
	ShowItemsTab int    `json:"showItemsTab"`
}
