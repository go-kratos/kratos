package http

import (
	"go-common/app/admin/main/esports/model"
	bm "go-common/library/net/http/blademaster"
)

func addDetail(c *bm.Context) {
	v := new(model.MatchDetail)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.AddDetail(c, v))
}

func editDetail(c *bm.Context) {
	v := new(model.MatchDetail)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.EditDetail(c, v))
}

func forbidDetail(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int   `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.ForbidDetail(c, v.ID, v.State))
}

func onLine(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int64 `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.UpOnline(c, v.ID, v.State))
}

func listDetail(c *bm.Context) {
	var (
		list []*model.MatchDetail
		cnt  int64
		err  error
	)
	v := new(struct {
		MaID int64 `json:"ma_id,omitempty" form:"ma_id" validate:"required"`
		Pn   int64 `form:"pn" validate:"min=0" default:"1"`
		Ps   int64 `form:"ps" validate:"min=0,max=30" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if list, cnt, err = esSvc.ListDetail(c, v.MaID, v.Pn, v.Ps); err != nil {
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
