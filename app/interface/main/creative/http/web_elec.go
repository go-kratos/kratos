package http

import (
	"go-common/app/interface/main/creative/model/elec"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func webUserElec(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	elecUser, err := elecSvc.UserInfo(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(elecUser, nil)
}

func webUserElecUpdate(c *bm.Context) {
	params := c.Request.Form
	stateStr := params.Get("state")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	state, err := strconv.ParseInt(stateStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", stateStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	elecUser, err := elecSvc.UserUpdate(c, mid, int8(state), ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(elecUser, nil)
}

func webArcElecUpdate(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	stateStr := params.Get("state")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.ParseInt(stateStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", stateStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, elecSvc.ArcUpdate(c, mid, aid, int8(state), ip))
}

func webElecNotify(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	notify, err := elecSvc.Notify(c, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(notify, nil)
}

func webElecStatus(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	st, err := elecSvc.Status(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(st, nil)
}

func webElecUpStatus(c *bm.Context) {
	params := c.Request.Form
	spdayStr := params.Get("display_specialday")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	spday, err := strconv.ParseInt(spdayStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", spdayStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, elecSvc.UpStatus(c, mid, int(spday), ip))
}

func webElecRecentRank(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	sizeStr := params.Get("size")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	size, err := strconv.ParseInt(sizeStr, 10, 8)
	if err != nil || size == 0 || size > 100 { //返回条数 （最大100，不传默认50）
		size = 50
	}
	recRank, err := elecSvc.RecentRank(c, mid, size, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string][]*elec.Rank{
		"list": recRank,
	}, nil)
}

func webElecCurrentRank(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	curRank, err := elecSvc.CurrentRank(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string][]*elec.Rank{
		"list": curRank,
	}, nil)
}

func webElecTotalRank(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tolRank, err := elecSvc.TotalRank(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string][]*elec.Rank{
		"list": tolRank,
	}, nil)
}

func webElecDailyBill(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	bg := params.Get("begin")
	end := params.Get("end")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 100 {
		ps = 100
	}
	bill, err := elecSvc.DailyBill(c, mid, pn, ps, bg, end, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(bill, nil)
}

func webElecBalance(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	bal, err := elecSvc.Balance(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(bal, nil)
}

func webRemarkList(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	bg := params.Get("begin")
	end := params.Get("end")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 12 {
		ps = 12
	}
	bill, err := elecSvc.RemarkList(c, mid, pn, ps, bg, end, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(bill, nil)
}

func webRemarkDetail(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := c.Request.Form
	idStr := params.Get("id")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rm, err := elecSvc.RemarkDetail(c, mid, id, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(rm, nil)
}

func webRemark(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("id")
	msg := params.Get("msg")
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", idStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	st, err := elecSvc.Remark(c, mid, id, msg, ak, ck, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(st, nil)
}

func webRecentElec(c *bm.Context) {
	req := c.Request
	params := req.Form
	ip := metadata.String(c, metadata.RemoteIP)
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if mid == 0 {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	pn, _ := strconv.Atoi(pageStr)
	if pn <= 0 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps <= 0 || ps > 20 {
		ps = 20
	}
	var rec *elec.RecentElecList
	elecStat, _ := elecSvc.UserState(c, mid, ip, ak, ck)
	if elecStat != nil && elecStat.State == "2" {
		rec, _ = elecSvc.RecentElec(c, mid, pn, ps, ip)
	}
	if rec == nil {
		rec = recl
	}
	c.JSON(rec, nil)
}
