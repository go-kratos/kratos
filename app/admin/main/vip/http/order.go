package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func orderList(c *bm.Context) {
	arg := new(model.ArgPayOrder)
	if err := c.Bind(arg); err != nil {
		return
	}
	res, count, err := vipSvc.OrderList(c, arg)
	info := new(model.PageInfo)
	info.Count = int(count)
	info.Item = res
	info.CurrentPage = arg.PN
	c.JSON(info, err)
}

func refund(c *bm.Context) {
	arg := new(struct {
		OrderNo      string  `form:"order_no" validate:"required"`
		RefundAmount float64 `form:"refund_amount" validate:"required"`
	})
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.Refund(c, arg.OrderNo, username.(string), arg.RefundAmount))
}
