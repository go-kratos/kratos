package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"go-common/library/time"
)

const (
	_platform  = 12
	_isDeleted = 1
)

func versionList(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*model.Version
		count int64
		ver   = req.Get("ver")
		page  = atoi(req.Get("page"))
		size  = 20
	)
	if page == 0 {
		page = 1
	}
	db := tvSrv.DBShow.Where("plat=?", _platform).Where("state!=?", _isDeleted)
	if ver != "" {
		db = db.Where("version=?", ver)
	}
	db.Model(&model.Version{}).Count(&count)
	if err = db.Model(&model.Version{}).Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &model.VersionPager{
		TotalCount: count,
		Pn:         page,
		Ps:         size,
		Items:      items,
	}
	c.JSON(pager, nil)
}

func versionInfo(c *bm.Context) {
	var (
		req = c.Request.Form
		vid = parseInt(req.Get("id"))
		err error
	)
	exist := model.Version{}
	if err = tvSrv.DBShow.Where("id=?", vid).Where("state!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(exist, nil)
}

func saveVersion(c *bm.Context) {
	var (
		req = c.Request.PostForm

		vid = parseInt(req.Get("id"))
		err error
	)
	exist := model.Version{}
	if err = tvSrv.DBShow.Where("id=?", vid).Where("state!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	alert, simple := validateVerPostData(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	if err := tvSrv.DBShow.Model(&model.Version{}).Where("id=?", vid).Update(simple).Error; err != nil {
		log.Error("tvSrv.saveVersion error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err := tvSrv.DBShow.Model(&model.Version{}).Where("id=?", vid).Update(map[string]int8{"state": simple.State}).Error; err != nil {
		log.Error("tvSrv.saveVersion error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addVersion(c *bm.Context) {
	var (
		err error
	)
	alert, simple := validateVerPostData(c)
	if alert != "" {
		renderErrMsg(c, ecode.RequestErr.Code(), alert)
		return
	}
	if err = tvSrv.DBShow.Create(simple).Error; err != nil {
		log.Error("tvSrv.addVersion error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func versionDel(c *bm.Context) {
	var (
		req = c.Request.PostForm

		vid = parseInt(req.Get("id"))
		err error
	)
	exist := model.Version{}
	if err = tvSrv.DBShow.Where("id=?", vid).Where("state!=?", _isDeleted).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err := tvSrv.DBShow.Model(&model.Version{}).Where("id=?", vid).Update(map[string]int{"state": _isDeleted}).Error; err != nil {
		log.Error("tvSrv.versionDel error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func validateVerPostData(c *bm.Context) (alert string, simple *model.Version) {
	var (
		req     = c.Request.PostForm
		plat    = atoi(req.Get("plat"))
		version = req.Get("version")
		build   = atoi(req.Get("build"))
		ptime   = time.Time(parseInt(req.Get("ptime")))
		state   = atoi(req.Get("state"))
		desc    = req.Get("description")
	)
	if plat == 0 {
		alert = "平台不能为空"
		return
	}
	if version == "" {
		alert = "版本号不能为空"
		return
	}
	if build == 0 {
		alert = "build号不能为空"
		return
	}
	if int64(ptime) == 0 {
		alert = "发布时间不能为空"
		return
	}
	if desc == "" {
		alert = "描述不能为空"
		return
	}
	return "", &model.Version{Plat: int8(plat), Description: desc, Version: version, Build: build, State: int8(state), Ptime: ptime}
}
