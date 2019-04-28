package http

import (
	"encoding/json"
	"strconv"

	"go-common/app/admin/main/manager/conf"
	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func rankGroups(c *bm.Context) {
	form := c.Request.Form
	pn, _ := strconv.Atoi(form.Get("pn"))
	if pn < 1 {
		pn = 1
	}
	ps, _ := strconv.Atoi(form.Get("ps"))
	if ps < 1 || ps > conf.Conf.Cfg.RankGroupMaxPs {
		ps = conf.Conf.Cfg.RankGroupMaxPs
	}
	groups, total, err := mngSvc.RankGroups(c, pn, ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int{
		"page":     pn,
		"pagesize": ps,
		"total":    total,
	}
	c.JSONMap(map[string]interface{}{
		"pager": page,
		"data":  groups,
	}, err)
}

func rankGroup(c *bm.Context) {
	form := c.Request.Form
	gid, _ := strconv.ParseInt(form.Get("id"), 10, 64)
	if gid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("id unnarmal (%d)", gid)
		return
	}
	c.JSON(mngSvc.RankGroup(c, gid))
}

func addRankGroup(c *bm.Context) {
	form := c.Request.Form
	name := form.Get("name")
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("name is empty")
		return
	}
	desc := form.Get("desc")
	auths, _ := xstr.SplitInts(form.Get("auths"))
	g := &model.RankGroup{Name: name, Desc: desc}
	c.JSON(mngSvc.AddRankGroup(c, g, auths))
}

func updateRankGroup(c *bm.Context) {
	form := c.Request.Form
	gid, _ := strconv.ParseInt(form.Get("id"), 10, 64)
	if gid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("id unnarmal (%d)", gid)
		return
	}
	name := form.Get("name")
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("name is empty")
		return
	}
	desc := form.Get("desc")
	auths, _ := xstr.SplitInts(form.Get("auths"))
	g := &model.RankGroup{ID: gid, Name: name, Desc: desc}
	c.JSON(nil, mngSvc.UpdateRankGroup(c, g, auths))
}

func delRankGroup(c *bm.Context) {
	form := c.Request.Form
	gid, _ := strconv.ParseInt(form.Get("id"), 10, 64)
	if gid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("id unnarmal (%d)", gid)
		return
	}
	c.JSON(nil, mngSvc.DelRankGroup(c, gid))
}

func addRankUser(c *bm.Context) {
	form := c.Request.Form
	uid, _ := strconv.ParseInt(form.Get("uid"), 10, 64)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("uid unnarmal (%d)", uid)
		return
	}
	c.JSON(nil, mngSvc.AddRankUser(c, uid))
}

func rankUsers(c *bm.Context) {
	form := c.Request.Form
	pn, _ := strconv.Atoi(form.Get("pn"))
	if pn < 1 {
		pn = 1
	}
	ps, _ := strconv.Atoi(form.Get("ps"))
	if ps < 1 || ps > conf.Conf.Cfg.RankGroupMaxPs {
		ps = conf.Conf.Cfg.RankGroupMaxPs
	}
	un := form.Get("username")
	users, total, err := mngSvc.RankUsers(c, pn, ps, un)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int{
		"page":     pn,
		"pagesize": ps,
		"total":    total,
	}
	c.JSONMap(map[string]interface{}{
		"pager": page,
		"data":  users,
	}, err)
}

func saveRankUser(c *bm.Context) {
	form := c.Request.Form
	uid, _ := strconv.ParseInt(form.Get("uid"), 10, 64)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("uid unnormal (%d)", uid)
		return
	}
	ranks := form.Get("ranks")
	if ranks == "" {
		c.JSON(nil, ecode.RequestErr)
		log.Error("ranks is empty")
		return
	}
	rs := make(map[string]int)
	if err := json.Unmarshal([]byte(ranks), &rs); err != nil {
		c.JSON(nil, ecode.RequestErr)
		log.Error("ranks unnormal (%s) error(%v)", ranks, err)
		return
	}
	rMap := make(map[int64]int)
	for g, r := range rs {
		gid, _ := strconv.ParseInt(g, 10, 64)
		rMap[gid] = r
	}
	c.JSON(nil, mngSvc.SaveRankUser(c, uid, rMap))
}

func delRankUser(c *bm.Context) {
	form := c.Request.Form
	uid, _ := strconv.ParseInt(form.Get("uid"), 10, 64)
	if uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		log.Error("uid unnormal (%d)", uid)
		return
	}
	c.JSON(nil, mngSvc.DelRankUser(c, uid))
}
