package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/region"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	_emptyShowItems = []*region.ShowItem{}
)

// regions get region data
func regions(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	language := params.Get("lang")
	ver := params.Get("ver")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	data, version, err := regionSvc.Regions(c, plat, build, ver, mobiApp, device, language)
	if err == ecode.NotModified {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
		"ver":  version,
	}
	returnDataJSON(c, res, 1, nil)
}

// regions get region data
func regionsList(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	language := params.Get("lang")
	entrance := params.Get("entrance")
	ver := params.Get("ver")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	data, version, err := regionSvc.RegionsList(c, plat, build, ver, mobiApp, device, language, entrance)
	if err == ecode.NotModified {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
		"ver":  version,
	}
	returnDataJSON(c, res, 1, nil)
}

// regions get region data
func regionsIndex(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	language := params.Get("lang")
	ver := params.Get("ver")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	data, version, err := regionSvc.NewRegionList(c, plat, build, ver, mobiApp, device, language)
	if err == ecode.NotModified {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
		"ver":  version,
	}
	returnDataJSON(c, res, 1, nil)
}

// regionShow region show
func regionShow(c *bm.Context) {
	header := c.Request.Header
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	buildStr := params.Get("build")
	channel := params.Get("channel")
	network := params.Get("network")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get(_headerBuvid)
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	adExtra := params.Get("ad_extra")
	// GetAudit
	if audit, ok := regionSvc.Audit(c, mobiApp, plat, build, rid, true); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.Show(c, plat, rid, build, mid, channel, buvid, network, mobiApp, device, adExtra)
		returnJSON(c, data, nil)
	}
}

// regionChildShow region childShow
func regionChildShow(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	tagIDStr := params.Get("tag_id")
	buildStr := params.Get("build")
	channel := params.Get("channel")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		mid   int64
		tagID int
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if tagIDStr != "" {
		if tagID, err = strconv.Atoi(tagIDStr); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// GetAudit
	if audit, ok := regionSvc.AuditChild(c, mobiApp, "default", plat, build, rid, tagID); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.ChildShow(c, plat, mid, rid, tagID, build, channel, mobiApp, time.Now())
		returnJSON(c, data, nil)
	}
}

// regionChildListShow region childlistShow
func regionChildListShow(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("rid")
	tagIDStr := params.Get("tag_id")
	mobiApp := params.Get("mobi_app")
	pnStr := params.Get("pn")
	// psStr := params.Get("ps")
	orderStr := params.Get("order")
	buildStr := params.Get("build")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		return
	}
	var tagID int
	if tagIDStr != "" {
		if tagID, err = strconv.Atoi(tagIDStr); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	// ps, err := strconv.Atoi(psStr)
	// if err != nil || ps > 60 || ps <= 0 {
	ps := 20
	// }
	if pn*ps > 400 {
		returnJSON(c, _emptyShowItems, nil)
		return
	}
	order := ""
	switch orderStr {
	case "view":
		order = "click"
	case "reply":
		order = "scores"
	case "danmaku":
		order = "dm"
	case "favorite":
		order = "stow"
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	platform := params.Get("platform")
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// GetAudit
	if audit, ok := regionSvc.AuditChildList(c, mobiApp, order, plat, build, rid, tagID, pn, ps); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.ChildListShow(c, plat, rid, tagID, pn, ps, build, mid, order, platform, mobiApp, device)
		returnJSON(c, data, nil)
	}
}

// regionChildListShow region childlistShow
func regionShowDynamic(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("rid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	rid, err := strconv.Atoi(ridStr)
	mobiApp := params.Get("mobi_app")
	mobiApp = model.MobiAPPBuleChange(mobiApp)
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	device := params.Get("device")
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps > 50 || ps <= 0 {
		ps = 50
	}
	if pn*ps > 200 {
		returnJSON(c, _emptyShowItems, nil)
		return
	}
	plat := model.Plat(mobiApp, device)
	data := regionSvc.ShowDynamic(c, plat, build, rid, pn, ps)
	returnJSON(c, data, nil)
}

func regionDynamic(c *bm.Context) {
	header := c.Request.Header
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	buildStr := params.Get("build")
	channel := params.Get("channel")
	network := params.Get("network")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get(_headerBuvid)
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	adExtra := params.Get("ad_extra")
	// GetAudit
	if audit, ok := regionSvc.Audit(c, mobiApp, plat, build, rid, true); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.Dynamic(c, plat, rid, build, mid, channel, buvid, network, mobiApp, device, adExtra, time.Now())
		returnJSON(c, data, nil)
	}
}

func regionDynamicList(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	pullStr := params.Get("pull")
	ctimeStr := params.Get("ctime")
	buildStr := params.Get("build")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		log.Error("pullStr(%s) error(%v)", pullStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// ctime
	ctime, err := strconv.ParseInt(ctimeStr, 10, 64)
	if err != nil || ctime < 0 {
		ctime = 0
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// GetAudit
	if _, ok := regionSvc.Audit(c, mobiApp, plat, build, rid, false); ok {
		data := map[string]interface{}{}
		returnJSON(c, data, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.DynamicList(c, plat, rid, pull, ctime, mid, time.Now())
		returnJSON(c, data, nil)
	}
}

func regionDynamicChild(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	buildStr := params.Get("build")
	tagIDStr := params.Get("tag_id")
	// check params
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	var tagID int
	if tagIDStr != "" {
		if tagID, err = strconv.Atoi(tagIDStr); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	// GetAudit
	if audit, ok := regionSvc.AuditChild(c, mobiApp, "", plat, build, rid, tagID); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.DynamicChild(c, plat, rid, tagID, build, mid, mobiApp, time.Now())
		returnJSON(c, data, nil)
	}
}

func regionDynamicChildList(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	ridStr := params.Get("rid")
	pullStr := params.Get("pull")
	ctimeStr := params.Get("ctime")
	buildStr := params.Get("build")
	tagIDStr := params.Get("tag_id")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		log.Error("build(%s) error(%v)", buildStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		log.Error("ridStr(%s) error(%v)", ridStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		log.Error("pullStr(%s) error(%v)", pullStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// ctime
	ctime, err := strconv.ParseInt(ctimeStr, 10, 64)
	if err != nil || ctime < 0 {
		ctime = 0
	}
	var tagID int
	if tagIDStr != "" {
		if tagID, err = strconv.Atoi(tagIDStr); err != nil {
			log.Error("tagId(%s) error(%v)", tagID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	// GetAudit
	if audit, ok := regionSvc.AuditChild(c, mobiApp, "", plat, build, rid, tagID); ok {
		returnJSON(c, audit, nil)
	} else {
		mobiApp = model.MobiAPPBuleChange(mobiApp)
		data := regionSvc.DynamicListChild(c, plat, rid, tagID, build, pull, ctime, mid, time.Now())
		returnJSON(c, data, nil)
	}
}
