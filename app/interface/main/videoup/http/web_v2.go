package http

import (
	"encoding/json"
	"io/ioutil"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func webV2Add(c *bm.Context) {
	var (
		aid       int64
		ap        = &archive.ArcParam{}
		cp        = &archive.ArcParam{}
		err       error
		validated bool
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
	defer func() {
		cp.Aid = ap.Aid
		cp.Mid = ap.Mid
		cp.IPv6 = ap.IPv6
		cp.UpFrom = ap.UpFrom
		build, buvid := getBuildInfo(c)
		vdpSvc.SendArchiveLog(aid, build, buvid, "add", "web/v2", cp, err)
	}()
	validated, err = vdpSvc.Validate(c, ap.Geetest, "web", mid)
	if validated || err == ecode.CreativeGeetestAPIErr {
		ap.Mid = mid
		ap.UpFrom = archive.UpFromWeb
		aid, err = vdpSvc.WebAdd(c, mid, ap, validated)
		if err != nil {
			log.Error("vdpSvc.WebAdd Err mid(%+d)|ap(%+v)|err(%+v)", mid, ap, err)
			c.JSON(nil, err)
			return
		}
		c.JSON(map[string]interface{}{
			"aid": aid,
		}, nil)
	} else {
		c.JSON(nil, ecode.CreativeGeetestErr)
	}
}
