package http

import (
	"go-common/app/interface/main/account/model"
	v1 "go-common/app/service/main/vip/api"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func bindInfoByMid(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	a := new(model.ArgBindInfo)
	if err = c.Bind(a); err != nil {
		return
	}
	a.Mid = mid.(int64)
	a.AppID = vipmol.EleAppID
	c.JSON(vipSvc.BindInfoByMid(c, a))
}

func createAssociateOrder(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	a := new(model.ArgCreateAssociateOrder)
	if err = c.Bind(a); err != nil {
		return
	}
	a.Mid = mid.(int64)
	a.IP = metadata.String(c, metadata.RemoteIP)
	c.JSON(vipSvc.CreateAssociateOrder(c, a))
}

func associatePanel(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	a := new(vipmol.ArgAssociatePanel)
	if err = c.Bind(a); err != nil {
		return
	}
	a.Mid = mid.(int64)
	a.IP = metadata.String(c, metadata.RemoteIP)
	res := new(struct {
		PriceList []*v1.AssociatePanelInfo `json:"price_list"`
	})
	res.PriceList, err = vipSvc.AssociatePanel(c, a)
	if res.PriceList == nil {
		res.PriceList = []*v1.AssociatePanelInfo{}
	}
	c.JSON(res, err)
}

func redpackets(c *bm.Context) {
	c.JSON(vipSvc.EleRedPackages(c))
}

func specailfoods(c *bm.Context) {
	c.JSON(vipSvc.EleSpecailFoods(c))
}

func actlimit(ctx *bm.Context) {
	var mid int64
	midi, exists := ctx.Get("mid")
	if exists {
		mid = midi.(int64)
	}
	if err := vipSvc.ActivityTimeLimit(mid); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}

func iplimit(ctx *bm.Context) {
	req := ctx.Request
	params := req.Form
	sappkey := params.Get("appkey")
	ip := metadata.String(ctx, metadata.RemoteIP)
	if err := vipSvc.ActivityWhiteIPLimit(sappkey, ip); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}

func openlimit(ctx *bm.Context) {
	req := ctx.Request
	params := req.Form
	outOpenID := params.Get("out_open_id")
	if err := vipSvc.ActivityWhiteOutOpenIDLimit(outOpenID); err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
}
