package http

import (
	"context"
	"strconv"
	"time"

	appMdl "go-common/app/interface/main/creative/model/app"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/article"
	"go-common/app/interface/main/creative/model/danmu"
	"go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/main/creative/model/elec"
	"go-common/app/interface/main/creative/model/operation"
	"go-common/app/interface/main/creative/model/order"
	"go-common/app/interface/main/creative/model/search"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func creatorMy(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mf, err := accSvc.MyInfo(c, mid, ip, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"myinfo":   mf,
		"viewinfo": whiteSvc.Viewinfo(mf),
	}, nil)
}

func creatorIndex(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	// get data
	var (
		stat            *data.Stat
		elecStat        *elec.UserState
		elecBal         *elec.Balance
		elecEarnings    *elec.Earnings
		archives        []*archive.ArcVideoAudit
		banner          []*operation.BannerCreator
		replies         *search.Replies
		articles        []*article.Meta
		artStat         artmdl.UpStat
		dmRecent        *danmu.DmRecent
		creatorDataShow *data.CreatorDataShow
		g               = &errgroup.Group{}
		ctx             = context.TODO()
	)
	g.Go(func() error {
		stat, _ = dataSvc.NewStat(ctx, mid, ip)
		if stat != nil {
			stat.Day30 = nil
			stat.Arcs = nil
		}
		return nil
	})
	g.Go(func() error {
		if arcs, err := arcSvc.WebArchives(ctx, mid, 0, "", "", "is_pubing,pubed,not_pubed", ip, 1, 2, 0); err == nil && arcs.Archives != nil {
			archives = arcs.Archives
		} else {
			archives = make([]*archive.ArcVideoAudit, 0)
		}
		return nil
	})
	g.Go(func() error {
		elecStat, _ = elecSvc.UserState(ctx, mid, ip, ak, ck)
		elecEarnings = &elec.Earnings{}
		if elecStat != nil && elecStat.State == "2" {
			elecEarnings.State = 1
			elecBal, _ = elecSvc.Balance(ctx, mid, ip)
			if elecBal != nil && elecBal.Wallet != nil {
				elecEarnings.Balance = elecBal.Wallet.SponsorBalance //充电数量
			}
			if elecBal != nil && elecBal.BpayAcc != nil {
				elecEarnings.Brokerage = elecBal.BpayAcc.Brokerage //贝壳数量
			}
		}
		return nil
	})
	g.Go(func() error {
		_, banner, _ = operSvc.AppBanner(ctx)
		return nil
	})
	g.Go(func() error {
		if arts, err := artSvc.Articles(ctx, mid, 1, 2, 0, 0, 0, ip); err == nil && arts != nil && len(arts.Articles) != 0 {
			articles = arts.Articles
		} else {
			articles = []*article.Meta{}
		}
		return nil
	})
	g.Go(func() error {
		replies, _ = replySvc.AppIndexReplies(ctx, ak, ck, mid, 0, 0, 0, search.All, 0, "", "", "", ip, 1, 10)
		if replies == nil {
			replies = &search.Replies{}
		}
		return nil
	})
	g.Go(func() error {
		artStat, _ = artSvc.ArticleStat(ctx, mid, ip)
		return nil
	})
	g.Go(func() error {
		dmRecent, _ = danmuSvc.Recent(ctx, mid, 1, 2, ip)
		return nil
	})
	g.Wait()
	creatorDataShow = &data.CreatorDataShow{}
	if len(archives) > 0 {
		creatorDataShow.Archive = 1
	}
	if len(articles) > 0 {
		creatorDataShow.Article = 1
	}
	c.JSON(map[string]interface{}{
		"archives":        archives,
		"archive_stat":    stat,
		"elec_earnings":   elecEarnings,
		"order_earnings":  &order.OasisEarnings{},
		"growth_earnings": &order.GrowthEarnings{},
		"banner":          banner,
		"articles":        articles,
		"replies":         replies,
		"article_stat":    artStat,
		"data_show":       creatorDataShow,
		"danmu":           dmRecent.List,
	}, nil)
}

func creatorArchives(c *bm.Context) {
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	order := params.Get("order")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	class := params.Get("class")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	// check params
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 50 {
		ps = 10
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	arc, err := arcSvc.WebArchives(c, mid, int16(tid), keyword, order, class, ip, pn, ps, 0)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arc, nil)
}

func creatorEarnings(c *bm.Context) {
	params := c.Request.Form
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		elecStat     *elec.UserState
		elecBal      *elec.Balance
		elecEarnings *elec.Earnings
		g            = &errgroup.Group{}
		ctx          = context.TODO()
	)
	g.Go(func() error {
		elecStat, _ = elecSvc.UserState(ctx, mid, ip, ak, ck)
		elecEarnings = &elec.Earnings{}
		if elecStat != nil && elecStat.State == "2" {
			elecEarnings.State = 1
			elecBal, _ = elecSvc.Balance(c, mid, ip)
			if elecBal != nil && elecBal.Wallet != nil {
				elecEarnings.Balance = elecBal.Wallet.SponsorBalance //充电数量
			}
			if elecBal != nil && elecBal.BpayAcc != nil {
				elecEarnings.Brokerage = elecBal.BpayAcc.Brokerage //贝壳数量
			}
		}
		return nil
	})
	g.Wait()
	cw := appMdl.EarningsCopyWriter{
		Elec:   "每月6日结算为贝壳，6-10日可在PC上进行提现",
		Growth: "每月1日结算为贝壳，6-10日可在PC上进行提现",
		Oasis:  "请在PC上进行接单、投稿操作",
	}
	c.JSON(map[string]interface{}{
		"elec_earnings":   elecEarnings,
		"order_earnings":  &order.OasisEarnings{},
		"growth_earnings": &order.GrowthEarnings{},
		"copywriter":      cw,
	}, nil)
}

func creatorViewArc(c *bm.Context) {
	req := c.Request
	params := c.Request.Form
	aidStr := params.Get("aid")
	tidStr := params.Get("typeid")
	title := params.Get("title")
	filename := params.Get("filename")
	desc := params.Get("desc")
	cover := params.Get("cover")
	ip := metadata.String(c, metadata.RemoteIP)
	cookie := req.Header.Get("cookie")
	ak := params.Get("access_key")
	// check user
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
	av, err := arcSvc.View(c, mid, aid, ip, archive.PlatformAndroid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if av == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid < 0 {
		tid = 0
	}
	ptags, _ := dataSvc.TagsWithChecked(c, mid, uint16(tid), title, filename, desc, cover, archive.TagPredictFromAPP)
	elecArc, _ := elecSvc.ArchiveState(c, aid, mid, ip)
	elecStat, _ := elecSvc.UserState(c, mid, ip, ak, cookie)
	if elecArc == nil {
		elecArc = &elec.ArcState{}
	}
	c.JSON(map[string]interface{}{
		"archive":     av.Archive,
		"videos":      av.Videos,
		"predict_tag": ptags,
		"arc_elec":    elecArc,
		"user_elec":   elecStat,
	}, nil)
}

func creatorVideoQuit(c *bm.Context) {
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
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	vq, err := dataSvc.AppVideoQuitPoints(c, cid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(vq, nil)
}

func creatorBanner(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn <= 0 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 0 || ps > 20 {
		ps = 20
	}
	oper, err := operSvc.CreatorOperationList(c, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(oper, nil)
}

func creatorArchiveData(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
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
	arcStat, err := dataSvc.ArchiveStat(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	_, vds, _ := arcSvc.Videos(c, mid, aid, ip)
	arcStat.Videos = vds
	c.JSON(arcStat, nil)
}

func creatorDelArc(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
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
	c.JSON(nil, arcSvc.Del(c, mid, aid, ip))
}

func creatorArcTagInfo(c *bm.Context) {
	params := c.Request.Form
	tagNameStr := params.Get("tag_name")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if len(tagNameStr) == 0 {
		log.Error("tagNameStr len zero (%s)", tagNameStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	code, msg := arcSvc.TagCheck(c, mid, tagNameStr)
	c.JSON(map[string]interface{}{
		"code": code,
		"msg":  msg,
	}, nil)
}

func creatorReplyList(c *bm.Context) {
	req := c.Request
	params := req.Form
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	isReport, err := strconv.Atoi(params.Get("is_report"))
	if err != nil {
		isReport = 0
	}
	tp, err := strconv.Atoi(params.Get("type"))
	if err != nil {
		tp = 1
	}
	oid, err := strconv.ParseInt(params.Get("oid"), 10, 64)
	if err != nil {
		oid = 0
	}
	pn, err := strconv.Atoi(params.Get("pn"))
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(params.Get("ps"))
	if err != nil || ps <= 10 || ps > 100 {
		ps = 10
	}
	p := &search.ReplyParam{
		Ak:          params.Get("access_key"),
		Ck:          c.Request.Header.Get("cookie"),
		OMID:        mid,
		OID:         oid,
		Pn:          pn,
		Ps:          ps,
		IP:          metadata.String(c, metadata.RemoteIP),
		IsReport:    int8(isReport),
		Type:        int8(tp),
		FilterCtime: params.Get("filter"),
		Kw:          params.Get("keyword"),
		Order:       params.Get("order"),
	}
	replies, err := replySvc.Replies(c, p)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"data": replies.Result,
		"pager": map[string]int{
			"pn":    p.Pn,
			"ps":    p.Ps,
			"total": replies.Total,
		},
	}, nil)
}

func creatorPre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	lang := params.Get("lang")
	if lang != "en" {
		lang = "ch"
	}
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mf, err := accSvc.MyInfo(c, mid, ip, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	mf.Commercial = arcSvc.AllowCommercial(c, mid)
	tpl, _ := tplSvc.Templates(c, mid)
	c.JSON(map[string]interface{}{
		"uploadinfo": whiteSvc.UploadInfoForCreator(mf, mid),
		"typelist":   arcSvc.AppTypes(c, lang),
		"myinfo":     mf,
		"template":   tpl,
	}, nil)
}

func creatorPredictTag(c *bm.Context) {
	params := c.Request.Form
	tidStr := params.Get("typeid")
	title := params.Get("title")
	filename := params.Get("filename")
	desc := params.Get("desc")
	cover := params.Get("cover")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid < 0 {
		tid = 0
	}
	ptags, _ := dataSvc.TagsWithChecked(c, mid, uint16(tid), title, filename, desc, cover, archive.TagPredictFromAPP)
	c.JSON(ptags, nil)
}

func creatorDataArchive(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	req := c.Request
	params := req.Form
	tyStr := params.Get("type")
	tmidStr := params.Get("tmid")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// check params
	ty, err := strconv.Atoi(tyStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, ok := data.IncrTy(int8(ty)); !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		arcStat  []*data.ThirtyDay
		archives []*archive.ArcVideoAudit
		show     int
		g        = &errgroup.Group{}
		ctx      = context.TODO()
	)
	log.Info("creatorDataArchive mid(%d) type(%d) access", mid, ty)
	g.Go(func() error {
		arcStat, err = dataSvc.ThirtyDayArchive(ctx, mid, int8(ty))
		return nil
	})
	g.Go(func() error {
		if arcs, err := arcSvc.WebArchives(ctx, mid, 0, "", "", "is_pubing,pubed,not_pubed", ip, 1, 2, 0); err == nil && arcs.Archives != nil {
			archives = arcs.Archives
		} else {
			archives = make([]*archive.ArcVideoAudit, 0)
		}
		return nil
	})
	g.Wait()
	if len(archives) > 0 {
		show = 1
	}
	if arcStat == nil {
		log.Info("creatorDataArchive mid(%d) type(%d) arcStat nil", mid, ty)
	}
	c.JSON(map[string]interface{}{
		"archive_stat": arcStat,
		"show":         show,
	}, nil)
}

func creatorDataArticle(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	req := c.Request
	params := req.Form
	tmidStr := params.Get("tmid")
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	tmid, _ := strconv.ParseInt(tmidStr, 10, 64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	var (
		artStat  []*artmdl.ThirtyDayArticle
		articles []*article.Meta
		show     int
		g        = &errgroup.Group{}
		ctx      = context.TODO()
	)
	g.Go(func() error {
		artStat, _ = dataSvc.ThirtyDayArticle(ctx, mid, ip)
		return nil
	})
	g.Go(func() error {
		if arts, err := artSvc.Articles(ctx, mid, 1, 2, 0, 0, 0, ip); err == nil && arts != nil && len(arts.Articles) != 0 {
			articles = arts.Articles
		} else {
			articles = []*article.Meta{}
		}
		return nil
	})
	g.Wait()
	if len(articles) > 0 {
		show = 1
	}
	c.JSON(map[string]interface{}{
		"article_stat": artStat,
		"show":         show,
	}, nil)
}

func creatorDescFormat(c *bm.Context) {
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
