package http

import (
	"strconv"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func aiConfig(c *bm.Context) {
	c.JSON(svc.AiConfig(c), nil)
}

func aiWhite(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error
		pnStr  = params.Get("pn")
		list   = make([]*model.AiWhite, 0)
		page   = &model.Page{}
		pn, ps int
		total  int64
		datas  struct {
			Data []*model.AiWhite `json:"data"`
			Page *model.Page      `json:"page"`
		}
	)
	ps = 50
	if pn, err = strconv.Atoi(pnStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if list, total, err = svc.AiWhite(c, pn, ps); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	page.Num = pn
	page.Size = ps
	page.Total = total
	datas.Data = list
	datas.Page = page
	c.JSON(datas, nil)
}

func aiWhiteAdd(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error
		midstr = params.Get("mid")
		mid    int64
	)
	if mid, err = strconv.ParseInt(midstr, 10, 64); mid <= 0 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AiWhiteAdd(c, mid))
}
func aiWhiteEdit(c *bm.Context) {
	var (
		params     = c.Request.Form
		err        error
		midStr     = params.Get("mid")
		stateStr   = params.Get("state")
		mid, state int64
	)
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if state, err = strconv.ParseInt(stateStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AiWhiteEdit(c, mid, int8(state)))
}

func aiCaseScore(c *bm.Context) {
	var (
		params  = c.Request.Form
		content = params.Get("content")
	)
	c.JSON(svc.AiScore(c, content))
}

func aiCaseAdd(c *bm.Context) {
	var (
		params        = c.Request.Form
		err           error
		sourceStr     = params.Get("source")
		content       = params.Get("content")
		typeStr       = params.Get("type")
		aiCase        = &model.AiCase{}
		source, type1 int64
	)
	source, err = strconv.ParseInt(sourceStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	type1, err = strconv.ParseInt(typeStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if content == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aiCase.Source = int8(source)
	aiCase.Content = content
	aiCase.Type = int8(type1)
	c.JSON(nil, svc.AiCaseAdd(c, aiCase))
}
