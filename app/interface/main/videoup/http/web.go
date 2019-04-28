package http

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func webAdd(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
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
	ap.UpFrom = archive.UpFromWeb
	// code, msg, data := vdpSvc.WebFilterArcParam(c, ap, ip)
	// if code != 0 {
	// 	c.Error = ecode.VideoupFieldFilterForbid
	// 	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
	// 		"code":    code,
	// 		"message": msg,
	// 		"data":    data,
	// 	}))
	// 	return
	// }
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.IPv6 = ap.IPv6
		cp.UpFrom = ap.UpFrom
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", "web", cp, err)
	}()
	aid, err = vdpSvc.WebAdd(c, mid, ap, false)
	if err != nil {
		log.Error("vdpSvc.WebAdd Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": aid,
	}, nil)
}

func webEdit(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
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
	ap.UpFrom = archive.UpFromWeb
	// code, msg, data := vdpSvc.WebFilterArcParam(c, ap, ip)
	// if code != 0 {
	// 	c.Error = ecode.VideoupFieldFilterForbid
	// 	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
	// 		"code":    code,
	// 		"message": msg,
	// 		"data":    data,
	// 	}))
	// 	return
	// }
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.UpFrom = ap.UpFrom
		cp.IPv6 = ap.IPv6
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "edit", "web", cp, err)
	}()
	err = vdpSvc.WebEdit(c, ap, mid)
	if err != nil {
		log.Error("vdpSvc.WebEdit Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": ap.Aid,
	}, nil)
}

func webUpCover(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	cover := c.Request.Form.Get("cover")
	c.Request.Form.Del("cover")
	ss := strings.Split(cover, ",")
	if len(ss) != 2 || len(ss[1]) == 0 {
		log.Error("cover(%s) format error", cover)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", ss[1], err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ftype := http.DetectContentType(bs)
	if ftype != "image/jpeg" && ftype != "image/png" && ftype != "image/webp" {
		log.Error("file type not allow file type(%s)", ftype)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	url, err := vdpSvc.WebUpCover(c, ftype, bs, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"url": url,
	}, nil)
}

func webFilter(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err error
	)
	// msg := c.Request.Form.Get("msg")
	// if len(msg) != 0 {
	// 	_, err = vdpSvc.WebSingleFilter(c, msg)
	// }
	c.JSON(nil, err)
}

// staffTitleFilter 过滤联合投稿职能
func webStaffTitleFilter(c *bm.Context) {
	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		err error
	)
	title := c.Request.Form.Get("title")
	if len(title) != 0 {
		var hit []string
		_, hit, err = vdpSvc.WebSingleFilter(c, title)
		if len(hit) > 0 {
			err = ecode.VideoupStaffTitleFilter
		}
	}
	c.JSON(nil, err)
}

func webCmAdd(c *bm.Context) {
	var (
		aid int64
		ap  = &archive.ArcParam{}
		cp  = &archive.ArcParam{}
		err error
	)

	midI, _ := c.Get("mid")
	mid, ok := midI.(int64)
	if !ok || mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
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
	ap.UpFrom = archive.UpFromCM
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.IPv6 = ap.IPv6
		cp.UpFrom = ap.UpFrom
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", "cm", cp, err)
	}()
	aid, err = vdpSvc.WebCmAdd(c, mid, ap)
	if err != nil {
		log.Error("vdpSvc.WebCmAdd Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"aid": aid,
	}, nil)
}
