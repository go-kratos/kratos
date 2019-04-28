package http

import (
	"strings"

	"go-common/app/service/main/point/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func pointInfo(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	c.JSON(svc.PointInfo(c, mid.(int64)))
}

func pointInfoInner(c *bm.Context) {
	m := new(model.ArgMid)
	if err := c.Bind(m); err != nil {
		return
	}
	c.JSON(svc.PointInfo(c, m.Mid))
}

func pointAddByBp(c *bm.Context) {
	var (
		err error
		p   int64
	)
	if err = checkAuth(c); err != nil {
		c.JSON(nil, err)
	}
	pa := new(model.ArgPointAdd)
	if err = c.Bind(pa); err != nil {
		log.Error("point add by bp bind %+v", err)
		return
	}
	if p, err = svc.PointAddByBp(c, pa); err != nil {
		log.Error("point add by bp(%+v) faild(%+v)", pa, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"point": p,
	}, nil)
}

func pointConsume(c *bm.Context) {
	var (
		err    error
		status int8
	)
	if err = checkAuth(c); err != nil {
		c.JSON(nil, err)
	}
	arg := new(model.ArgPointConsume)
	if err = c.Bind(arg); err != nil {
		log.Error("point consume bind %+v", err)
		return
	}
	if arg.Point <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = svc.ConsumePoint(c, arg); err != nil {
		log.Error("point consume(%+v) faild(%+v)", arg, err)
	}
	c.JSON(map[string]interface{}{
		"status": status,
	}, nil)
}

func pointHistory(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
		m      = new(model.ArgPointHistory)
	)
	if err := c.Bind(m); err != nil {
		return
	}
	phs, total, cursor, err := svc.PointHistory(c, mid.(int64), m.Cursor, m.PS)
	data := make(map[string]interface{})
	data["total"] = total
	data["phs"] = phs
	data["cursor"] = cursor
	c.JSON(data, err)
}

func oldPointHistory(c *bm.Context) {
	var (
		m = new(model.ArgOldPointHistory)
	)
	if err := c.Bind(m); err != nil {
		return
	}
	phs, total, err := svc.OldPointHistory(c, m.Mid, m.PN, m.PS)
	data := make(map[string]interface{})
	data["total"] = total
	data["phs"] = phs
	c.JSON(data, err)
}

func configs(c *bm.Context) {
	c.JSON(svc.AllConfig(c), nil)
}

func config(c *bm.Context) {
	arg := new(model.ArgConfig)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ChangeType == 0 {
		arg.ChangeType = model.Contract
	}
	c.JSON(svc.Config(c, int(arg.ChangeType), arg.Mid, arg.Bp))
}

func checkAuth(c *bm.Context) (err error) {
	req := c.Request
	params := req.Form
	sappkey := params.Get("appkey")
	if len(sappkey) == 0 || !strings.Contains(whiteAppkeys, sappkey) {
		err = ecode.AccessDenied
		return
	}
	return
}

func pointAdd(c *bm.Context) {
	var (
		err    error
		status int8
	)
	if err = checkAuth(c); err != nil {
		c.JSON(nil, err)
	}
	arg := new(model.ArgPoint)
	if err = c.Bind(arg); err != nil {
		log.Error("point add bind %+v", err)
		return
	}
	if arg.Point <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = svc.AddPoint(c, arg); err != nil {
		log.Error("point add(%+v) faild(%+v)", arg, err)
	}
	c.JSON(map[string]interface{}{
		"status": status,
	}, nil)
}
