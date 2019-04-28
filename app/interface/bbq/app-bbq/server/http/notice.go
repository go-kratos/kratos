package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

func noticeNum(c *bm.Context) {
	var mid int64
	midValue, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid = midValue.(int64)
	c.JSON(srv.GetNoticeNum(c, mid))
}

func noticeOverview(c *bm.Context) {
	var mid int64
	midValue, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid = midValue.(int64)
	c.JSON(srv.NoticeOverview(c, mid))
}

func noticeList(c *bm.Context) {
	var mid int64
	midValue, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid = midValue.(int64)

	req := &v1.NoticeListRequest{}
	if err := c.Bind(req); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	req.Mid = mid
	c.JSON(srv.NoticeList(c, req))
}
