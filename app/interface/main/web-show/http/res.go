package http

import (
	"strconv"

	rsmdl "go-common/app/interface/main/web-show/model/resource"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid = "Buvid"
	_buvid       = "buvid3"
)

func resources(c *bm.Context) {
	arg := new(rsmdl.ArgRess)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid, arg.Sid, arg.Buvid = device(c)
	data, count, err := resSvc.Resources(c, arg)
	if err != nil {
		log.Error("resSvc.Resource error(%v)", err)
		c.JSON(nil, ecode.Degrade)
		return
	}
	c.JSONMap(map[string]interface{}{
		"count": count,
		"data":  data,
	}, nil)
}

func resource(c *bm.Context) {
	arg := new(rsmdl.ArgRes)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid, arg.Sid, arg.Buvid = device(c)
	data, count, err := resSvc.Resource(c, arg)
	if err != nil {
		log.Error("resSvc.Resource error(%v)", err)
		c.JSON(nil, ecode.Degrade)
		return
	}
	c.JSONMap(map[string]interface{}{
		"count": count,
		"data":  data,
	}, nil)
}

func relation(c *bm.Context) {
	arg := new(rsmdl.ArgAid)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Mid, arg.Sid, arg.Buvid = device(c)
	c.JSON(resSvc.Relation(c, arg))
}

func advideo(c *bm.Context) {
	arg := new(rsmdl.ArgAid)
	if err := c.Bind(arg); err != nil {
		return
	}
	midTemp, ok := c.Get("mid")
	if !ok {
		log.Info("mid not exist")
		arg.Mid = 0
	} else {
		arg.Mid = midTemp.(int64)
	}
	c.JSON(resSvc.VideoAd(c, arg), nil)
}

func urlMonitor(c *bm.Context) {
	params := c.Request.Form
	pfStr := params.Get("pf")
	pf, _ := strconv.Atoi(pfStr)
	c.JSON(resSvc.URLMonitor(c, pf), nil)
}

func device(c *bm.Context) (mid int64, sid, buvid string) {
	midTemp, ok := c.Get("mid")
	buvid = c.Request.Header.Get(_headerBuvid)
	if buvid == "" {
		cookie, _ := c.Request.Cookie(_buvid)
		if cookie != nil {
			buvid = cookie.Value
		}
	}
	if !ok {
		if sidCookie, err := c.Request.Cookie("sid"); err == nil {
			sid = sidCookie.Value
		}
	} else {
		mid = midTemp.(int64)
	}
	return
}
