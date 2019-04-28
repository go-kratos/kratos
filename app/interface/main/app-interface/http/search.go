package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/search"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid = "Buvid"
	_keyWordLen  = 50
)

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
	parent := params.Get("parent_mode")
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
	isQuery, _ := strconv.Atoi(params.Get("is_org_query"))
	plat := model.Plat(mobiApp, device)
	c.JSON(srcSvr.Search(c, mid, zoneid, mobiApp, device, platform, buvid, keyword, duration, order, filtered, lang, fromSource, recommend, parent, plat, rid, highlight, build, pn, ps, isQuery, checkOld(plat, build), time.Now()))
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
	case "1":
		typeV = "season"
	case "2":
		typeV = "upper"
	case "3":
		typeV = "movie"
	case "4":
		typeV = "live_room"
	case "5":
		typeV = "live_user"
	case "6":
		typeV = "article"
	case "7":
		typeV = "season2"
	case "8":
		typeV = "movie2"
	case "9":
		typeV = "tag"
	case "10":
		typeV = "video"
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(srcSvr.SearchByType(c, mid, zoneid, mobiApp, device, platform, buvid, typeV, keyword, filtered, order, plat, build, highlight, categoryID, userType, orderSort, pn, ps, checkOld(plat, build), time.Now()))
}

func searchLive(c *bm.Context) {
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
	order := params.Get("order")
	platform := params.Get("platform")
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
	plat := model.Plat(mobiApp, device)
	switch sType {
	case "4":
		if (model.IsAndroid(plat) && build > search.SearchLiveAllAndroid) || (model.IsIPhone(plat) && build > search.SearchLiveAllIOS) || model.IsIPad(plat) || model.IsIPhoneB(plat) {
			typeV = "live_all"
		} else {
			typeV = "live_room"
		}
	case "5":
		typeV = "live_user"
	}
	if typeV == "live_all" {
		c.JSON(srcSvr.SearchLiveAll(c, mid, mobiApp, platform, buvid, device, typeV, keyword, order, build, pn, ps))
	} else {
		c.JSON(srcSvr.SearchLive(c, mid, mobiApp, platform, buvid, device, typeV, keyword, order, build, pn, ps))
	}
}

// ip string, limit int
func hotSearch(c *bm.Context) {
	var (
		mid   int64
		build int
		limit int
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	buvid := header.Get("Buvid")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if limit, err = strconv.Atoi(params.Get("limit")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srcSvr.HotSearch(c, buvid, mid, build, limit, mobiApp, device, platform, time.Now()), nil)
}

// suggest search suggest data.
func suggest(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	term := params.Get("keyword")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	c.JSON(srcSvr.Suggest(c, mid, buvid, term, build, mobiApp, device, time.Now()), nil)
}

// suggest2 search suggest data from new api.
func suggest2(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	mobiApp := params.Get("mobi_app")
	term := params.Get("keyword")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	buvid := header.Get(_headerBuvid)
	platform := params.Get("platform")
	c.JSON(srcSvr.Suggest2(c, mid, platform, buvid, term, build, mobiApp, time.Now()), nil)
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
	device := params.Get("device")
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
	c.JSON(srcSvr.Suggest3(c, mid, platform, buvid, term, device, build, highlight, mobiApp, time.Now()), nil)
}

func checkOld(plat int8, build int) bool {
	const (
		_oldAndroid = 513000
		_oldIphone  = 6060
	)
	return (model.IsIPhone(plat) && build <= _oldIphone) || (model.IsAndroid(plat) && build <= _oldAndroid)
}

func searchUser(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	keyword := params.Get("keyword")
	filtered := params.Get("filtered")
	order := params.Get("order")
	fromSource := params.Get("from_source")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	userType, _ := strconv.Atoi(params.Get("user_type"))
	highlight, _ := strconv.Atoi(params.Get("highlight"))
	if order == "" {
		order = "totalrank"
	}
	if order != "totalrank" && order != "fans" && order != "level" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	orderSort, _ := strconv.Atoi(params.Get("order_sort"))
	if orderSort != 1 {
		orderSort = 0
	}
	if fromSource == "" {
		fromSource = "dynamic_uname"
	}
	if fromSource != "dynamic_uname" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, _ := strconv.Atoi(params.Get("pn"))
	if pn < 1 {
		pn = 1
	}
	ps, _ := strconv.Atoi(params.Get("ps"))
	if ps < 1 || ps > 20 {
		ps = 20
	}
	buvid := header.Get(_headerBuvid)
	c.JSON(srcSvr.User(c, mid, buvid, mobiApp, device, platform, keyword, filtered, order, fromSource, highlight, build, userType, orderSort, pn, ps, time.Now()), nil)
}

func recommend(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	platform := params.Get("platform")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	from, _ := strconv.Atoi(params.Get("from"))
	show, _ := strconv.Atoi(params.Get("show"))
	buvid := header.Get("Buvid")
	c.JSON(srcSvr.Recommend(c, mid, build, from, show, buvid, platform, mobiApp, device))
}

func defaultWords(c *bm.Context) {
	var (
		build int
		mid   int64
		err   error
	)
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	platform := params.Get("platform")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	from, _ := strconv.Atoi(params.Get("from"))
	buvid := header.Get("Buvid")
	c.JSON(srcSvr.DefaultWords(c, mid, build, from, buvid, platform, mobiApp, device))
}

func recommendNoResult(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		build  int
		mid    int64
		err    error
	)
	platform := params.Get("platform")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get("Buvid")
	keyword := params.Get("keyword")
	if keyword == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, _ := strconv.Atoi(params.Get("pn"))
	if pn < 1 {
		pn = 1
	}
	ps, _ := strconv.Atoi(params.Get("ps"))
	if ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(srcSvr.RecommendNoResult(c, platform, mobiApp, device, buvid, keyword, build, pn, ps, mid))
}

func resource(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		build  int
		mid    int64
		err    error
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	network := params.Get("network")
	buvid := header.Get("Buvid")
	adExtra := params.Get("ad_extra")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(srcSvr.Resource(c, mobiApp, device, network, buvid, adExtra, build, plat, mid))
}

func recommendPre(c *bm.Context) {
	var (
		params = c.Request.Form
		header = c.Request.Header
		build  int
		mid    int64
		err    error
	)
	platform := params.Get("platform")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	buvid := header.Get("Buvid")
	ps, _ := strconv.Atoi(params.Get("ps"))
	if ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(srcSvr.RecommendPre(c, platform, mobiApp, device, buvid, build, ps, mid))
}

func searchEpisodes(c *bm.Context) {
	var (
		params    = c.Request.Form
		mid, ssID int64
		err       error
	)
	if ssID, err = strconv.ParseInt(params.Get("season_id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if ssID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srcSvr.SearchEpisodes(c, mid, ssID))
}
