package http

import (
	"go-common/app/interface/main/space/model"
	bm "go-common/library/net/http/blademaster"
)

func tagSub(c *bm.Context) {
	v := new(struct {
		TagID int64 `form:"tag_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.TagSub(c, mid, v.TagID))
}

func tagCancelSub(c *bm.Context) {
	v := new(struct {
		TagID int64 `form:"tag_id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.TagCancelSub(c, mid, v.TagID))
}

func tagSubList(c *bm.Context) {
	var (
		mid   int64
		tags  []*model.Tag
		total int
		err   error
	)
	v := new(struct {
		Vmid int64 `form:"vmid" validate:"min=1"`
		Pn   int   `form:"pn" default:"1" validate:"min=1"`
		Ps   int   `form:"ps" default:"100" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if tags, total, err = spcSvc.TagSubList(c, mid, v.Vmid, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	type data struct {
		Tags []*model.Tag `json:"tags"`
		*model.Page
	}
	c.JSON(&data{Tags: tags, Page: &model.Page{Pn: v.Pn, Ps: v.Ps, Total: total}}, nil)
}
