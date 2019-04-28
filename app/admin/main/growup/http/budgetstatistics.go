package http

import (
	"net/http"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func budgetDayStatistics(c *bm.Context) {
	v := new(struct {
		Type  int `form:"type"`
		From  int `form:"from" default:"0" validate:"min=0"`
		Limit int `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, infos, err := svr.BudgetDayStatistics(c, v.Type, v.From, v.Limit)
	if err != nil {
		log.Error("s.budgetDayStatistics error(%v)", err)
		c.JSON(nil, err)
		return
	}

	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": infos,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func budgetDayGraph(c *bm.Context) {
	v := new(struct {
		Type int `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	info, err := svr.BudgetDayGraph(c, v.Type)
	if err != nil {
		log.Error("s.budgetDayGraph error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": info,
	}))
}

func budgetMonthStatistics(c *bm.Context) {
	v := new(struct {
		Type  int `form:"type"`
		From  int `form:"from" default:"0" validate:"min=0"`
		Limit int `form:"limit" default:"20" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, infos, err := svr.BudgetMonthStatistics(c, v.Type, v.From, v.Limit)
	if err != nil {
		log.Error("s.budgetMonthStatistics BudgetMonthStatistics error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": infos,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}
