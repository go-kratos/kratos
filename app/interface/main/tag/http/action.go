package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func likeArcTag(c *bm.Context) {
	var (
		aid int64
		tid int64
		err error
	)
	params := c.Request.Form
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	aidStr := params.Get("aid")
	tidStr := params.Get("tag_id")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.Like(c, mid, aid, tid, time.Now()); err != nil {
		log.Error("svr.Like(%d,%d) error(%v)", aid, tid, err)
	}
	c.JSON(nil, err)
}

func hateArcTag(c *bm.Context) {
	var (
		aid int64
		tid int64
		err error
	)
	params := c.Request.Form
	mid, _ := c.Get("mid")
	aidStr := params.Get("aid")
	tidStr := params.Get("tag_id")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.Hate(c, mid.(int64), aid, tid, time.Now()); err != nil {
		log.Error("svr.Hate(%d,%d) error(%v)", aid, tid, err)
	}
	c.JSON(nil, err)
}
