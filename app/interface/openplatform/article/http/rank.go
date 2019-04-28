package http

import (
	"strconv"

	artmdl "go-common/app/interface/openplatform/article/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func rankCategories(c *bm.Context) {
	c.JSON(artSrv.RankCategories(c), nil)
}

func ranks(c *bm.Context) {
	var (
		request  = c.Request
		params   = request.Form
		cid, mid int64
	)
	cid, _ = strconv.ParseInt(params.Get("cid"), 10, 64)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data, note, err := artSrv.Ranks(c, cid, mid, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if data == nil {
		data = []*artmdl.RankMeta{}
	}
	res := make(map[string]interface{})
	res["data"] = data
	res["note"] = note
	res["message"] = "ok"
	c.JSONMap(res, nil)
}
