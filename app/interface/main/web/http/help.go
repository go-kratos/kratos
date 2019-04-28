package http

import (
	"go-common/app/interface/main/web/conf"
	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func helpList(c *bm.Context) {
	var (
		rs  []*model.HelpList
		err error
	)
	v := new(struct {
		PTypeID string `form:"parentTypeId" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if rs, err = webSvc.HelpList(c, v.PTypeID); err != nil {
		c.JSON(nil, ecode.Degrade)
		return
	}
	c.JSON(rs, nil)
}

func helpDetail(c *bm.Context) {
	var (
		total  int
		detail []*model.HelpDeatil
		list   []*model.HelpList
		err    error
	)
	v := new(struct {
		PTypeID string `form:"questionTypeId" validate:"required"`
		KeyFlag int    `form:"keyFlag" default:"1" validate:"min=1"`
		FID     string `form:"fId"`
		Pn      int    `form:"pn" default:"1" validate:"min=1"`
		Ps      int    `form:"ps" default:"15" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Ps > conf.Conf.Rule.MaxHelpPageSize {
		v.Ps = conf.Conf.Rule.MaxHelpPageSize
	}
	if detail, list, total, err = webSvc.HelpDetail(c, v.FID, v.PTypeID, v.KeyFlag, v.Pn, v.Ps); err != nil {
		c.JSON(nil, ecode.Degrade)
		log.Error("webSvc.HelpDetail(%s,%d,%d,%d) error(%v)", v.PTypeID, v.KeyFlag, v.Pn, v.Ps, err)
		return
	}
	data := make(map[string]interface{}, 2)
	rsDetail := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	rsDetail["items"] = detail
	rsDetail["page"] = page
	data["detail"] = rsDetail
	data["list"] = list
	c.JSON(data, nil)
}

func helpSearch(c *bm.Context) {
	var (
		total int
		list  []*model.HelpDeatil
		err   error
	)
	v := new(struct {
		PTypeID  string `form:"questionTypeId" default:"-1"`
		KeyWords string `form:"keyWords" validate:"required"`
		KeyFlag  int    `form:"keyFlag" default:"1" validate:"min=1"`
		Pn       int    `form:"pn" default:"1" validate:"min=1"`
		Ps       int    `form:"ps" default:"15" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Ps > conf.Conf.Rule.MaxHelpPageSize {
		v.Ps = conf.Conf.Rule.MaxHelpPageSize
	}
	if list, total, err = webSvc.HelpSearch(c, v.PTypeID, v.KeyWords, v.KeyFlag, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		log.Error("webSvc.HelpDetail(%s,%d,%d,%d) error(%v)", v.KeyWords, v.KeyFlag, v.Pn, v.Ps, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}
