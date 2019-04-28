package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/article"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func creatorArticlePre(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var (
		isAuthor int
		reason   string
	)
	ia, err := artSvc.IsAuthor(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if ia {
		isAuthor = 1
	} else {
		isAuthor = 0
		reason = "您还未开通专栏权限，请先在PC上进行申请"
	}
	c.JSON(map[string]interface{}{
		"is_author":  isAuthor,
		"reason":     reason,
		"submit_url": "https://member.bilibili.com/article-text/mobile",
	}, nil)
}

func creatorArticleList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	sortStr := params.Get("sort")
	groupStr := params.Get("group")
	categoryStr := params.Get("category")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	sort, err := strconv.Atoi(sortStr)
	if err != nil || sort < 0 {
		sort = 0
	}
	group, err := strconv.Atoi(groupStr)
	if err != nil || group < 0 {
		group = 0
	}
	category, err := strconv.Atoi(categoryStr)
	if err != nil || category < 0 {
		category = 0
	}
	arts, err := artSvc.Articles(c, mid, int(pn), int(ps), sort, group, category, ip)
	if ecode.Cause(err) == ecode.ArtCreationNoPrivilege {
		err = nil
	}
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if arts == nil {
		arts = &article.ArtList{
			Articles: []*article.Meta{},
			Type:     &model.CreationArtsType{},
			Page:     &model.ArtPage{},
		}
	}
	c.JSON(arts, nil)
}

func creatorArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	art, err := artSvc.View(c, aid, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(art, nil)
}

func creatorDelArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSvc.DelArticle(c, aid, mid, ip))
}

func creatorWithDrawArticle(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, artSvc.WithDrawArticle(c, aid, mid, ip))
}

func creatorDraftList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	arts, err := artSvc.Drafts(c, mid, pn, ps, ip)
	if ecode.Cause(err) == ecode.ArtCreationNoPrivilege {
		err = nil
	}
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if arts == nil {
		arts = &article.DraftList{
			Drafts: []*article.Meta{},
			Page:   &model.ArtPage{},
		}
	}
	c.JSON(arts, nil)
}
