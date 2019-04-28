package http

//assist 创作中心协管相关

import (
	"go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/xstr"
	"net/http"
	"strconv"
	"time"
)

func assists(c *bm.Context) {
	midStr := c.Request.Form.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assists, err := assSvc.Assists(c, mid)
	if err != nil {
		log.Error("assistSvc.Assists(%v) error(%v)", assists, mid)
		return
	}
	c.JSON(assists, nil)
}

func assistsMids(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assmidsStr := params.Get("assmids")
	assmids, err := xstr.SplitInts(assmidsStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assmidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(assmids) > 20 {
		log.Error("assmids(%d) number gt 20", len(assmids))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	asByMids, err := assSvc.AssistsMidsTotal(c, mid, assmids)
	if err != nil {
		log.Error("assistSvc.AssistsMidsTotal(%v), mids(%v), assmids(%v), error(%v)", asByMids, mid, assmids, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(asByMids, nil)
}

func assistInfo(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typeStr := params.Get("type")
	tp, err := strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assist, err := assSvc.Assist(c, mid, assistMid, tp)
	if err != nil {
		c.JSON(nil, err)
		log.Error("assSvc.Assist(%s) error(%v)|mid(%d)|assistMid(%d)", assist, err, mid, assistMid)
		return
	}
	c.JSON(assist, nil)
}

func assistIDs(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ids, err := assSvc.AssistIDs(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ids, nil)
}

func assistAdd(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	// mid
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil || assistMid == 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = assSvc.AddAssist(c, mid, assistMid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"mid":        mid,
		"assist_mid": assistMid,
	}, nil)
}

func assistDel(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := assSvc.DelAssist(c, mid, assistMid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"mid":        mid,
		"assist_mid": assistMid,
	}, nil)
}

// exist to be assist from follower
func assistExit(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := assSvc.Exit(c, mid, assistMid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"mid":        mid,
		"assist_mid": assistMid,
	}, nil)
}

func assistLogs(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	var (
		err                                               error
		mid, assistMid, ps, pn, total, bgnCtime, endCtime int64
		assistLogs                                        []*assist.Log
	)
	mid, err = strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, _ = strconv.ParseInt(assistMidStr, 10, 64)
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	ps, err = strconv.ParseInt(psStr, 10, 64)
	if err != nil || ps <= 10 {
		ps = 10
	}
	pn, err = strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		pn = 1
	}
	bgnCtimeStr := params.Get("stime")
	bgnCtime, err = strconv.ParseInt(bgnCtimeStr, 10, 64)
	if err != nil || bgnCtime <= 0 {
		bgnCtime = time.Now().Add(-time.Hour * 72).Unix()
	}
	endCtimeStr := params.Get("etime")
	endCtime, err = strconv.ParseInt(endCtimeStr, 10, 64)
	if err != nil || endCtime <= 0 {
		endCtime = time.Now().Unix()
	}
	formatedBgnCtime := time.Unix(bgnCtime, 0)
	formatedEndCtime := time.Unix(endCtime, 0)
	assistLogs, err = assSvc.Logs(c, mid, assistMid, formatedBgnCtime, formatedEndCtime, int((pn-1)*ps), int(ps))
	if err != nil {
		log.Error("assistSvc.AssistLogs(%v) error(%v)", assistLogs, err)
		return
	}
	total, err = assSvc.LogCnt(c, mid, assistMid, formatedBgnCtime, formatedEndCtime)
	if err != nil {
		log.Error("assSvc.LogCnt: mid (%d),assistMid (%d),bgnctime (%v),endctime (%v):error(%v)", mid, assistMid, formatedBgnCtime, formatedEndCtime, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    assistLogs,
		"pager": map[string]int64{
			"current": pn,
			"size":    ps,
			"total":   total,
		},
	}))
}

func assistLogAdd(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tpStr := params.Get("type")
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	actStr := params.Get("action")
	act, err := strconv.ParseInt(actStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", act, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	subIDStr := params.Get("subject_id")
	subID, err := strconv.ParseInt(subIDStr, 10, 64)
	if err != nil || subID <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", subIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	objIDStr := params.Get("object_id")
	if len(objIDStr) == 0 {
		log.Error("objIDStr length eq zero(%s)", objIDStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	detail := params.Get("detail")
	if len(detail) == 0 {
		log.Error("detail len is zero (%s) error(%v)", detail, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = assSvc.AddLog(c, mid, assistMid, tp, act, subID, objIDStr, detail); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"mid":        mid,
		"assist_mid": assistMid,
		"type":       tp,
		"action":     act,
		"subject_id": subID,
		"object_id":  objIDStr,
	}, nil)
}

func assistLogCancel(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	logIDStr := params.Get("log_id")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil || logID <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", logIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, _ := strconv.ParseInt(assistMidStr, 10, 64)
	if assistMid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := assSvc.CancelLog(c, mid, assistMid, logID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"mid":        mid,
		"assist_mid": assistMid,
		"log_id":     logID,
	}, nil)
}

func assistLogInfo(c *bm.Context) {
	params := c.Request.Form
	logIDStr := params.Get("log_id")
	logID, err := strconv.ParseInt(logIDStr, 10, 64)
	if err != nil || logID <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", logIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, _ := strconv.ParseInt(assistMidStr, 10, 64)
	if assistMid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	logInfo, err := assSvc.LogInfo(c, logID, mid, assistMid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("assSvc.Assist(%s) error(%v)|logId(%d)|mid(%d)|assistMid(%d)", logInfo, err, logID, mid, assistMid)
		return
	}
	c.JSON(logInfo, nil)
}

func assistUps(c *bm.Context) {
	params := c.Request.Form
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", assistMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil || ps <= 20 {
		ps = 20
	}
	assistUpsPager, err := assSvc.AssistUps(c, assistMid, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "",
		"data":    assistUpsPager.Data,
		"pager":   assistUpsPager.Pager,
	}))
}

func assistLogObj(c *bm.Context) {
	params := c.Request.Form
	objIDStr := params.Get("object_id")
	objID, err := strconv.ParseInt(objIDStr, 10, 64)
	if err != nil || objID <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", objIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid <= 0 {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tpStr := params.Get("type")
	tp, err := strconv.ParseInt(tpStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	actStr := params.Get("action")
	act, _ := strconv.ParseInt(actStr, 10, 64)
	if act < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", act, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	logInfo, err := assSvc.LogObj(c, mid, objID, tp, act)
	if err != nil {
		c.JSON(nil, err)
		log.Error("assSvc.LogObj(%s) error(%v)|mid(%d)|logId(%d)|tp(%d)|act(%d)", logInfo, err, mid, objID, tp, act)
		return
	}
	c.JSON(logInfo, nil)
}
