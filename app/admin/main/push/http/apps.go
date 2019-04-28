package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func appList(c *bm.Context) {
	var (
		err   error
		items []*model.App
	)
	if err = pushSrv.DB.Find(&items).Error; err != nil {
		log.Error("appList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(items, nil)
}

func addApp(c *bm.Context) {
	app := new(model.App)
	if err := c.Bind(app); err != nil {
		return
	}
	if !pushSrv.DB.Where("name=?", app.Name).First(&model.App{}).RecordNotFound() {
		log.Warn("addApp(%+v) repeat", app)
		c.JSON(nil, ecode.PushRecordRepeatErr)
		return
	}
	if err := pushSrv.DB.Create(app).Error; err != nil {
		log.Error("addApp(%s) error(%v)", app.Name, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func appInfo(c *bm.Context) {
	var (
		req  = c.Request.Form
		info = &model.App{}
	)
	id, _ := strconv.ParseInt(req.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.First(info, id).Error; err != nil {
		log.Error("appInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

func saveApp(c *bm.Context) {
	app := new(model.App)
	if err := c.Bind(app); err != nil {
		return
	}
	if err := pushSrv.DB.Model(&model.App{ID: app.ID}).Updates(map[string]interface{}{"name": app.Name, "push_limit_user": app.PushLimitUser}).Error; err != nil {
		log.Error("saveApp(%d,%s) error(%v)", app.ID, app.Name, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func delApp(c *bm.Context) {
	req := c.Request.Form
	id, _ := strconv.ParseInt(req.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Model(&model.App{ID: id}).Update("dtime", time.Now().Unix()).Error; err != nil {
		log.Error("delApps(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
