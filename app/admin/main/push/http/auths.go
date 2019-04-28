package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func authList(c *bm.Context) {
	var (
		req   = c.Request.Form
		auths []*model.Auth
	)
	appID, _ := strconv.ParseInt(req.Get("app_id"), 10, 64)
	if err := pushSrv.DB.Model(&model.App{ID: appID}).Related(&auths).Error; err != nil {
		log.Error("authList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(auths, nil)
}

func addAuth(c *bm.Context) {
	auth := new(model.Auth)
	if err := c.Bind(auth); err != nil {
		return
	}
	if err := pushSrv.DB.Create(auth).Error; err != nil {
		log.Error("addAuth(%+v) error(%v)", auth, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func authInfo(c *bm.Context) {
	auth := new(model.Auth)
	if err := c.Bind(auth); err != nil {
		return
	}
	if err := pushSrv.DB.First(auth, auth.ID).Error; err != nil {
		log.Error("authInfo(%d) error(%v)", auth.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(auth, nil)
}

func saveAuth(c *bm.Context) {
	auth := new(model.Auth)
	if err := c.Bind(auth); err != nil {
		return
	}
	if err := pushSrv.DB.Model(&model.Auth{ID: auth.ID}).Update(auth).Error; err != nil {
		log.Error("saveAuth(%+v) error(%v)", auth, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func delAuth(c *bm.Context) {
	id, _ := strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Model(&model.Auth{ID: id}).Update("dtime", time.Now().Unix()).Error; err != nil {
		log.Error("delAuth(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
