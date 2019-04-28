package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// addPa 申请保护弹幕.
func addPa(c *bm.Context) {
	var (
		str   string
		err   error
		cid   int64
		dmids []int64
	)
	mid, _ := c.Get("mid")
	params := c.Request.Form
	str = params.Get("cid")
	cid, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	str = params.Get("dmids")
	if dmids, err = xstr.SplitInts(str); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.AddProtectApply(c, mid.(int64), cid, dmids)
	c.JSON(nil, err)

}

// recall 弹幕撤回
func recall(c *bm.Context) {
	var (
		str  string
		msg  string
		err  error
		cid  int64
		dmid int64
	)
	mid, _ := c.Get("mid")
	params := c.Request.Form
	str = params.Get("cid")
	cid, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmid, err = strconv.ParseInt(params.Get("dmid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	msg, err = dmSvc.Recall(c, mid.(int64), cid, dmid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.Recall(%v,%d,%d) error(%v)", mid, cid, dmid, err)
		return
	}
	res := map[string]interface{}{}
	if msg != "" {
		res["message"] = msg
	}
	c.JSONMap(res, err)
}

// midHash 获取用户mid hash.
func midHash(c *bm.Context) {
	var err error
	mid, _ := c.Get("mid")
	hash, err := dmSvc.MidHash(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.MidHash(%d) error(%v)", mid.(int64), err)
		return
	}
	res := map[string]interface{}{}
	res["data"] = map[string]interface{}{
		"hash": hash,
	}
	c.JSONMap(res, err)
}

// transfer 弹幕转移
func transfer(c *bm.Context) {
	p := c.Request.Form
	toCid, err := strconv.ParseInt(p.Get("to"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	offset, err := strconv.ParseFloat(p.Get("offset"), 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := c.Get("mid")
	if err = dmSvc.CheckExist(c, mid.(int64), toCid); err != nil {
		c.JSON(nil, err)
		return
	}
	fromCids, err := xstr.SplitInts(p.Get("from"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, cid := range fromCids {
		if cid == toCid {
			c.JSON(nil, ecode.DMTransferSame)
			return
		}
		if err = dmSvc.CheckExist(c, mid.(int64), cid); err != nil {
			c.JSON(nil, err)
			return
		}
		if err = dmSvc.TransferJob(c, mid.(int64), cid, toCid, offset); err != nil {
			log.Error("dmSvc.TransferJob(mid:%d,from:%d,to:%d,offset:%v) error(%v)", mid.(int64), cid, toCid, offset, err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, err)
}

func transferList(c *bm.Context) {
	var (
		cid int64
		p   = c.Request.Form
	)
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	l, err := dmSvc.TransferList(c, cid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(l, err)
}

func transferRetry(c *bm.Context) {
	var (
		id  int64
		err error
		p   = c.Request.Form
	)
	id, err = strconv.ParseInt(p.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := c.Get("mid")
	err = dmSvc.TransferRetry(c, id, mid.(int64))
	c.JSON(nil, err)
}
