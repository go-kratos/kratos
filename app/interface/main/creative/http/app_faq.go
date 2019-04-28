package http

import (
	"context"
	"go-common/app/interface/main/creative/model/faq"
	mMdl "go-common/app/interface/main/creative/model/music"
	resMdl "go-common/app/interface/main/creative/model/resource"
	resmdl "go-common/app/service/main/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"
	"strconv"
)

func appH5FaqEditor(c *bm.Context) {
	var (
		total         int
		items         []*faq.Detail
		err           error
		pn            = 1
		ps            = 20
		keyFlag       = 1
		faqQuesTypeID = faq.PhoneFaqQuesTypeID
	)
	params := c.Request.Form
	device := params.Get("device")
	if device == "ipad" || device == "pad" {
		log.Warn("openpad faqSvc.Detail (%s,%s,%s)", faqQuesTypeID, device, faq.PadFaqQuesTypeID)
		faqQuesTypeID = faq.PadFaqQuesTypeID
	}
	if items, total, err = faqSvc.Detail(c, faqQuesTypeID, keyFlag, pn, ps); err != nil {
		log.Error("faqSvc.Detail(%s,%d,%d) error(%v)", faqQuesTypeID, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	detail := map[string]interface{}{
		"items": items,
		"page": map[string]int{
			"num":   pn,
			"size":  ps,
			"total": total,
		},
	}
	c.JSON(map[string]interface{}{
		"detail": detail,
	}, nil)
}

func appCooperate(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	idStr := params.Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	c.JSONMap(map[string]interface{}{
		"data": musicSvc.Cooperate(c, id, mid),
	}, nil)
}

func appCooperatePre(c *bm.Context) {
	var (
		banners []*resmdl.Assignment
		coos    []*mMdl.Cooperate
		build   int
		err     error
		g       = &errgroup.Group{}
		ctx     = context.TODO()
	)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	header := c.Request.Header
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platStr := params.Get("platform")
	network := params.Get("network")
	buvid := header.Get("Buvid")
	adExtra := params.Get("ad_extra")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := resMdl.Plat(mobiApp, device)
	g.Go(func() error {
		coos = musicSvc.CooperatePre(ctx, mid, platStr, build)
		return nil
	})
	g.Go(func() error {
		banners, _ = resSvc.CooperateBanner(ctx, mobiApp, device, network, buvid, adExtra, build, plat, mid, false)
		return nil
	})
	g.Wait()
	c.JSON(map[string]interface{}{
		"coos":    coos,
		"banners": banners,
	}, nil)
}
