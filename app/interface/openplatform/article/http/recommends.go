package http

import (
	"strconv"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const (
	_headerBuvid     = "Buvid"
	_buvid           = "buvid3"
	_recommendRegion = 0
	_rankPage        = 9
)

func recommends(c *bm.Context) {
	var (
		mid    int64
		rs     []*model.RecommendArtWithLike
		err    error
		params = c.Request.Form
		aids   []int64
		sky    *model.SkyHorseResp
	)
	if aids, _ = xstr.SplitInts(params.Get("aids")); len(aids) == 0 {
		aids, _ = xstr.SplitInts(params.Get("adis")) //兼容ios客户端bug
	}
	cid, _ := strconv.ParseInt(params.Get("cid"), 10, 64)
	sort, _ := strconv.Atoi(params.Get("sort"))
	pn, _ := strconv.ParseInt(params.Get("pn"), 10, 64)
	ps, _ := strconv.ParseInt(params.Get("ps"), 10, 64)
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
		ps = conf.Conf.Article.MaxRecommendPsSize
	}
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	plat := model.Plat(mobiApp, device)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	pageType, _ := strconv.Atoi(params.Get("from"))
	if pageType == 0 {
		// 分区页
		pageType = 2
	}
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	buvid := buvid(c)
	if rs, sky, err = artSrv.SkyHorse(c, cid, int(pn), int(ps), aids, sort, mid, build, buvid, plat); err != nil {
		dao.PromError("推荐接口")
		log.Error("service.Recommends(%d) error(%+v)", mid, err)
		c.JSON(nil, err)
		return
	}
	var as []*model.Meta
	for _, r := range rs {
		as = append(as, &r.Meta)
	}
	artSrv.RecommendInfoc(mid, plat, pageType, int(cid), build, buvid, metadata.String(c, metadata.RemoteIP), as, false, time.Now(), pn, sky)
	if rs == nil {
		rs = []*model.RecommendArtWithLike{}
	}
	res := make(map[string]interface{})
	res["aids_len"] = conf.Conf.Article.RecommendAidLen
	res["data"] = rs
	c.JSONMap(res, nil)
}

func home(c *bm.Context) {
	var (
		mid    int64
		rs     *model.RecommendHome
		err    error
		params = c.Request.Form
		ip     = metadata.String(c, metadata.RemoteIP)
		aids   []int64
		sky    *model.SkyHorseResp
	)
	aids, _ = xstr.SplitInts(params.Get("aids"))
	pn, _ := strconv.ParseInt(params.Get("pn"), 10, 64)
	ps, _ := strconv.ParseInt(params.Get("ps"), 10, 64)
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
		ps = conf.Conf.Article.MaxRecommendPsSize
	}
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	plat := model.Plat(mobiApp, device)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	pageType, _ := strconv.Atoi(params.Get("from"))
	if pageType == 0 {
		// 首页tab
		pageType = 8
	}
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	buvid := buvid(c)
	if rs, sky, err = artSrv.RecommendHome(c, int8(plat), build, int(pn), int(ps), aids, mid, ip, time.Now(), buvid); err != nil {
		dao.PromError("推荐接口")
		log.Error("service.Recommends(%d) error(%+v)", mid, err)
		c.JSON(nil, ecode.Degrade)
		return
	}
	var as []*model.Meta
	for _, r := range rs.Articles {
		as = append(as, &r.Meta)
	}
	artSrv.RecommendInfoc(mid, plat, pageType, _recommendRegion, build, buvid, metadata.String(c, metadata.RemoteIP), as, false, time.Now(), pn, sky)
	if len(rs.Ranks) > 0 && pn == 1 {
		var as []*model.Meta
		for _, r := range rs.Ranks {
			as = append(as, r.Meta)
		}
		artSrv.RecommendInfoc(mid, plat, _rankPage, 0, build, buvid, metadata.String(c, metadata.RemoteIP), as, false, time.Now(), pn, nil)
	}
	res := make(map[string]interface{})
	res["aids_len"] = conf.Conf.Article.RecommendAidLen
	res["data"] = rs
	c.JSONMap(res, nil)
}

func recommendsPlus(c *bm.Context) {
	var (
		mid    int64
		rs     *model.RecommendPlus
		err    error
		params = c.Request.Form
		aids   []int64
		sky    *model.SkyHorseResp
	)
	if aids, _ = xstr.SplitInts(params.Get("aids")); len(aids) == 0 {
		aids, _ = xstr.SplitInts(params.Get("adis")) //兼容ios客户端bug
	}
	cid, _ := strconv.ParseInt(params.Get("cid"), 10, 64)
	sort, _ := strconv.Atoi(params.Get("sort"))
	pn, _ := strconv.ParseInt(params.Get("pn"), 10, 64)
	ps, _ := strconv.ParseInt(params.Get("ps"), 10, 64)
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
		ps = conf.Conf.Article.MaxRecommendPsSize
	}
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	plat := model.Plat(mobiApp, device)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	pageType, _ := strconv.Atoi(params.Get("from"))
	if pageType == 0 {
		// 分区页
		pageType = 2
	}
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	buvid := buvid(c)
	if rs, sky, err = artSrv.RecommendPlus(c, cid, int8(plat), build, int(pn), int(ps), aids, mid, time.Now(), sort, buvid); err != nil {
		c.JSON(nil, err)
		dao.PromError("推荐接口")
		log.Error("service.Recommends(%d) error(%+v)", mid, err)
		return
	}
	var as []*model.Meta
	for _, r := range rs.Articles {
		as = append(as, &r.Meta)
	}
	artSrv.RecommendInfoc(mid, plat, pageType, int(cid), build, buvid, metadata.String(c, metadata.RemoteIP), as, false, time.Now(), pn, sky)
	if len(rs.Ranks) > 0 && pn == 1 {
		var as []*model.Meta
		for _, r := range rs.Ranks {
			as = append(as, r.Meta)
		}
		artSrv.RecommendInfoc(mid, plat, _rankPage, 0, build, buvid, metadata.String(c, metadata.RemoteIP), as, false, time.Now(), pn, nil)
	}
	if rs.Articles == nil {
		rs.Articles = []*model.RecommendArtWithLike{}
	}
	res := make(map[string]interface{})
	res["aids_len"] = conf.Conf.Article.RecommendAidLen
	res["data"] = rs
	c.JSONMap(res, nil)
}

func allRecommends(c *bm.Context) {
	v := new(struct {
		Pn int `form:"pn" validate:"min=1"`
		Ps int `form:"ps" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	count, recommends, err := artSrv.AllRecommends(c, v.Pn, v.Ps)
	c.JSON(map[string]interface{}{
		"total": count,
		"list":  recommends,
	}, err)
}
