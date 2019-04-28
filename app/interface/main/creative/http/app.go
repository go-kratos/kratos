package http

import (
	"context"
	"fmt"
	"go-common/app/interface/main/creative/model/account"
	accmdl "go-common/app/interface/main/creative/model/account"
	"go-common/app/interface/main/creative/model/activity"
	"go-common/app/interface/main/creative/model/app"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/main/creative/model/elec"
	"go-common/app/interface/main/creative/model/faq"
	"go-common/app/interface/main/creative/model/message"
	"go-common/app/interface/main/creative/model/newcomer"
	resMdl "go-common/app/interface/main/creative/model/resource"
	"go-common/app/interface/main/creative/model/search"
	"go-common/app/interface/main/creative/model/watermark"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"strconv"
	"strings"
	"time"
)

func appIndex(c *bm.Context) {
	var (
		err                        error
		rec                        *elec.RecentElecList
		recList                    []*elec.RecentElec
		stat                       *data.Stat
		archives                   []*archive.OldArchiveVideoAudit
		arcs                       *search.Result
		replies                    *search.Replies
		elecStat                   *elec.UserState
		dataStat                   *data.AppStatList
		topBanners, academyBanners []*resMdl.Banner
		portal                     []*app.Portal
		topMsgs                    []*message.Message
		appTasks                   *newcomer.AppIndexNewcomer
		artStat                    artmdl.UpStat
		artAuthor, build, coop     int
		g                          = &errgroup.Group{}
		ctx                        = context.TODO()
	)
	req := c.Request
	header := c.Request.Header
	params := req.Form
	ip := metadata.String(c, metadata.RemoteIP)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platStr := params.Get("platform")
	network := params.Get("network")
	buvid := header.Get("Buvid")
	adExtra := params.Get("ad_extra")
	if coop, _ = strconv.Atoi(params.Get("coop")); coop > 0 {
		coop = 1
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	resID, _ := strconv.Atoi(params.Get("resource_id"))
	resMdlPlat := resMdl.Plat(mobiApp, device)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
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
	g.Go(func() error {
		stat, _ = dataSvc.NewStat(ctx, mid, ip)
		if stat != nil {
			stat.Day30 = nil
			stat.Arcs = nil
		}
		return nil
	})
	g.Go(func() error {
		if arcs, err = arcSvc.Archives(ctx, mid, 0, "", "", "is_pubing,pubed,not_pubed", ip, 1, 2, coop); err == nil && arcs.OldArchives != nil {
			archives = arcs.OldArchives
		} else {
			archives = make([]*archive.OldArchiveVideoAudit, 0)
		}
		return nil
	})
	g.Go(func() error {
		replies, _ = replySvc.AppIndexReplies(ctx, ak, ck, mid, 0, 0, 0, search.All, resMdlPlat, "", "", "", ip, 1, 10)
		if replies == nil {
			replies = &search.Replies{}
		}
		return nil
	})
	g.Go(func() error {
		elecStat, _ = elecSvc.UserState(ctx, mid, ip, ak, ck)
		if elecStat != nil && elecStat.State == "2" {
			rec, _ = elecSvc.RecentElec(ctx, mid, 1, 2, ip)
			if rec != nil {
				recList = rec.List
			}
		}
		return nil
	})
	g.Go(func() error {
		academyBanners, _ = resSvc.AcademyBanner(ctx, mobiApp, device, network, metadata.String(c, metadata.RemoteIP), buvid, adExtra, build, resID, resMdlPlat, mid, false)
		return nil
	})
	g.Go(func() error {
		topBanners, _ = resSvc.TopBanner(ctx, mobiApp, device, network, metadata.String(c, metadata.RemoteIP), buvid, adExtra, build, resID, resMdlPlat, mid, false)
		isAuthor, _ := artSvc.IsAuthor(ctx, mid, ip)
		if isAuthor {
			artAuthor = 1
		} else {
			artAuthor = 0
		}
		portal, _ = appSvc.Portals(ctx, mid, artAuthor, build, app.PortalIntro, platStr, resMdlPlat)
		dataStat, _ = dataSvc.AppStat(ctx, mid)
		if platStr == "android" && build < 510007 {
			dataStat.Show = 0 // android < 510007 close
		}
		return nil
	})
	g.Go(func() error {
		artStat, _ = artSvc.ArticleStat(ctx, mid, ip)
		return nil
	})
	g.Go(func() error {
		topMsgs, _ = appSvc.TopMsg(ctx, mid, build, platStr, mobiApp, ak, ck, ip)
		return nil
	})
	g.Go(func() error {
		appTasks, _ = newcomerSvc.AppIndexNewcomer(ctx, mid, platStr)
		return nil
	})
	g.Wait()
	up := &app.Up{
		Art: artAuthor,
	}
	if len(archives) > 0 {
		up.Arc = 1
	}
	c.JSON(map[string]interface{}{
		"stat":             stat,
		"tasks":            appTasks,
		"archives":         archives,
		"replies":          replies,
		"user_elec":        elecStat,
		"recent_elec_rank": recList,
		"banner":           topBanners,
		"aca_banner":       academyBanners,
		"portal_list":      portal,
		"data_stat":        dataStat,
		"art_stat":         artStat,
		"up":               up,
		"tip":              operSvc.NoticeStr,
		"top_acts":         arcSvc.TopAct(),
		"top_msgs":         topMsgs,
		"block_intros":     appSvc.BlockIntros(build, platStr),
	}, nil)
}

func appArcView(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	lang := params.Get("lang")
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	mobiApp := params.Get("mobi_app")
	plat := params.Get("platform")
	// check user
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
	var (
		arc                       *archive.Archive
		av                        *archive.ArcVideo
		video                     []*archive.Video
		white                     int
		elecStat                  *elec.UserState
		elecArc                   *elec.ArcState
		activities                []*activity.Activity
		appFmt                    []*archive.AppFormat
		mf                        *account.MyInfo
		wm                        *watermark.Watermark
		g                         = &errgroup.Group{}
		ctx                       = context.TODO()
		faqs                      = make(map[string]*faq.Faq)
		lotteryCheck, lotteryBind bool
		recFriends                []*accmdl.Friend
	)
	g.Go(func() error {
		if av, err = arcSvc.View(ctx, mid, aid, ip, plat); err == nil && av != nil {
			arc = av.Archive
			video = av.Videos
			if (mobiApp == "android" && build < 511001) || (mobiApp == "ios" && build < 6011) {
				av.Archive.Desc = archive.ShortDesc(av.Archive.Desc)
			}
			activities = arcSvc.Activities(ctx)
			if arc.MissionID > 0 {
				protectTagBeforeMission(ctx, arc, plat, build)
				actInList := false
				for _, v := range activities {
					if v.ID == arc.MissionID {
						actInList = true
						break
					}
				}
				if !actInList {
					activities = append(activities, &activity.Activity{
						ID:   arc.MissionID,
						Name: arc.MissionName,
					})
				}
			}
		}
		return nil
	})
	g.Go(func() error {
		elecStat, _ = elecSvc.UserState(ctx, mid, ip, ak, ck)
		return nil
	})
	g.Go(func() error {
		elecArc, _ = elecSvc.ArchiveState(ctx, aid, mid, ip)
		if elecArc == nil {
			elecArc = &elec.ArcState{}
		}
		return nil
	})
	g.Go(func() error {
		appFmt, _ = arcSvc.AppFormats(ctx)
		return nil
	})
	g.Go(func() error {
		mf, _ = accSvc.MyInfo(ctx, mid, ip, time.Now())
		_, white = whiteSvc.UploadInfoForMainApp(mf, plat, mid)
		return nil
	})
	g.Go(func() error {
		wm, _ = wmSvc.WaterMark(ctx, mid)
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
		lotteryBind, _ = dymcSvc.LotteryNotice(ctx, aid, mid)
		return nil
	})
	g.Go(func() error {
		recFriends, _ = accSvc.RecFollows(ctx, mid)
		return nil
	})
	g.Wait()
	if arc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]interface{}{
		"archive":     arc,
		"videos":      video,
		"user_elec":   elecStat,
		"arc_elec":    elecArc,
		"typelist":    arcSvc.AppTypes(c, lang),
		"activities":  activities,
		"desc_format": appFmt,
		"dpub":        arcSvc.Dpub(),
		"rules":       arcSvc.EditRules(c, white, arc.State, lotteryBind),
		"watermark":   wm,
		"tip":         vsSvc.AppManagerTip,
		"cus_tip":     vsSvc.CusManagerTip,
		// common data
		"camera_cfg":  appSvc.CameraCfg,
		"module_show": arcSvc.AppModuleShowMap(mid, lotteryCheck),
		"icons":       appSvc.Icons(),
		"faqs":        faqs,
		"rec_friends": recFriends,
	}, nil)
}

func appArcDel(c *bm.Context) {
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
	// check params
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// del
	c.JSON(nil, arcSvc.Del(c, mid, aid, ip))
}

func appArchives(c *bm.Context) {
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	order := params.Get("order")
	tidStr := params.Get("tid")
	kw := params.Get("keyword")
	class := params.Get("class")
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		pn, ps, coop int
		mid, tid     int64
	)
	midStr, _ := c.Get("mid")
	if mid = midStr.(int64); mid < 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if pn, _ = strconv.Atoi(pageStr); pn <= 0 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(psStr); ps <= 0 || ps > 50 {
		ps = 10
	}
	if tid, _ = strconv.ParseInt(tidStr, 10, 16); tid < 0 {
		tid = 0
	}
	if coop, _ = strconv.Atoi(params.Get("coop")); coop > 0 {
		coop = 1
	}
	arc, err := arcSvc.Archives(c, mid, int16(tid), kw, order, class, ip, pn, ps, coop)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if arc != nil {
		arc.Tip = operSvc.CreativeStr
	}
	c.JSON(arc, nil)
}

func appReplyList(c *bm.Context) {
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
	resMdlPlat := resMdl.Plat(params.Get("mobi_app"), params.Get("device"))
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
		ResMdlPlat:  resMdlPlat,
	}
	replies, err := replySvc.Replies(c, p)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"pager": map[string]int{
			"current": pn,
			"size":    ps,
			"total":   replies.Total,
		},
		"data": replies.Result,
	}, nil)
}

func appUpInfo(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
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
	v, err := accSvc.UpInfo(c, mid, ip)
	if err != nil {
		log.Error("memberSvc.UpInfo(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(v, nil)
}

func appPre(c *bm.Context) {
	req := c.Request
	params := req.Form
	ip := metadata.String(c, metadata.RemoteIP)
	buildStr := params.Get("build")
	platStr := params.Get("platform")
	tmidStr := params.Get("tmid")
	device := params.Get("device")
	mobiApp := params.Get("mobi_app")
	resMdlPlat := resMdl.Plat(mobiApp, device)
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
		err                       error
		up                        *account.UpInfo
		portal                    []*app.Portal
		build, artAuthor, showCre int
		g                         = &errgroup.Group{}
		ctx                       = context.TODO()
		academyIntro              = &app.AcademyIntro{}
		actIntro                  = &app.ActIntro{}
		mf                        *account.MyInfo
	)
	if buildStr != "" {
		build, err = strconv.Atoi(buildStr)
		if err != nil {
			log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	g.Go(func() error {
		mf, err = accSvc.MyInfo(ctx, mid, ip, time.Now())
		if err != nil {
			log.Info("accSvc.MyInfo (%d) err(%v)", mid, err)
		}
		return nil
	})
	g.Go(func() error {
		up, _ = accSvc.UpInfo(ctx, mid, ip)
		if up != nil && account.IsUper(up) {
			showCre = 1
		}
		if showCre != 1 {
			if isAuthor, err := upSvc.ArcUpInfo(ctx, mid, ip); err == nil && isAuthor == 1 {
				showCre = 1
			}
		}
		return nil
	})
	g.Go(func() error {
		isAuthor, _ := artSvc.IsAuthor(ctx, mid, ip)
		if isAuthor {
			artAuthor = 1
		} else {
			artAuthor = 0
		}
		portal, _ = appSvc.Portals(ctx, mid, artAuthor, build, app.PortalNotice, platStr, resMdlPlat)
		return nil
	})
	g.Wait()
	entranceGuidance := "成为UP主，分享你的创作"
	if platStr == "android" {
		entranceGuidance = "投稿"
	}
	if showCre == 0 {
		ccURL := appSvc.H5Page().CreativeCollege
		if (build >= 5350000 && platStr == "android") || (build >= 8240 && platStr == "ios") {
			ccURL = fmt.Sprintf("%s?from=my&navhide=1", ccURL)
		} else {
			ccURL = fmt.Sprintf("%s?from=my", ccURL)
		}
		academyIntro = &app.AcademyIntro{
			Show:  1,
			Title: "创作学院",
			URL:   ccURL,
		}
	}
	uploadinfo, _ := whiteSvc.UploadInfoForMainApp(mf, platStr, mid)
	if uploadinfo["info"] == 1 {
		actIntro = &app.ActIntro{
			Show:  1,
			Title: "热门活动",
			URL:   appSvc.H5Page().HotAct,
		}
	}
	c.JSON(map[string]interface{}{
		"creative": map[string]interface{}{
			"portal_list": portal,
			"show":        showCre,
		},
		"entrance": map[string]interface{}{
			"guidance": entranceGuidance,
			"show":     1,
		},
		"academy": academyIntro,
		"act":     actIntro,
	}, nil)
}

func appBanner(c *bm.Context) {
	params := c.Request.Form
	psStr := params.Get("ps")
	pnStr := params.Get("pn")
	tp := params.Get("type")
	_, ok := c.Get("mid")
	if !ok {
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
	oper, err := operSvc.AppOperationList(c, pn, ps, tp)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(oper, nil)
}

func appNewPre(c *bm.Context) {
	req := c.Request
	params := req.Form
	ip := metadata.String(c, metadata.RemoteIP)
	midStr := params.Get("mid")
	tmidStr := params.Get("tmid")
	var (
		err                             error
		mid, tmid                       int64
		up                              *account.UpInfo
		isUp, showCreative, showAcademy int
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt midStr(%s) error(%v)", midStr, err)
		c.JSON(nil, err)
		return
	}
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if tmidStr != "" {
		if tmid, err = strconv.ParseInt(tmidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt tmidStr(%s) error(%v)", tmidStr, err)
			c.JSON(nil, err)
			return
		}
		if tmid > 0 && dataSvc.IsWhite(mid) {
			mid = tmid
		}
	}
	up, _ = accSvc.UpInfo(c, mid, ip)
	if up != nil && account.IsUper(up) {
		isUp = 1
	}
	if isUp != 1 {
		if isAuthor, err := upSvc.ArcUpInfo(c, mid, ip); err == nil && isAuthor == 1 {
			isUp = 1
		}
	}
	if isUp == 0 {
		showAcademy = 1
	}
	if dataSvc.IsForbidVideoup(mid) {
		showCreative = 0
	} else {
		showCreative = 1
	}
	c.JSON(map[string]interface{}{
		"academy": showAcademy,
		"show":    showCreative,
		"is_up":   isUp,
	}, nil)
}

func protectTagBeforeMission(c context.Context, arc *archive.Archive, plat string, build int) {
	if ((plat == "android" && build < 5260000) || (plat == "ios" && build <= 6680)) && len(arc.Tag) > 0 {
		var mp *activity.Protocol
		mp, _ = arcSvc.MissionProtocol(c, arc.MissionID)
		if len(mp.Tags) > 0 {
			sTags := strings.Split(mp.Tags, ",")
			missionTagMap := make(map[string]int8)
			for _, tag := range sTags {
				missionTagMap[tag] = 0
			}
			arcTags := strings.Split(arc.Tag, ",")
			var splitedArcTags []string
			for _, arcTag := range arcTags {
				if _, ok := missionTagMap[arcTag]; !ok {
					splitedArcTags = append(splitedArcTags, arcTag)
				}
			}
			arc.Tag = strings.Join(splitedArcTags, ",")
		}
	}
}

func appSimpleArcVideos(c *bm.Context) {
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	order := params.Get("order")
	kw := params.Get("keyword")
	class := params.Get("class")
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		pn, ps, tid int
	)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if pn, _ = strconv.Atoi(pageStr); pn <= 0 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(psStr); ps <= 0 || ps > 50 {
		ps = 10
	}
	if tid, _ = strconv.Atoi(params.Get("tid")); tid <= 0 {
		tid = 0
	}
	mid = converTmid(c, mid)
	savs, err := arcSvc.SimpleArcVideos(c, mid, int16(tid), kw, order, class, ip, pn, ps, 0)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(savs, nil)
}

func appTaskBind(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	// check white list
	//if task := whiteSvc.TaskWhiteList(mid); task != 1 {
	//	c.JSON(nil, ecode.RequestErr)
	//	return
	//}
	id, err := newcomerSvc.TaskBind(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}
