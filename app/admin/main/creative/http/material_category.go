package http

import (
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/material"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"

	"github.com/jinzhu/gorm"
)

func searchMCategory(c *bm.Context) {
	var (
		req     = c.Request.Form
		err     error
		items   []*material.Category
		count   int64
		name    = req.Get("name")
		page    = atoi(req.Get("page"))
		typeStr = req.Get("type")
		size    = 20
	)
	if page == 0 {
		page = 1
	}
	db := svc.DB.Where("state!=?", material.StateOff)

	if typeStr != "" {
		db = db.Where("type=?", atoi(typeStr))
	}
	if name != "" {
		db = db.Where("name=?", name)
	}
	db.Model(&material.Category{}).Count(&count)
	if err = db.Model(&material.Category{}).Order("rank ASC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &material.CategoryPager{
		Items: items,
		Pager: &material.Pager{Num: page, Size: size, Total: count},
	}
	c.JSON(pager, nil)
}

func category(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &material.CategoryParam{}
	if err = svc.DB.Where("id=?", id).First(&m).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]*material.CategoryParam{
		"data": m,
	}, nil)
}

func editMCategory(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	uid, uname := getUIDName(c)
	m := &material.CategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := material.Category{}
	if err = svc.DB.Where("id=?", id).Where("state!=?", material.StateOff).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	m.UID = uid
	if err := svc.DB.Model(&material.Category{}).Where("id=?", id).Update(m).Update(map[string]interface{}{"type": m.Type}).Update(map[string]interface{}{"new": m.New}).Error; err != nil {
		log.Error("svc.editMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMaterialTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "edit", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func indexMCategory(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.IndexParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := material.Category{}
	if err = svc.DB.Where("id=?", m.ID).Where("state!=?", material.StateOff).First(&exist).Error; err != nil {
		log.Error("svc.indexMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	exist1 := material.Category{}
	if err = svc.DB.Where("id=?", m.SwitchID).Where("state!=?", material.StateOff).First(&exist1).Error; err != nil {
		log.Error("svc.indexMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	m.UID = uid
	if err := svc.DB.Model(&material.Category{}).Where("id=?", m.ID).Update(map[string]int64{"rank": m.SwitchIndex}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err := svc.DB.Model(&material.Category{}).Where("id=?", m.SwitchID).Update(map[string]int64{"rank": m.Index}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}

	svc.SendMusicLog(c, logcli.LogClientArchiveMaterialTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "index", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func addMCategory(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &material.CategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		log.Error("svc.addCategory bind error(%v)", err)
		return
	}
	m.UID = uid
	//Category.type 跟素材type一致
	exist := &material.Category{}
	if err = svc.DB.Where("state!=?", material.StateOff).Where("name=?", m.Name).Where("type=?", m.Type).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(nil, err)
		return
	}
	if exist.ID > 0 {
		c.JSON(map[string]int64{
			"id": exist.ID,
		}, nil)
		return
	}
	if err = svc.DB.Create(m).Error; err != nil {
		log.Error("svc.addMCategory  Create error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMaterialTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: m.Name})
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func delMCategory(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	exist := material.Category{}
	if err = svc.DB.Where("id=?", id).Where("state!=?", material.StateOff).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	if err := svc.DB.Model(material.Category{}).Where("id=?", id).Update(map[string]int{"state": material.StateOff}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.delMCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMaterialTypeCategory, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "del", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}
