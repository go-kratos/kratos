package http

import (
	"go-common/app/interface/main/account/model"
	cardv1 "go-common/app/service/main/card/api/grpc/v1"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func userCard(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(cardSvc.UserCard(c, mid.(int64)))
}

func cardInfo(c *bm.Context) {
	var err error
	arg := new(model.ArgCardID)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(cardSvc.Card(c, arg.ID))
}

func cardHots(c *bm.Context) {
	c.JSON(cardSvc.CardHots(c))
}

func cardGroups(c *bm.Context) {
	var mid int64
	midi, exists := c.Get("mid")
	if exists {
		mid = midi.(int64)
	}
	c.JSON(cardSvc.AllGroup(c, mid))
}

func cardsByGid(c *bm.Context) {
	var err error
	arg := new(model.ArgGroupID)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(cardSvc.CardsByGid(c, arg.ID))
}

func equip(c *bm.Context) {
	var err error
	arg := new(model.ArgCardID)
	if err = c.Bind(arg); err != nil {
		return
	}
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(nil, cardSvc.Equip(c, &cardv1.ModelArgEquip{Mid: mid.(int64), CardId: arg.ID}))
}

func demount(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(nil, cardSvc.Demount(c, mid.(int64)))
}
