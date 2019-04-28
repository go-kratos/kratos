package http

import (
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strings"
)

func synonymList(c *bm.Context) {
	var (
		err   error
		total int64
		st    []*model.SynonymTag
		param = new(model.ParamSynonymList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Pn < 1 {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if st, total, err = svc.SynonymList(c, param.Keyword, param.Pn, param.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["synonym"] = st
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	c.JSON(data, nil)
}

func synonymEdit(c *bm.Context) {
	var (
		err      error
		username string
		// mid      int64
		param = new(model.ParamSynonymEdit)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, username = managerInfo(c)
	if username == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.SynonymAdd(c, username, param.TName, param.Adverb))
}

func synonymInfo(c *bm.Context) {
	var (
		err      error
		stagInfo *model.SynonymInfo
		param    = new(struct {
			Tid int64 `form:"tid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if stagInfo, err = svc.SynonymInfo(c, param.Tid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(stagInfo, nil)
}

func synonymDel(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamSynonymDel)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if len(param.Adverb) <= 0 {
		err = svc.SynonymDelete(c, param.Tid)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.RemoveSynonymSon(c, param.Tid, param.Adverb))
}

func synonymIsExist(c *bm.Context) {
	var (
		err   error
		tag   *model.Tag
		param = new(model.ParamSynonymExist)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if strings.Compare(param.TName, param.Adverb) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tag, err = svc.SynonymIsExist(c, param.Adverb); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 3)
	data["id"] = tag.ID
	data["name"] = tag.Name
	data["state"] = tag.State
	c.JSON(data, nil)
}
