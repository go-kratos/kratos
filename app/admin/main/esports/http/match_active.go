package http

import (
	"go-common/app/admin/main/esports/model"
	bm "go-common/library/net/http/blademaster"
)

func addAct(c *bm.Context) {
	v := new(model.ParamMA)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(esSvc.AddAct(c, v))
}

func editAct(c *bm.Context) {
	v := new(model.ParamMA)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(esSvc.EditAct(c, v))
}

func forbidAct(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int   `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.ForbidAct(c, v.ID, v.State))
}

func listAct(c *bm.Context) {
	var (
		list []*model.MatchModule
		cnt  int64
		err  error
	)
	v := new(struct {
		Mid int64 `form:"mid"`
		Pn  int64 `form:"pn" validate:"min=0" default:"1"`
		Ps  int64 `form:"ps" validate:"min=0,max=30" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if list, cnt, err = esSvc.ListAct(c, v.Mid, v.Pn, v.Ps); err != nil {
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
