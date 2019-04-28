package http

import (
	"encoding/json"
	"go-common/app/service/openplatform/ticket-sales/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"net/http"

	"go-common/library/net/http/blademaster/binding"
)

func settleCompare(c *bm.Context) {
	req := &model.GetSettleOrdersRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err = binding.Form.Bind(c.Request, req); err != nil {
		c.Render(http.StatusOK, render.MapJSON{
			"errno": 1,
			"msg":   err.Error(),
		})
		return
	}
	data, err := svc.GetSettleOrders(c, req.Date, req.Ref == 1, req.ExtParams, req.PageSize)
	if err != nil {
		c.Render(http.StatusOK, render.MapJSON{
			"errno": 1,
			"msg":   err.Error(),
		})
		return
	}
	c.Render(http.StatusOK, render.MapJSON{
		"errno": 0,
		"data":  data,
	})
}

func settleRepush(c *bm.Context) {
	var req interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = svc.RepushSettleOrders(c, req)
	if err != nil {
		c.Render(http.StatusOK, render.MapJSON{
			"errno": 1,
			"msg":   err.Error(),
		})
	}
	c.Render(http.StatusOK, render.MapJSON{
		"errno": 0,
		"msg":   "将在5分钟内重推",
	})
}
