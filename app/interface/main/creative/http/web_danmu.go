package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func webListDmPurchases(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	danmus, err := danmuSvc.AdvDmPurchaseList(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(danmus, nil)
}

func webPassDmPurchase(c *bm.Context) {
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	danmuIDStr := params.Get("id")
	danmuID, err := strconv.ParseInt(danmuIDStr, 10, 64)
	if err != nil || danmuID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.PassDmPurchase(c, mid, danmuID, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDenyDmPurchase(c *bm.Context) {
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	danmuIDStr := params.Get("id")
	danmuID, err := strconv.ParseInt(danmuIDStr, 10, 64)
	if err != nil || danmuID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.DenyDmPurchase(c, mid, danmuID, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webCancelDmPurchase(c *bm.Context) {
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	danmuIDStr := params.Get("id")
	danmuID, err := strconv.ParseInt(danmuIDStr, 10, 64)
	if err != nil || danmuID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.CancelDmPurchase(c, mid, danmuID, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDmList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	idStr := params.Get("id")
	pool := params.Get("pool")
	aidStr := params.Get("aid")
	midStr := params.Get("mid")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
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
	if err != nil || ps <= 0 {
		ps = 20
	}
	ip := metadata.String(c, metadata.RemoteIP)
	order := "progress"
	orderStr := params.Get("order")
	if orderStr != order {
		order = "ctime"
	}
	list, err := danmuSvc.List(c, mid, aid, id, pn, ps, order, pool, midStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(list, nil)
}

func webDmEdit(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("cid")
	dmidsStr := params.Get("dmids")
	stateStr := params.Get("state")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
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
		return
	}
	c.JSON(nil, nil)
}

func webDmDistri(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	params := c.Request.Form
	cidStr := params.Get("cid")
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	v, err := arcSvc.Video(c, mid, aid, cid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	list := make(map[int64]int64)
	if v.From == "vupload" {
		list, err = danmuSvc.Distri(c, mid, cid, ip)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(map[string]interface{}{
		"duration": v.Duration * 1000,
		"type":     v.From,
		"list":     list,
	}, nil)
}

func webDmTransfer(c *bm.Context) {
	params := c.Request.Form
	fromCIDStr := params.Get("from_cid")
	toCIDStr := params.Get("to_cid")
	offsetStr := params.Get("offset")
	offset, err := strconv.ParseFloat(offsetStr, 64)
	if err != nil {
		offset = float64(0)
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	fromCID, err := strconv.ParseInt(fromCIDStr, 10, 64)
	if err != nil || fromCID <= 0 {
		err = ecode.RequestErr
		return
	}
	toCID, err := strconv.ParseInt(toCIDStr, 10, 64)
	if err != nil || toCID <= 0 {
		err = ecode.RequestErr
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	err = danmuSvc.Transfer(c, mid, fromCID, toCID, offset, ak, ck, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDmUpPool(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("cid")
	dmidsStr := params.Get("dmids")
	poolStr := params.Get("pool")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		err = ecode.RequestErr
		return
	}
	pool, err := strconv.Atoi(poolStr)
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
	err = danmuSvc.UpPool(c, id, dmids, int8(pool), mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDmReportCheck(c *bm.Context) {
	var err error
	params := c.Request.Form
	cidStr := params.Get("cid")
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	dmidStr := params.Get("dmid")
	dmid, err := strconv.ParseInt(dmidStr, 10, 64)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	opStr := params.Get("op")
	op, err := strconv.ParseInt(opStr, 10, 64)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.DmReportCheck(c, mid, cid, dmid, op, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDmProtectArchive(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	data, err := danmuSvc.DmProtectArchive(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func webDmProtectList(c *bm.Context) {
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	sortStr := params.Get("sort")
	if sortStr != "playtime" {
		sortStr = "ctime"
	}
	pnStr := params.Get("page")
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn <= 0 {
		pn = 1
	}
	aidStr := params.Get("aid")
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	vlist, err := danmuSvc.DmProtectList(c, mid, pn, aidStr, sortStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(vlist, nil)
}

func webDmProtectOper(c *bm.Context) {
	params := c.Request.Form
	statusStr := params.Get("status")
	status, err := strconv.ParseInt(statusStr, 10, 64)
	if err != nil || status != 1 {
		status = 0
	}
	idsStr := params.Get("ids")
	_, err = xstr.SplitInts(idsStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	err = danmuSvc.DmProtectOper(c, mid, status, idsStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webDmReport(c *bm.Context) {
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn <= 0 {
		pn = 1
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil || ps <= 0 {
		ps = 20
	}
	aidStr := params.Get("aid")
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	data, err := danmuSvc.DmReportList(c, mid, pn, ps, aidStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func webUserMid(c *bm.Context) {
	params := c.Request.Form
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	nameStr := params.Get("name")
	ip := metadata.String(c, metadata.RemoteIP)
	mid, err := danmuSvc.UserMid(c, nameStr, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"mid": mid,
	}, nil)
}
