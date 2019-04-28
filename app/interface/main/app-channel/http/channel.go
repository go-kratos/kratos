package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/app-channel/model"
	"go-common/app/interface/main/app-channel/model/channel"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_headerBuvid = "Buvid"
)

func index(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		cidInt int64
		header = c.Request.Header
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	idxStr := params.Get("idx")
	pullStr := params.Get("pull")
	loginEventStr := params.Get("login_event")
	cidStr := params.Get("channel_id")
	displayIDStr := params.Get("display_id")
	cname := params.Get("channel_name")
	qn, _ := strconv.Atoi(params.Get("qn"))
	fnver, _ := strconv.Atoi(params.Get("fnver"))
	fnval, _ := strconv.Atoi(params.Get("fnval"))
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	displayID, err := strconv.Atoi(displayIDStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get(_headerBuvid)
	// disid := header.Get(_headerDisplayID)
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		pull = true
	}
	loginEvent, err := strconv.Atoi(loginEventStr)
	if err != nil {
		loginEvent = 0
	}
	idx, err := strconv.ParseInt(idxStr, 10, 64)
	if err != nil || idx < 0 {
		idx = 0
	}
	if cidInt, _ = strconv.ParseInt(cidStr, 10, 64); cidInt == 0 && cname == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	now := time.Now()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data, err := channelSvc.Index(c, mid, cidInt, idx, plat, mobiApp, device, buvid, cname, build, loginEvent, displayID, qn, fnver, fnval, pull, now)
	c.JSON(data, err)
}

func index2(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		cidInt int64
		header = c.Request.Header
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	idxStr := params.Get("idx")
	pullStr := params.Get("pull")
	loginEventStr := params.Get("login_event")
	cidStr := params.Get("channel_id")
	displayIDStr := params.Get("display_id")
	cname := params.Get("channel_name")
	qn, _ := strconv.Atoi(params.Get("qn"))
	fnver, _ := strconv.Atoi(params.Get("fnver"))
	fnval, _ := strconv.Atoi(params.Get("fnval"))
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	displayID, err := strconv.Atoi(displayIDStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get(_headerBuvid)
	// disid := header.Get(_headerDisplayID)
	pull, err := strconv.ParseBool(pullStr)
	if err != nil {
		pull = true
	}
	loginEvent, err := strconv.Atoi(loginEventStr)
	if err != nil {
		loginEvent = 0
	}
	idx, err := strconv.ParseInt(idxStr, 10, 64)
	if err != nil || idx < 0 {
		idx = 0
	}
	if cidInt, _ = strconv.ParseInt(cidStr, 10, 64); cidInt == 0 && cname == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	now := time.Now()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data, err := channelSvc.Index2(c, mid, cidInt, idx, plat, mobiApp, device, buvid, cname, metadata.String(c, metadata.RemoteIP), build, loginEvent, displayID, qn, fnver, fnval, pull, now)
	c.JSON(data, err)
}

func tab(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
		cidInt int64
	)
	cidStr := params.Get("channel_id")
	cname := params.Get("channel_name")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if cidInt, _ = strconv.ParseInt(cidStr, 10, 64); cidInt == 0 && cname == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	plat := model.Plat(mobiApp, device)
	c.JSON(channelSvc.Tab(c, cidInt, mid, cname, plat))
}

func subscribeAdd(c *bm.Context) {
	var (
		params  = c.Request.Form
		mid     int64
		cidInt  int64
		err     error
		fromInt int
		now     = time.Now()
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	cidStr := params.Get("channel_id")
	fromStr := params.Get("from")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if cidInt, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fromInt, _ = strconv.Atoi(fromStr)
	c.JSON(nil, channelSvc.SubscribeAdd(c, mid, cidInt, now))
	channelSvc.OperationInfoc(mobiApp, device, "add", build, fromInt, cidInt, mid, now)
}

func subscribeCancel(c *bm.Context) {
	var (
		params  = c.Request.Form
		mid     int64
		cidInt  int64
		err     error
		fromInt int
		now     = time.Now()
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	cidStr := params.Get("channel_id")
	fromStr := params.Get("from")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if cidInt, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fromInt, _ = strconv.Atoi(fromStr)
	c.JSON(nil, channelSvc.SubscribeCancel(c, mid, cidInt, now))
	channelSvc.OperationInfoc(mobiApp, device, "cannel", build, fromInt, cidInt, mid, now)
}

func subscribeUpdate(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
	)
	cidStr := params.Get("channel_ids")
	if cidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(nil, channelSvc.SubscribeUpdate(c, mid, cidStr))
}

func list(c *bm.Context) {
	param := &channel.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.MID = midInter.(int64)
	}
	plat := model.Plat(param.MobiApp, param.Device)
	limit := 18 //频道聚合页需要展示最多18个我的订阅
	c.JSON(channelSvc.List(c, param.MID, plat, param.Build, limit, param.Ver, param.MobiApp, param.Device, param.Lang))
}

func subscribe(c *bm.Context) {
	param := &channel.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.MID = midInter.(int64)
	}
	c.JSON(channelSvc.Subscribe(c, param.MID, 0))
}

func discover(c *bm.Context) {
	param := &channel.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.MID = midInter.(int64)
	}
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(channelSvc.Discover(c, param.ID, param.MID, plat))
}

func category(c *bm.Context) {
	param := &channel.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(channelSvc.Category(c, plat))
}

func square(c *bm.Context) {
	param := &channel.ParamSquare{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.MID = midInter.(int64)
	}
	buvid := c.Request.Header.Get(_headerBuvid)
	plat := model.Plat(param.MobiApp, param.Device)
	c.JSON(channelSvc.Square(c, param.MID, plat, param.Build, param.LoginEvent, param.MobiApp, param.Device, param.Lang, buvid))
}

func mysub(c *bm.Context) {
	param := &channel.Param{}
	if err := c.Bind(param); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		param.MID = midInter.(int64)
	}
	c.JSON(channelSvc.Mysub(c, param.MID, 0))
}

func tablist(c *bm.Context) {
	var (
		params = c.Request.Form
		mid    int64
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	cidStr := params.Get("channel_id")
	displayIDStr := params.Get("display_id")
	cname := params.Get("channel_name")
	fromStr := params.Get("display_id")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	displayID, err := strconv.Atoi(displayIDStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	cidInt, _ := strconv.ParseInt(cidStr, 10, 64)
	if cidInt == 0 && cname == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fromInt, err := strconv.Atoi(fromStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	now := time.Now()
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	data, err := channelSvc.TabList(c, cidInt, mid, cname, mobiApp, displayID, build, fromInt, plat, now)
	c.JSON(data, err)
}
