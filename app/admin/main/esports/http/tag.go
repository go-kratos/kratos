package http

import (
	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func tagInfo(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(esSvc.TagInfo(c, v.ID))
}

func tagList(c *bm.Context) {
	var (
		list []*model.Tag
		cnt  int64
		err  error
	)
	v := new(struct {
		Pn int64 `form:"pn" validate:"min=0"`
		Ps int64 `form:"ps" validate:"min=0,max=30"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	if v.Ps == 0 {
		v.Ps = 20
	}
	if list, cnt, err = esSvc.TagList(c, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int64{
		"num":   v.Pn,
		"size":  v.Ps,
		"count": cnt,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}

func addTag(c *bm.Context) {
	v := new(model.Tag)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.AddTag(c, v))
}

func editTag(c *bm.Context) {
	v := new(model.Tag)
	if err := c.Bind(v); err != nil {
		return
	}
	if v.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.EditTag(c, v))
}

func forbidTag(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"min=1"`
		State int   `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.ForbidTag(c, v.ID, v.State))
}
