package http

import (
	"strconv"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid       = "Buvid"
	_headerDisplayID   = "Display-ID"
	_headerDeviceID    = "Device-ID"
	_androidFnvalBuild = 5325000
	_iosFnvalBuild     = 8160
	_iosQnBuildGt      = 8170
	_iosQnBuildLt      = 8190
	_androidQnBuildLt  = 5335000
	_qn480             = 32
)

func feedIndex(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	network := params.Get("network")
	buildStr := params.Get("build")
	idxStr := params.Get("idx")
	pullStr := params.Get("pull")
	styleStr := params.Get("style")
	loginEventStr := params.Get("login_event")
	openEvent := params.Get("open_event")
	bannerHash := params.Get("banner_hash")
	adExtra := params.Get("ad_extra")
	qnStr := params.Get("qn")
	interest := params.Get("interest")
	flushStr := params.Get("flush")
	autoplayCard, _ := strconv.Atoi(params.Get("autoplay_card"))
	fnver, _ := strconv.Atoi(params.Get("fnver"))
	fnval, _ := strconv.Atoi(params.Get("fnval"))
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	if (model.IsAndroid(plat) && build <= _androidFnvalBuild) || (model.IsIOSNormal(plat) && build <= _iosFnvalBuild) {
		fnval = 0
	}
	style, _ := strconv.Atoi(styleStr)
	flush, _ := strconv.Atoi(flushStr)
	// get audit data, if check audit hit.
	is, ok := feedSvc.Audit(c, mobiApp, plat, build)
	if ok {
		c.JSON(is, nil)
		return
	}
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	dvcid := header.Get(_headerDeviceID)
	// page
	idx, err := strconv.ParseInt(idxStr, 10, 64)
	if err != nil || idx < 0 {
		idx = 0
	}
	// pull default
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		pull = true
	}
	// login event
	loginEvent, err := strconv.Atoi(loginEventStr)
	if err != nil {
		loginEvent = 0
	}
	// qn
	qn, _ := strconv.Atoi(qnStr)
	now := time.Now()
	// index
	data, userFeature, isRcmd, newUser, code, feedclean, autoPlayInfoc, err := feedSvc.Index(c, mid, plat, build, buvid, network, mobiApp, device, platform, openEvent, loginEvent, idx, pull, now, bannerHash, adExtra, qn, interest, style, flush, fnver, fnval, autoplayCard)
	res := map[string]interface{}{
		"data": data,
		"config": map[string]interface{}{
			"feed_clean_abtest": feedclean,
		},
	}
	c.JSONMap(res, err)
	if err != nil {
		return
	}
	// infoc
	items := make([]*ai.Item, 0, len(data))
	for _, item := range data {
		items = append(items, item.AI)
	}
	feedSvc.IndexInfoc(c, mid, plat, build, buvid, disid, "/x/feed/index", userFeature, style, code, items, isRcmd, pull, newUser, now, "", dvcid, network, flush, autoPlayInfoc, 0)
}

func feedUpper(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check page
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 {
		ps = 20
	} else if ps > 200 {
		ps = 200
	}
	plat := model.Plat(mobiApp, device)
	now := time.Now()
	uas, _ := feedSvc.Upper(c, mid, plat, build, pn, ps, now)
	data := map[string]interface{}{}
	if len(uas) != 0 {
		data["item"] = uas
	} else {
		data["item"] = []struct{}{}
	}
	uls, count := feedSvc.UpperLive(c, mid)
	if len(uls) != 0 {
		data["live"] = struct {
			Item  []*feed.Item `json:"item"`
			Count int          `json:"count"`
			Conut int          `json:"conut"`
		}{uls, count, count}
	}
	c.JSON(data, nil)
}

func feedUpperArchive(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check page
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 {
		ps = 20
	} else if ps > 200 {
		ps = 200
	}
	plat := model.Plat(mobiApp, device)
	now := time.Now()
	uas, _ := feedSvc.UpperArchive(c, mid, plat, build, pn, ps, now)
	data := map[string]interface{}{}
	if len(uas) != 0 {
		data["item"] = uas
	} else {
		data["item"] = []struct{}{}
	}
	c.JSON(data, nil)
}

func feedUpperBangumi(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check page
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 {
		ps = 20
	} else if ps > 200 {
		ps = 200
	}
	plat := model.Plat(mobiApp, device)
	now := time.Now()
	uas, _ := feedSvc.UpperBangumi(c, mid, plat, build, pn, ps, now)
	data := map[string]interface{}{}
	if len(uas) != 0 {
		data["item"] = uas
	} else {
		data["item"] = []struct{}{}
	}
	c.JSON(data, nil)
}

func feedUpperArticle(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check page
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 {
		ps = 20
	} else if ps > 200 {
		ps = 200
	}
	plat := model.Plat(mobiApp, device)
	now := time.Now()
	uas, _ := feedSvc.UpperArticle(c, mid, plat, build, pn, ps, now)
	data := map[string]interface{}{}
	if len(uas) != 0 {
		data["item"] = uas
	} else {
		data["item"] = []struct{}{}
	}
	c.JSON(data, nil)
}

func feedUnreadCount(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	// get params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	total, feedCount, articleCount := feedSvc.UnreadCount(c, mid, plat, build, time.Now())
	c.JSON(struct {
		Total   int `json:"total"`
		Count   int `json:"count"`
		Article int `json:"article"`
	}{total, feedCount, articleCount}, nil)
}

func feedDislike(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	gt := params.Get("goto")
	if gt == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, _ := strconv.ParseInt(params.Get("id"), 10, 64)
	reasonID, _ := strconv.ParseInt(params.Get("reason_id"), 10, 64)
	cmreasonID, _ := strconv.ParseInt(params.Get("cm_reason_id"), 10, 64)
	feedbackID, _ := strconv.ParseInt(params.Get("feedback_id"), 10, 64)
	upperID, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	rid, _ := strconv.ParseInt(params.Get("rid"), 10, 64)
	tagID, _ := strconv.ParseInt(params.Get("tag_id"), 10, 64)
	adcb := params.Get("ad_cb")
	buvid := header.Get(_headerBuvid)
	if buvid == "" && mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, feedSvc.Dislike(c, mid, id, buvid, gt, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, time.Now()))
}

func feedDislikeCancel(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	gt := params.Get("goto")
	if gt == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, _ := strconv.ParseInt(params.Get("id"), 10, 64)
	reasonID, _ := strconv.ParseInt(params.Get("reason_id"), 10, 64)
	cmreasonID, _ := strconv.ParseInt(params.Get("cm_reason_id"), 10, 64)
	feedbackID, _ := strconv.ParseInt(params.Get("feedback_id"), 10, 64)
	upperID, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	rid, _ := strconv.ParseInt(params.Get("rid"), 10, 64)
	tagID, _ := strconv.ParseInt(params.Get("tag_id"), 10, 64)
	adcb := params.Get("ad_cb")
	buvid := header.Get(_headerBuvid)
	if buvid == "" && mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, feedSvc.DislikeCancel(c, mid, id, buvid, gt, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, time.Now()))
}

func feedUpperRecent(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	aidStr := params.Get("param")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	upperStr := params.Get("vmid")
	upperID, err := strconv.ParseInt(upperStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(struct {
		Item []*feed.Item `json:"item"`
	}{feedSvc.UpperRecent(c, mid, upperID, aid, time.Now())}, nil)
}

func feedIndexTab(c *bm.Context) {
	var (
		id      int64
		items   []*feed.Item
		isBnj   bool
		bnjDays int
		cover   string
		err     error
		mid     int64
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	now := time.Now()
	idStr := params.Get("id")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	if id, _ = strconv.ParseInt(idStr, 10, 64); id <= 0 {
		c.JSON(struct {
			Tab []*operate.Menu `json:"tab"`
		}{feedSvc.Menus(c, plat, build, now)}, nil)
		return
	}
	items, cover, isBnj, bnjDays, err = feedSvc.Actives(c, id, mid, now)
	c.JSON(struct {
		Cover   string       `json:"cover"`
		IsBnj   bool         `json:"is_bnj,omitempty"`
		BnjDays int          `json:"bnj_days,omitempty"`
		Item    []*feed.Item `json:"item"`
	}{cover, isBnj, bnjDays, items}, err)
}

func feedIndex2(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	header := c.Request.Header
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	dvcid := header.Get(_headerDeviceID)
	param := &feed.IndexParam{}
	// get params
	if err := c.Bind(param); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, ok := cdm.Columnm[param.Column]
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// 兼容老的style逻辑，3为新单列
	style := int(cdm.Columnm[param.Column])
	if style == 1 {
		style = 3
	}
	// check params
	plat := model.Plat(param.MobiApp, param.Device)
	// get audit data, if check audit hit.
	if data, ok := feedSvc.Audit2(c, param.MobiApp, plat, param.Build, param.Column); ok {
		c.JSON(struct {
			Item []card.Handler `json:"items"`
		}{Item: data}, nil)
		return
	}
	if (model.IsAndroid(plat) && param.Build <= _androidFnvalBuild) || (model.IsIOSNormal(plat) && param.Build <= _iosFnvalBuild) {
		param.Fnval = 0
	}
	if (model.IsAndroid(plat) && param.Build > _androidFnvalBuild && param.Build < _androidQnBuildLt) || (model.IsIOSNormal(plat) && param.Build > _iosQnBuildGt && param.Build <= _iosQnBuildLt) || param.Qn <= 0 {
		param.Qn = _qn480
	}
	now := time.Now()
	// index
	plat = model.PlatAPPBuleChange(plat)
	data, config, infc, err := feedSvc.Index2(c, buvid, mid, plat, param, style, now)
	c.JSON(struct {
		Item   []card.Handler `json:"items"`
		Config *feed.Config   `json:"config"`
	}{Item: data, Config: config}, err)
	if err != nil {
		return
	}
	// infoc
	items := make([]*ai.Item, 0, len(data))
	for _, item := range data {
		items = append(items, item.Get().Rcmd)
	}
	feedSvc.IndexInfoc(c, mid, plat, param.Build, buvid, disid, "/x/feed/index", infc.UserFeature, style, infc.Code, items, infc.IsRcmd, param.Pull, infc.NewUser, now, "", dvcid, param.Network, param.Flush, infc.AutoPlayInfoc, param.DeviceType)
}

func feedIndexTab2(c *bm.Context) {
	var (
		id      int64
		items   []card.Handler
		isBnj   bool
		bnjDays int
		cover   string
		err     error
		mid     int64
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	now := time.Now()
	idStr := params.Get("id")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	forceHost, _ := strconv.Atoi(params.Get("force_host"))
	buildStr := params.Get("build")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	if id, _ = strconv.ParseInt(idStr, 10, 64); id <= 0 {
		c.JSON(struct {
			Tab []*operate.Menu `json:"tab"`
		}{feedSvc.Menus(c, plat, build, now)}, nil)
		return
	}
	items, cover, isBnj, bnjDays, err = feedSvc.Actives2(c, id, mid, mobiApp, plat, build, forceHost, now)
	c.JSON(struct {
		Cover   string         `json:"cover"`
		IsBnj   bool           `json:"is_bnj,omitempty"`
		BnjDays int            `json:"bnj_days,omitempty"`
		Item    []card.Handler `json:"items"`
	}{cover, isBnj, bnjDays, items}, err)
}

func feedIndexConverge(c *bm.Context) {
	var (
		mid   int64
		title string
		cover string
		uri   string
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	param := &feed.ConvergeParam{}
	if err := c.Bind(param); err != nil {
		return
	}
	plat := model.Plat(param.MobiApp, param.Device)
	if (model.IsAndroid(plat) && param.Build <= _androidFnvalBuild) || (model.IsIOSNormal(plat) && param.Build <= _iosFnvalBuild) {
		param.Fnval = 0
	}
	data, converge, err := feedSvc.Converge(c, mid, plat, param, time.Now())
	if converge != nil {
		title = converge.Title
		cover = converge.Cover
		uri = converge.URI
	}
	c.JSON(struct {
		Items []card.Handler `json:"items,omitempty"`
		Title string         `json:"title,omitempty"`
		Cover string         `json:"cover,omitempty"`
		Param string         `json:"param,omitempty"`
		URI   string         `json:"uri,omitempty"`
	}{Items: data, Title: title, Cover: cover, Param: strconv.FormatInt(param.ID, 10), URI: uri}, err)
}
