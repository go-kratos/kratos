package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func clientAdd(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)

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
	ap.UpFrom = archive.UpFromWindows
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.IPv6 = ap.IPv6
		cp.UpFrom = ap.UpFrom
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", "windows", cp, err)
	}()
	aid, err = vdpSvc.ClientAdd(c, mid, ap)
	if err != nil {
		c.JSON(nil, err)
		log.Error("addErr err(%+v)|ap(%+v)", err, ap)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": aid,
	}, nil)
}

func clientEdit(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)

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
	aid = ap.Aid
	ap.Mid = mid
	ap.UpFrom = archive.UpFromWindows
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.UpFrom = ap.UpFrom
		cp.IPv6 = ap.IPv6
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "edit", "windows", cp, err)
	}()
	err = vdpSvc.ClientEdit(c, ap, mid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("editErr err(%+v)|ap(%+v)", err, ap)
		return
	}

	c.JSON(map[string]interface{}{
		"aid": ap.Aid,
	}, nil)
}

func clientUpCover(c *bm.Context) {
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
	url, err := vdpSvc.ClientUpCover(c, ftype, bs, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url": url,
	}, nil)
}
