package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// arcView archive with page info.
func arcView(c *bm.Context) {
	params := c.Request.Form
	// check params
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.View3(c, aid))
}

// arcViews archives with page info by aids
func arcViews(c *bm.Context) {
	params := c.Request.Form
	aidsStr := params.Get("aids")
	// check params
	aids, err := xstr.SplitInts(aidsStr)
	if err != nil {
		log.Error("query aids(%s) split error(%v)", aidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aids) > 20 {
		log.Error("query aids(%s) too long, appkey(%s)", aidsStr, params.Get("appkey"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.Views3(c, aids))
}

// arcPage get pages by aid
func arcPage(c *bm.Context) {
	params := c.Request.Form
	// check params
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.Page3(c, aid))
}

// video get video by aid & cid.
func video(c *bm.Context) {
	params := c.Request.Form
	// check params
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cidStr := params.Get("cid")
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.Video3(c, aid, cid))
}

// description get description by aid & cid.
func description(c *bm.Context) {
	params := c.Request.Form
	// check params
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.Description(c, aid))
}
