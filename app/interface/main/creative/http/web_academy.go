package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func webAcademyTags(c *bm.Context) {
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	// check params
	tgs, _ := acaSvc.TagList(c)
	c.JSON(tgs, nil)
}

func webAcademyArchives(c *bm.Context) {
	params := c.Request.Form
	tidsStr := params.Get("tids")
	bsStr := params.Get("business")
	pageStr := params.Get("pn")
	psStr := params.Get("ps")
	keyword := params.Get("keyword")
	order := params.Get("order")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	var (
		tids []int64
		err  error
	)
	// check params
	if tidsStr != "" {
		if tids, err = xstr.SplitInts(tidsStr); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", tidsStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	bs, err := strconv.Atoi(bsStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", bsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", pageStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", psStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn == 0 {
		pn = 1
	}
	if ps > 20 {
		ps = 20
	}
	aca := &academy.EsParam{
		Tid:      tids,
		Business: bs,
		Pn:       pn,
		Ps:       ps,
		Keyword:  keyword,
		Order:    order,
		IP:       ip,
	}
	arcs, err := acaSvc.ArchivesWithES(c, aca)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arcs, nil)
}

func webAddFeedBack(c *bm.Context) {
	params := c.Request.Form
	category := params.Get("category")
	course := params.Get("course")
	suggest := params.Get("suggest")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	id, err := acaSvc.AddFeedBack(c, category, course, suggest, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}
