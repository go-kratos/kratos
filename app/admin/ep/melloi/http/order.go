package http

import (
	"strconv"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryOrder(c *bm.Context) {
	qor := model.QueryOrderRequest{}
	if err := c.BindWith(&qor, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	if err := qor.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryOrder(&qor))
}

func addOrder(c *bm.Context) {
	order := model.Order{}
	if err := c.BindWith(&order, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.AddOrder(&order))
}

func updateOrder(c *bm.Context) {
	order := model.Order{}
	if err := c.BindWith(&order, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.UpdateOrder(&order))
}

func delOrder(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.DelOrder(v.ID))
}

func addOrderReport(c *bm.Context) {
	report := new(model.OrderReport)
	if err := c.BindWith(&report, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}

	nameStr := ""
	if username, err := c.Request.Cookie("username"); err == nil || username != nil {
		nameStr = username.Value
	}
	c.JSON(nil, srv.AddReport(nameStr, report))
}

func queryOrderReport(c *bm.Context) {
	params := c.Request.Form
	orderID := params.Get("order_id")
	oid, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.QueryReportByOrderID(oid))
}
