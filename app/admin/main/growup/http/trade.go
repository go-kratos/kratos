package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/app/admin/main/growup/model"

	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func goodsSync(c *blademaster.Context) {
	var arg = new(struct {
		GoodsType int `form:"goods_type" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	eff, err := svr.SyncGoods(c, arg.GoodsType)
	if err != nil {
		log.Error("svr.SyncGoods err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(eff, nil)
}

func goodsUpdate(c *blademaster.Context) {
	var arg = new(struct {
		ID       int64 `form:"id" validate:"required"`
		Discount int   `form:"discount" validate:"min=1,max=100"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	eff, err := svr.UpdateGoodsInfo(c, arg.Discount, arg.ID)
	if err != nil {
		log.Error("svr.UpdateGoodsInfo err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(eff, nil)
}

func goodsDisplay(c *blademaster.Context) {
	arg := new(struct {
		IDs     []int64 `form:"ids,split" validate:"required"`
		Display int     `form:"display" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	eff, err := svr.UpdateGoodsDisplay(c, model.DisplayStatus(arg.Display), arg.IDs)
	if err != nil {
		log.Error("svr.UpdateGoodsDisplay err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(eff, nil)
}

func goodsList(c *blademaster.Context) {
	arg := new(struct {
		From  int `form:"from" validate:"min=0" default:"0"`
		Limit int `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	total, list, err := svr.GoodsList(c, arg.From, arg.Limit)
	if err != nil {
		log.Error("svr.GoodsList err(%v)", err)
		c.JSON(nil, err)
		return
	}
	renderPagRes(list, total, arg.Limit)(c)
}

func orderList(c *blademaster.Context) {
	arg := new(model.OrderQueryArg)
	if err := c.Bind(arg); err != nil {
		return
	}
	total, list, err := svr.OrderList(c, arg, arg.From, arg.Limit)
	if err != nil {
		log.Error("svr.OrderList err(%v)", err)
		c.JSON(nil, err)
		return
	}
	renderPagRes(list, total, arg.Limit)(c)
}

func renderPagRes(list interface{}, total int64, ps int) func(c *blademaster.Context) {
	return func(c *blademaster.Context) {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    0,
			"message": "0",
			"data":    list,
			"paging": map[string]interface{}{
				"page_size": ps,
				"total":     total,
			},
		}))
	}
}

func orderExport(c *blademaster.Context) {
	arg := new(model.OrderQueryArg)
	if err := c.Bind(arg); err != nil {
		return
	}
	content, err := svr.OrderExport(c, arg, arg.From, arg.Limit)
	if err != nil {
		log.Error("svr.OrderExport err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "creative_order"),
	})
}

func orderStatistics(c *blademaster.Context) {
	arg := new(model.OrderQueryArg)
	if err := c.Bind(arg); err != nil {
		return
	}
	res, err := svr.OrderStatistics(c, arg)
	if err != nil {
		log.Error("svr.OrderStatistics err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
