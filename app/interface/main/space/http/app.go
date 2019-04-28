package http

import (
	"go-common/app/interface/main/space/model"
	bm "go-common/library/net/http/blademaster"
)

func appIndex(c *bm.Context) {
	v := new(model.AppIndexArg)
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		v.Mid = midInter.(int64)
	}
	c.JSON(spcSvc.AppIndex(c, v))
}

func appPlayedGame(c *bm.Context) {
	var (
		mid   int64
		list  []*model.AppGame
		count int
		err   error
	)
	v := new(struct {
		VMid     int64  `form:"mid" validate:"min=1"`
		Platform string `form:"platform" default:"android" validate:"required"`
		Pn       int    `form:"pn" default:"1" validate:"min=1"`
		Ps       int    `form:"ps" default:"20" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	list, count, err = spcSvc.AppPlayedGame(c, mid, v.VMid, v.Platform, v.Pn, v.Ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"pn":    v.Pn,
		"ps":    v.Ps,
		"count": count,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}

func appTopPhoto(c *bm.Context) {
	var mid int64
	v := new(struct {
		Vmid     int64  `form:"mid" validate:"min=1"`
		Platform string `form:"platform" default:"android"`
		Device   string `form:"device"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	imgURL := spcSvc.AppTopPhoto(c, mid, v.Vmid, v.Platform, v.Device)
	c.JSON(&struct {
		ImgURL string `json:"img_url"`
	}{ImgURL: imgURL}, nil)
}
