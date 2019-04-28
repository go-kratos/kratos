package http

import (
	mdlSearch "go-common/app/interface/main/tv/model/search"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid = "Buvid"
	_keyWordLen  = 50
)

// searchAll all search .
func searchAll(c *bm.Context) {
	var (
		err    error
		v      = new(mdlSearch.UserSearch)
		header = c.Request.Header
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Keyword == "" || len([]rune(v.Keyword)) > _keyWordLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Order == "" {
		v.Order = "totalrank"
	}
	if v.Page < 1 {
		v.Page = 1
	}
	if v.Pagesize < 1 || v.Pagesize > 20 {
		v.Pagesize = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		v.MID = midInter.(int64)
	}
	v.Buvid = header.Get(_headerBuvid)
	c.JSON(secSvc.SearchAll(c, v))
}

// userSearch search by user .
func userSearch(c *bm.Context) {
	var (
		err    error
		v      = new(mdlSearch.UserSearch)
		header = c.Request.Header
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Order == "" {
		v.Order = "totalrank"
	}
	if v.Order != "totalrank" && v.Order != "fans" && v.Order != "level" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.OrderSort != 1 {
		v.OrderSort = 0
	}
	if v.FromSource == "" {
		v.FromSource = "app_search"
	}
	if v.Page < 1 {
		v.Page = 1
	}
	if v.Pagesize < 1 || v.Pagesize > 20 {
		v.Pagesize = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		v.MID = midInter.(int64)
	}
	v.Buvid = header.Get(_headerBuvid)
	c.JSON(secSvc.UserSearch(c, v))
}

// pgcSearch search pgc opera and film .
func pgcSearch(c *bm.Context) {
	var (
		err    error
		v      = new(mdlSearch.UserSearch)
		header = c.Request.Header
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Keyword == "" || len([]rune(v.Keyword)) > _keyWordLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Page < 1 {
		v.Page = 1
	}
	if v.Pagesize < 1 || v.Pagesize > 20 {
		v.Pagesize = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		v.MID = midInter.(int64)
	}
	v.Buvid = header.Get(_headerBuvid)
	c.JSON(secSvc.PgcSearch(c, v))
}
