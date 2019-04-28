package http

import (
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func listArticles(c *bm.Context) {
	var (
		id, mid int64
	)
	id, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(artSrv.ListArticles(c, id, mid))
}

func webListArticles(c *bm.Context) {
	var (
		id, mid int64
	)
	id, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(artSrv.WebListArticles(c, id, mid))
}

func listInfo(c *bm.Context) {
	var (
		aid int64
	)
	aid, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.ListInfo(c, aid))
}

func upLists(c *bm.Context) {
	var (
		upMid, sort int64
		err         error
		lists       model.UpLists
	)
	upMid, _ = strconv.ParseInt(c.Request.Form.Get("mid"), 10, 64)
	if upMid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sort, _ = strconv.ParseInt(c.Request.Form.Get("sort"), 10, 64)
	if lists, err = artSrv.UpLists(c, upMid, int8(sort)); err != nil {
		c.JSON(nil, err)
		return
	}
	if lists.Lists == nil {
		lists.Lists = []*model.List{}
	}
	c.JSON(lists, err)
}

func refreshLists(c *bm.Context) {
	var (
		id int64
	)
	id, _ = strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.RefreshList(c, id))
}

func rebuildAllListReadCount(c *bm.Context) {
	artSrv.RebuildAllListReadCount(c)
	c.JSON(nil, nil)
}
