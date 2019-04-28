package model

import (
	account "go-common/app/service/main/account/model"
)

// NavResp  struct of nav api response
type NavResp struct {
	IsLogin bool `json:"isLogin"`
	//AccessStatus  int    `json:"accessStatus"`
	//DueRemark     string `json:"dueRemark"`
	EmailVerified int    `json:"email_verified"`
	Face          string `json:"face"`
	LevelInfo     struct {
		Cur     int         `json:"current_level"`
		Min     int         `json:"current_min"`
		NowExp  int         `json:"current_exp"`
		NextExp interface{} `json:"next_exp"`
	} `json:"level_info"`
	Mid            int64   `json:"mid"`
	MobileVerified int     `json:"mobile_verified"`
	Coins          float64 `json:"money"`
	Moral          float32 `json:"moral"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"officialVerify"`
	Pendant        account.PendantInfo `json:"pendant"`
	Scores         int                 `json:"scores"`
	Uname          string              `json:"uname"`
	VipDueDate     int64               `json:"vipDueDate"`
	VipStatus      int                 `json:"vipStatus"`
	VipType        int                 `json:"vipType"`
	VipPayType     int32               `json:"vip_pay_type"`
	Wallet         *Wallet             `json:"wallet"`
	HasShop        bool                `json:"has_shop"`
	ShopURL        string              `json:"shop_url"`
	AllowanceCount int                 `json:"allowance_count"`
}

// FailedNavResp struct of failed nav response
type FailedNavResp struct {
	IsLogin bool `json:"isLogin"`
}

// Wallet struct.
type Wallet struct {
	Mid           int64   `json:"mid"`
	BcoinBalance  float32 `json:"bcoin_balance"`
	CouponBalance float32 `json:"coupon_balance"`
	CouponDueTime int64   `json:"coupon_due_time"`
}
