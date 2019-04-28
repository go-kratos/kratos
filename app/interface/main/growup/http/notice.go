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
		Platform int `form:"platform"`
		From     int `form:"from" default:"0" validate:"min=0"`
		Limit    int `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, total, err := svc.GetNotices(c, v.Type, v.Platform, v.From, v.Limit)
	if err != nil {
		log.Error("growup svc.GetNotices error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
		"paging": map[string]int64{
			"page_size": int64(v.Limit),
			"total":     total,
		},
	}))
}

func latestNotice(c *bm.Context) {
	v := new(struct {
		Platform int `form:"platform" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	notice, err := svc.LatestNotice(c, v.Platform)
	if err != nil {
		log.Error("svc.LatestNotice error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(notice, nil)
}
