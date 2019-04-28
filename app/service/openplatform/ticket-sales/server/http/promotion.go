package http

import (
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
)

func getPromo(c *bm.Context) {
	arg := new(rpcV1.PromoID)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.GetPromo(c, arg))
}

func createPromo(c *bm.Context) {
	arg := new(rpcV1.CreatePromoRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.CreatePromo(c, arg))
}

func operatePromo(c *bm.Context) {
	arg := new(rpcV1.OperatePromoRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.OperatePromo(c, arg))
}

func editPromo(c *bm.Context) {
	arg := new(rpcV1.EditPromoRequest)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(svc.EditPromo(c, arg))
}
