package http

import (
	"strings"

	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func preBusinss(b *business.Business) (invalid bool) {
	b.Name = common.FilterBusinessName(strings.TrimSpace(b.Name))
	b.Desc = strings.TrimSpace(b.Desc)
	emails := strings.Split(b.Developer, ",")
	if len(emails) > 5 || b.Name == "" || b.Desc == "" || b.TP <= 0 || b.TP > 4 {
		invalid = true
		return
	}

	return
}

// addBusiness .
func addBusiness(c *bm.Context) {
	b := &business.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	if preBusinss(b) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	b.UID = uid(c)
	c.JSON(srv.AddBusiness(c, b))
}

// updateBusiness .
func updateBusiness(c *bm.Context) {
	b := &business.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.ID <= 0 || preBusinss(b) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	b.UID = uid(c)
	c.JSON(nil, srv.UpdateBusiness(c, b))
}

// setBusiness .
func setBusinessState(c *bm.Context) {
	b := &business.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	b.UID = uid(c)
	c.JSON(nil, srv.SetBusinessState(c, b))
}

// getBusiness .
func getBusiness(c *bm.Context) {
	b := &business.Business{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := srv.Business(c, b)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(res, nil)
}

// getBusinessList .
func getBusinessList(c *bm.Context) {
	ids := getAccessBiz(c)
	res, err := srv.BusinessList(c, ids, false)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(res, nil)
}

// getBusinessList .
func getBusinessEnable(c *bm.Context) {
	ids := getAccessBiz(c)
	res, err := srv.BusinessList(c, ids, true)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

// addBizCFG
func addBizCFG(c *bm.Context) {
	b := &business.BizCFG{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.BusinessID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err := srv.AddBizCFG(c, b)
	if err != nil {
		httpCode(c, err.Error(), ecode.RequestErr)
		return
	}
	c.JSON(id, nil)
}

// updateBizCFG
func updateBizCFG(c *bm.Context) {
	b := &business.BizCFG{}
	if err := c.Bind(b); err != nil {
		return
	}
	if b.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.UpdateBizCFG(c, b))
}

// listBizCFGs
func listBizCFGs(c *bm.Context) {
	b := &business.BizCFG{}
	if err := c.Bind(b); err != nil {
		return
	}
	res, err := srv.ListBizCFGs(c, b.BusinessID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(res, nil)
}

// 保留字配置
func reserveCFG(c *bm.Context) {
	opt := new(struct {
		BizID int64 `form:"business_id" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.ReserveCFG(c, opt.BizID))
}
