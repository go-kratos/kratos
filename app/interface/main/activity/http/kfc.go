package http

import (
	bm "go-common/library/net/http/blademaster"
)

func kfcInfo(c *bm.Context) {
	p := new(struct {
		ID int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(kfcSvc.KfcInfo(c, p.ID, mid))
}

func kfcUse(c *bm.Context) {
	p := new(struct {
		CouponCode string `form:"coupon_code" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	if len([]rune(p.CouponCode)) == 12 {
		kfcSvc.KfcUse(c, p.CouponCode)
	}
	c.JSON(200, nil)
}

func deliverKfc(c *bm.Context) {
	p := new(struct {
		ID  int64 `form:"id" validate:"min=1"`
		Mid int64 `form:"mid" validate:"min=1"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, kfcSvc.DeliverKfc(c, p.ID, p.Mid))
}
