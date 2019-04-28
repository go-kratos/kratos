package http

import (
	"strconv"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func article(c *bm.Context) {
	var (
		mid          int64
		pn, ps, sort int
		ok           bool
		err          error
	)
	params := c.Request.Form
	midStr := params.Get("mid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	sortStr := params.Get("sort")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxArticlePs {
		ps = conf.Conf.Rule.MaxArticlePs
	}
	if sortStr != "" {
		if sort, ok = model.ArticleSortType[sortStr]; !ok {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		sort = artmdl.FieldDefault
	}
	c.JSON(spcSvc.Article(c, mid, pn, ps, sort))
}
