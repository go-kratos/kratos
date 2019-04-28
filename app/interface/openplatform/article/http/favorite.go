package http

import (
	"strconv"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func addFavorite(c *bm.Context) {
	var (
		err           error
		mid, aid, fid int64
		params        = c.Request.Form
		ip            = metadata.String(c, metadata.RemoteIP)
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	aidStr := params.Get("id")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fidStr := params.Get("fid")
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	meta, err := artSrv.ArticleMeta(c, aid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if meta == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	err = artSrv.AddFavorite(c, mid, aid, fid, ip)
	artSrv.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(meta.Author.Mid, 10), "", strconv.FormatInt(mid, 10), aidStr, "article", infoc.ActionFav, fidStr)
	c.JSON(nil, err)
}

func delFavorite(c *bm.Context) {
	var (
		err           error
		mid, aid, fid int64
		params        = c.Request.Form
		ip            = metadata.String(c, metadata.RemoteIP)
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	aidStr := params.Get("id")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fidStr := params.Get("fid")
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if err = artSrv.DelFavorite(c, mid, aid, fid, ip); err != nil {
		dao.PromError("删除收藏")
		log.Error("artSrv.DelFavorite(%d,%d,%d) error(%+v)", mid, aid, fid, err)
	}
	c.JSON(nil, err)
}

func favorites(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		pn, ps   int
		favs     []*artmdl.Favorite
		page     *artmdl.Page
		params   = c.Request.Form
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	fidStr := params.Get("fid")
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pnStr := params.Get("pn")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	psStr := params.Get("ps")
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Article.CreationDefaultSize
	} else if ps > conf.Conf.Article.CreationMaxSize {
		ps = conf.Conf.Article.CreationMaxSize
	}
	if favs, page, err = artSrv.ValidFavs(c, mid, fid, pn, ps, ip); err != nil {
		c.JSON(nil, err)
		return
	}
	type data struct {
		Favorites []*artmdl.Favorite `json:"favorites"`
		Page      *artmdl.Page       `json:"page"`
	}
	c.JSON(&data{
		Favorites: favs,
		Page:      page,
	}, nil)
}

func allFavorites(c *bm.Context) {
	var (
		err      error
		mid, fid int64
		pn, ps   int
		favs     []*artmdl.Favorite
		page     *artmdl.Page
		params   = c.Request.Form
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	fidStr := params.Get("fid")
	if fidStr != "" {
		if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil || fid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pnStr := params.Get("pn")
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	psStr := params.Get("ps")
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = conf.Conf.Article.CreationDefaultSize
	} else if ps > conf.Conf.Article.CreationMaxSize {
		ps = conf.Conf.Article.CreationMaxSize
	}
	if favs, page, err = artSrv.Favs(c, mid, fid, pn, ps, ip); err != nil {
		c.JSON(nil, err)
		return
	}
	type data struct {
		Favorites []*artmdl.Favorite `json:"favorites"`
		Page      *artmdl.Page       `json:"page"`
	}
	c.JSON(&data{
		Favorites: favs,
		Page:      page,
	}, nil)
}
