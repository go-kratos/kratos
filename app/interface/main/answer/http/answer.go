package http

import (
	"fmt"
	"strconv"
	"strings"

	"go-common/app/interface/main/answer/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func localized(c *bm.Context) string {
	langs := detectLocalizedWeb(c)
	if len(langs) == 0 {
		return model.LangZhCN
	}
	switch langs[0] {
	case model.LangZhTW:
		return model.LangZhTW
	case model.LangZhHK:
		return model.LangZhTW
	default:
		return model.LangZhCN
	}
}

func checkBirthDay(c *bm.Context) {
	var mid, _ = c.Get("mid")
	if ok := svc.CheckBirthday(c, mid.(int64)); !ok {
		c.JSON(nil, ecode.MemberBirthdayInfoIsNull)
		return
	}
	c.JSON(nil, nil)
}

// checkPro check second step answers
func checkPro(c *bm.Context) {
	var (
		err     error
		ids     []int64
		params  = c.Request.Form
		qIds    = params.Get("qs_ids")
		mid, _  = c.Get("mid")
		ansHash = make(map[int64]string)
	)
	qidArr := strings.Split(qIds, ",")
	for _, qid := range qidArr {
		id, _ := strconv.ParseInt(qid, 10, 64)
		ansHash[id] = params.Get("ans_hash_" + qid)
		ids = append(ids, id)
	}
	hid, err := svc.ProCheck(c, mid.(int64), ids, ansHash, localized(c))
	if err != nil {
		log.Error("svc.ProCheck(%d,%s,%+v) error(%+v)", mid.(int64), xstr.JoinInts(ids), ansHash, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(fmt.Sprintf(model.ProPassed, hid), nil)
}

// checkBase check first step answers
func checkBase(c *bm.Context) {
	var (
		err     error
		ids     []int64
		params  = c.Request.Form
		qIds    = params.Get("qs_ids")
		mid, _  = c.Get("mid")
		ansHash = make(map[int64]string)
	)
	qidsArr := strings.Split(qIds, ",")
	for _, qid := range qidsArr {
		id, _ := strconv.ParseInt(qid, 10, 64)
		ansHash[id] = params.Get("ans_hash_" + qid)
		ids = append(ids, id)
	}
	req, err := svc.CheckBase(c, mid.(int64), ids, ansHash, localized(c))
	if err != nil {
		log.Error("svc.BaseCheck(%d,%s,%+v) error(%+v)", mid.(int64), xstr.JoinInts(ids), ansHash, err)
		c.JSON(nil, err)
		return
	}
	res := make(map[string]interface{})
	if req != nil && len(req.QidList) > 0 {
		res["next"] = false
		res["ids"] = req.QidList
		c.JSON(res, nil)
		return
	}
	res["next"] = true
	c.JSON(res, nil)
}

// checkExtra extra question check.
func checkExtra(c *bm.Context) {
	var (
		ids     []int64
		params  = c.Request.Form
		qIds    = params.Get("qs_ids")
		mid, ok = c.Get("mid")
		ansHash = make(map[int64]string)
		ua      = c.Request.Header.Get("User-Agent")
		refer   = c.Request.Header.Get("Referer")
	)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	qidsArr := strings.Split(qIds, ",")
	for _, qid := range qidsArr {
		id, _ := strconv.ParseInt(qid, 10, 64)
		h := params.Get("ans_hash_" + qid)
		if h != "" {
			ansHash[id] = params.Get("ans_hash_" + qid)
			ids = append(ids, id)
		}
	}
	buvid := c.Request.Header.Get("Buvid")
	if buvid == "" {
		cookie, _ := c.Request.Cookie("buvid3")
		if cookie != nil {
			buvid = cookie.Value
		}
	}
	c.JSON(nil, svc.ExtraCheck(c, mid.(int64), ids, ansHash, ua, localized(c), refer, buvid))
}

// getBaseQ get first step questions
func baseQus(c *bm.Context) {
	var (
		mid, ok = c.Get("mid")
		mobile  = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svc.BaseQ(c, mid.(int64), localized(c), mobile))
}

// getProType get second step question types
func proType(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svc.ProType(c, mid.(int64), localized(c)))
}

// getQstByType get second step questions
func proQus(c *bm.Context) {
	var (
		params  = c.Request.Form
		mid, ok = c.Get("mid")
		tIdsStr = params.Get("type_ids")
		mobile  = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svc.ConvertProQues(c, mid.(int64), tIdsStr, localized(c), mobile))
}

// extraQus extra question.
func extraQus(c *bm.Context) {
	var (
		mid, ok = c.Get("mid")
		mobile  = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svc.ConvertExtraQs(c, mid.(int64), localized(c), mobile))
}

func cool(c *bm.Context) {
	var (
		err    error
		mid    int64
		hid    int64
		params = c.Request.Form
		hidStr = params.Get("id")
	)
	if midI, ok := c.Get("mid"); ok {
		mid = midI.(int64)
	}
	if hid, err = strconv.ParseInt(hidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.Cool(c, hid, mid))
}

func extraScore(c *bm.Context) {
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(svc.ExtraScore(c, mid.(int64)))
}
