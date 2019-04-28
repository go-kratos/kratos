package http

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// get splash
func transcode(c *bm.Context) {
	v := new(model.ReqTransode)
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(nil, auditSvc.Transcode(v))
}

// get splash
func hotword(c *bm.Context) {
	hotword := gobSvc.Hotword
	if hotword == nil {
		c.JSON(nil, ecode.ServiceUnavailable)
		return
	}
	c.JSON(hotword, nil)
}

// get splash
func splash(c *bm.Context) {
	v := new(struct {
		Channel string `form:"channel" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(gobSvc.PickSph(v.Channel))
}

func favorites(c *bm.Context) {
	v := new(model.FormFav)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.AccessKey != "" {
		if mid, ok := c.Get("mid"); ok { // if not logged in, not request follow
			c.JSON(favSvc.Favorites(c, v.ToReq(mid.(int64))))
			return
		}
	}
	c.JSON(nil, ecode.NoLogin)
}

func favAct(c *bm.Context) {
	v := new(model.FormFavAct)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.AccessKey != "" {
		if mid, ok := c.Get("mid"); ok { // if not logged in, not request follow
			c.JSON(nil, favSvc.FavAct(c, v.ToReq(mid.(int64))))
			return
		}
	}
	c.JSON(nil, ecode.NoLogin)
}

func applyPGC(c *bm.Context) {
	v := new(model.ReqApply)
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(nil, auditSvc.ApplyPGC(c, v))
}

func labels(c *bm.Context) {
	v := new(struct {
		CatType  int `form:"cat_type" validate:"required,min=1,max=2"`
		Category int `form:"category" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(gobSvc.Labels(c, v.CatType, v.Category))
}

func region(c *bm.Context) {
	var (
		err error
		res []*model.Region
		m   = make(map[string]interface{})
	)
	if res, err = tvSvc.Regions(c); err != nil {
		return
	}
	m["mtime"] = tvSvc.MaxTime
	m["data"] = res
	c.JSONMap(m, nil)
}
