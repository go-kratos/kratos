package http

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/space"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_crop = "@750w_250h_1c"
)

type userAct struct {
	Client   string `json:"client"`
	Buvid    string `json:"buvid"`
	Mid      int64  `json:"mid"`
	Time     int64  `json:"time"`
	From     string `json:"from"`
	Build    string `json:"build"`
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Action   string `json:"action"`
	ActionID string `json:"action_id"`
	Extra    string `json:"extra"`
}

func spaceAll(c *bm.Context) {
	var (
		mid    int64
		vmid   int64
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	platform := params.Get("platform")
	device := params.Get("device")
	buildStr := params.Get("build")
	name := params.Get("name")
	// check params
	if vmid, _ = strconv.ParseInt(params.Get("vmid"), 10, 64); vmid < 1 && name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInt, ok := c.Get("mid"); ok {
		mid = midInt.(int64)
	}
	if build, err = strconv.Atoi(buildStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	plat := model.Plat(mobiApp, device)
	space, err := spaceSvr.Space(c, mid, vmid, plat, build, pn, ps, platform, device, mobiApp, name, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if model.IsIPhone(plat) && space.Space != nil && space.Space.ImgURL != "" {
		space.Space.ImgURL = space.Space.ImgURL + _crop
	}
	space.Relation = compRealtion(space.Relation, mobiApp, build, device)
	c.JSON(space, nil)
	// for ai big data
	userActPub.Send(context.TODO(), strconv.FormatInt(mid, 10), &userAct{
		Client:   mobiApp,
		Buvid:    c.Request.Header.Get("Buvid"),
		Mid:      mid,
		Time:     time.Now().Unix(),
		From:     params.Get("from"),
		Build:    buildStr,
		ItemID:   space.Card.Mid,
		ItemType: "mid",
		Action:   "space",
	})
}

func upArchive(c *bm.Context) {
	var (
		pn, ps int
	)
	params := c.Request.Form
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.UpArcs(c, vmid, pn, ps, time.Now()), nil)
}

func myinfo(c *bm.Context) {
	mid, err := authSvc.AuthToken(c)
	if err != nil {
		shouldChangeError := false
		params := c.Request.Form
		mobiApp := params.Get("mobi_app")
		device := params.Get("device")
		build, _ := strconv.Atoi(params.Get("build"))
		plat := model.Plat(mobiApp, device)
		if model.IsIPhone(plat) && build > config.LoginBuild.Iphone {
			shouldChangeError = true
		}
		if shouldChangeError && ecode.Cause(err).Equal(ecode.NoLogin) {
			c.JSON(nil, ecode.UserLoginInvalid)
			return
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(accSvr.Myinfo(c, mid))
}

func mine(c *bm.Context) {
	params := &space.MineParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		params.Mid = midInter.(int64)
	}
	plat := model.Plat(params.MobiApp, params.Device)
	c.JSON(accSvr.Mine(c, params.Mid, params.Platform, params.Filtered, params.Build, plat))
}

func mineIpad(c *bm.Context) {
	params := &space.MineParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		params.Mid = midInter.(int64)
	}
	plat := model.Plat(params.MobiApp, params.Device)
	if model.IsIPad(plat) {
		plat = model.PlatIPad
	}
	c.JSON(accSvr.MineIpad(c, params.Mid, params.Platform, params.Filtered, params.Build, plat))
}

func upArticle(c *bm.Context) {
	var (
		pn, ps int
	)
	params := c.Request.Form
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.UpArticles(c, vmid, pn, ps), nil)
}

func contribute(c *bm.Context) {
	var (
		vmid   int64
		build  int
		pn, ps int
		err    error
	)
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	// check params
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, _ = strconv.ParseInt(params.Get("vmid"), 10, 64); vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(spaceSvr.Contribute(c, plat, build, vmid, pn, ps, time.Now()))
}

func contribution(c *bm.Context) {
	var (
		vmid         int64
		build        int
		maxID, minID int64
		size         int
		err          error
	)
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	maxIDStr := params.Get("max_id")
	minIDStr := params.Get("min_id")
	// check params
	if build, err = strconv.Atoi(params.Get("build")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, _ = strconv.ParseInt(params.Get("vmid"), 10, 64); vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if maxIDStr != "" {
		if maxID, err = strconv.ParseInt(maxIDStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if minIDStr != "" {
		if minID, err = strconv.ParseInt(minIDStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if size, _ = strconv.Atoi(params.Get("size")); size < 1 || size > 20 {
		size = 20
	}
	cursor, err := model.NewCursor(maxID, minID, size)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("%+v", err)
		return
	}
	plat := model.Plat(mobiApp, device)
	c.JSON(spaceSvr.Contribution(c, plat, build, vmid, cursor, time.Now()))
}

func upContribute(c *bm.Context) {
	var (
		vmid  int64
		attrs *space.Attrs
		items []*space.Item
		err   error
	)
	params := c.Request.Form
	vmidStr := params.Get("vmid")
	attrsStr := params.Get("attrs")
	itemsStr := params.Get("items")
	// check params
	if vmid, _ = strconv.ParseInt(vmidStr, 10, 64); vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(attrsStr), &attrs); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("json.Unmarshal(%s) error(%v)", attrsStr, err)
		return
	}
	if err = json.Unmarshal([]byte(itemsStr), &items); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("json.Unmarshal(%s) error(%v)", itemsStr, err)
		return
	}
	c.JSON(spaceSvr.AddContribute(c, vmid, attrs, items), nil)
}

func bangumi(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
	)
	params := c.Request.Form
	if midInt, ok := c.Get("mid"); ok {
		mid = midInt.(int64)
	}
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.Bangumi(c, mid, vmid, nil, pn, ps), nil)
}

func community(c *bm.Context) {
	var (
		pn, ps int
	)
	params := c.Request.Form
	ak := params.Get("access_key")
	platform := params.Get("platform")
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.Community(c, vmid, pn, ps, ak, platform), nil)
}

func coinArc(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
	)
	params := c.Request.Form
	if midInt, ok := c.Get("mid"); ok {
		mid = midInt.(int64)
	}
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.CoinArcs(c, mid, vmid, nil, pn, ps), nil)
}

func likeArc(c *bm.Context) {
	var (
		mid    int64
		pn, ps int
	)
	params := c.Request.Form
	if midInt, ok := c.Get("mid"); ok {
		mid = midInt.(int64)
	}
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.LikeArcs(c, mid, vmid, nil, pn, ps), nil)
}
func report(c *bm.Context) {
	params := c.Request.Form
	ak := params.Get("access_key")
	reason := params.Get("reason")
	mid, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	if mid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spaceSvr.Report(c, mid, reason, ak))
}

func clips(c *bm.Context) {
	var (
		pos, ps int
	)
	params := c.Request.Form
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pos, _ = strconv.Atoi(params.Get("offset")); pos < 0 {
		pos = 0
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.Clip(c, vmid, pos, ps), nil)
}

func albums(c *bm.Context) {
	var (
		pos, ps int
	)
	params := c.Request.Form
	// check params
	vmid, _ := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pos, _ = strconv.Atoi(params.Get("offset")); pos < 0 {
		pos = 0
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 20 {
		ps = 20
	}
	c.JSON(spaceSvr.Album(c, vmid, pos, ps), nil)
}

// checkPay pay movie or bangumi
func compRealtion(rel int, mobiApp string, build int, device string) (r int) {
	const (
		_upAndroid = 505000
		_banIphone = 5550
		_banIPad   = 10450
	)
	switch mobiApp {
	case "android", "android_G":
		if build < _upAndroid && rel == -1 {
			return -999
		}
	case "iphone", "iphone_G":
		if build < _banIphone && rel == -1 {
			return -999
		}
	case "ipad", "ipad_G":
		if build <= _banIPad && rel == -1 {
			return -999
		}
	}
	return rel
}
