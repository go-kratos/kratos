package http

import (
	"go-common/app/service/main/usersuit/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func medalInfo(c *bm.Context) {
	var (
		err  error
		info *model.MedalInfo
		arg  = new(model.ArgMIDNID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if info, err = usersuitSvc.MedalInfo(c, arg.MID, arg.NID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

func medalGet(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMIDNID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = usersuitSvc.MedalGet(c, arg.MID, arg.NID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func medalCheck(c *bm.Context) {
	var (
		err  error
		info *model.MedalCheck
		arg  = new(model.ArgMIDNID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if info, err = usersuitSvc.MedalCheck(c, arg.MID, arg.NID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

func medalActivated(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(usersuitSvc.MedalActivated(c, arg.MID))
}

func medalMy(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(usersuitSvc.MedalMyInfo(c, arg.MID))
}

func medalAllInfo(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(usersuitSvc.MedalAllInfo(c, arg.MID))
}

func medalPopup(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(usersuitSvc.MedalPopup(c, arg.MID))
}

func medalInstall(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMedalInstall)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, usersuitSvc.MedalInstall(c, arg.Mid, arg.Nid, arg.IsActivated))
}

func medalUser(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgMID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(usersuitSvc.MedalUserInfo(c, arg.MID, metadata.String(c, metadata.RemoteIP)))
}
