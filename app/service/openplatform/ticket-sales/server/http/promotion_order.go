package http

import (
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
)

func checkCreatePromoOrder(c *bm.Context) {
	arg := new(rpcV1.CheckCreatePromoOrderRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.CheckCreateStatus(c, arg))
}

func createPromoOrder(c *bm.Context) {
	arg := new(rpcV1.CreatePromoOrderRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.CreatePromoOrder(c, arg))
}

func payNotify(c *bm.Context) {
	arg := new(rpcV1.OrderID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.PromoPayNotify(c, arg))
}

func cancelOrder(c *bm.Context) {
	arg := new(rpcV1.OrderID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.CancelOrder(c, arg))
}

func checkIssue(c *bm.Context) {
	arg := new(rpcV1.OrderID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.CheckIssue(c, arg))
}

func finishIssue(c *bm.Context) {
	arg := new(rpcV1.FinishIssueRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.FinishIssue(c, arg))
}
