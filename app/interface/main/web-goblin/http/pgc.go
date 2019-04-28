package http

import (
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func pgcfull(c *bm.Context) {
	var (
		err   error
		items interface{}
	)
	v := new(struct {
		Pn     int64  `form:"pn"          validate:"min=1"`
		Ps     int64  `form:"ps"          validate:"min=1,max=50"`
		Tp     int    `form:"season_type" validate:"min=1,max=5"`
		Source string `form:"bsource"     validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !accessCard(_pgcFull, v.Source) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if items, err = srvWeb.PgcFull(c, v.Tp, v.Pn, v.Ps, v.Source); err != nil {
		c.JSON(nil, err)
		return
	}
	if items == nil {
		items = struct{}{}
	}
	data := make(map[string]interface{}, 5)
	data["app_name"] = _appName
	data["package_name"] = _packageName
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["source"] = v.Source
	data["shortvideos"] = items
	c.JSONMap(data, nil)
}

func pgcincre(c *bm.Context) {
	var (
		err   error
		items interface{}
	)
	v := new(struct {
		Pn        int64  `form:"pn"          validate:"min=1"`
		Ps        int64  `form:"ps"          validate:"min=1,max=50"`
		Tp        int    `form:"season_type" validate:"min=1,max=5"`
		StartTime int64  `form:"start_ts"    validate:"required"`
		EndTime   int64  `form:"end_ts"      validate:"required"`
		Source    string `form:"bsource"     validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !accessCard(_pgcIncre, v.Source) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if items, err = srvWeb.PgcIncre(c, v.Tp, v.Pn, v.Ps, v.StartTime, v.EndTime, v.Source); err != nil {
		c.JSON(nil, err)
		return
	}
	if items == nil {
		items = struct{}{}
	}
	data := make(map[string]interface{}, 5)
	data["app_name"] = _appName
	data["package_name"] = _packageName
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["source"] = v.Source
	data["shortvideos"] = items
	c.JSONMap(data, nil)
}
