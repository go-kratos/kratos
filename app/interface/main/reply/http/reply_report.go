package http

import (
	"strconv"
	"strings"
	"text/template"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// reportReply report a reply.
func reportReply(c *bm.Context) {
	params := c.Request.Form
	mid, _ := c.Get("mid")
	oidsStr := params.Get("oid")
	rpsStr := params.Get("rpid")
	tpStr := params.Get("type")
	reaStr := params.Get("reason")
	cont := params.Get("content")
	platform := params.Get("platform")
	buildStr := params.Get("build")
	buvid := c.Request.Header.Get("buvid")
	var build int64
	var err error
	if buildStr != "" {
		build, err = strconv.ParseInt(buildStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", buildStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	// check params
	oids, err := xstr.SplitInts(oidsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIds, err := xstr.SplitInts(rpsStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpsStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) != len(rpIds) {
		log.Warn("oids(%s) not equal rpids(%s)", oidsStr, rpsStr)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	reaTmp, err := strconv.ParseInt(reaStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", reaStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rea := int8(reaTmp)
	if rea == model.ReportReasonOther {
		cont = strings.TrimSpace(cont)
		cl := len([]rune(cont))
		if 200 < cl || cl < 2 {
			log.Warn("content(%s) length %d, max 200, min 2", cont, cl)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
		cont = template.HTMLEscapeString(cont)
	} else {
		cont = ""
	}
	for i := 0; i < len(oids); i++ {
		cd, err := rpSvr.AddReport(c, mid.(int64), oids[i], rpIds[i], int8(tp), rea, cont, platform, build, buvid)
		if err != nil {
			var data map[string]int
			if err == ecode.ReplyReportDeniedAsCD {
				data = map[string]int{"ttl": cd}
			} else {
				log.Warn("rpSvr.AddReport(%d, %d, %d, %d, %d, %s, %s, %d, %s) error(%v)", mid, oids[i], rpIds[i], tp, rea, cont, err, platform, build, buvid)
			}
			c.JSON(data, err)
			return
		}
	}
	c.JSON(nil, nil)
}

func reportRelated(c *bm.Context) {
	var (
		mid    int64
		escape = true
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	rpidStr := params.Get("root")
	tpStr := params.Get("type")
	midIf, ok := c.Get("mid")
	if ok {
		mid = midIf.(int64)
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(oid:%d) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpid, err := strconv.ParseInt(rpidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(rpid:%d) error(%v)", rpidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(type:%d) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	// check android and ios appkey
	// if mobile, no html escape
	if isMobile(params) {
		escape = false
	}
	sub, root, rels, err := rpSvr.ReportRelated(c, mid, oid, rpid, int8(tp), escape)
	if err != nil {
		log.Warn("rpSvr.ReportRelated(%d,%d,%d) error(%v)", oid, rpid, tp, err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"root":    root,
		"related": rels,
		"upper": map[string]interface{}{
			"mid": sub.Mid,
		},
	}
	c.JSON(data, err)
}

func reportSndReply(c *bm.Context) {
	var (
		mid    int64
		escape = true
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	rpidStr := params.Get("root")
	tpStr := params.Get("type")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	midIf, ok := c.Get("mid")
	if ok {
		mid = midIf.(int64)
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(oid:%d) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpid, err := strconv.ParseInt(rpidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(rpid:%d) error(%v)", rpidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(type:%d) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(pn:%d) error(%v)", pnStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(ps:%d) error(%v)", psStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	// check android and ios appkey
	// if mobile, no html escape
	if isMobile(params) {
		escape = false
	}
	sub, root, rs, err := rpSvr.ReportReply(c, mid, oid, rpid, int8(tp), int(pn), int(ps), escape)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["page"] = map[string]int{
		"num":   int(pn),
		"size":  int(ps),
		"count": root.Count,
	}
	data["upper"] = map[string]interface{}{
		"mid": sub.Mid,
	}
	data["root"] = root
	data["replies"] = rs
	c.JSON(data, err)
}
