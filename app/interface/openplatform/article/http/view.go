package http

import (
	"strconv"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func view(c *bm.Context) {
	var (
		id     int64
		err    error
		art    *artmdl.Article
		params = c.Request.Form
	)
	idStr := params.Get("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if art, err = artSrv.Article(c, id); err != nil {
		c.JSON(nil, err)
		return
	}
	if art == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(art, err)
}

func addView(c *bm.Context) {
	var (
		mid    int64
		params = c.Request.Form
		page   = params.Get("page")
		from   = params.Get("from")
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	if page == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if from == "" {
		from = "unknow"
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	plat := artmdl.Plat(mobiApp, device)
	build := params.Get("build")
	buvid := buvid(c)
	// for tianma mainCard -> 7
	if from == "mainCard" {
		from = "7"
	}
	ua := c.Request.Header.Get("User-Agent")
	referer := c.Request.Header.Get("Referer")
	artSrv.ShowInfoc(ip, time.Now(), buvid, mid, plat, page, from, build, ua, referer)
	c.JSON(nil, nil)
}

func viewInfo(c *bm.Context) {
	var (
		id      int64
		mid     int64
		data    *artmdl.ViewInfo
		request = c.Request
		params  = request.Form
		ip      = metadata.String(c, metadata.RemoteIP)
		err     error
	)
	idStr := params.Get("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	cheat := cheatInfo(c, mid, id)
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	plat := artmdl.Plat(mobiApp, device)
	from := params.Get("from")
	if data, err = artSrv.ViewInfo(c, mid, id, ip, cheat, plat, from); err != nil {
		c.JSON(nil, err)
		return
	}
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	buvid := buvid(c)
	if from == "articleSlideShow" {
		data.Pre, data.Next = artSrv.ViewList(c, id, buvid, "articleSlide", ip, build, plat, mid)
	} else {
		data.Pre, data.Next = artSrv.ViewList(c, id, buvid, from, ip, build, plat, mid)
	}
	// for tianma mainCard -> 7
	if from == "mainCard" {
		from = "7"
	}
	if from != "articleSlide" {
		ua := c.Request.Header.Get("User-Agent")
		artSrv.ViewInfoc(mid, plat, build, "doc", from, buvid, id, time.Now(), ua)
		artSrv.AIViewInfoc(mid, plat, build, "doc", from, buvid, id, time.Now(), ua)
	}
	c.JSON(data, nil)
}

func list(c *bm.Context) {
	var (
		mid     int64
		request = c.Request
		params  = request.Form
		pn, ps  int64
	)
	midStr := params.Get("mid")
	mid, _ = strconv.ParseInt(midStr, 10, 64)
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pnStr := params.Get("pn")
	pn, _ = strconv.ParseInt(pnStr, 10, 64)
	if pn <= 0 {
		pn = 1
	}
	psStr := params.Get("ps")
	ps, _ = strconv.ParseInt(psStr, 10, 64)
	if ps <= 0 {
		ps = 20
	} else if ps > conf.Conf.Article.MaxUpperListPsSize {
		ps = conf.Conf.Article.MaxUpperListPsSize
	}
	c.JSON(artSrv.UpArtMetasAndLists(c, mid, int(pn), int(ps), artmdl.FieldDefault))
}

func earlyArticles(c *bm.Context) {
	var (
		err    error
		aid    int64
		params = c.Request.Form
	)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(artSrv.MoreArts(c, aid))
}

func moreArts(c *bm.Context) {
	var (
		err      error
		aid, mid int64
		params   = c.Request.Form
	)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(artSrv.Mores(c, aid, mid))
}

func cheatInfo(c *bm.Context, mid, id int64) (res *artmdl.CheatInfo) {
	req := c.Request
	params := req.Form
	res = &artmdl.CheatInfo{
		Mid:   strconv.FormatInt(mid, 10),
		Cvid:  strconv.FormatInt(id, 10),
		Refer: req.Header.Get("Referer"),
		UA:    req.Header.Get("User-Agent"),
		Ts:    strconv.FormatInt(time.Now().Unix(), 10),
		IP:    metadata.String(c, metadata.RemoteIP),
	}
	if csid, err := req.Cookie("sid"); err == nil {
		res.Sid = csid.Value
	}
	if params.Get("access_key") == "" {
		res.Client = infoc.ClientWeb
		if ck, err := req.Cookie("buvid3"); err == nil {
			res.Buvid = ck.Value
		}
	} else {
		if params.Get("platform") == "ios" {
			if params.Get("device") == "pad" {
				res.Client = infoc.ClientIpad
			} else {
				res.Client = infoc.ClientIphone
			}
		} else if params.Get("platform") == "android" {
			res.Client = infoc.ClientAndroid
		}
		res.Buvid = req.Header.Get("buvid")
		res.Build = params.Get("build")
	}
	return
}

func actInfo(c *bm.Context) {
	var (
		request = c.Request
		params  = request.Form
	)
	mobiApp := params.Get("mobi_app")
	plat := artmdl.Plat(mobiApp, "")
	c.JSON(artSrv.ActInfo(c, plat))
}
