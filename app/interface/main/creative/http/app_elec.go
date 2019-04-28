package http

import (
	"go-common/app/interface/main/creative/model/elec"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

var (
	cb   = &elec.ChargeBill{}
	recl = &elec.RecentElecList{}
)

func appElecBill(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 20 {
		ps = 20
	}
	bal, err := elecSvc.Balance(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var bil *elec.ChargeBill
	elecStat, _ := elecSvc.UserState(c, mid, ip, ak, ck)
	if elecStat != nil && elecStat.State == "2" {
		bil, _ = elecSvc.AppDailyBill(c, mid, pn, ps, ip)
	}
	if bil == nil {
		bil = cb
	}
	c.JSON(map[string]interface{}{
		"balance": bal,
		"bill":    bil,
	}, nil)
}

func appElecRecentRank(c *bm.Context) {
	req := c.Request
	params := req.Form
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 20 {
		ps = 20
	}
	var rec *elec.RecentElecList
	elecStat, _ := elecSvc.UserState(c, mid, ip, ak, ck)
	if elecStat != nil && elecStat.State == "2" {
		rec, _ = elecSvc.RecentElec(c, mid, pn, ps, ip)
	}
	if rec == nil {
		rec = recl
	}
	c.JSON(rec, nil)
}
