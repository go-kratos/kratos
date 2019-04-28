package http

import (
	"go-common/app/service/main/account/model"
	bm "go-common/library/net/http/blademaster"
)

// info
func info(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	info, err := accSvc.Info(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

// infoByName
func infoByName(c *bm.Context) {
	p := new(model.ParamNames)
	if err := c.Bind(p); err != nil {
		return
	}
	infos, err := accSvc.InfosByName(c, p.Names)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(infos, nil)
}

// infos
func infos(c *bm.Context) {
	p := new(model.ParamMids)
	if err := c.Bind(p); err != nil {
		return
	}
	infos, err := accSvc.Infos(c, p.Mids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(infos, nil)
}

// card
func card(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	card, err := accSvc.Card(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(card, nil)
}

// cards
func cards(c *bm.Context) {
	p := new(model.ParamMids)
	if err := c.Bind(p); err != nil {
		return
	}
	cards, err := accSvc.Cards(c, p.Mids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(cards, nil)
}

// vip
func vip(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	v, err := accSvc.Vip(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(v, nil)
}

func vips(c *bm.Context) {
	p := new(model.ParamMids)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(accSvc.Vips(c, p.Mids))
}

// profile
func profile(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	pfl, err := accSvc.Profile(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pfl, nil)
}

// profileWithStat
func profileWithStat(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	pfl, err := accSvc.ProfileWithStat(c, p.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pfl, nil)
}

// privacy
func privacy(c *bm.Context) {
	p := new(model.ParamMid)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(accSvc.Privacy(c, p.Mid))
}
