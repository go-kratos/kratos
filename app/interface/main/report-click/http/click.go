package http

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/interface/main/report-click/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_platWeb       = "0"
	_platH5        = "1"
	_platOuter     = "2"
	_platIos       = "3"
	_platAndroid   = "4"
	_platAndroidTV = "5"
)

var _expireCookie = time.Date(2022, time.November, 10, 23, 0, 0, 0, time.UTC)

// webClick write the archive data.
func webClick(c *bm.Context) {
	err := flashPlay(c, _platWeb, conf.Conf.Click.WebSecret)
	c.JSON(nil, err)
}

// outerClick write the archive data.
func outerClick(c *bm.Context) {
	err := flashPlay(c, _platOuter, conf.Conf.Click.OutSecret)
	c.JSON(nil, err)
}

// iosClick write the archive data.
func iosClick(c *bm.Context) {
	err := mobilePlay(c, conf.Conf.Click.AesKey, conf.Conf.Click.AesIv, conf.Conf.Click.AesSalt, _platIos)
	c.JSON(nil, err)
}

// androidClick write the archive data.
func androidClick(c *bm.Context) {
	err := mobilePlay(c, conf.Conf.Click.AesKey, conf.Conf.Click.AesIv, conf.Conf.Click.AesSalt, _platAndroid)
	c.JSON(nil, err)
}

// android2Click write the archive data.
func android2Click(c *bm.Context) {
	err := mobilePlay(c, conf.Conf.Click.AesKey2, conf.Conf.Click.AesIv2, conf.Conf.Click.AesSalt2, _platAndroid)
	c.JSON(nil, err)
}

// androidTV == android2Click write the archive data.
func androidTV(c *bm.Context) {
	err := mobilePlay(c, conf.Conf.Click.AesKey2, conf.Conf.Click.AesIv2, conf.Conf.Click.AesSalt2, _platAndroidTV)
	c.JSON(nil, err)
}

// outerClickH5  h5 outer click same to flash plat .
func outerClickH5(c *bm.Context) {
	c.JSON(nil, h5Play(c, _platOuter))
}

// h5Click write the archive data.
func h5Click(c *bm.Context) {
	c.JSON(nil, h5Play(c, _platH5))
}

// webH5Click write the archive data.
func webH5Click(c *bm.Context) {
	c.JSON(nil, h5Play(c, _platWeb))
}

// flashPlay.
func flashPlay(c *bm.Context, plat, secret string) (err error) {
	var (
		buvid  string
		mid    int64
		ck     *http.Cookie
		params = c.Request.Form
		unix   = time.Now()
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if ck, err = c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	userAgent := c.Request.Header.Get("User-Agent")
	refer := c.Request.Header.Get("Referer")
	sign := params.Get("sign")
	if err = clickSvr.FlashSigned(params, secret, unix); err != nil {
		log.Error("clickSvr.FlashSigned() error(%v)", err)
		if err == ecode.ClickQuerySignErr {
			log.Warn("click sign error(%s,%s,%s,%s)", sign, refer, userAgent, c.Request.Header.Get("Origin"))
		}
		return
	}
	midStr := params.Get("mid")
	if mid != 0 && midStr != strconv.FormatInt(mid, 10) {
		log.Warn("flashPlay stat mid(%d) not equal stat mid(%s)", mid, midStr)
		return
	}
	midStr = strconv.FormatInt(mid, 10)
	aid := params.Get("aid")
	var cookieSid string
	if ck, err := c.Request.Cookie("sid"); err == nil {
		cookieSid = ck.Value
	}
	typeID := params.Get("type")
	subType := params.Get("sub_type")
	sid := params.Get("sid")
	epid := params.Get("epid")
	// service.
	ip := metadata.String(c, metadata.RemoteIP)
	clickSvr.Play(c, plat, aid, params.Get("cid"), params.Get("part"),
		midStr, params.Get("lv"), params.Get("ftime"), params.Get("stime"),
		params.Get("did"), ip, userAgent, buvid, cookieSid, refer, typeID, subType, sid, epid, "", "", "", "", "", "")
	return
}

// mobilePlay.
func mobilePlay(c *bm.Context, aesKey, aesIv, aesSalt, plat string) (err error) {
	// check params.
	req := c.Request
	unix := time.Now()
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll error(%v)", err)
		err = ecode.ServerErr
		return
	}
	req.Body.Close()
	bs, err = clickSvr.Decrypt(bs, aesKey, aesIv)
	if err != nil {
		log.Error("clickSvr.Decrypt(%s) error(%d)", bs, err)
		return
	}
	p, err := clickSvr.Verify(bs, aesSalt, unix)
	if err != nil {
		log.Error("clickSvr.Verify(%s) error(%d)", bs, err)
		return
	}
	req.Form = p // for log
	typeID := p.Get("type")
	subType := p.Get("sub_type")
	sid := p.Get("sid")
	epid := p.Get("epid")
	playMode := p.Get("play_mode")
	platform := p.Get("platform")
	device := p.Get("device")
	mobiAapp := p.Get("mobi_app")
	autoPlay := p.Get("auto_play")
	ap, _ := strconv.ParseInt(autoPlay, 10, 64)
	// service.
	aidStr := p.Get("aid")
	var (
		accessKey string
		midStr    = p.Get("mid")
		noAccess  bool
	)
	paasMid, _ := strconv.ParseInt(midStr, 10, 64)
	accessKey = p.Get("access_key")
	if paasMid > 0 && accessKey == "" {
		noAccess = true
	}
	if accessKey != "" {
		c.Request.Form.Set("access_key", accessKey)
		authSvc.User(c)
		mid, ok := c.Get("mid")
		if !ok {
			log.Warn("idfSvc.Access() access_key", accessKey)
			if paasMid > 0 {
				noAccess = true
			}
		}
		if mid != nil {
			midStr = strconv.FormatInt(mid.(int64), 10)
		}
	}
	refer := c.Request.Header.Get("Referer")
	userAgent := c.Request.Header.Get("User-Agent")
	if noAccess {
		userAgent = userAgent + " (no_accesskey)"
	}
	if ap == 1 || ap == 2 { // abandon the logic that transforms the plat to 6/7/8/9, keep the plat and modify the UA
		userAgent += " (inline_play_begin)"
	}
	buvid := req.Header.Get("buvid")
	var cookieSid string
	if ck, err := c.Request.Cookie("sid"); err == nil {
		cookieSid = ck.Value
	}
	ip := metadata.String(c, metadata.RemoteIP)
	clickSvr.Play(c, plat, aidStr, p.Get("cid"), p.Get("part"), midStr, p.Get("lv"),
		p.Get("ftime"), p.Get("stime"), p.Get("did"), ip, userAgent, buvid,
		cookieSid, refer, typeID, subType, sid, epid, playMode, platform, device, mobiAapp, autoPlay, "")
	return
}

func h5Play(c *bm.Context, plat string) (err error) {
	var (
		buvid  string
		mid    int64
		ck     *http.Cookie
		params = c.Request.Form
		unix   = time.Now()
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if ck, err = c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	// check params.
	st := params.Get("stime")
	stime, err := strconv.ParseInt(st, 10, 64)
	if err != nil {
		err = ecode.ClickQueryFormatErr
		return
	}
	if unix.Unix()-stime > 60 {
		err = ecode.ClickServerTimeout
		return
	}
	typeID := params.Get("type")
	subType := params.Get("sub_type")
	sid := params.Get("sid")
	epid := params.Get("epid")
	var ft string
	// check cookie did
	var (
		did string
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	cookie, err := c.Request.Cookie("rpdid")
	if err != nil || cookie.Value == "" {
		did = clickSvr.GenDid(ip, unix)
		http.SetCookie(c.Writer, &http.Cookie{Name: "rpdid", Value: did, Path: "/", Domain: ".bilibili.com", Expires: _expireCookie})
		err = nil
	} else {
		did = cookie.Value
	}
	_, ft = clickSvr.CheckDid(did)
	if ft == "" {
		log.Error("ft null ft:%s,did:%s", ft, did)
		return
	}
	midStr := params.Get("mid")
	if mid != 0 && midStr != strconv.FormatInt(mid, 10) {
		log.Warn("h5 stat mid(%d) not equal stat mid(%s)", mid, midStr)
		return
	}
	midStr = strconv.FormatInt(mid, 10)
	aid := params.Get("aid")
	userAgent := c.Request.Header.Get("User-Agent")
	var cookieSid string
	if ck, err := c.Request.Cookie("sid"); err == nil {
		cookieSid = ck.Value
	}
	refer := c.Request.Header.Get("Referer")
	// service.
	clickSvr.Play(c, plat, aid, params.Get("cid"), params.Get("part"),
		midStr, params.Get("lv"), ft, params.Get("stime"), did,
		ip, userAgent, buvid, cookieSid, refer, typeID, subType, sid, epid, "", "", "", "", "", "")
	return
}
