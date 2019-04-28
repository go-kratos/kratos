package http

import (
	"strconv"

	"go-common/app/interface/main/space/conf"
	"go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func riderList(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
		rider  *assist.AssistUpsPager
		err    error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxRiderPs {
		ps = conf.Conf.Rule.MaxRiderPs
	}
	if rider, err = spcSvc.RiderList(c, mid, pn, ps); err != nil {
		log.Error("spcSvc.RiderList(%d,%d,%d) error(%v)", mid, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"pn":    pn,
		"ps":    ps,
		"count": int(rider.Pager.Total),
	}
	data["page"] = page
	data["list"] = rider.Data
	c.JSON(data, nil)
}

func exitRider(c *bm.Context) {
	var (
		mid, upMid int64
		err        error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	upMidStr := params.Get("up_mid")
	if upMid, err = strconv.ParseInt(upMidStr, 10, 64); err != nil || upMid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.ExitRider(c, mid, upMid))
}
