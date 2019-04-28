package http

import (
	"context"
	"strconv"
	"time"

	accmdl "go-common/app/interface/main/creative/model/account"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/faq"
	"go-common/app/interface/main/creative/model/watermark"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func appArcDescFormat(c *bm.Context) {
	params := c.Request.Form
	ip := metadata.String(c, metadata.RemoteIP)
	typeidStr := params.Get("typeid")
	cpStr := params.Get("copyright")
	lang := params.Get("lang")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	typeid, err := strconv.ParseInt(typeidStr, 10, 16)
	if typeid < 0 || err != nil {
		typeid = 0
	}
	copyright, err := strconv.ParseInt(cpStr, 10, 16)
	if copyright <= 0 || err != nil {
		copyright = archive.CopyrightReprint
	}
	desc, length, err := arcSvc.DescFormatForApp(c, typeid, copyright, lang, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"desc_length": length,
		"desc_format": desc,
	}, nil)
}

func appArchivePre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	plat := params.Get("platform")
	lang := params.Get("lang")
	if lang != "en" {
		lang = "ch"
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err          error
		mf           *accmdl.MyInfo
		wm           *watermark.Watermark
		g            = &errgroup.Group{}
		ctx          = context.TODO()
		faqs         = make(map[string]*faq.Faq)
		recFriends   []*accmdl.Friend
		lotteryCheck bool
	)
	g.Go(func() error {
		mf, err = accSvc.MyInfo(ctx, mid, ip, time.Now())
		if err != nil {
			log.Info("accSvc.MyInfo (%d) err(%v)", mid, err)
		}
		if mf != nil {
			mf.Commercial = arcSvc.AllowCommercial(c, mid)
		}
		return nil
	})
	g.Go(func() error {
		wm, err = wmSvc.WaterMark(ctx, mid)
		if err != nil {
			log.Info("wmSvc.WaterMark (%d) err(%+v) WaterMark(%+v)", mid, err, wm)
		}
		if len(wm.URL) == 0 {
			wm.State = 1
		}
		return nil
	})
	g.Go(func() error {
		faqs = faqSvc.Pre(ctx)
		return nil
	})
	g.Go(func() error {
		lotteryCheck, _ = dymcSvc.LotteryUserCheck(ctx, mid)
		return nil
	})
	g.Go(func() error {
		recFriends, _ = accSvc.RecFollows(ctx, mid)
		return nil
	})
	g.Wait()
	uploadinfo, _ := whiteSvc.UploadInfoForMainApp(mf, plat, mid)
	mf.DymcLottery = lotteryCheck
	c.JSON(map[string]interface{}{
		"uploadinfo": uploadinfo,
		"typelist":   arcSvc.AppTypes(c, lang),
		"myinfo":     mf,
		"arctip":     arcSvc.ArcTip,
		"activities": arcSvc.Activities(c),
		"watermark":  wm,
		"fav":        arcSvc.Fav(c, mid),
		"tip":        vsSvc.AppManagerTip,
		"cus_tip":    vsSvc.CusManagerTip,
		// common data
		"camera_cfg":  appSvc.CameraCfg,
		"module_show": arcSvc.AppModuleShowMap(mid, lotteryCheck),
		"icons":       appSvc.Icons(),
		"faqs":        faqs,
		"rec_friends": recFriends,
	}, nil)
}
