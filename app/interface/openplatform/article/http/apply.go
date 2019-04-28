package http

import (
	"go-common/app/interface/openplatform/article/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func applyInfo(c *bm.Context) {
	var (
		mid int64
	)
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(artSrv.ApplyInfo(c, mid))
}

func apply(c *bm.Context) {
	var (
		mid               int64
		request           = c.Request
		params            = request.Form
		content, category string
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	content = params.Get("content")
	if int64(len([]rune(content))) > conf.Conf.Article.MaxApplyContentLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	category = params.Get("category")
	if int64(len([]rune(category))) > conf.Conf.Article.MaxApplyCategoryLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSrv.Apply(c, mid, content, category))
}

func isAuthor(c *bm.Context) {
	var (
		mid, mediaID   int64
		err            error
		author, forbid bool
		id             int64
		level          = true
		canEdit        = true
	)
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	mediaIDStr := c.Request.Form.Get("media_id")
	if mediaIDStr != "" {
		mediaID, err = strconv.ParseInt(mediaIDStr, 10, 64)
		if err != nil || mediaID < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mediaID > 0 {
		if level, err = artSrv.LevelRequired(c, mid); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if id, err = artSrv.MediaArticle(c, mediaID, mid); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if id > 0 && artSrv.EditTimes(c, id) <= 0 {
			canEdit = false
		}
	}
	if author, forbid, err = artSrv.IsAuthor(c, mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{"is_author": author, "forbid": forbid, "level": level, "id": id, "can_edit": canEdit}, nil)
}

func addAuthor(c *bm.Context) {
	var (
		mid int64
	)
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	c.JSON(nil, artSrv.AddAuthor(c, mid))
}
