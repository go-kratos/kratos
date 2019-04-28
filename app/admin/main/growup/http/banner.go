package http

import (
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func banners(c *bm.Context) {
	v := new(struct {
		From  int64 `form:"from" validate:"min=0" default:"0"`
		Limit int64 `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, bs, err := svr.Banners(c, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.Banners error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    bs,
		"paging": map[string]int64{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func addBanner(c *bm.Context) {
	v := new(struct {
		Image   string `form:"image" validate:"required"`
		Link    string `form:"link" validate:"required"`
		StartAt int64  `form:"start_at" validate:"required"`
		EndAt   int64  `form:"end_at" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	dup, err := svr.AddBanner(c, v.Image, v.Link, v.StartAt, v.EndAt)
	if err != nil {
		log.Error("growup svr.AddBanner error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    dup,
	}))
}

func editBanner(c *bm.Context) {
	v := new(struct {
		ID      int64  `form:"id" validate:"required"`
		Image   string `form:"image" validate:"required"`
		Link    string `form:"link" validate:"required"`
		StartAt int64  `form:"start_at" validate:"required"`
		EndAt   int64  `form:"end_at" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	dup, err := svr.EditBanner(c, v.ID, v.StartAt, v.EndAt, v.Image, v.Link)
	if err != nil {
		log.Error("growup svr.AddBanner error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    dup,
	}))
}

func off(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.Off(c, time.Now().Unix(), v.ID)
	if err != nil {
		log.Error("growup svr.AddBanner error(%v)", err)
	}
	c.JSON(nil, err)
}
