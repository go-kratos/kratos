package http

import (
	"io/ioutil"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/feedback/conf"
	"go-common/app/interface/main/feedback/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"
)

const (
	_timeFromat = "2006-01-02 15:04:05"
)

// addReply
func addReply(c *bm.Context) {
	header := c.Request.Header
	params := c.Request.Form
	// params
	buvid := header.Get("Buvid")
	system := params.Get("system")
	version := params.Get("version")
	midStr := params.Get("mid")
	content := params.Get("content")
	imgURL := params.Get("img_url")
	logURL := params.Get("log_url")
	device := params.Get("device")
	channel := params.Get("channel")
	entrance := params.Get("entrance")
	netState := params.Get("net_state")
	netOperator := params.Get("net_operator")
	agencyArea := params.Get("agency_area")
	platform := params.Get("platform")
	browser := params.Get("browser")
	qq := params.Get("qq")
	mobiApp := params.Get("mobi_app")
	email := params.Get("email")
	tagStr := params.Get("tag_id")
	// check params
	if buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if version == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		mid int64
		err error
	)
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if content == "" && imgURL == "" && logURL == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tagID, _ := strconv.ParseInt(tagStr, 10, 64)
	c.JSON(feedbackSvr.AddReply(c, mid, tagID, buvid, system, version, mobiApp, filtered(content), imgURL, logURL, device, channel, entrance, netState, netOperator, agencyArea, platform, browser, qq, email, time.Now()))
}

func addWebReply(c *bm.Context) {
	params := c.Request.Form
	// system := params.Get("system")
	version := params.Get("version")
	midStr := params.Get("mid")
	sidStr := params.Get("session_id")
	content := params.Get("content")
	imgURL := params.Get("img_url")
	logURL := params.Get("log_url")
	// device := params.Get("device")
	buvid := params.Get("buvid")
	// channel := params.Get("channel")
	netState := params.Get("net_state")
	netOperator := params.Get("net_operator")
	agencyArea := params.Get("agency_area")
	platform := params.Get("platform")
	browser := params.Get("browser")
	qq := params.Get("qq")
	email := params.Get("email")
	tagStr := params.Get("tag_id")
	aidStr := params.Get("aid")
	var (
		mid int64
		sid int64
		err error
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if content == "" && imgURL == "" && logURL == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	sid, _ = strconv.ParseInt(sidStr, 10, 64)
	tagID, _ := strconv.ParseInt(tagStr, 10, 64)
	if platform == "" {
		platform = "ugc"
	}
	c.JSON(feedbackSvr.AddWebReply(c, mid, sid, tagID, aidStr, filtered(content), imgURL, netState, netOperator, agencyArea, platform, version, buvid, browser, qq, email, time.Now()))
}

//replys
func replys(c *bm.Context) {
	header := c.Request.Header
	params := c.Request.Form
	buvid := header.Get("Buvid")
	system := params.Get("system")
	version := params.Get("version")
	midStr := params.Get("mid")
	platform := params.Get("platform")
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	entrance := params.Get("entrance")
	// check params
	if system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if version == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		mid int64
		err error
	)
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if mid == 0 && buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 1 || ps > conf.Conf.Feedback.ReplysNum {
		ps = conf.Conf.Feedback.ReplysNum
	}
	rs, isEndReply, err := feedbackSvr.Replys(c, buvid, platform, mobiApp, device, system, version, entrance, mid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data":         rs,
		"is_end_reply": isEndReply,
	}
	c.JSONMap(res, err)
}

func sessions(c *bm.Context) {
	params := c.Request.Form
	tagid := params.Get("tag_id")
	platform := params.Get("platform")
	stateStr := params.Get("state")
	midStr := params.Get("mid")
	start := params.Get("start")
	end := params.Get("end")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		etime time.Time
	)
	if start == "" {
		start = "00-00-00 00:00:00"
	}
	stime, _ := time.Parse(_timeFromat, start)
	etime, _ = time.Parse(_timeFromat, end)
	if end == "" {
		etime = time.Now()
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps < 0 || ps > 10 {
		ps = 10
	}
	total, sessions, err := feedbackSvr.Sessions(c, mid, stateStr, tagid, platform, stime, etime, ps, pn)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data":  sessions,
		"total": total,
	}
	c.JSONMap(res, err)
}

func sessionsClose(c *bm.Context) {
	params := c.Request.Form
	sidStr := params.Get("session_id")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, feedbackSvr.UpdateSessionState(c, 3, sid))
}

func replyTag(c *bm.Context) {
	params := c.Request.Form
	tp := params.Get("type") // NOTE: player
	platform := params.Get("platform")
	c.JSON(model.Tags[tp][platform], nil)
}

func ugcTag(c *bm.Context) {
	params := c.Request.Form
	tp := params.Get("type")
	platform := params.Get("platform")
	mold, err := strconv.Atoi(tp)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(feedbackSvr.Tags(c, mid, mold, platform))
}

func webReply(c *bm.Context) {
	params := c.Request.Form
	sidStr := params.Get("session_id")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(feedbackSvr.WebReplys(c, sid, mid))
}

func filtered(content string) string {
	v := make([]rune, 0, len(content))
	for _, c := range content {
		if c != utf8.RuneError {
			v = append(v, c)
		}
	}
	return string(v)
}

func playerCheck(c *bm.Context) {
	var (
		params                                     = c.Request.Form
		header                                     = c.Request.Header
		platform, ipChangeTimes                    int
		mid, checkTime, aid, connectSpeed, ioSpeed int64
		region, school                             string
		err                                        error
	)
	platformStr := params.Get("platform")
	if platform = model.FormPlatForm(platformStr); platform == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if region = params.Get("region"); region == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if school = params.Get("school"); school == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr := params.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	checkTimeStr := params.Get("check_time")
	if checkTime, err = strconv.ParseInt(checkTimeStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, _ = strconv.ParseInt(params.Get("aid"), 10, 64)
	ipChangeTimesStr := params.Get("ip_change_times")
	if ipChangeTimes, err = strconv.Atoi(ipChangeTimesStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	connectSpeedStr := params.Get("connect_speed")
	if connectSpeed, err = strconv.ParseInt(connectSpeedStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ioSpeedStr := params.Get("io_speed")
	if ioSpeed, err = strconv.ParseInt(ioSpeedStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, feedbackSvr.PlayerCheck(c, platform, ipChangeTimes, mid, checkTime, aid, connectSpeed, ioSpeed, region, school, header.Get("X-Cache-Server-Addr")))
}

// addReplyH5 add from H5
func addReplyH5(c *bm.Context) {
	c.Request.ParseMultipartForm(model.MaxUploadSize)
	header := c.Request.Header
	params := c.Request.MultipartForm
	buvid := header.Get("Buvid")
	var system, version, midStr, content, logURL, device, channel, entrance, netState, netOperator, agencyArea, platform, browser, qq, mobiApp, email, tagStr string
	// params
	if len(params.Value["system"]) > 0 {
		system = params.Value["system"][0]
	}
	if len(params.Value["version"]) > 0 {
		version = params.Value["version"][0]
	}
	if len(params.Value["mid"]) > 0 {
		midStr = params.Value["mid"][0]
	}
	if len(params.Value["content"]) > 0 {
		content = params.Value["content"][0]
	}
	if len(params.Value["log_url"]) > 0 {
		logURL = params.Value["log_url"][0]
	}
	if len(params.Value["device"]) > 0 {
		device = params.Value["device"][0]
	}
	if len(params.Value["channel"]) > 0 {
		channel = params.Value["channel"][0]
	}
	if len(params.Value["entrance"]) > 0 {
		entrance = params.Value["entrance"][0]
	}
	if len(params.Value["net_state"]) > 0 {
		netState = params.Value["net_state"][0]
	}
	if len(params.Value["net_operator"]) > 0 {
		netOperator = params.Value["net_operator"][0]
	}
	if len(params.Value["agency_area"]) > 0 {
		agencyArea = params.Value["agency_area"][0]
	}
	if len(params.Value["platform"]) > 0 {
		platform = params.Value["platform"][0]
	}
	if len(params.Value["browser"]) > 0 {
		browser = params.Value["browser"][0]
	}
	if len(params.Value["qq"]) > 0 {
		qq = params.Value["qq"][0]
	}
	if len(params.Value["mobi_app"]) > 0 {
		mobiApp = params.Value["mobi_app"][0]
	}
	if len(params.Value["email"]) > 0 {
		email = params.Value["email"][0]
	}
	if len(params.Value["buvid"]) > 0 && buvid == "" {
		buvid = params.Value["buvid"][0]
	}
	if len(params.Value["tag_id"]) > 0 {
		tagStr = params.Value["tag_id"][0]
	}
	// check params
	if buvid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if system == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if version == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		mid int64
		err error
	)
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	imgs := c.Request.MultipartForm.File["imgs"]
	var (
		imgURL string
		urls   []string
		mutex  = sync.Mutex{}
	)
	g, ctx := errgroup.WithContext(c)
	for k, fh := range imgs {
		if k == conf.Conf.Feedback.ImgLimit {
			break
		}
		var (
			img      multipart.File
			url      string
			fileName string
			fileTpye string
			body     []byte
		)
		if img, err = fh.Open(); err != nil {
			log.Error("H5 addReply Open %s failed", fh.Filename)
			err = nil
			continue
		}
		defer img.Close()
		fileName = fh.Filename
		fileTpye = strings.TrimPrefix(path.Ext(fileName), ".")
		if body, err = ioutil.ReadAll(img); err != nil {
			log.Error("H5 addReply ioutil.ReadAll %s failed", fh.Filename)
			err = nil
			continue
		}
		g.Go(func() (err error) {
			if url, err = feedbackSvr.Upload(ctx, "", fileTpye, time.Now(), body); err != nil {
				log.Error("H5 addReply Upload %s failed", fh.Filename)
				err = nil
				return
			}
			mutex.Lock()
			urls = append(urls, url)
			mutex.Unlock()
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(urls) > 0 {
		imgURL = strings.Join(urls, ";")
	}
	if content == "" && imgURL == "" && logURL == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tagID, _ := strconv.ParseInt(tagStr, 10, 64)
	c.JSON(feedbackSvr.AddReply(c, mid, tagID, buvid, system, version, mobiApp, filtered(content), imgURL, logURL, device, channel, entrance, netState, netOperator, agencyArea, platform, browser, qq, email, time.Now()))
}
