package http

import (
	"go-common/app/interface/main/space/model"
	bm "go-common/library/net/http/blademaster"
)

func bangumiList(c *bm.Context) {
	var (
		mid   int64
		list  []*model.Bangumi
		count int
		err   error
	)
	v := new(struct {
		Vmid int64 `form:"vmid" validate:"min=1"`
		Pn   int   `form:"pn" default:"1" validate:"min=1"`
		Ps   int   `form:"ps" default:"15" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if list, count, err = spcSvc.BangumiList(c, mid, v.Vmid, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	type data struct {
		List []*model.Bangumi `json:"list"`
		*model.Page
	}
	c.JSON(&data{List: list, Page: &model.Page{Pn: v.Pn, Ps: v.Ps, Total: count}}, nil)
}

func bangumiConcern(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	v := new(struct {
		SeasonID int64 `form:"season_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, spcSvc.BangumiConcern(c, mid, v.SeasonID))
}

func bangumiUnConcern(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	v := new(struct {
		SeasonID int64 `form:"season_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, spcSvc.BangumiUnConcern(c, mid, v.SeasonID))
}
