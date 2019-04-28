package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func notices(c *bm.Context) {
	v := new(struct {
		Type     int `form:"type"`
		Status   int `form:"status"`
		Platform int `form:"platform"`
		From     int `form:"from" default:"0" validate:"min=0"`
		Limit    int `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, notices, err := svr.Notices(c, v.Type, v.Status, v.Platform, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.Notices error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    notices,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func insertNotice(c *bm.Context) {
	v := new(struct {
		Title    string `form:"title" validate:"required"`
		Type     int    `form:"type" validate:"required"`
		Platform int    `form:"platform" validate:"required"`
		Link     string `form:"link" validate:"required"`
		Status   int    `form:"status" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.InsertNotice(c, v.Title, v.Type, v.Platform, v.Link, v.Status)
	if err != nil {
		log.Error("growup svr.Notices error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateNotice(c *bm.Context) {
	v := new(struct {
		ID       int64  `form:"id"`
		Title    string `form:"title"`
		Type     int    `form:"type"`
		Platform int    `form:"platform"`
		Link     string `form:"link"`
		Status   int    `form:"status"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateNotice(c, v.Type, v.Platform, v.Title, v.Link, v.ID, v.Status)
	if err != nil {
		log.Error("growup svr.Notices error(%v)", err)
	}
	c.JSON(nil, err)
}
