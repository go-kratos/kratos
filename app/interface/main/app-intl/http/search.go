package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-intl/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const _keyWordLen = 50

func searchAll(c *bm.Context) {
	var (
		build  int
		mid    int64
		pn, ps int
		err    error
	)
	params := c.Request.Form
	header := c.Request.Header
	// params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ridStr := params.Get("rid")
	keyword := params.Get("keyword")
	highlightStr := params.Get("highlight")
	lang := params.Get("lang")
	duration := params.Get("duration")
	order := params.Get("order")
	filtered := params.Get("filtered")
	platform := params.Get("platform")
	zoneidStr := params.Get("zoneid")
	fromSource := params.Get("from_source")
	recommend := params.Get("recommend")
	// header
	buvid := header.Get("Buvid")
	// check params
	if keyword == "" || len([]rune(keyword)) > _keyWordLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	zoneid, _ := strconv.ParseInt(zoneidStr, 10, 64)
	rid, _ := strconv.Atoi(ridStr)
	highlight, _ := strconv.Atoi(highlightStr)
	// page and size
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	switch order {
	case "default", "":
		order = "totalrank"
	case "view":
		order = "click"
	case "danmaku":
		order = "dm"
	}
	if duration == "" {
		duration = "0"
	}
	if recommend == "" || recommend != "1" {
		recommend = "0"
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(searchSvc.Search(c, mid, zoneid, mobiApp, device, platform, buvid, keyword, duration, order, filtered, lang, fromSource, recommend, plat, rid, highlight, build, pn, ps, time.Now()))
}

func searchByType(c *bm.Context) {
	var (
		build  int
		mid    int64
		pn, ps int
		typeV  string
		err    error
	)
	params := c.Request.Form
	header := c.Request.Header
	// params
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	sType := params.Get("type")
	keyword := params.Get("keyword")
	filtered := params.Get("filtered")
	zoneidStr := params.Get("zoneid")
	order := params.Get("order")
	platform := params.Get("platform")
	highlightStr := params.Get("highlight")
	categoryIDStr := params.Get("category_id")
	userTypeStr := params.Get("user_type")
	orderSortStr := params.Get("order_sort")
	// header
	buvid := header.Get("Buvid")
	if keyword == "" || len([]rune(keyword)) > _keyWordLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	userType, _ := strconv.Atoi(userTypeStr)
	orderSort, _ := strconv.Atoi(orderSortStr)
	zoneid, _ := strconv.ParseInt(zoneidStr, 10, 64)
	categoryID, _ := strconv.Atoi(categoryIDStr)
	highlight, _ := strconv.Atoi(highlightStr)
	// page and size
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	switch sType {
	case "2":
		typeV = "upper"
	case "6":
		typeV = "article"
	case "7":
		typeV = "season2"
	case "8":
		typeV = "movie2"
	case "9":
		typeV = "tag"
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(searchSvc.SearchByType(c, mid, zoneid, mobiApp, device, platform, buvid, typeV, keyword, filtered, order, plat, build, highlight, categoryID, userType, orderSort, pn, ps, time.Now()))
}

// suggest3 search suggest data from newest api.
func suggest3(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	mobiApp := params.Get("mobi_app")
	term := params.Get("keyword")
	highlight, _ := strconv.Atoi(params.Get("highlight"))
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	platform := params.Get("platform")
	c.JSON(searchSvc.Suggest3(c, mid, platform, buvid, term, build, highlight, mobiApp, time.Now()), nil)
}
