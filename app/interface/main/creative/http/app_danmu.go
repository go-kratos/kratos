package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func appDmList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	idStr := params.Get("id")
	pool := params.Get("pool")
	aidStr := params.Get("aid")
	midIStr := params.Get("mid")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn <= 0 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 0 || pn > 10 {
		ps = 10
	}
	ip := metadata.String(c, metadata.RemoteIP)
	order := "progress"
	orderStr := params.Get("order")
	if orderStr != order {
		order = "ctime"
	}
	list, err := danmuSvc.List(c, mid, aid, id, pn, ps, order, pool, midIStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(list, nil)
}

func appDmEdit(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("cid")
	dmidsStr := params.Get("dmids")
	stateStr := params.Get("state")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		err = ecode.RequestErr
		return
	}
	state, err := strconv.Atoi(stateStr)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	var dmids []int64
	if dmids, err = xstr.SplitInts(dmidsStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("xstr.SplitInts(dmidsStr %v) err(%v)", dmidsStr, err)
		return
	}
	for _, dmid := range dmids {
		if dmid <= 0 {
			c.JSON(nil, ecode.RequestErr)
			log.Error("dmids range err (dmid %d) err(%v)", dmid, err)
			return
		}
	}
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.Edit(c, mid, id, int8(state), dmids, ip)
	if err != nil {
		c.JSON(nil, err)
	}
	c.JSON(nil, nil)
}

func appDmRecent(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		pn, ps int64
		err    error
	)
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps > 10 || ps < 1 {
		ps = 10
	}
	recent, err := danmuSvc.Recent(c, mid, pn, ps, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(recent, nil)
}

func appDmEditBatch(c *bm.Context) {
	params := c.Request.Form
	paramsStr := params.Get("params")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	err := danmuSvc.EditBatch(c, mid, paramsStr, ip)
	if err != nil {
		c.JSON(nil, err)
	}
	c.JSON(nil, nil)
}
