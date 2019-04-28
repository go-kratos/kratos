package http

import (
	"strconv"
	"strings"
	"text/template"

	"go-common/app/interface/main/reply/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_remarkLength = 200
)

func adminSubjectMid(c *bm.Context) {
	var (
		err  error
		tp   int64
		oid  int64
		mid  int64
		adid int64
	)
	params := c.Request.Form
	tpStr := params.Get("type")
	oidStr := params.Get("oid")
	midStr := params.Get("mid")
	adidStr := params.Get("adid")
	remark := params.Get("remark")
	if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", midStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	err = rpSvr.AdminSubjectMid(c, adid, mid, oid, int8(tp), remark)
	c.JSON(nil, err)
}

func adminSubRegist(c *bm.Context) {
	var (
		err   error
		oid   int64
		tp    int64
		state int64
		mid   int64
	)

	params := c.Request.Form
	midStr := params.Get("mid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	stateStr := params.Get("state")
	appkey := params.Get("appkey")
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", midStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}

	if stateStr != "" {
		if state, err = strconv.ParseInt(stateStr, 10, 8); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", stateStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if tpStr != "" {
		if tp, err = strconv.ParseInt(tpStr, 10, 8); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	err = rpSvr.AdminSubRegist(c, oid, mid, int8(tp), int8(state), appkey)
	c.JSON(nil, err)
}

func adminSubject(c *bm.Context) {

	params := c.Request.Form
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
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
	sub, err := rpSvr.AdminGetSubject(c, oid, int8(tp))
	if err != nil {
		log.Warn("rpSvr.AdminGetSubjectState(oid%d,tp,%d)error(%v)", oid, int8(tp))
		c.JSON(nil, err)
		return
	}
	c.JSON(sub, nil)
}

// adminModifySubject modify subject state.
func adminSubjectState(c *bm.Context) {
	var (
		err error
		mid int64
	)

	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	stateStr := params.Get("state")
	midStr := params.Get("mid")
	remark := params.Get("remark")
	// check params
	remark = strings.TrimSpace(remark)
	rml := len([]rune(remark))
	if rml > _remarkLength {
		log.Warn("remark(%s) length %d, max %d", remark, rml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
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
	state, err := strconv.ParseInt(stateStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", stateStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", midStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	err = rpSvr.AdminSubjectState(c, adid, oid, mid, int8(tp), int8(state), remark)
	c.JSON(nil, err)
}

func adminAuditSub(c *bm.Context) {

}

// adminPassReply pass reply normal.
func adminPassReply(c *bm.Context) {
	var err error
	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	remark := params.Get("remark")
	// check params
	remark = strings.TrimSpace(remark)
	rml := len([]rune(remark))
	if rml > _remarkLength {
		log.Warn("remark(%s) length %d, max %d", remark, rml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	// remark = template.JSEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
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
	if len(oids) != len(rpIDs) {
		log.Warn("oids:%s rpIDs:%s ", oids, rpIDs)
		return
	}
	for i := 0; i < len(oids); i++ {
		err = rpSvr.AdminPass(c, adid, oids[i], rpIDs[i], int8(tp), remark)
	}
	c.JSON(nil, err)
}

// adminRecoverReply recover reply normal.
func adminRecoverReply(c *bm.Context) {
	var err error
	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	remark := params.Get("remark")
	// check params
	remark = strings.TrimSpace(remark)
	rml := len([]rune(remark))
	if rml > _remarkLength {
		log.Warn("remark(%s) length %d, max %d", remark, rml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpIDStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
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
	err = rpSvr.AdminRecover(c, adid, oid, rpID, int8(tp), remark)
	c.JSON(nil, err)
}

// adminEditReply edit reply content by admin.
func adminEditReply(c *bm.Context) {
	var err error
	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	msg := params.Get("message")
	remark := params.Get("remark")
	// check params
	msg = strings.TrimSpace(msg)
	ml := len([]rune(msg))
	if conf.Conf.Reply.MaxConLen < ml || ml < conf.Conf.Reply.MinConLen {
		log.Warn("content(%s) length %d, max %d, min %d", msg, ml, conf.Conf.Reply.MaxConLen, conf.Conf.Reply.MinConLen)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	msg = template.HTMLEscapeString(msg)
	remark = strings.TrimSpace(remark)
	rml := len([]rune(remark))
	if rml > _remarkLength {
		log.Warn("remark(%s) length %d, max %d", remark, rml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpIDStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
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
	err = rpSvr.AdminEdit(c, adid, oid, rpID, int8(tp), msg, remark)
	c.JSON(nil, err)
}

// report reason
const (
	TypeReportReasonOthers = 0
)

// _delReply delete reply, this call be limited to this file.
func _delReply(c *bm.Context, isReport bool) {
	var (
		err     error
		reason  int64
		freason int64
	)

	params := c.Request.Form
	adidStr := params.Get("adid")
	adname := params.Get("adname")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	morStr := params.Get("moral")
	notStr := params.Get("notify")
	remark := params.Get("remark")
	ftimeStr := params.Get("ftime")
	auditStr := params.Get("audit")
	reasonStr := params.Get("reason")
	content := params.Get("reason_content")
	freasonStr := params.Get("freason")
	// check params
	remark = strings.TrimSpace(remark)
	ml := len([]rune(remark))
	if ml > _remarkLength {
		log.Warn("remark (%s) length %d, max %d", remark, ml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tps, err := xstr.SplitInts(tpStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	moral, err := strconv.ParseInt(morStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	notify, err := strconv.ParseBool(notStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if reasonStr != "" {
		if reason, err = strconv.ParseInt(reasonStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", reasonStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	if freasonStr != "" {
		if freason, err = strconv.ParseInt(freasonStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", freasonStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}
	ftime, _ := strconv.ParseInt(ftimeStr, 10, 64)
	audit, _ := strconv.ParseInt(auditStr, 10, 8)
	if len(oids) == 0 || len(rpIDs) != len(oids) || len(tps) != len(oids) {
		log.Warn("admin del reply oids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if reason != TypeReportReasonOthers {
		content = ""
	}
	for i := 0; i < len(oids); i++ {
		if isReport {
			err = rpSvr.AdminDeleteByReport(c, adid, oids[i], rpIDs[i], ftime, int8(tps[i]), int(moral), notify, adname, remark, int8(audit), int8(reason), content, int8(freason))
		} else {
			err = rpSvr.AdminDelete(c, adid, oids[i], rpIDs[i], ftime, int8(tps[i]), int(moral), notify, adname, remark, int8(reason), int8(freason))
		}
	}
	c.JSON(nil, err)
}

// adminDelReply delete reply by admin.
func adminDelReply(c *bm.Context) {
	_delReply(c, false)
}

// adminDelReplyByReport delete a report reply by admin.
func adminDelReplyByReport(c *bm.Context) {
	_delReply(c, true)
}

func adminReportStateSet(c *bm.Context) {

	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	stateStr := params.Get("state")
	var adid int64
	if adidStr != "" {
		var err error
		adid, err = strconv.ParseInt(adidStr, 10, 64)
		if err != nil {
			log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
			err = ecode.RequestErr
			c.JSON(nil, err)
			return
		}
	}

	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tps, err := xstr.SplitInts(tpStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) == 0 || len(rpIDs) != len(oids) || len(tps) != len(oids) {
		log.Warn("admin set report oids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	state, err := strconv.ParseInt(stateStr, 10, 8)
	if err != nil {
		log.Warn("admin set report soids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		err = rpSvr.AdminReportStateSet(c, adid, oids[i], rpIDs[i], int8(tps[i]), int8(state))
	}
	c.JSON(nil, err)
}

func adminTransferReport(c *bm.Context) {

	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	auditStr := params.Get("audit")
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.SplitInts(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.SplitInts(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tps, err := xstr.SplitInts(tpStr)
	if err != nil {
		log.Warn("strconv.SplitInts(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) == 0 || len(rpIDs) != len(oids) || len(tps) != len(oids) {
		log.Warn("admin del reply oids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	audit, _ := strconv.ParseInt(auditStr, 10, 8)
	for i := 0; i < len(oids); i++ {
		err = rpSvr.AdminReportTransfer(c, adid, oids[i], rpIDs[i], int8(tps[i]), int8(audit))
	}
	c.JSON(nil, err)
}

// adminIgnoreReport ignore a report.
func adminIgnoreReport(c *bm.Context) {
	var err error
	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	auditStr := params.Get("audit")
	remark := params.Get("remark")
	// check params
	remark = strings.TrimSpace(remark)
	ml := len([]rune(remark))
	if ml > _remarkLength {
		log.Warn("remark (%s) length %d, max %d", remark, ml, 255)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tps, err := xstr.SplitInts(tpStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	audit, _ := strconv.ParseInt(auditStr, 10, 8)
	if len(oids) == 0 || len(rpIDs) != len(oids) || len(tps) != len(oids) {
		log.Warn("admin del reply oids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		err = rpSvr.AdminReportIgnore(c, adid, oids[i], rpIDs[i], int8(tps[i]), int8(audit), remark)
	}
	c.JSON(nil, err)
}

// adminAddTopReply add top reply
func adminAddTopReply(c *bm.Context) {

	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	actStr := params.Get("action")
	// check params
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oid, err := strconv.ParseInt(oidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpID, err := strconv.ParseInt(rpIDStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
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
	act, err := strconv.ParseInt(actStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(actStr :%s) err(%v)", actStr, err)
		act = 1
	}
	err = rpSvr.AdminAddTop(c, adid, oid, rpID, int8(tp), int8(act))
	c.JSON(nil, err)
}

func adminReportRecover(c *bm.Context) {
	var err error
	params := c.Request.Form
	adidStr := params.Get("adid")
	oidStr := params.Get("oid")
	tpStr := params.Get("type")
	rpIDStr := params.Get("rpid")
	auditStr := params.Get("audit")
	remark := params.Get("remark")
	// check params
	remark = strings.TrimSpace(remark)
	rml := len([]rune(remark))
	if rml > _remarkLength {
		log.Warn("remark(%s) length %d, max %d", remark, rml, _remarkLength)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	remark = template.HTMLEscapeString(remark)
	adid, err := strconv.ParseInt(adidStr, 10, 64)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", adidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	oids, err := xstr.SplitInts(oidStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", oidStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	rpIDs, err := xstr.SplitInts(rpIDStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", rpIDStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	tps, err := xstr.SplitInts(tpStr)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	audit, err := strconv.ParseInt(auditStr, 10, 8)
	if err != nil {
		log.Warn("strconv.ParseInt(%s) error(%v)", tpStr, err)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(oids) == 0 || len(rpIDs) != len(oids) || len(tps) != len(oids) {
		log.Warn("admin del reply oids: %v, rpIDs: %v, tps: %v", oids, rpIDs, tps)
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	for i := 0; i < len(oids); i++ {
		err = rpSvr.AdminReportRecover(c, adid, oids[i], rpIDs[i], int8(tps[i]), int8(audit), remark)
	}
	c.JSON(nil, err)
}
