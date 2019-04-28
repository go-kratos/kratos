package http

import (
	"context"
	"strconv"

	"go-common/app/interface/main/creative/model/data"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func converTmid(c *bm.Context, mid int64) (retMid int64) {
	tmidStr := c.Request.Form.Get("tmid")
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		retMid = tmid
	} else {
		retMid = mid
	}
	return
}

func appDataArc(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	arcStat, err := dataSvc.ArchiveStat(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	_, vds, _ := arcSvc.Videos(c, mid, aid, ip)
	arcStat.Videos = vds
	c.JSON(arcStat, nil)
}

func appDataVideoQuit(c *bm.Context) {
	params := c.Request.Form
	cidStr := params.Get("cid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	// check params
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	pts, err := dataSvc.AppVideoQuitPoints(c, cid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pts, nil)
}

func appFan(c *bm.Context) {
	req := c.Request
	params := req.Form
	tyStr := params.Get("type")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err error
		ty  int
		fan *data.AppFan
	)
	ty, err = strconv.Atoi(tyStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", tyStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	fan, err = dataSvc.UpFansAnalysisForApp(c, mid, ty, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(fan, nil)
}

func appFanRank(c *bm.Context) {
	req := c.Request
	params := req.Form
	tyStr := params.Get("type")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err error
		ty  int
		rk  map[string][]*data.RankInfo
	)
	ty, err = strconv.Atoi(tyStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", tyStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	rk, err = dataSvc.FanRankApp(c, mid, ty, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rk, nil)
}

func appOverView(c *bm.Context) {
	req := c.Request
	params := req.Form
	tyStr := params.Get("type")
	ip := metadata.String(c, metadata.RemoteIP)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err error
		ty  int
		ao  *data.AppOverView
	)
	ty, err = strconv.Atoi(tyStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", tyStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	ao, err = dataSvc.OverView(c, mid, int8(ty), ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ao, nil)
}

func appArchiveAnalyze(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	arcStat, err := dataSvc.ArchiveAnalyze(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	_, vds, _ := arcSvc.Videos(c, mid, aid, ip)
	arcStat.Videos = vds
	c.JSON(arcStat, nil)
}

func appVideoRetention(c *bm.Context) {
	params := c.Request.Form
	cidStr := params.Get("cid")
	ip := metadata.String(c, metadata.RemoteIP)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid = converTmid(c, mid)
	pts, err := dataSvc.VideoRetention(c, cid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pts, nil)
}

func appDataArticle(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		artStat artmdl.UpStat
		artIncr []*artmdl.ThirtyDayArticle
		g       = &errgroup.Group{}
		ctx     = context.TODO()
	)
	mid = converTmid(c, mid)
	g.Go(func() error {
		artStat, _ = artSvc.ArticleStat(ctx, mid, ip)
		return nil
	})
	g.Go(func() error {
		artIncr, _ = dataSvc.ThirtyDayArticle(ctx, mid, ip)
		return nil
	})
	g.Wait()
	c.JSON(map[string]interface{}{
		"stat": artStat,
		"incr": artIncr,
	}, nil)
}
