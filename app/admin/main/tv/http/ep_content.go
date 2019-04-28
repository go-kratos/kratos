package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func createEP(c *bm.Context) {
	var (
		err error
		epc = &model.TVEpContent{}
	)
	if err = c.Bind(epc); err != nil {
		return
	}
	exist := model.Content{}
	tvSrv.DB.Where("epid=?", epc.ID).First(&exist)
	if exist.ID <= 0 { // data not exist, brand new data
		if err = tvSrv.DB.Create(epc.ToContent(true)).Error; err != nil {
			log.Error("tvSrv.createEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if err = tvSrv.DB.Create(epc).Error; err != nil {
			log.Error("tvSrv.createEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
	} else {
		// data exists, but was deleted
		if exist.IsDeleted == 1 {
			if err = tvSrv.EpDel(epc.ID, 0); err != nil {
				c.JSON(nil, err)
				return
			}
		}
		// data exist, so we update
		if err := tvSrv.DB.Model(&model.TVEpContent{}).Where("id=?", epc.ID).Update(epc).Error; err != nil {
			log.Error("tvSrv.modifyEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if err := tvSrv.DB.Model(&model.TVEpContent{}).Where("id=?", epc.ID).Update(map[string]string{"long_title": epc.LongTitle, "title": epc.Title, "cover": epc.Cover}).Error; err != nil {
			log.Error("tvSrv.modifyEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
		lic := epc.ToContent(false)
		if err := tvSrv.DB.Model(&model.Content{}).Where("epid=?", epc.ID).Update(lic).Error; err != nil {
			log.Error("tvSrv.modifyEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if err := tvSrv.DB.Model(&model.Content{}).Where("epid=?", epc.ID).Update(map[string]string{"title": lic.Title, "subtitle": lic.Subtitle, "desc": lic.Desc, "cover": lic.Cover}).Error; err != nil {
			log.Error("tvSrv.modifyEP error(%v)", err)
			c.JSON(nil, err)
			return
		}
		if err := tvSrv.EpCheck(&exist); err != nil { // manage ep's check status
			c.JSON(nil, err)
			return
		}
	}
	renderErrMsg(c, 0, "0")
	log.Info("createEP success(%v)", epc)
}

func removeEP(c *bm.Context) {
	var (
		req = c.Request.PostForm
		err error
		id  = parseInt(req.Get("id"))
	)
	exist := model.Content{}
	if err = tvSrv.DB.Where("epid=?", id).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.EpDel(id, 1); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.EpCheck(&exist); err != nil { // manage ep's check status
		c.JSON(nil, err)
		return
	}
	log.Info("removeEP success(%v)", exist)
	renderErrMsg(c, 0, "0")
}

func createSeason(c *bm.Context) {
	var (
		err error
		eps = &model.TVEpSeason{}
	)
	if err = c.Bind(eps); err != nil {
		return
	}
	if err = tvSrv.SnUpdate(c, eps); err != nil {
		c.JSON(nil, err)
		return
	}
	renderErrMsg(c, 0, "0")
}

func removeSeason(c *bm.Context) {
	var (
		req = c.Request.PostForm
		err error
		id  = parseInt(req.Get("id"))
	)
	exist := model.TVEpSeason{}
	if err = tvSrv.DB.Where("id=?", id).Where("is_deleted=?", 0).First(&exist).Error; err != nil {
		renderErrMsg(c, ecode.RequestErr.Code(), "Data not exist")
		return
	}
	if err = tvSrv.SeasonRemove(&exist); err != nil {
		c.JSON(nil, err)
		return
	}
	log.Info("removeSeason success(%v)", exist)
	renderErrMsg(c, 0, "0")
}

func act(c *bm.Context, action int) {
	param := new(struct {
		CID int64 `form:"cid" validate:"required"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.EpAct(c, param.CID, action))
}

func online(c *bm.Context) {
	act(c, 1)
}

func hidden(c *bm.Context) {
	act(c, 0)
}
