package http

import (
	"go-common/app/interface/main/creative/model/appeal"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
	"strings"
)

func webAppealContact(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	cookie := c.Request.Header.Get("cookie")
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	ct, err := apSvc.PhoneEmail(c, cookie, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]string{
		"phone": ct.TelPhone,
		"email": ct.Email,
	}, nil)
}

func webAppealList(c *bm.Context) {
	params := c.Request.Form
	state := params.Get("state")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	// check
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 10
	}
	all, open, closed, aps, err := apSvc.List(c, mid, pn, ps, state, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"pager": map[string]int{
			"current":      pn,
			"size":         ps,
			"total":        all,
			"open_count":   open,
			"closed_count": closed,
		},
		"appeals": aps,
	}, nil)
}

func webAppealDetail(c *bm.Context) {
	params := c.Request.Form
	apidStr := params.Get("apid")
	// check params
	apid, err := strconv.ParseInt(apidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ap, err := apSvc.Detail(c, mid, apid, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		log.Error("apSvc.Detail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if ap == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(ap, nil)
}

func webAppealDown(c *bm.Context) {
	params := c.Request.Form
	apidStr := params.Get("apid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	apid, err := strconv.ParseInt(apidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	var is bool
	if is, err = checkStateAndMID(c, mid, apid, ip); err != nil {
		log.Error("checkStateAndMID error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if !is {
		log.Error("checkStateAndMID not your appeal (%v)", is)
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(nil, apSvc.State(c, mid, apid, appeal.StateUserClosed, ip))
}

func webAppealAdd(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	content := params.Get("content")
	qq := params.Get("qq")
	pics := params.Get("pics")
	phone := params.Get("phone")
	email := params.Get("email")
	typeidStr := params.Get("typeid")
	title := params.Get("title")
	desc := params.Get("desc")
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", typeidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ap := &appeal.BusinessAppeal{
		BusinessTypeID:  tid,
		BusinessMID:     mid,
		BusinessTitle:   title,
		BusinessContent: desc,
	}
	apid, err := apSvc.Add(c, mid, aid, qq, phone, email, content, strings.Replace(pics, ";", ",", -1), ip, ap)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if apid > 0 {
		var is bool
		if is, err = checkStateAndMID(c, mid, apid, ip); err != nil {
			log.Error("checkStateAndMID error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if !is {
			log.Error("checkStateAndMID not your appeal (%v)", is)
			c.JSON(nil, ecode.NothingFound)
			return
		}
		c.JSON(nil, apSvc.Reply(c, mid, apid, appeal.ReplySystemEvent, appeal.ReplyMsg, "", ip))
	}
}

func webAppealReply(c *bm.Context) {
	params := c.Request.Form
	apidStr := params.Get("apid")
	content := params.Get("content")
	pics := params.Get("pics")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	apid, err := strconv.ParseInt(apidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", apidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var is bool
	if is, err = checkStateAndMID(c, mid, apid, ip); err != nil {
		log.Error("checkStateAndMID error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if !is {
		log.Error("checkStateAndMID not your appeal (%v)", is)
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(nil, apSvc.Reply(c, mid, apid, appeal.ReplyUserEvent, content, pics, ip))
}

func webAppealStar(c *bm.Context) {
	params := c.Request.Form
	apidStr := params.Get("apid")
	starStr := params.Get("star")
	ip := metadata.String(c, metadata.RemoteIP)
	// check params
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	apid, err := strconv.ParseInt(apidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", apidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	star, err := strconv.ParseInt(starStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", starStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if star < 0 || star > 3 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var is bool
	if is, err = checkMID(c, mid, apid, ip); err != nil {
		log.Error("star checkMID error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if !is {
		log.Error("star checkMID not your appeal (%v)", is)
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(nil, apSvc.Star(c, mid, apid, star, ip))
}

func checkStateAndMID(c *bm.Context, mid, apid int64, ip string) (is bool, err error) {
	ap, err := apSvc.Detail(c, mid, apid, ip)
	if err != nil || ap == nil {
		return
	}
	if appeal.IsClosed(ap.State) {
		err = ecode.NothingFound
		return
	}
	if ap.Mid == mid {
		is = true
	}
	return
}

func checkMID(c *bm.Context, mid, apid int64, ip string) (is bool, err error) {
	ap, err := apSvc.Detail(c, mid, apid, ip)
	if err != nil || ap == nil {
		return
	}
	if ap.Mid == mid {
		is = true
	}
	return
}
