package http

import (
	"strconv"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func hotspotArts(c *bm.Context) {
	var (
		mid    int64
		rs     []*model.MetaWithLike
		hot    *model.Hotspot
		err    error
		params = c.Request.Form
		aids   []int64
	)
	cid, _ := strconv.ParseInt(params.Get("id"), 10, 64)
	sort, _ := strconv.Atoi(params.Get("sort"))
	pn, _ := strconv.ParseInt(params.Get("pn"), 10, 64)
	ps, _ := strconv.ParseInt(params.Get("ps"), 10, 64)
	if cid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn <= 0 {
		pn = 1
	}
	if pn > conf.Conf.Article.MaxRecommendPnSize {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps <= 0 {
		ps = 20
	} else if ps > conf.Conf.Article.MaxRecommendPsSize {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if hot, rs, err = artSrv.HotspotArts(c, cid, int(pn), int(ps), aids, int8(sort), mid); err != nil {
		c.JSON(nil, err)
		dao.PromError("热点运营接口")
		log.Error("service.HotspotArts(%d) error(%+v)", mid, err)
		return
	}
	res := make(map[string]interface{})
	res["aids_len"] = conf.Conf.Article.RecommendAidLen
	if rs == nil {
		rs = []*model.MetaWithLike{}
	}
	res["data"] = &model.HotspotResp{
		Hotspot:  hot,
		Articles: rs,
	}
	c.JSONMap(res, nil)
}
