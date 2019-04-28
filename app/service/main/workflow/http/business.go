package http

import (
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// UpBusinessExtra update business extra
func upBusinessExtra(c *bm.Context) {
	ap := new(struct {
		Cid      int32  `form:"cid" validate:"required"`
		Mid      int64  `form:"mid" validate:"required"`
		Business int8   `form:"business" validate:"required"`
		Key      string `form:"key" validate:"required"`
		Val      string `form:"val" validate:"required"`
	})
	if err := c.BindWith(ap, binding.FormPost); err != nil {
		return
	}
	c.JSON(nil, wkfSvc.UpBusinessExtraV2(c, ap.Cid, ap.Mid, ap.Business, ap.Key, ap.Val))
}

// BusinessExtra get business extra
func businessExtra(c *bm.Context) {
	ap := new(struct {
		Cid      int32 `form:"cid" validate:"required"`
		Mid      int64 `form:"mid" validate:"required"`
		Business int8  `form:"business" validate:"required"`
	})
	if err := c.Bind(ap); err != nil {
		return
	}
	c.JSON(wkfSvc.BusinessExtraV2(c, ap.Cid, ap.Mid, ap.Business))
}
