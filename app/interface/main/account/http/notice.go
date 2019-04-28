package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func notice(c *bm.Context) {
	c.JSON(struct{}{}, nil)
}

func closeNotice(c *bm.Context) {
	c.JSON(struct{}{}, nil)
}

func noticeV2(c *bm.Context) {
	var (
		params = c.Request.Form
		pf     string
		build  int64
	)
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uuid := params.Get("uuid")
	if uuid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pf = params.Get("platform")
	build, _ = strconv.ParseInt(params.Get("build"), 10, 64)
	n, err := memberSvc.NoticeV2(c, mid.(int64), uuid, pf, build)
	if err != nil {
		log.Error("memberSvc.NoticeV2(%d, %s) error(%v)", mid, uuid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(n, nil)
}

func closeNoticeV2(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	uuid := c.Request.Form.Get("uuid")
	if uuid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := memberSvc.CloseNoticeV2(c, mid.(int64), uuid)
	if err != nil {
		log.Error("memberSvc.CloseNoticeV2(%d, %s) error(%v)", mid, uuid, err)
		c.JSON(nil, err)
	}
	c.JSON(nil, nil)
}
