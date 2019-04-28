package http

import (
	"time"

	"go-common/app/interface/main/web-goblin/conf"
	webmdl "go-common/app/interface/main/web-goblin/model/web"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_lmt      = 10000
	_pgcFull  = "pgcfull"
	_pgcIncre = "pgcincre"
	_ugcFull  = "ugcfull"
	_ugcIncre = "ugcincre"
)

func ugcfull(c *bm.Context) {
	var (
		err   error
		items []*webmdl.Mi
	)
	v := new(struct {
		Pn     int64  `form:"pn"  validate:"min=1"`
		Ps     int64  `form:"ps"  validate:"min=1,max=50"`
		Source string `form:"bsource"  validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !accessCard(_ugcFull, v.Source) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	source := "?bsource=" + v.Source
	if items, err = srvWeb.UgcFull(c, v.Pn, v.Ps, source); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 5)
	data["app_name"] = _appName
	data["package_name"] = _packageName
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["source"] = v.Source
	data["shortvideos"] = items
	c.JSONMap(data, nil)
}

func ugcincre(c *bm.Context) {
	var (
		err  error
		item []*webmdl.Mi
	)
	v := new(struct {
		Pn      int    `form:"pn"       validate:"min=1"`
		Ps      int    `form:"ps"       validate:"min=1,max=50"`
		StartTs int64  `form:"start_ts" validate:"required"`
		EndTs   int64  `form:"end_ts"   validate:"required"`
		Source  string `form:"bsource"   validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !accessCard(_ugcIncre, v.Source) {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if v.StartTs >= v.EndTs || v.EndTs-v.StartTs >= conf.Conf.OutSearch.Rspan {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Ps*v.Pn > _lmt {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	source := "?bsource=" + v.Source
	if item, err = srvWeb.UgcIncre(c, v.Pn, v.Ps, v.StartTs, v.EndTs, source); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 5)
	data["app_name"] = _appName
	data["package_name"] = _packageName
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["source"] = v.Source
	data["shortvideos"] = item
	c.JSONMap(data, nil)
}

func accessCard(flag, arg string) bool {
	var (
		b = false
		m = map[string][]string{
			_pgcFull:  conf.Conf.OutSearch.AcPgcFull,
			_pgcIncre: conf.Conf.OutSearch.AcPgcIncre,
			_ugcFull:  conf.Conf.OutSearch.AcUgcFull,
			_ugcIncre: conf.Conf.OutSearch.AcUgcIncre,
		}
	)
	for _, v := range m[flag] {
		if v == arg {
			b = true
			break
		}
	}
	return b
}
