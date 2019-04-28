package http

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/view"
	resource "go-common/app/service/main/resource/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_viewPath     = "/x/v2/view"
	_viewPagePath = "/x/v2/view/page"
	_headerBuvid  = "Buvid"
)

var (
	_rate    = []int64{564, 1328, 2592, 4192}
	_formats = map[int]string{1: "mp4", 2: "hdmp4", 3: "flv", 4: "flv"}
	_dislike = []*view.Dislike{
		{
			ID:   5,
			Name: "标题党/封面党",
		},
		{
			ID:   6,
			Name: "内容质量差",
		},
		{
			ID:   7,
			Name: "内容/封面令人不适",
		},
		{
			ID:   8,
			Name: "营销广告",
		},
	}
)

// viewIndex view handler
func viewIndex(c *bm.Context) {
	var (
		mid, aid, movieID int64
		err               error
		parentMode        int
	)
	params := c.Request.Form
	header := c.Request.Header
	// get params
	aidStr := params.Get("aid")
	movieidStr := params.Get("movie_id")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ak := params.Get("access_key")
	buildStr := params.Get("build")
	from := params.Get("from")
	trackid := params.Get("trackid")
	network := params.Get("network")
	adExtra := params.Get("ad_extra")
	parentModeStr := params.Get("parent_mode")
	parentMode, _ = strconv.Atoi(parentModeStr)
	qnStr := params.Get("qn")
	qn, _ := strconv.Atoi(qnStr)
	fnverStr := params.Get("fnver")
	fnver, _ := strconv.Atoi(fnverStr)
	fnvalStr := params.Get("fnval")
	fnval, _ := strconv.Atoi(fnvalStr)
	forceHost, _ := strconv.Atoi(params.Get("force_host"))
	spmid := params.Get("spmid")
	fromSpmid := params.Get("from_spmid")
	// check params
	if aidStr == "" && movieidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aidStr != "" && aidStr != "0" {
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else if movieidStr != "" && movieidStr != "0" {
		if movieID, err = strconv.ParseInt(movieidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if aid < 1 && movieID < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get("Buvid")
	disid := header.Get("Display-ID")
	cdnIP := header.Get("X-Cache-Server-Addr")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	autoplay, _ := strconv.Atoi(params.Get("autoplay"))
	now := time.Now()
	// view
	ip := metadata.String(c, metadata.RemoteIP)
	viewSvr.ViewInfoc(mid, int(plat), trackid, aidStr, ip, _viewPath, buildStr, buvid, disid, from, now, err, autoplay, spmid, fromSpmid)
	data, err := viewSvr.View(c, mid, aid, movieID, plat, build, qn, fnver, fnval, forceHost, parentMode, ak, mobiApp, device, buvid, cdnIP, network, adExtra, from, now)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data.Dislikes = _dislike
	compMeta(data, build)
	if mid == 0 && data.Duration > 360 {
		data.Paster, _ = viewSvr.Paster(c, plat, resource.VdoAdsTypeNologin, aidStr, strconv.Itoa(int(data.TypeID)), buvid)
	}
	c.JSON(data, nil)
	viewSvr.RelateInfoc(mid, aid, int(plat), trackid, buildStr, buvid, disid, ip, _viewPath, data.ReturnCode, data.UserFeature, from, data.Relates, now, data.IsRec)
}

// viewPage view page handler.
func viewPage(c *bm.Context) {
	var (
		mid, aid int64
		build    int
		err      error
	)
	params := c.Request.Form
	header := c.Request.Header
	// get params
	aidStr := params.Get("aid")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	ak := params.Get("access_key")
	buildStr := params.Get("build")
	from := params.Get("from")
	trackid := params.Get("trackid")
	spmid := params.Get("spmid")
	fromSpmid := params.Get("from_spmid")
	// check params
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if build, err = strconv.Atoi(buildStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get("Buvid")
	disid := header.Get("Display-ID")
	cdnIP := header.Get("X-Cache-Server-Addr")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	autoplay, _ := strconv.Atoi(params.Get("autoplay"))
	ip := metadata.String(c, metadata.RemoteIP)
	now := time.Now()
	// view page
	viewSvr.ViewInfoc(mid, int(plat), trackid, aidStr, ip, _viewPagePath, buildStr, buvid, disid, from, now, err, autoplay, spmid, fromSpmid)
	data, err := viewSvr.ViewPage(c, mid, aid, 0, plat, build, ak, mobiApp, device, cdnIP, false, now)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	compMeta(data, build)
	c.JSON(data, nil)
}

// videoShot video shot .
func videoShot(c *bm.Context) {
	var aid, cid int64
	params := c.Request.Form
	if aid, _ = strconv.ParseInt(params.Get("aid"), 10, 64); aid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid, _ = strconv.ParseInt(params.Get("cid"), 10, 64); cid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(viewSvr.Shot(c, aid, cid))
}

// addShare add a share.
func addShare(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	buvid := c.Request.Header.Get("Buvid")
	from := params.Get("from")
	build := params.Get("build")
	// check
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	share, isReport, upID, err := viewSvr.AddShare(c, aid, mid, metadata.String(c, metadata.RemoteIP))
	c.JSON(struct {
		Aid   int64 `json:"aid"`
		Count int   `json:"count"`
	}{aid, share}, err)
	if err != nil {
		return
	}
	// for ai big data
	sendUserAct(mid, mobiApp, buvid, from, build, aidStr, "av", "share")
	if collector != nil && isReport {
		collector.InfoAntiCheat2(c, strconv.FormatInt(upID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), "av", "share", "")
	}
}

func addCoin(c *bm.Context) {
	var (
		mid, aid, upID int64
		avType         int64
		actLike        = "cointolike"
		actCoin        = "coin"
	)
	params := c.Request.Form
	// check
	mobiApp := params.Get("mobi_app")
	buvid := c.Request.Header.Get("Buvid")
	from := params.Get("from")
	build := params.Get("build")
	aidStr := params.Get("aid")
	upIDStr := params.Get("upid")
	ak := params.Get("access_key")
	multiStr := params.Get("multiply")
	selectLikeStr := params.Get("select_like")
	selectLike, _ := strconv.Atoi(selectLikeStr)
	aid, _ = strconv.ParseInt(aidStr, 10, 64)
	upID, _ = strconv.ParseInt(upIDStr, 10, 64)
	multiply, err := strconv.ParseInt(multiStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if avType, _ = strconv.ParseInt(params.Get("avtype"), 10, 64); avType == 0 {
		avType = 1
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// for ai big data
	sendUserAct(mid, mobiApp, buvid, from, build, aidStr, "av", actCoin)
	if selectLike == 1 {
		sendUserAct(mid, mobiApp, buvid, from, build, aidStr, "av", actLike)
	}
	prompt, like, err := viewSvr.AddCoin(c, aid, mid, upID, avType, multiply, ak, selectLike)
	c.JSON(struct {
		Prompt bool `json:"prompt,omitempty"`
		Like   bool `json:"like"`
	}{Prompt: prompt, Like: like}, err)
	if err != nil {
		return
	}
	if collector != nil {
		var itemType = "archive"
		if avType == 2 {
			itemType = "article"
		}
		collector.InfoAntiCheat2(c, strconv.FormatInt(upID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), itemType, actCoin, "")
		if like {
			collector.InfoAntiCheat2(c, strconv.FormatInt(upID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), "av", actLike, "")
		}
	}
}

func like(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	// check
	mobiApp := params.Get("mobi_app")
	buvid := c.Request.Header.Get("Buvid")
	from := params.Get("from")
	build := params.Get("build")
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	like, err := strconv.Atoi(params.Get("like"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if like != 0 && like != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	upperID, toast, err := viewSvr.Like(c, aid, mid, int8(like))
	c.JSON(struct {
		Toast string `json:"toast"`
	}{Toast: toast}, err)
	if err != nil {
		return
	}
	// for ai big data
	action := "like"
	if like == 1 {
		action = "like_cancel"
	}
	sendUserAct(mid, mobiApp, buvid, from, build, aidStr, "av", action)
	if collector != nil {
		collector.InfoAntiCheat2(c, strconv.FormatInt(upperID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), "av", action, "")
	}
}

func dislike(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	buvid := c.Request.Header.Get("Buvid")
	from := params.Get("from")
	build := params.Get("build")
	// check
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dislike, err := strconv.Atoi(params.Get("dislike"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if dislike != 0 && dislike != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	upperID, err := viewSvr.Dislike(c, aid, mid, int8(dislike))
	c.JSON(nil, err)
	if err != nil {
		return
	}
	// for ai big data
	action := "dislike"
	if dislike == 1 {
		action = "dislike_cancel"
	}
	sendUserAct(mid, mobiApp, buvid, from, build, aidStr, "av", action)
	if collector != nil {
		collector.InfoAntiCheat2(c, strconv.FormatInt(upperID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), "av", action, "")
	}
}

func addFav(c *bm.Context) {
	var (
		mid  int64
		vmid int64
		fid  int64
		aid  int64
		err  error
	)
	params := c.Request.Form
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	// check
	if aid, err = strconv.ParseInt(params.Get("aid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, err = strconv.ParseInt(params.Get("vmid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(params.Get("fid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ak := params.Get("access_key")
	prompt, err := viewSvr.AddFav(c, mid, vmid, []int64{fid}, aid, ak)
	c.JSON(struct {
		Prompt bool `json:"prompt,omitempty"`
	}{Prompt: prompt}, err)
}

// adDislike ad dislike
func adDislike(c *bm.Context) {
	var mid int64
	params := c.Request.Form
	header := c.Request.Header
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	gt := params.Get("goto")
	id, _ := strconv.ParseInt(params.Get("id"), 10, 64)
	reasonID, _ := strconv.ParseInt(params.Get("reason_id"), 10, 64)
	cmreasonID, _ := strconv.ParseInt(params.Get("cm_reason_id"), 10, 64)
	rid, _ := strconv.ParseInt(params.Get("rid"), 10, 64)
	tagID, _ := strconv.ParseInt(params.Get("tag_id"), 10, 64)
	adcb := params.Get("ad_cb")
	buvid := header.Get("Buvid")
	if buvid == "" && mid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// for ad data
	err := dislikePub.Send(context.TODO(), strconv.FormatInt(mid, 10), &cmDislike{
		ID:         id,
		Buvid:      buvid,
		Goto:       gt,
		Mid:        mid,
		ReasonID:   reasonID,
		CMReasonID: cmreasonID,
		UpperID:    0,
		Rid:        rid,
		TagID:      tagID,
		ADCB:       adcb,
	})
	c.JSON(nil, err)
}

func compMeta(v *view.View, build int) {
	if build == 5800 || build == 508000 {
		for _, vp := range v.Pages {
			metas := make([]*view.Meta, 0, 4)
			for i, r := range _rate {
				meta := &view.Meta{
					Quality: i + 1,
					Size:    int64(float64(r*v.Duration) * 1.1 / 8.0),
					Format:  _formats[i+1],
				}
				metas = append(metas, meta)
			}
			vp.Metas = metas
		}
	}
}

// vipPlayURL get big-member token.
func vipPlayURL(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	params := c.Request.Form
	aid, err := strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	cid, _ := strconv.ParseInt(params.Get("cid"), 10, 64)
	if aid == 0 || cid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(viewSvr.VipPlayURL(c, aid, cid, mid))
}

// follow check if follow.
func follow(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	params := c.Request.Form
	vmid, err := strconv.ParseInt(params.Get("vmid"), 10, 64)
	if err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(viewSvr.Follow(c, vmid, mid))
}

func upperRecmd(c *bm.Context) {
	var (
		mid    int64
		vmid   int64
		header = c.Request.Header
		params = c.Request.Form
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if vmid, err = strconv.ParseInt(params.Get("vmid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	plat := model.Plat(mobiApp, device)
	buvid := header.Get(_headerBuvid)
	data, err := viewSvr.UpperRecmd(c, plat, platform, mobiApp, device, buvid, build, mid, vmid)
	c.JSON(data, err)
}

func likeTriple(c *bm.Context) {
	var (
		mid       int64
		actTriple = "triplelike"
	)
	params := &view.TripleParam{}
	if err := c.Bind(params); err != nil {
		return
	}
	buvid := c.Request.Header.Get("Buvid")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	// for ai big data
	sendUserAct(mid, params.MobiApp, buvid, params.From, params.Build, strconv.FormatInt(params.AID, 10), "av", actTriple)
	triple, err := viewSvr.LikeTriple(c, params.AID, mid, params.Ak)
	c.JSON(triple, err)
	if err != nil {
		return
	}
	if collector != nil && triple.Anticheat {
		collector.InfoAntiCheat2(c, strconv.FormatInt(triple.UpID, 10), strconv.FormatInt(params.AID, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(params.AID, 10), "av", actTriple, "")
	}
}

func sendUserAct(mid int64, mobiApp, buvid, from, build, itemID, itemType, action string) {
	userActPub.Send(context.TODO(), strconv.FormatInt(mid, 10), &userAct{
		Client:   mobiApp,
		Buvid:    buvid,
		Mid:      mid,
		Time:     time.Now().Unix(),
		From:     from,
		Build:    build,
		ItemID:   itemID,
		ItemType: itemType,
		Action:   action,
	})
}
