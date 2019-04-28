package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func like(c *bm.Context) {
	var (
		id, mid, likeType  int64
		params             = c.Request.Form
		idStr, likeTypeStr string
	)
	idStr = params.Get("id")
	likeTypeStr = params.Get("type")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	likeType, _ = strconv.ParseInt(likeTypeStr, 10, 64)
	if (id <= 0) || (likeType == 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get mid
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	meta, err := artSrv.ArticleMeta(c, id)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if meta == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	err = artSrv.Like(c, mid, id, int(likeType))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	likeStr := "like"
	if likeType == 2 {
		likeStr = "like_cancel"
	}
	artSrv.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(meta.Author.Mid, 10), "", strconv.FormatInt(mid, 10), idStr, "article", likeStr, "")
	c.JSON(nil, nil)
}
