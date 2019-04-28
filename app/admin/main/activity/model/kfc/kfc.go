package kfc

import "go-common/library/time"

// ListParams .
type ListParams struct {
	CouponCode string `form:"coupon_code"`
	Mid        int64  `form:"mid"`
	Pn         int    `form:"pn" default:"0" validate:"min=0"`
	Ps         int    `form:"ps" default:"15" validate:"min=1"`
}

//BnjKfcCoupon def
type BnjKfcCoupon struct {
	ID         int64     `json:"id"  gorm:"column:id"`
	Mid        int64     `json:"mid"  gorm:"column:mid"`
	CouponCode string    `json:"coupon_code"  gorm:"column:coupon_code"`
	Desc       string    `json:"desc"  gorm:"column:desc"`
	State      int       `json:"state"  gorm:"column:state"`
	DeleteTime time.Time `json:"delete_time"  gorm:"column:delete_time" time_format:"2006-01-02 15:04:05"`
	Ctime      time.Time `json:"ctime"  gorm:"column:ctime" time_format:"2006-01-02 15:04:05"`
	Mtime      time.Time `json:"mtime"  gorm:"column:mtime" time_format:"2006-01-02 15:04:05"`
}

// TableName BnjKfcCoupon def
func (BnjKfcCoupon) TableName() string {
	return "bnj_kfc_coupon"
}
