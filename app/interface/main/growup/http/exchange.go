package http

import (
	"net/http"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func goodsState(c *bm.Context) {
	data, err := svc.GoodsState(c)
	if err != nil {
		log.Error("growup svc.GoodsState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func goodsShow(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.GoodsShow(c, mid)
	if err != nil {
		log.Error("growup svc.GoodsShow error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func goodsRecord(c *bm.Context) {
	v := new(struct {
		Page int `form:"page" default:"1" validate:"min=1"`
		Size int `form:"size" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, total, err := svc.GoodsRecord(c, mid, v.Page, v.Size)
	if err != nil {
		log.Error("growup svc.GoodsRecord error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
		"paging": map[string]int{
			"page_size": v.Size,
			"total":     total,
		},
	}))
}

func goodsBuy(c *bm.Context) {
	var err error
	v := new(struct {
		ProductID string `form:"product_id" validate:"required"`
		GoodsType int    `form:"goods_type" validate:"required"`
		Price     int64  `form:"price" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if err = svc.GoodsBuy(c, mid, v.ProductID, v.GoodsType, v.Price); err != nil {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}
