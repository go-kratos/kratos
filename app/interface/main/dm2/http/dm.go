package http

import (
	"math"
	"net/http"
	"strconv"
	xtime "time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/ip"
	"go-common/library/net/metadata"
	"go-common/library/time"
)

func httpCode(err error) (code int) {
	switch err {
	case ecode.NotModified:
		code = http.StatusNotModified
	case ecode.RequestErr:
		code = http.StatusBadRequest
	case ecode.NothingFound:
		code = http.StatusNotFound
	case ecode.ServiceUnavailable:
		code = http.StatusServiceUnavailable
	default:
		code = http.StatusInternalServerError
	}
	return
}

func dmXML(c *bm.Context) {
	var (
		p           = c.Request.Form
		comp        = p.Get("comp")
		contentType = "text/xml"
	)
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	data, err := dmSvc.DMXML(c, model.SubTypeVideo, oid)
	if err != nil {
		c.AbortWithStatus(httpCode(err))
		log.Error("dmSvc.XML(%d) error(%v)", oid, err)
		return
	}
	c.Writer.Header().Set("Content-Encoding", "deflate")
	c.Writer.Header().Set("Last-Modified", xtime.Now().Format(http.TimeFormat))
	if comp == "0" {
		c.Writer.Header().Set("Content-Encoding", "none")
		if data, err = dmSvc.Gzdeflate(data); err != nil {
			log.Error("dmSvc.Gzdeflate(%d) error(%v)", oid, err)
			c.AbortWithStatus(httpCode(err))
			return
		}
	}
	c.Bytes(http.StatusOK, contentType, data)
}

func dmSeg(c *bm.Context) {
	var (
		plat        int32
		mid         int64
		contentType = "application/octet-stream"
		p           = c.Request.Form
	)
	iMid, ok := c.Get("mid")
	if ok {
		mid = iMid.(int64)
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ps, err := strconv.ParseInt(p.Get("ps"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	platform, err := strconv.ParseInt(p.Get("plat"), 10, 64)
	if err != nil {
		plat = model.PlatUnknow
	} else {
		plat = int32(platform)
	}
	data, err := dmSvc.DMSeg(c, int32(tp), plat, mid, aid, oid, ps)
	if err != nil {
		c.AbortWithStatus(httpCode(err))
		return
	}
	if len(data) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Bytes(http.StatusOK, contentType, data)
}

func dmSegV2(c *bm.Context) {
	var (
		plat int32
		mid  int64
		p    = c.Request.Form
	)
	iMid, ok := c.Get("mid")
	if ok {
		mid = iMid.(int64)
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.ParseInt(p.Get("pn"), 10, 64)
	if err != nil || pn <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	platform, err := strconv.ParseInt(p.Get("plat"), 10, 64)
	if err != nil {
		plat = model.PlatUnknow
	} else {
		plat = int32(platform)
	}
	c.JSON(dmSvc.DMSegV2(c, int32(tp), mid, aid, oid, pn, plat))
}

func dm(c *bm.Context) {
	p := c.Request.Form
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(dmSvc.DM(c, int32(tp), aid, oid))
}

func ajaxDM(c *bm.Context) {
	var (
		p    = c.Request.Form
		msgs = make([]string, 0)
	)
	app := p.Get("mobi_app")
	if app == "android" || app == "iphone" {
		c.JSON(msgs, nil)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(dmSvc.AjaxDM(c, aid))
}

// validDMStyle 验证弾幕pool and mode.
func validDMStyle(pool, mode int32) (valid bool) {
	switch pool {
	case model.PoolNormal, model.PoolSubtitle:
		if mode == model.ModeRolling || mode == model.ModeBottom || mode == model.ModeTop ||
			mode == model.ModeReverse || mode == model.ModeSpecial {
			valid = true
		}
	case model.PoolSpecial:
		if mode == model.ModeCode || mode == model.ModeBAS {
			valid = true
		}
	}
	return
}

//dm post
func dmPost(c *bm.Context) {
	var (
		plat = int64(model.PlatUnknow)
		p    = c.Request.Form
		now  = xtime.Now().Unix()
	)
	mid, _ := c.Get("mid")
	msg := p.Get("msg")
	typ, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil || int32(typ) != model.SubTypeVideo {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	progress, err := strconv.ParseInt(p.Get("progress"), 10, 32) // NOTE 老接口过来的弹幕时间为秒
	if err != nil || progress < 0 || progress > math.MaxInt32 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	color, err := strconv.ParseInt(p.Get("color"), 10, 64)
	if err != nil || color < 0 || color > math.MaxInt32 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	fontsize, err := strconv.ParseInt(p.Get("fontsize"), 10, 32)
	if err != nil || fontsize <= 0 || fontsize > 127 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pool, err := strconv.ParseInt(p.Get("pool"), 10, 32)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mode, err := strconv.ParseInt(p.Get("mode"), 10, 32)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !validDMStyle(int32(pool), int32(mode)) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rnd, err := strconv.ParseInt(p.Get("rnd"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	platStr := p.Get("plat")
	if platStr != "" {
		if plat, err = strconv.ParseInt(platStr, 10, 32); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	dm := &model.DM{
		Type:     int32(typ),
		Oid:      oid,
		Mid:      mid.(int64),
		Progress: int32(progress),
		Pool:     int32(pool),
		State:    model.StateNormal,
		Ctime:    time.Time(now),
		Mtime:    time.Time(now),
		Content: &model.Content{
			FontSize: int32(fontsize),
			Color:    color,
			IP:       int64(ip.InetAtoN(metadata.String(c, metadata.RemoteIP))),
			Mode:     int32(mode),
			Plat:     int32(plat),
			Msg:      msg,
			Ctime:    time.Time(now),
			Mtime:    time.Time(now),
		},
	}
	if dm.Pool == model.PoolSpecial {
		dm.ContentSpe = &model.ContentSpecial{
			Msg:   msg,
			Ctime: time.Time(now),
			Mtime: time.Time(now),
		}
	}
	if err = dmSvc.Post(c, dm, aid, rnd); err != nil {
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"dmid": dm.ID,
	}
	c.JSON(data, nil)
}

// judgeDM dm judge list.
func judgeDM(c *bm.Context) {
	p := c.Request.Form
	cid, err := strconv.ParseInt(p.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmid, err := strconv.ParseInt(p.Get("dmid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data, err := dmSvc.JudgeDms(c, 1, cid, dmid)
	if err != nil {
		log.Error("dmSvc.JudgeDms(cid:%d,dmid:%d) error(%v)", cid, dmid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func dmAdvert(c *bm.Context) {
	p := c.Request.Form
	arg := &model.ADReq{
		ClientIP: metadata.String(c, metadata.RemoteIP),
		Buvid:    c.Request.Header.Get("Buvid"),
		MobiApp:  p.Get("mobi_app"),
		ADExtra:  p.Get("ad_extra"),
	}
	if mid, ok := c.Get("mid"); ok {
		arg.Mid = mid.(int64)
	}
	typ, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil || int32(typ) != model.SubTypeVideo {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(p.Get("aid"), 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Aid = aid
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Oid = oid
	build, err := strconv.ParseInt(p.Get("build"), 10, 64)
	if err != nil || build <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Build = build
	c.JSON(dmSvc.DMAdvert(c, arg))
}
