package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
)

func appEdit(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)
	mobi := c.Request.Form.Get("mobi_app")
	build := c.Request.Form.Get("build")
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.Request.Body.Close()
	if mobi == "iphone" && (build == "4470" || build == "4430") {
		err = iosUnmarshal(bs, ap)
		err = iosUnmarshal(bs, cp)
	} else {
		err = json.Unmarshal(bs, ap)
		err = json.Unmarshal(bs, cp)
	}
	if err != nil {
		c.JSON(nil, ecode.VideoupParamErr)
		return
	}
	aid = ap.Aid
	ap.Mid = mid
	// 老的只编辑基础信息， 不需要按照平台自动转换
	ap.UpFrom = archive.UpFromAPP
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.UpFrom = ap.UpFrom
		cp.IPv6 = ap.IPv6
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "edit", "app", cp, err)
	}()
	err = vdpSvc.AppEdit(c, ap, mid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("editErr err(%+v)|ap(%+v)", err, ap)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": ap.Aid,
	}, nil)
}

func iosUnmarshal(bs []byte, ap *archive.ArcParam) (err error) {
	var apStr = struct {
		Aid       string `json:"aid"`
		Copyright string `json:"copyright"`
		Cover     string `json:"cover"`
		Title     string `json:"title"`
		TypeID    string `json:"tid"`
		Tag       string `json:"tag"`
		Desc      string `json:"desc"`
		MissionID string `json:"mission_id"`
		OpenElec  string `json:"open_elec"`
		DTime     string `json:"dtime"`
	}{}
	if err = json.Unmarshal(bs, &apStr); err != nil {
		return
	}
	aid, err := strconv.ParseInt(apStr.Aid, 10, 64)
	if err != nil {
		return
	}
	copyright, err := strconv.ParseInt(apStr.Copyright, 10, 8)
	if err != nil {
		return
	}
	typeID, err := strconv.ParseInt(apStr.TypeID, 10, 16)
	if err != nil {
		return
	}
	missionID, err := strconv.ParseInt(apStr.MissionID, 10, 10)
	if err != nil {
		return
	}
	openElec, err := strconv.ParseInt(apStr.OpenElec, 10, 8)
	if err != nil {
		return
	}
	dtime, err := strconv.ParseInt(apStr.DTime, 10, 64)
	if err != nil {
		return
	}
	ap.Aid = aid
	ap.Copyright = int8(copyright)
	ap.TypeID = int16(typeID)
	ap.MissionID = int(missionID)
	ap.OpenElec = int8(openElec)
	ap.DTime = time.Time(dtime)
	ap.Cover = apStr.Cover
	ap.Title = apStr.Title
	ap.Tag = apStr.Tag
	ap.Desc = apStr.Desc
	return
}

func appUpCover(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("c.Request().FormFile(\"file\") error(%v) | ", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Error("ioutil.ReadAll(c.Request().Body) error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	if ftype != "image/jpeg" && ftype != "image/png" && ftype != "image/webp" {
		log.Error("filetype not allow file type(%s)", ftype)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	url, err := vdpSvc.AppUpCover(c, ftype, bs, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url": url,
	}, nil)
}

func appAdd(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)
	form := c.Request.Form
	ar := &archive.AppRequest{
		Build:    form.Get("build"),
		MobiApp:  form.Get("mobi_app"),
		Platform: form.Get("platform"),
		Device:   form.Get("device"),
	}
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.Request.Body.Close()
	err = json.Unmarshal(bs, ap)
	err = json.Unmarshal(bs, cp)
	if err != nil {
		c.JSON(nil, ecode.VideoupParamErr)
		return
	}
	ap.Aid = 0
	ap.Mid = mid
	ap.UpFrom = appUpFrom(ar.Platform, ar.Device)
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.UpFrom = ap.UpFrom
		cp.IPv6 = ap.IPv6
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", ar.Platform, cp, err)
	}()
	if err == nil {
		aid, err = vdpSvc.AppAdd(c, mid, ap, ar)
	}
	if err != nil {
		log.Error("AppAdd Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": aid,
	}, nil)
}

func appEditFull(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)
	form := c.Request.Form
	ar := &archive.AppRequest{
		Build:    form.Get("build"),
		MobiApp:  form.Get("mobi_app"),
		Platform: form.Get("platform"),
		Device:   form.Get("device"),
	}
	buildStr := form.Get("build")
	buildNum, _ := strconv.ParseInt(buildStr, 10, 64)
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)
	ap = &archive.ArcParam{}
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.Request.Body.Close()
	err = json.Unmarshal(bs, ap)
	err = json.Unmarshal(bs, cp)
	if err != nil {
		c.JSON(nil, ecode.VideoupParamErr)
		return
	}
	aid = ap.Aid
	ap.Mid = mid
	ap.UpFrom = appUpFrom(ar.Platform, ar.Device)
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.UpFrom = ap.UpFrom
		cp.IPv6 = ap.IPv6
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "edit", ar.Platform, cp, err)
	}()
	err = vdpSvc.AppEditFull(c, ap, mid, buildNum, ar)
	if err != nil {
		log.Error("AppEditFull Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": ap.Aid,
	}, nil)
}

func appUpFrom(platfrom, device string) (res int8) {
	if platfrom == "ios" {
		if device == "pad" {
			res = archive.UpFromIpad
		} else {
			res = archive.UpFromAPPiOS
		}
	} else if platfrom == "android" {
		res = archive.UpFromAPPAndroid
	} else {
		res = archive.UpFromAPP
	}
	return
}
