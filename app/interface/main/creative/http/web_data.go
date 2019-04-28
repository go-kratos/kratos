package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/data"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func webVideoQuitPoints(c *bm.Context) {
	params := c.Request.Form
	cidStr := params.Get("cid")
	tmidStr := params.Get("tmid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	data, err := dataSvc.VideoQuitPoints(c, cid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func webArchive(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := dataSvc.ArchiveStat(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func webArticleData(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := artSvc.ArticleStat(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func base(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// get data
	uv, _ := dataSvc.ViewerBase(c, mid)
	ua, _ := dataSvc.ViewerArea(c, mid)
	if len(uv) == 0 {
		uv = nil
	}
	if len(ua) == 0 {
		ua = nil
	}
	c.JSON(map[string]interface{}{
		"viewer_base": uv,
		"viewer_area": ua,
		"period":      data.Tip(),
	}, nil)
}

func trend(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// get data
	ut, _ := dataSvc.CacheTrend(c, mid)
	if len(ut) == 0 {
		ut = nil
	}
	c.JSON(ut, nil)
}

func action(c *bm.Context) {
	params := c.Request.Form
	month := params.Get("month")
	tmidStr := params.Get("tmid")
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// get data
	rfd, _ := dataSvc.RelationFansDay(c, mid)
	rfdh, _ := dataSvc.RelationFansHistory(c, mid, month)
	rfm, _ := dataSvc.RelationFansMonth(c, mid)
	uah, _ := dataSvc.ViewerActionHour(c, mid)
	if len(rfd) == 0 {
		rfd = nil
	}
	if len(rfdh) == 0 {
		rfdh = nil
	}
	if len(rfm) == 0 {
		rfm = nil
	}
	if len(uah) == 0 {
		uah = nil
	}
	c.JSON(map[string]interface{}{
		"relation_fans_day":     rfd,
		"relation_fans_history": rfdh,
		"relation_fans_month":   rfm,
		"viewer_action_hour":    uah,
	}, nil)
}

func survey(c *bm.Context) {
	params := c.Request.Form
	tyStr := params.Get("type")
	tmidStr := params.Get("tmid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	ty, err := strconv.Atoi(tyStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, ok := data.IncrTy(int8(ty)); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	data, _ := dataSvc.UpIncr(c, mid, int8(ty), ip)
	if len(data) == 0 {
		data = nil
	}
	c.JSON(data, nil)
}

func pandect(c *bm.Context) {
	params := c.Request.Form
	tyStr := params.Get("type")
	tmidStr := params.Get("tmid")
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	ty, err := strconv.Atoi(tyStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, ok := data.IncrTy(int8(ty)); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	data, _ := dataSvc.ThirtyDayArchive(c, mid, int8(ty))
	if len(data) == 0 {
		data = nil
	}
	c.JSON(data, nil)
}

func webFan(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	// check params
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.UpFansAnalysisForWeb(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func webPlaySource(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.UpPlaySourceAnalysis(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func webArcPlayAnalysis(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	cpStr := params.Get("copyright")
	midI, ok := c.Get("mid")
	ip := metadata.String(c, metadata.RemoteIP)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	cp, err := strconv.Atoi(cpStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.UpArcPlayAnalysis(c, mid, cp, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func webArtThirtyDay(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	tyStr := params.Get("type")

	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}

	ty, err := strconv.Atoi(tyStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := data.CheckType(byte(ty)); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.ArtThirtyDay(c, mid, byte(ty))
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func webArtRank(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	tyStr := params.Get("type")

	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}

	ty, err := strconv.Atoi(tyStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ok := data.CheckType(byte(ty)); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.ArtRank(c, mid, byte(ty))
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}

func webArtReadAnalysis(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")

	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}

	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}

	data, err := dataSvc.ArtReadAnalysis(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(data, nil)
}
