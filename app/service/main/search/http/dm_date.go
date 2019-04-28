package http

import (
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func dmDate(c *bm.Context) {
	var (
		err error
		sp  = &model.DmDateParams{
			Bsp: &model.BasicSearchParams{},
		}
		res *model.SearchResult
	)
	if err = c.Bind(sp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = c.Bind(sp.Bsp); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err = svr.DmDate(c, sp)
	if err != nil {
		log.Error("srv.DmDate(%v) error(%v)", sp, err)
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(res, err)
}
