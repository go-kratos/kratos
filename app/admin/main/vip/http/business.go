package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_defpn = 1
	_defps = 10
)

func business(c *bm.Context) {
	var (
		err error
		r   *model.VipBusinessInfo
	)
	arg := new(model.ArgID)
	if err = c.Bind(arg); err != nil {
		return
	}
	if r, err = vipSvc.BusinessInfo(c, int(arg.ID)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(r, nil)
}

func updateBusiness(c *bm.Context) {
	var (
		err error
		arg = new(model.VipBusinessInfo)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vipSvc.UpdateBusinessInfo(c, arg))
}

func addBusiness(c *bm.Context) {
	var (
		err error
		arg = new(model.VipBusinessInfo)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.AddBusinessInfo(c, arg))
}

func businessList(c *bm.Context) {
	var (
		infos []*model.VipBusinessInfo
		total int64
		err   error
		arg   = new(model.ArgPage)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if arg.Pn == 0 {
		arg.Pn = _defpn
	}
	if arg.Ps == 0 {
		arg.Ps = _defps
	}
	if infos, total, err = vipSvc.BusinessList(c, arg.Pn, arg.Ps, arg.Status); err != nil {
		c.JSON(nil, err)
		return
	}
	res := new(struct {
		Data  []*model.VipBusinessInfo `json:"data"`
		Total int64                    `json:"total"`
	})
	res.Data = infos
	res.Total = total
	c.JSON(res, nil)
}
