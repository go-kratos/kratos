package model

import "go-common/library/time"

//Promotion 活动表
type Promotion struct {
	PromoID       int64     `json:"promo_id"`
	Type          int16     `json:"type"`
	ItemID        int64     `json:"item_id"`
	SKUID         int64     `json:"sku_id"`
	Extra         int64     `json:"extra"`
	ExpireSec     int64     `json:"expire_sec"`
	SKUCount      int64     `json:"sku_count"`
	Amount        int64     `json:"amount"`
	BuyerCount    int64     `json:"buyer_count"`
	BeginTime     int64     `json:"begin_time"`
	EndTime       int64     `json:"end_time"`
	Status        int16     `json:"status"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
	PrivSKUID     int64     `json:"priv_sku_id"`
	UsableCoupons string    `json:"usable_coupons"`
}

//PromotionGroup 拼团表
type PromotionGroup struct {
	PromoID    int64     `json:"promo_id"`
	GroupID    int64     `json:"group_id"`
	UID        int64     `json:"uid"`
	OrderCount int64     `json:"order_count"`
	Status     int16     `json:"status"`
	ExpireAt   int64     `json:"expire_at"`
	Ctime      time.Time `json:"ctime"`
	Mtime      time.Time `json:"mtime"`
}

//PromotionOrder 拼团订单表
type PromotionOrder struct {
	PromoID  int64     `json:"promo_id"`
	GroupID  int64     `json:"group_id"`
	OrderID  int64     `json:"order_id"`
	IsMaster int16     `json:"is_master"`
	UID      int64     `json:"uid"`
	Status   int16     `json:"status"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
	SKUID    int64     `json:"sku_id"`
}
