package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func addAppeal(c *bm.Context) {
	var (
		err     error
		params  = c.Request.Form
		btid    int64
		bid     int64
		midI, _ = c.Get("mid")
		btidStr = params.Get("business_typeid")
		bidStr  = params.Get("bid")
		reason  = params.Get("reason")
	)
	if btid, err = strconv.ParseInt(btidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt err(err) %v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if bid, err = strconv.ParseInt(bidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt err(err) %v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = creditSvc.AddAppeal(c, btid, bid, midI.(int64), reason); err != nil {
		log.Error("creditSvc.AddAppeal error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func appealStatus(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
		mid    int64
		bid    int64
		bidStr = params.Get("bid")
	)
	if bid, err = strconv.ParseInt(bidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt err(err) %v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if ok {
		mid, _ = midI.(int64)
	}
	state, err := creditSvc.AppealState(c, mid, bid)
	if err != nil {
		log.Error("creditSvc.AppealState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	// 未申述过 true, 已申诉过flase
	c.JSON(!state, nil)
}
