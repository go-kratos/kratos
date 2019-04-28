package http

import (
	bm "go-common/library/net/http/blademaster"
)

func orderList(c *bm.Context) {
	args := new(struct {
		Mid          int64  `form:"mid"`
		OrderNo      string `form:"order_no"`
		Status       int8   `form:"status"`
		PaymentStime int64  `form:"payment_stime"`
		PaymentEtime int64  `form:"payment_etime"`
		PageNum      int64  `form:"pn" default:"1"`
		PageSize     int64  `form:"ps" default:"20"`
	})
	if err := c.Bind(args); err != nil {
		return
	}

	c.JSON(tvSrv.OrderList(args.Mid, args.PageNum, args.PageSize, args.PaymentStime, args.PaymentEtime, args.Status, args.OrderNo))

}
