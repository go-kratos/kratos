package http

import (
	"strconv"

	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func pageList(c *bm.Context) {
	var (
		aid   int64
		err   error
		pages []*api.Page
	)
	aidStr := c.Request.Form.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pages, err = playSvr.PageList(c, aid); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(pages) == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(pages, nil)
}

func videoShot(c *bm.Context) {
	v := new(struct {
		Aid   int64 `form:"aid" validate:"min=1"`
		Cid   int64 `form:"cid"`
		Index bool  `form:"index"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(playSvr.VideoShot(c, v.Aid, v.Cid, v.Index))
}

func playURLToken(c *bm.Context) {
	var (
		aid, cid, mid int64
		err           error
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cid, _ = strconv.ParseInt(params.Get("cid"), 10, 64)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(playSvr.PlayURLToken(c, mid, aid, cid))
}
