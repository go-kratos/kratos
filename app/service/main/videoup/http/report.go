package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// viewArchive get archive info.
func arcReport(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	aidStr := params.Get("aid")
	tpStr := params.Get("type")
	reason := params.Get("reason")
	pics := params.Get("pics")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp < 0 || tp > 9 {
		log.Error("type(%d) or question empty error(%v)", tp, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err = vdpSvc.ArcReport(c, mid, aid, int8(tp), reason, pics, time.Now())
	if err != nil {
		log.Error(" vdpSvc.ArcReport(%d,%d,%d,%s,%s) error(%v)", mid, aid, int8(tp), reason, pics, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}
