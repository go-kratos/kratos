package http

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	accmdl "go-common/app/interface/main/creative/model/account"
	arcMdl "go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/danmu"
	elecMdl "go-common/app/interface/main/creative/model/elec"
	"go-common/app/interface/main/creative/model/order"
	porderM "go-common/app/interface/main/creative/model/porder"
	"go-common/app/interface/main/creative/model/tag"
	"go-common/app/interface/main/creative/model/watermark"
	accSvcModel "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

func webViewArc(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	hidStr := params.Get("history")
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		elecArc  *elecMdl.ArcState
		archive  *arcMdl.Archive
		videos   []*arcMdl.Video
		staffs   []*arcMdl.StaffView
		sMids    []int64
		sMap     map[int64]*accSvcModel.Info
		wm       *watermark.Watermark
		g        = &errgroup.Group{}
		ctx      = context.TODO()
		subtitle *danmu.SubtitleSubjectReply
	)
	g.Go(func() error {
		elecArc, _ = elecSvc.ArchiveState(ctx, aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		var (
			av *arcMdl.ArcVideo
		)
		if av, err = arcSvc.View(ctx, mid, aid, ip, arcMdl.PlatformWeb); err == nil && av != nil {
			archive = av.Archive
			videos = av.Videos
			staffs = av.Archive.Staffs
		} else {
			archive = &arcMdl.Archive{}
			videos = []*arcMdl.Video{}
			staffs = []*arcMdl.StaffView{}
		}
		return err
	})
	g.Go(func() error {
		wm, _ = wmSvc.WaterMark(c, mid)
		return nil
	})
	g.Go(func() error {
		subtitle, _ = danmuSvc.SubView(c, aid, ip)
		return nil
	})
	if err = g.Wait(); err != nil {
		c.JSON(nil, err)
		return
	}
	for _, v := range staffs {
		sMids = append(sMids, v.ApMID)
	}
	if sMap, err = accSvc.Infos(c, sMids, ip); err != nil {
		log.Error("accSvc.Infos(%v) error(%v)", sMids, err)
		err = nil
	}
	for k, v := range staffs {
		if m, ok := sMap[v.ApMID]; ok {
			staffs[k].ApName = m.Name
		}
	}
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err == nil && hid > 0 {
		history, err := arcSvc.HistoryView(c, mid, hid, ip)
		if err == nil && history.Mid > 0 {
			archive.Title = history.Title
			archive.Desc = history.Content
			archive.Cover = history.Cover
			videos = history.Video
		}
	}
	c.JSON(map[string]interface{}{
		"archive":   archive,
		"videos":    videos,
		"staffs":    staffs,
		"arc_elec":  elecArc,
		"watermark": wm,
		"tip":       vsSvc.WebManagerTip,
		"subtitle":  subtitle,
	}, nil)
}

func webArchives(c *bm.Context) {
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	order := params.Get("order")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	class := params.Get("status")
	tmidStr := params.Get("tmid")
	ip := metadata.String(c, metadata.RemoteIP)
	var tmid int64
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, ok := midI.(int64)
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if tmidStr != "" {
		tmid, _ = strconv.ParseInt(tmidStr, 10, 64)
	}
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// check params
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 20 {
		ps = 10
	}
	coop, _ := strconv.Atoi(params.Get("coop"))
	if coop > 0 {
		coop = 1
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	arc, err := arcSvc.WebArchives(c, mid, int16(tid), keyword, order, class, ip, pn, ps, coop)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arc, nil)
}

func webStaffApplies(c *bm.Context) {
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	tidStr := params.Get("tid")
	keyword := params.Get("keyword")
	state := params.Get("state")
	tmidStr := params.Get("tmid")
	//check user
	var tmid int64
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)
	if tmidStr != "" {
		tmid, _ = strconv.ParseInt(tmidStr, 10, 64)
	}
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	// check params
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 20 {
		ps = 10
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	arc, err := arcSvc.ApplySearch(c, mid, int16(tid), keyword, state, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arc, nil)
}

func webViewPre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	lang := params.Get("lang")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	var (
		err                        error
		mf                         *accmdl.MyInfo
		g                          = &errgroup.Group{}
		ctx                        = context.TODO()
		industryList, showtypeList []*porderM.Config
		orders                     = make([]*order.Order, 0)
		videoJam                   *arcMdl.VideoJam
		dymcLottery                bool
		wm                         *watermark.Watermark
		prePay                     map[string]interface{}
		staffConf                  struct {
			TypeList []*conf.StaffTypeConf `json:"typelist"`
			Titles   []*tag.StaffTitle     `json:"titles"`
		}
	)
	g.Go(func() error {
		mf, err = accSvc.MyInfo(ctx, mid, ip, time.Now())
		if err != nil {
			log.Info("accSvc.MyInfo (%d) err(%v)", mid, err)
		}
		if mf != nil {
			mf.Commercial = arcSvc.AllowCommercial(ctx, mid)
		}
		return nil
	})
	g.Go(func() error {
		industryList, _ = adSvc.IndustryList(ctx)
		showtypeList, _ = adSvc.ShowList(ctx)
		return nil
	})
	g.Go(func() error {
		if arcSvc.AllowOrderUps(mid) {
			orders, _ = arcSvc.ExecuteOrders(ctx, mid, ip)
		}
		return nil
	})
	g.Go(func() error {
		videoJam, _ = arcSvc.VideoJam(ctx, ip)
		return nil
	})
	g.Go(func() error {
		dymcLottery, _ = dymcSvc.LotteryUserCheck(ctx, mid)
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
		prePay, err = paySvc.Pre(ctx, mid)
		return nil
	})
	g.Wait()
	if mf != nil {
		mf.DymcLottery = dymcLottery
	}
	staffConf.Titles = arcSvc.StaffTitles(ctx)
	staffConf.TypeList = staffSvc.TypeConfig()
	c.JSON(map[string]interface{}{
		"prepay":        prePay,
		"video_jam":     videoJam,
		"typelist":      arcSvc.Types(c, lang),
		"activities":    arcSvc.Activities(c),
		"myinfo":        mf,
		"orders":        orders,
		"industry_list": industryList,
		"showtype_list": showtypeList,
		"watermark":     wm,
		"fav":           arcSvc.Fav(c, mid),
		"tip":           vsSvc.WebManagerTip,
		"staff_conf":    staffConf,
	}, nil)
}

func webDelArc(c *bm.Context) {
	req := c.Request
	ip := metadata.String(c, metadata.RemoteIP)
	parmes := req.Form
	challenge := parmes.Get("geetest_challenge")
	validate := parmes.Get("geetest_validate")
	seccode := parmes.Get("geetest_seccode")
	success := req.Form.Get("geetest_success")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(parmes.Get("aid"), 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%d) error(%v)", aid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	successi, err := strconv.Atoi(success)
	if err != nil {
		successi = 1
	}
	if gtSvc.Validate(c, challenge, validate, seccode, "web", ip, successi, mid) {
		c.JSON(nil, arcSvc.Del(c, mid, aid, ip))
	} else {
		c.JSON(nil, ecode.CreativeGeetestErr)
	}
}

func webVideos(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	archive, videos, err := arcSvc.Videos(c, mid, aid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"archive": archive,
		"videos":  videos,
	}, nil)
}

func webDescFormat(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	typeidStr := params.Get("typeid")
	cpStr := params.Get("copyright")
	lang := params.Get("lang")
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	copyright, err := strconv.ParseInt(cpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if copyright != 1 && copyright != 2 {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	format, err := arcSvc.DescFormat(c, typeid, copyright, lang, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(format, nil)
}

func webViewPoints(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	aidStr := params.Get("aid")
	cidStr := params.Get("cid")
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midStr.(int64)
	if !ok {
		log.Error("mid(%s) to int64 not ok", midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	vps, err := arcSvc.WebViewPoints(c, aid, cid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(vps, nil)
}

/*func webViewPointsEdit(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	aidStr := params.Get("aid")
	cidStr := params.Get("cid")
	midStr, ok := c.Get("mid")
	pointStr := params.Get("points")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, ok := midStr.(int64)
	if !ok {
		log.Error("mid(%s) to int64 not ok", midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", cidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := arcSvc.WebViewPointsEdit(c, aid, cid, mid, vp, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
*/

func webUserSearch(c *bm.Context) {
	var (
		err      error
		searchUp []*accmdl.SearchUp
		accCards map[int64]*accSvcModel.Card
		relas    map[int64]int
		mids     []int64
		g        = &errgroup.Group{}
	)
	searchUp = make([]*accmdl.SearchUp, 0)
	accCards = make(map[int64]*accSvcModel.Card)
	relas = make(map[int64]int)
	params := c.Request.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	kw := params.Get("kw")
	ip := metadata.String(c, metadata.RemoteIP)
	inMid, _ := strconv.ParseInt(kw, 10, 64)
	nickMid, err := danmuSvc.UserMid(c, kw, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if nickMid != 0 && nickMid != mid {
		mids = append(mids, nickMid)
	}
	if inMid != 0 && inMid != nickMid && inMid != mid {
		mids = append(mids, inMid)
	}
	if len(mids) == 0 {
		c.JSON(map[string]interface{}{
			"users": []*accmdl.SearchUp{},
		}, nil)
		return
	}
	g.Go(func() error {
		if accCards, err = accSvc.Cards(c, mids, ip); err != nil {
			log.Error("accSvc.Infos(%v) error(%v)", mids, err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		if relas, err = accSvc.FRelations(c, mid, mids, ip); err != nil {
			log.Error("accSvc.Relations(%d,%v) error(%v)", mid, mids, err)
			return err
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		c.JSON(nil, err)
		return
	}

	for _, v := range accCards {
		up := &accmdl.SearchUp{
			Mid:     v.Mid,
			Name:    v.Name,
			Face:    v.Face,
			Silence: v.Silence,
		}
		if rela, ok := relas[v.Mid]; ok {
			up.Relation = rela
			if rela >= 128 {
				up.IsBlock = true
			}
		}
		searchUp = append(searchUp, up)
	}
	c.JSON(map[string]interface{}{
		"users": searchUp,
	}, nil)
}
