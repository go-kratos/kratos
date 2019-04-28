package http

import (
	"strconv"

	"go-common/app/interface/openplatform/article/conf"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func meta(c *bm.Context) {
	var (
		err    error
		aid    int64
		am     *artmdl.Meta
		params = c.Request.Form
	)
	idStr := params.Get("id")
	if aid, err = strconv.ParseInt(idStr, 10, 64); err != nil || aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if am, err = artSrv.ArticleMeta(c, aid); err != nil {
		c.JSON(nil, err)
		return
	} else if am == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(am, nil)
}

func metas(c *bm.Context) {
	var (
		err    error
		aids   []int64
		ams    map[int64]*artmdl.Meta
		params = c.Request.Form
		mid    int64
		resIDs []int64
	)
	idsStr := params.Get("ids")
	midStr := params.Get("mid")
	mid, _ = strconv.ParseInt(midStr, 10, 64)
	if aids, err = xstr.SplitInts(idsStr); err != nil || len(aids) < 1 || len(aids) > conf.Conf.Article.MaxArticleMetas {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ams, err = artSrv.ArticleMetas(c, aids); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(ams) == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if mid > 0 {
		for _, artm := range ams {
			resIDs = append(resIDs, artm.ID)
		}
		likeRes, _ := artSrv.HadLikesByMid(c, mid, resIDs)
		for _, art := range ams {
			isLike := likeRes[art.ID]
			if isLike > 0 {
				art.IsLike = true
			}
		}
	}
	c.JSON(ams, nil)
}

func addCheatFilter(c *bm.Context) {
	var (
		err    error
		aid    int64
		lv     int
		params = c.Request.Form
	)
	idStr := params.Get("id")
	if aid, err = strconv.ParseInt(idStr, 10, 64); err != nil || aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	lv, err = strconv.Atoi(params.Get("lv"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.AddCheatFilter(c, aid, lv))
}

func delCheatFilter(c *bm.Context) {
	var (
		err    error
		aid    int64
		params = c.Request.Form
	)
	idStr := params.Get("id")
	if aid, err = strconv.ParseInt(idStr, 10, 64); err != nil || aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.DelCheatFilter(c, aid))
}
