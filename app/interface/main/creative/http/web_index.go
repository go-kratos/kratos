package http

import (
	"context"
	"strconv"

	"go-common/app/interface/main/creative/model/order"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func webIndexStat(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	// check user
	midI, ok := c.Get("mid")
	tmidStr := params.Get("tmid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	stat, err := dataSvc.NewStat(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	rfd, _ := dataSvc.RelationFansDay(c, mid)
	log.Info("dataSvc.NewStat(%+v) mid(%d)", stat, mid)
	c.JSON(map[string]interface{}{
		"total_click":       stat.Play,
		"total_dm":          stat.Dm,
		"total_reply":       stat.Comment,
		"total_fans":        stat.Fan,
		"total_fav":         stat.Fav,
		"total_like":        stat.Like,
		"total_share":       stat.Share,
		"incr_click":        stat.Play - stat.PlayLast,
		"incr_dm":           stat.Dm - stat.DmLast,
		"incr_reply":        stat.Comment - stat.CommentLast,
		"incr_fans":         stat.Fan - stat.FanLast,
		"inc_fav":           stat.Fav - stat.FavLast,
		"inc_like":          stat.Like - stat.LikeLast,
		"inc_share":         stat.Share - stat.ShareLast,
		"fan_recent_thirty": rfd,
	}, nil)
}

func webIndexTool(c *bm.Context) {
	params := c.Request.Form
	tyStr := params.Get("type")
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	// check params
	ty, err := strconv.ParseInt(tyStr, 10, 8)
	if err != nil || ty < 0 || ty > 2 {
		ty = 0
	}
	tool, err := operSvc.Tool(c, int8(ty))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tool, nil)
}

func webIndexFull(c *bm.Context) {
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	opers, err := operSvc.WebOperations(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(opers, nil)
}

func webIndexNotify(c *bm.Context) {
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	opers, err := operSvc.WebOperations(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if length := len(opers["creative"]); length > 0 {
		c.JSON(opers["creative"][length-1], nil)
	} else {
		c.JSON(nil, nil)
	}
}

func webIndexOper(c *bm.Context) {
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(operSvc.WebRelOperCache, nil)
}

func webIndexVersion(c *bm.Context) {
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(vsSvc.VersionMapCache, nil)
}

func webWhite(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	tmidStr := params.Get("tmid")
	var (
		err                 error
		tmid                int64
		orderAllow, orderUp int
		cmAd                *order.UpValidate
		g                   = &errgroup.Group{}
		ctx                 = context.TODO()
		task                int8
		ugcpay, staff       bool
		staffCount          int
	)
	mid, _ := midI.(int64)
	if tmidStr != "" {
		if tmid, err = strconv.ParseInt(tmidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tmidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if tmid > 0 && dataSvc.IsWhite(mid) {
			mid = tmid
		}
	}
	g.Go(func() error { //商单
		if ok = arcSvc.AllowOrderUps(mid); ok {
			orderAllow = 1
		}
		return nil
	})
	g.Go(func() error { //商单广告
		cmAd, _ = arcSvc.UpValidate(ctx, mid, ip)
		return nil
	})
	g.Go(func() error { //任务系统 0-关闭 1-开启
		task = whiteSvc.TaskWhiteList(mid)
		return nil
	})
	g.Go(func() error { //ugc付费白名单校验
		ugcpay, _ = paySvc.White(ctx, mid)
		return nil
	})
	g.Go(func() error { //联合投稿白名单
		staff, _ = upSvc.ShowStaff(ctx, mid)
		return nil
	})
	g.Go(func() error { //联合投稿staff合作按钮
		staffCount, _ = arcSvc.StaffValidate(ctx, mid)
		return nil
	})
	g.Wait()
	c.JSON(map[string]interface{}{
		"order": map[string]interface{}{
			"allow": orderAllow,
			"up":    orderUp,
		},
		"growup":      nil,
		"cm_ad":       cmAd,
		"task":        task,
		"ugcpay":      ugcpay,
		"staff":       staff,
		"staff_count": staffCount,
	}, nil)
}

func webIndexNewcomer(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	res, err := newcomerSvc.IndexNewcomer(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
