package http

import (
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_buvid = "buvid3"
)

func view(c *bm.Context) {
	var (
		aid, mid, cid int64
		cookieStr     string
		err           error
		rs            *model.View
	)
	cookieStr = c.Request.Header.Get("Cookie")
	aidStr := c.Request.Form.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	cidStr := c.Request.Form.Get("cid")
	if cidStr != "" {
		if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	cdnIP := c.Request.Header.Get("X-Cache-Server-Addr")
	if rs, err = webSvc.View(c, aid, cid, mid, cdnIP, cookieStr); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rs, nil)
}

func archiveStat(c *bm.Context) {
	var (
		aid int64
		err error
	)
	aidStr := c.Request.Form.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.ArchiveStat(c, aid))
}

func addShare(c *bm.Context) {
	var (
		aid, mid   int64
		err        error
		buvid, sid string
	)
	aidStr := c.Request.Form.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get mid
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if bdCookie, _ := c.Request.Cookie(_buvid); bdCookie != nil {
		buvid = bdCookie.Value
	}
	if sidCookie, _ := c.Request.Cookie("sid"); sidCookie != nil {
		sid = sidCookie.Value
	}
	c.JSON(webSvc.AddShare(c, aid, mid, c.Request.UserAgent(), c.Request.Referer(), c.Request.URL.Path, buvid, sid))
}

func description(c *bm.Context) {
	var (
		aid, page int64
		err       error
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pstr := params.Get("page")
	if pstr != "" {
		if page, err = strconv.ParseInt(pstr, 10, 64); err != nil || page < 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(webSvc.Description(c, aid, page))
}

func arcReport(c *bm.Context) {
	var (
		aid, mid, tp int64
		err          error
		params       = c.Request.Form
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aid")
	tpStr := params.Get("type")
	reason := params.Get("reason")
	pics := params.Get("pics")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, webSvc.ArcReport(c, mid, aid, tp, reason, pics))
}

func appealTags(c *bm.Context) {
	c.JSON(webSvc.AppealTags(c))
}

func arcAppeal(c *bm.Context) {
	var (
		mid int64
		err error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	data := make(map[string]string)
	for name := range params {
		switch name {
		case "tid":
			tidStr := params.Get("tid")
			if _, err = strconv.ParseInt(tidStr, 10, 64); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
			data["tid"] = tidStr
		case "aid":
			aidStr := params.Get("aid")
			if _, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
			data["oid"] = aidStr
		case "desc":
			desc := params.Get("desc")
			if desc == "" {
				c.JSON(nil, ecode.RequestErr)
				return
			}
			data["description"] = desc
		default:
			data[name] = params.Get(name)
		}
	}
	c.JSON(nil, webSvc.ArcAppeal(c, mid, data))
}

func authorRecommend(c *bm.Context) {
	var (
		aid int64
		err error
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.AuthorRecommend(c, aid))
}

func relatedArcs(c *bm.Context) {
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(webSvc.RelatedArcs(c, v.Aid))
}

func detail(c *bm.Context) {
	var (
		mid int64
		err error
		rs  *model.Detail
	)
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	cdnIP := c.Request.Header.Get("X-Cache-Server-Addr")
	if rs, err = webSvc.Detail(c, v.Aid, mid, cdnIP, c.Request.Header.Get("Cookie")); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rs, nil)
}

func arcUGCPay(c *bm.Context) {
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(webSvc.ArcUGCPay(c, mid, v.Aid))
}

func arcRelation(c *bm.Context) {
	v := new(struct {
		Aid int64 `form:"aid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(webSvc.ArcRelation(c, mid, v.Aid))
}
