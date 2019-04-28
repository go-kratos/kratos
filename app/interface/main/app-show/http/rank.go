package http

import (
	"strconv"

	"go-common/app/interface/main/app-show/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func rankAll(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	order := params.Get("order")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	plat := model.Plat(mobiApp, device)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > 100 || ps <= 0 {
		ps = 100
	}
	if ((pn-1)*ps)+1 > 100 {
		returnJSON(c, _emptyShowItems, nil)
		return
	}
	build, _ := strconv.Atoi(buildStr)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// GetAudit
	if audit, ok := rankSvc.Audit(c, mobiApp, order, plat, build, 0); ok {
		returnJSON(c, audit, nil)
	} else {
		data := rankSvc.RankShow(c, plat, 0, pn, ps, mid, order)
		returnJSON(c, data, nil)
	}
}

func rankRegion(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ridStr := params.Get("rid")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	plat := model.Plat(mobiApp, device)
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > 100 || ps <= 0 {
		ps = 100
	}
	if ((pn-1)*ps)+1 > 100 {
		returnJSON(c, _emptyShowItems, nil)
		return
	}
	build, _ := strconv.Atoi(buildStr)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// GetAudit
	if audit, ok := rankSvc.Audit(c, mobiApp, "all", plat, build, rid); ok {
		returnJSON(c, audit, nil)
	} else {
		data := rankSvc.RankShow(c, plat, rid, pn, ps, mid, "all")
		returnJSON(c, data, nil)
	}
}
