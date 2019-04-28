package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func businessList(c *bm.Context) {
	var (
		items   []*model.Business
		err     error
		count   int
		apps    []*model.App
		appsMap = make(map[int64]*model.App)
	)
	pager := new(model.Pager)
	if err = c.Bind(pager); err != nil {
		return
	}
	if err = pushSrv.DB.Offset((pager.Pn - 1) * pager.Ps).Limit(pager.Ps).Find(&items).Error; err != nil {
		log.Error("businessList(%d,%d) error(%v)", pager.Pn, pager.Ps, err)
		c.JSON(nil, err)
		return
	}
	if err = pushSrv.DB.Find(&apps).Error; err != nil {
		log.Error("businessList(%d,%d) error(%v)", pager.Pn, pager.Ps, err)
		c.JSON(nil, err)
		return
	}
	for _, app := range apps {
		appsMap[app.ID] = app
	}
	for _, item := range items {
		if appsMap[item.AppID] != nil {
			item.AppName = appsMap[item.AppID].Name
		}
	}
	if err = pushSrv.DB.Model(&model.Business{}).Count(&count).Error; err != nil {
		log.Error("businessList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := map[string]interface{}{
		"data": items,
		"pager": &model.Pager{
			Pn:    pager.Pn,
			Ps:    pager.Ps,
			Total: count,
		},
	}
	c.JSONMap(data, nil)
}

func addBusiness(c *bm.Context) {
	biz := new(model.Business)
	if err := c.Bind(biz); err != nil {
		return
	}
	biz.Token = model.RandomString(32)
	if !pushSrv.DB.Where("app_id=? and name=?", biz.AppID, biz.Name).First(&model.Business{}).RecordNotFound() {
		log.Warn("addBusiness(%+v) repeat", biz)
		c.JSON(nil, ecode.PushRecordRepeatErr)
		return
	}
	if err := pushSrv.DB.Create(biz).Error; err != nil {
		log.Error("addBusiness(%+v) error(%v)", biz, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func businessInfo(c *bm.Context) {
	var (
		req  = c.Request.Form
		info = &model.Business{}
	)
	id, _ := strconv.ParseInt(req.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.First(info, id).Error; err != nil {
		log.Error("businessInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(info, nil)
}

func saveBusiness(c *bm.Context) {
	biz := new(model.Business)
	if err := c.Bind(biz); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Omit("token", "mtime", "ctime", "dtime").Save(biz).Error; err != nil {
		log.Error("saveBusiness(%+v) error(%v)", biz, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func delBusiness(c *bm.Context) {
	id, _ := strconv.ParseInt(c.Request.Form.Get("id"), 10, 64)
	if id == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := pushSrv.DB.Model(&model.Business{ID: id}).Update("dtime", time.Now().Unix()).Error; err != nil {
		log.Error("delBusiness(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
