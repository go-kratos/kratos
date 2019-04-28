package http

import (
	"strconv"

	"go-common/app/admin/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func updateStatState(c *bm.Context) {
	var (
		params   = c.Request.Form
		state    int64
		id       int64
		operater string
		err      error
	)
	state, err = strconv.ParseInt(params.Get("state"), 10, 8)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	operater = params.Get("operater")
	if operater == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.UpdateState(c, int8(state), id, operater)
	if err != nil {
		log.Error("spySrv.UpdateState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func updateStatQuantity(c *bm.Context) {
	var (
		params   = c.Request.Form
		count    int64
		id       int64
		operater string
		err      error
	)
	count, err = strconv.ParseInt(params.Get("count"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	operater = params.Get("operater")
	if operater == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.UpdateStatQuantity(c, count, id, operater)
	if err != nil {
		log.Error("spySrv.UpdateStatQuantity error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func deleteStat(c *bm.Context) {
	var (
		params   = c.Request.Form
		id       int64
		operater string
		err      error
	)
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	operater = params.Get("operater")
	if operater == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.DeleteStat(c, 1, id, operater)
	if err != nil {
		log.Error("spySrv.DeleteStat error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func addRemark(c *bm.Context) {
	var (
		params   = c.Request.Form
		remark   string
		id       int64
		operater string
		err      error
	)
	remark = params.Get("remark")
	if err != nil || len(remark) > model.MaxRemarkLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	operater = params.Get("operater")
	if operater == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.AddLog2(c, &model.Log{
		RefID:   id,
		Name:    operater,
		Module:  model.UpdateStat,
		Context: remark,
	})
	if err != nil {
		log.Error("spySrv.AddLog2 error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func remarkList(c *bm.Context) {
	var (
		params = c.Request.Form
		id     int64
		err    error
		data   []*model.Log
	)
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err = spySrv.LogList(c, id, model.UpdateStat)
	if err != nil {
		log.Error("spySrv.logList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, err)
}

func statPage(c *bm.Context) {
	var (
		params = c.Request.Form
		id     int64
		mid    int64
		t      int64
		ps, pn int
		err    error
		data   *model.StatPage
	)
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	t, err = strconv.ParseInt(params.Get("type"), 10, 8)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if int8(t) == model.AccountType {
		mid = id
		id = 0
	}
	if ps, err = strconv.Atoi(params.Get("ps")); err != nil {
		ps = model.DefPs
	}
	if pn, err = strconv.Atoi(params.Get("pn")); err != nil {
		pn = model.DefPn
	}
	data, err = spySrv.StatPage(c, mid, id, int8(t), pn, ps)
	if err != nil {
		log.Error("spySrv.StatPage error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, err)
}
