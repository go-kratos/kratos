package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

var (
	_layout = "2006-01-02"
)

func statisGraph(c *bm.Context) {
	v := new(struct {
		Type    int64  `form:"type"`
		TagID   string `form:"tag_id" validate:"required"`
		Compare int    `form:"compare"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svr.StatisGraph(c, v.Type, v.TagID, v.Compare)
	if err != nil {
		log.Error("svr.StatisGraph error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func statisList(c *bm.Context) {
	v := new(struct {
		Type    int64  `form:"type"`
		TagID   string `form:"tag_id" validate:"required"`
		Compare int    `form:"compare"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svr.StatisList(c, v.Type, v.TagID, v.Compare)
	if err != nil {
		log.Error("svr.StatisList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func statisExport(c *bm.Context) {
	v := new(struct {
		Type    int64  `form:"type"`
		TagID   string `form:"tag_id" validate:"required"`
		Compare int    `form:"compare"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	content, err := svr.StatisExport(c, v.Type, v.TagID, v.Compare)
	if err != nil {
		log.Error("svr.StatisExport error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "statistics"),
	})
}

func ascList(c *bm.Context) {
	v := new(struct {
		Type     string  `form:"type" validate:"required"`
		Tags     []int64 `form:"tag_ids,split" validate:"required"`
		Date     string  `form:"date" validate:"required"`
		ScoreMin int     `form:"score_min"`
		ScoreMax int     `form:"score_max"`
		MID      int64   `form:"mid"`
		From     int     `form:"from" default:"0" validate:"min=0"`
		Limit    int     `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse(_layout, v.Date)
	if err != nil {
		return
	}
	total, data, err := svr.GetTrendAsc(c, v.Type, v.Tags, date, v.ScoreMin, v.ScoreMax, v.MID, v.From, v.Limit)
	if err != nil {
		log.Error("svr.GetTrendAsc error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func descList(c *bm.Context) {
	v := new(struct {
		Type     string  `form:"type" validate:"required"`
		Tags     []int64 `form:"tag_ids,split" validate:"required"`
		Date     string  `form:"date" validate:"required"`
		ScoreMin int     `form:"score_min"`
		ScoreMax int     `form:"score_max"`
		MID      int64   `form:"mid"`
		From     int     `form:"from" default:"0" validate:"min=0"`
		Limit    int     `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	date, err := time.Parse(_layout, v.Date)
	if err != nil {
		return
	}
	total, data, err := svr.GetTrendDesc(c, v.Type, v.Tags, date, v.ScoreMin, v.ScoreMax, v.MID, v.From, v.Limit)
	if err != nil {
		log.Error("svr.GetTrendDesc error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}
