package http

import (
	"io/ioutil"
	"net/http"

	"encoding/json"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func creatorEdit(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.CreatorParam{}
		err error
	)
	defer func() {
		//特例 creatorEdit
		ap.Title = cp.Title
		ap.Aid = cp.Aid
		ap.Tag = cp.Tag
		ap.Desc = cp.Desc
		ap.OpenElec = cp.OpenElec
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "edit", "creator", ap, err)
	}()
	midI, _ := c.Get("mid")
	mid, _ := midI.(int64)

	if err = c.Bind(cp); err != nil {
		err = ecode.VideoupParamErr
		return
	}
	if cp.Aid == 0 {
		err = ecode.VideoupParamErr
		return
	}
	aid = cp.Aid
	ap.Aid = cp.Aid
	ap.Title = cp.Title
	ap.Desc = cp.Desc
	ap.Tag = cp.Tag
	ap.OpenElec = cp.OpenElec

	err = vdpSvc.CreatorEdit(c, mid, cp)
	if err != nil {
		c.JSON(nil, err)
		log.Error("addErr err(%+v)|cp(%+v)", err, cp)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": cp.Aid,
	}, nil)
}

func creatorAdd(c *bm.Context) {
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
	ap.Mid = mid
	ap.UpFrom = archive.UpFromCreator
	defer func() {
		cp.IPv6 = ap.IPv6
		cp.Aid = ap.Aid
		cp.UpFrom = ap.UpFrom
		cp.Mid = ap.Mid
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", "creator", cp, err)
	}()
	if err == nil {
		aid, err = vdpSvc.CreatorAdd(c, mid, ap)
	}
	if err != nil {
		c.JSON(nil, err)
		log.Error("addErr err(%+v)|ap(%+v)", err, ap)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": aid,
	}, nil)
}

func creatorUpCover(c *bm.Context) {
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
	url, err := vdpSvc.CreatorUpCover(c, ftype, bs, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url": url,
	}, nil)
}
