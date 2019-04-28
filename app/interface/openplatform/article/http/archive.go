package http

import (
	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func archives(c *bm.Context) {
	var (
		err    error
		aids   []int64
		arcs   map[int64]*api.Arc
		params = c.Request.Form
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	idsStr := params.Get("ids")
	if aids, err = xstr.SplitInts(idsStr); err != nil || len(aids) < 1 || len(aids) > conf.Conf.Article.MaxArchives {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if arcs, err = artSrv.Archives(c, aids, ip); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(arcs) == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(arcs, err)
}
