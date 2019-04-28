package http

import (
	"strconv"

	"go-common/app/admin/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func userInfo(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		err    error
	)
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := spySrv.UserInfo(c, mid)
	if err != nil {
		log.Error("spySrv.UserInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, err)
}

func historyPage(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		pn, ps int
		err    error
		data   *model.HistoryPage
	)
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps, err = strconv.Atoi(params.Get("ps")); err != nil {
		ps = model.DefPs
	}
	if pn, err = strconv.Atoi(params.Get("pn")); err != nil {
		pn = model.DefPn
	}
	q := &model.HisParamReq{Mid: mid, Pn: pn, Ps: ps}
	data, err = spySrv.HisoryPage(c, q)
	if err != nil {
		log.Error("spySrv.HisoryPage error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, err)
}

func resetBase(c *bm.Context) {
	var (
		params = c.Request.Form
		name   = params.Get("name")
		mid    int64
		err    error
	)
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.ResetBase(c, mid, name)
	if err != nil {
		log.Error("spySrv.ResetBase error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func refreshBase(c *bm.Context) {
	var (
		params = c.Request.Form
		name   = params.Get("name")
		mid    int64
		err    error
	)
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err = spySrv.RefreshBase(c, mid, name)
	if err != nil {
		log.Error("spySrv.RefreshBase error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func resetEvent(c *bm.Context) {
	var (
		params = c.Request.Form
		name   = params.Get("name")
		mid    int64
		err    error
	)
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.ResetEvent(c, mid, name)
	if err != nil {
		log.Error("spySrv.ResetEvent error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func clearCount(c *bm.Context) {
	var (
		params = c.Request.Form
		name   = params.Get("name")
		mid    int64
		err    error
	)
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.ClearCount(c, mid, name)
	if err != nil {
		log.Error("spySrv.ClearCount error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}
