package http

import (
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func replySearch(c *bm.Context) {
	params := c.Request.Form
	appidStr := params.Get("appid")
	switch appidStr {
	case "reply_record":
		replyRecord(c)
	default:
		c.JSON(nil, ecode.RequestErr)
		return
	}
}

func replyRecord(c *bm.Context) {
	var (
		err error
		sp  = &model.ReplyRecordParams{
			Bsp: &model.BasicSearchParams{},
		}
		res    *model.SearchResult
		params = c.Request.Form
	)
	if params.Get("mid") == "" {
		log.Error("mid is required")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err = svr.ReplyRecord(c, sp)
	if err != nil {
		log.Error("svr.ArchiveCheck(%v,%d,%d) error(%v)", sp, sp.Bsp.Pn, sp.Bsp.Ps, err)
	}
	c.JSON(res, err)
}
