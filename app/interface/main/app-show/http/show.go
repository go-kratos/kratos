package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/show"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_headerBuvid     = "Buvid"
	_headerDisplayID = "Display-ID"
)

func shows(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	ak := params.Get("access_key")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// get audit data, if check audit hit.
	ss, ok := showSvc.Audit(c, mobiApp, plat, build)
	if ok {
		returnJSON(c, ss, nil)
		return
	}
	network := params.Get("network")
	ip := metadata.String(c, metadata.RemoteIP)
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// display
	ss = showSvc.Display(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, "hans", "", false, time.Now())
	returnJSON(c, ss, nil)
	// infoc
	if len(ss) == 0 {
		return
	}
	showSvc.Infoc(mid, plat, buvid, disid, ip, "/x/v2/show", ss[0].Body, time.Now())
}

func showsRegion(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	ak := params.Get("access_key")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// get audit data, if check audit hit.
	ss, ok := showSvc.Audit(c, mobiApp, plat, build)
	if !ok {
		ip := metadata.String(c, metadata.RemoteIP)
		buvid := header.Get(_headerBuvid)
		network := params.Get("network")
		var mid int64
		if midInter, ok := c.Get("mid"); ok {
			mid = midInter.(int64)
		}
		// display
		language := params.Get("lang")
		ss = showSvc.RegionDisplay(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, language, "", false, time.Now())
	}
	res := map[string]interface{}{
		"data": ss,
	}
	returnDataJSON(c, res, 25, nil)
}

func showsIndex(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	channel := params.Get("channel")
	ak := params.Get("access_key")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// get audit data, if check audit hit.
	ss, ok := showSvc.Audit(c, mobiApp, plat, build)
	if !ok {
		ip := metadata.String(c, metadata.RemoteIP)
		buvid := header.Get(_headerBuvid)
		network := params.Get("network")
		var mid int64
		if midInter, ok := c.Get("mid"); ok {
			mid = midInter.(int64)
		}
		// display
		language := params.Get("lang")
		adExtra := params.Get("ad_extra")
		ss = showSvc.Index(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, language, adExtra, false, time.Now())
	}
	res := map[string]interface{}{
		"data": ss,
	}
	returnDataJSON(c, res, 25, nil)
}

func showTemps(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
	)
	// get params
	mobiApp := params.Get("mobi_app")
	// check params
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	ip := metadata.String(c, metadata.RemoteIP)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// display
	data := showSvc.Display(c, mid, plat, 0, header.Get(_headerBuvid), "", ip, "", "wifi", mobiApp, device, "hans", "", true, time.Now())
	returnJSON(c, data, nil)
}

func showChange(c *bm.Context) {
	params := c.Request.Form
	header := c.Request.Header
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	randStr := params.Get("rand")
	buildStr := params.Get("build")
	// check params
	rand, err := strconv.Atoi(randStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", randStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// normal data
	ip := metadata.String(c, metadata.RemoteIP)
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	network := params.Get("network")
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	build, _ := strconv.Atoi(buildStr)
	// change
	sis := showSvc.Change(c, mid, build, plat, rand, buvid, ip, network, mobiApp, device)
	returnJSON(c, sis, nil)
	// infoc
	showSvc.Infoc(mid, plat, buvid, disid, ip, "/x/v2/show/change", sis, time.Now())
}

// showRegionChange
func showRegionChange(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	buildStr := params.Get("build")
	plat := model.Plat(mobiApp, device)
	// get params
	randStr := params.Get("rand")
	// check params
	rand, err := strconv.Atoi(randStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", randStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ridStr := params.Get("rid")
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, _ := strconv.Atoi(buildStr)
	data := showSvc.RegionChange(c, rid, rand, plat, build, mobiApp)
	returnJSON(c, data, nil)
}

// showBangumiChange
func showBangumiChange(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// get params
	randStr := params.Get("rand")
	// check params
	rand, err := strconv.Atoi(randStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", randStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := showSvc.BangumiChange(c, rand, plat)
	returnJSON(c, data, nil)
}

// showArticleChange
func showArticleChange(c *bm.Context) {
	data := []*show.Item{}
	returnJSON(c, data, nil)
}

// showDislike
func showDislike(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		mid    int64
	)
	// get params
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	idStr := params.Get("id")
	gt := params.Get("goto")
	if !model.IsGoto(gt) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// normal data
	ip := metadata.String(c, metadata.RemoteIP)
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	// parse id
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// change
	si := showSvc.Dislike(c, mid, plat, id, buvid, mobiApp, device, gt, ip)
	returnJSON(c, si, nil)
	// infoc
	showSvc.Infoc(mid, plat, buvid, disid, ip, "/x/v2/show/change/dislike", []*show.Item{si}, time.Now())
}

func showWidget(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	buildStr := params.Get("build")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var data []*show.Item
	if ss, ok := showSvc.AuditChild(c, mobiApp, plat, build); ok {
		if len(ss) > 3 {
			data = ss[:3]
		} else {
			data = ss
		}
		returnJSON(c, data, nil)
		return
	}
	data = showSvc.Widget(c, plat)
	returnJSON(c, data, nil)
}

// show live change
func showLiveChange(c *bm.Context) {
	params := c.Request.Form
	// get params
	randStr := params.Get("rand")
	ak := params.Get("access_key")
	rand, err := strconv.Atoi(randStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", randStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// change
	ip := metadata.String(c, metadata.RemoteIP)
	data := showSvc.LiveChange(c, mid, ak, ip, rand, time.Now())
	returnJSON(c, data, nil)
}

// popular hot tab popular
func popular(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		mid    int64
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	idxStr := params.Get("idx")
	loginEventStr := params.Get("login_event")
	lastParam := params.Get("last_param")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	loginEvent, err := strconv.Atoi(loginEventStr)
	if err != nil {
		loginEvent = 0
	}
	idx, err := strconv.ParseInt(idxStr, 10, 64)
	if err != nil || idx < 0 {
		idx = 0
	}
	now := time.Now()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get(_headerBuvid)
	// get audit data, if check audit hit.
	data, ok := showSvc.AuditFeed(c, mobiApp, plat, build)
	if !ok {
		data = showSvc.FeedIndex(c, mid, idx, plat, build, loginEvent, lastParam, mobiApp, device, buvid, now)
	}
	c.JSON(data, nil)
}

// popular hot tab popular
func popular2(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		mid    int64
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	idxStr := params.Get("idx")
	loginEventStr := params.Get("login_event")
	lastParam := params.Get("last_param")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	loginEvent, err := strconv.Atoi(loginEventStr)
	if err != nil {
		loginEvent = 0
	}
	idx, err := strconv.ParseInt(idxStr, 10, 64)
	if err != nil || idx < 0 {
		idx = 0
	}
	now := time.Now()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get(_headerBuvid)
	// get audit data, if check audit hit.
	data, ok := showSvc.AuditFeed2(c, mobiApp, plat, build)
	var (
		ver string
	)
	if !ok {
		data, ver, err = showSvc.FeedIndex2(c, mid, idx, plat, build, loginEvent, lastParam, mobiApp, device, buvid, now)
	}
	config := map[string]interface{}{
		"item_title": "当前热门",
	}
	c.JSONMap(map[string]interface{}{"data": data, "ver": ver, "config": config}, err)
}
