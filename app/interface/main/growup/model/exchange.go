package model

import (
	"go-common/library/time"
)

// GoodsInfo goods info
type GoodsInfo struct {
	ID            int64     `json:"-"`
	ProductID     string    `json:"product_id"`
	ResourceID    int64     `json:"-"`
	GoodsType     int       `json:"goods_type"`
	Discount      int       `json:"discount"`
	IsDisplay     int       `json:"-"`
	DisplayOnTime time.Time `json:"-"`
	ProductName   string    `json:"product_name"`  // 商品名称
	OriginPrice   int64     `json:"origin_price"`  // 实时成本, 单位分
	CurrentPrice  int64     `json:"current_price"` // 实时售价, 单位分
	Month         int32     `json:"month"`         //有效期
}

// GoodsOrder goods order
type GoodsOrder struct {
	MID        int64     `json:"-"`
	OrderNo    string    `json:"-"`
	OrderTime  time.Time `json:"order_time"`
	GoodsType  int       `json:"-"`
	GoodsID    string    `json:"-"`
	GoodsName  string    `json:"goods_name"`
	GoodsPrice int64     `json:"goods_price"`
	GoodsCost  int64     `json:"goods_cost"`
}
