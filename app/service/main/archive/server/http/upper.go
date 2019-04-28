package http

import (
	"strconv"

	"context"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func upperPassed(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// deal
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 || ps > 100 {
		ps = 20
	}
	as, err := arcSvc.UpperPassed3(c, mid, pn, ps)
	if err != nil {
		if ec := ecode.Cause(err); ec != ecode.NothingFound {
			log.Error("arcSvc.UpperPassed(%d) error(%d)", mid, err)
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(as, err)
}

// upperCount write the count of archives of Up.
func upperCount(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	count, err := arcSvc.UpperCount(c, mid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("arcSvc.UpperCount(%d) error(%d)", mid, err)
		return
	}
	var res struct {
		Count int `json:"count"`
	}
	res.Count = count
	c.JSON(res, nil)
}

// uppersCount uppers count
func uppersCount(c *bm.Context) {
	params := c.Request.Form
	midsStr := params.Get("mids")
	// check params
	mids, err := xstr.SplitInts(midsStr)
	if err != nil {
		log.Error("query mids(%s) split error(%v)", midsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(mids) > 20 {
		log.Error("query mids(%s) too long", midsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.UppersCount(c, mids))
}

// upperCache delete user cache.
func upperCache(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	action := params.Get("modifiedAttr")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if action == "updateUname" || action == "updateFace" || action == "updateByAdmin" {
		c.JSON(nil, arcSvc.UpperCache(context.TODO(), mid, action))
	}
}
