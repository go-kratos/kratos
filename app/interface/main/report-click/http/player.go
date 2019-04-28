package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/report-click/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// reportPlayer write the archive data.
func reportPlayer(c *bm.Context) {
	params := c.Request.Form
	header := c.Request.Header
	buvid := header.Get("Buvid")
	displayid := header.Get("Display-ID")
	ts := params.Get("ts")
	aid := params.Get("aid")
	cid := params.Get("cid")
	playedTime := params.Get("played_time")
	mid := params.Get("mid")
	moAp := params.Get("mobi_app")
	typeID := params.Get("type")
	subType := params.Get("sub_type")
	sid := params.Get("sid")
	epid := params.Get("epid")
	infocRealTime.Info(ts, buvid, displayid, mid, aid, cid, playedTime, strconv.FormatInt(time.Now().Unix(), 10), "1", moAp, "", typeID, subType, sid, epid, "")
	c.JSON(nil, nil)
}

// reportHeartbeat write the archive data.
func reportHeartbeat(c *bm.Context) {
	params := c.Request.Form
	header := c.Request.Header
	buvid := header.Get("Buvid")
	displayid := header.Get("Display-ID")
	sts := params.Get("start_ts")
	aid := params.Get("aid")
	if aid == "" {
		aid = params.Get("avid")
	}
	cid := params.Get("cid")
	playedTime := params.Get("played_time")
	mid := params.Get("mid")
	moAp := params.Get("mobi_app")
	typeID := params.Get("type")
	subType := params.Get("sub_type")
	sid := params.Get("sid")
	epid := params.Get("epid")
	playType := params.Get("play_type")
	if playType == "" {
		playType = params.Get("playtype")
	}
	infocRealTime.Info(sts, buvid, displayid, mid, aid, cid, playedTime, strconv.FormatInt(time.Now().Unix(), 10), "2", moAp, "", typeID, subType, sid, epid, playType)
	c.JSON(nil, nil)
}

func heartbeatMobile(c *bm.Context) {
	params := c.Request.Form
	header := c.Request.Header
	sts := params.Get("start_ts")
	build := params.Get("build")
	buvid := header.Get("Buvid")
	mobileApp := params.Get("mobi_app")
	device := params.Get("device")
	platform := params.Get("platform")
	session := params.Get("session")
	mid := params.Get("mid")
	aid := params.Get("aid")
	cid := params.Get("cid")
	sid := params.Get("sid")
	epid := params.Get("epid")
	tp := params.Get("type")
	subType := params.Get("sub_type")
	quality := params.Get("quality")
	totalTime := params.Get("total_time")
	pausedTime := params.Get("paused_time")
	playedTime := params.Get("played_time")
	videoDuration := params.Get("video_duration")
	playType := params.Get("play_type")
	networkType := params.Get("network_type")
	playProgressTimeLast := params.Get("last_play_progress_time")
	playProgressTimeMax := params.Get("max_play_progress_time")
	playMode := params.Get("play_mode")
	from := params.Get("from")
	epidStatus := params.Get("epid_status")
	playStatus := params.Get("play_status")
	userStatus := params.Get("user_status")
	actualPlayedTime := params.Get("actual_played_time")
	autoPlay := params.Get("auto_play")
	detailPlayTime := params.Get("detail_play_time")
	listPlayTime := params.Get("list_play_time")
	userAgent := c.Request.Header.Get("User-Agent")
	ts, err := strconv.ParseInt(sts, 10, 64)
	if err != nil || ts <= 0 {
		ts = time.Now().Unix()
		sts = strconv.FormatInt(ts, 10)
	}
	// NOTE /x/report//heartbeat/mobile auto_play = 2 ===> /x/report/click/android2 & ios
	// (自动播放的上报>> 2:天马feed流inline) 播放时长转成播放点击
	autoPlayInt, _ := strconv.ParseInt(autoPlay, 10, 64)
	fromInt, _ := strconv.ParseInt(from, 10, 64)
	videoDurInt, _ := strconv.ParseInt(videoDuration, 10, 64)
	playedTimeInt, _ := strconv.ParseInt(playedTime, 10, 64)
	var needCompens bool
	if (autoPlayInt == 2 || autoPlayInt == 1) && fromMap[fromInt] {
		userAgent = userAgent + " (inline_play_to_view)" // change from auto_play to inline_play_heartbeat, then to inline_play_to_view
		needCompens = true
	}
	if (autoPlayInt == 2 || autoPlayInt == 1) && fromInlineMap[fromInt] && playedTimeInt >= inlineDuration && (videoDurInt >= playedTimeInt) {
		userAgent += " (played_time_enough)" // new logic, if inline play more than 10s, count it also
		needCompens = true
	}
	if needCompens {
		var cookieSid, plat string
		if ck, err := c.Request.Cookie("sid"); err == nil {
			cookieSid = ck.Value
		}
		ip := metadata.String(c, metadata.RemoteIP)
		switch platform {
		case "android":
			plat = _platAndroid
		case "ios":
			plat = _platIos
		}
		clickSvr.Play(c, plat, aid, cid, params.Get("part"), mid, params.Get("lv"),
			"0", sts, buvid, ip, userAgent, buvid,
			cookieSid, c.Request.Header.Get("Referer"), tp,
			subType, sid, epid, playMode, platform, device, mobileApp, autoPlay, session)

		log.Warn("plat:%s,aid:%s,cid:%s,part:%s,mid:%s,lv:%s,0:%s,sts:%s,buvid:%s,ip:%s,userAgent:%s,"+
			"buvid:%s,cookieSid:%s,Referer:%s,tp:%s,subType:%s,sid:%s,epid:%s,playMode:%s,"+
			"platform:%s,device:%s,mobileApp:%s,autoPlay:%s,session:%s",
			plat, aid, cid, params.Get("part"), mid, params.Get("lv"),
			"0", sts, buvid, ip, userAgent, buvid,
			cookieSid, c.Request.Header.Get("Referer"), tp,
			subType, sid, epid, playMode, platform, device, mobileApp, autoPlay, session)
	}
	buildInt, _ := strconv.ParseInt(build, 10, 64)
	clickSvr.SuccReport(c, &model.SuccReport{ // record the success by different build
		MobiApp: mobileApp,
		Build:   buildInt,
	})
	infocStatistics.Info(sts, build, buvid, mobileApp, platform, session, mid, aid, cid, sid,
		epid, tp, subType, quality, totalTime, pausedTime, playedTime, videoDuration,
		playType, networkType, playProgressTimeLast, playProgressTimeMax, playMode, device, from, epidStatus, playStatus, userStatus, actualPlayedTime, autoPlay, detailPlayTime, listPlayTime)
	data := make(map[string]interface{}, 1)
	data["ts"] = ts
	c.JSON(data, nil)
}

// webHeartbeat write the archive data.
func webHeartbeat(c *bm.Context) {
	var (
		buvid, mid, term string
		params           = c.Request.Form
	)
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	sts := params.Get("start_ts")
	aid := params.Get("aid")
	cid := params.Get("cid")
	pause := params.Get("pause")
	playType := params.Get("play_type")
	if playType == "" {
		playType = params.Get("playtype")
	}
	playedTime := params.Get("played_time")
	if midI, ok := c.Get("mid"); ok {
		mid = strconv.FormatInt(midI.(int64), 10)
	}
	tp := params.Get("type")
	subType := params.Get("sub_type")
	sid := params.Get("sid")
	if sid == "" {
		sid = params.Get("seasonID")
	}
	epid := params.Get("epid")
	dt := params.Get("dt")
	if dt == "7" {
		// count m.bilibili.com visits times.
		term = "h5"
	} else {
		dt = "2"
		term = "web"
	}
	realtime := params.Get("realtime")
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	infocRealTime.Info(sts, buvid, "", mid, aid, cid, playedTime, ts, dt, term, pause, tp, subType, sid, epid, playType)
	clickSvr.Report(c, playedTime, cid, tp, subType, realtime, aid, mid, sid, epid, dt, ts)
	c.JSON(nil, nil)
}
