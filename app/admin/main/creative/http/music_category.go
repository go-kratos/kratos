package http

import (
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"

	"github.com/jinzhu/gorm"
)

func searchCategory(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*music.Category
		count int64
		pid   = atoi(req.Get("pid"))
		name  = req.Get("name")
		page  = atoi(req.Get("page"))
		sort  = atoi(req.Get("sort"))
		size  = 20
		order string
	)
	if page == 0 {
		page = 1
	}
	db := svc.DBArchive.Where("state!=?", music.MusicDelete)
	if pid != 0 {
		db = db.Where("pid=?", pid)
	}
	//pid 目前不做分级 pid=0
	if name != "" {
		db = db.Where("name=?", name)
	}
	db.Model(&music.Category{}).Count(&count)
	if sort == 1 {
		order = "camera_index"
	} else {
		order = "index"
	}
	if err = db.Model(&music.Category{}).Order(order).Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &music.CategoryPager{
		Items: items,
		Pager: &music.Pager{Num: page, Size: size, Total: count},
	}
	c.JSON(pager, nil)
}

func categoryInfo(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &music.CategoryParam{}
	if err = svc.DBArchive.Where("id=?", id).First(&m).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]*music.CategoryParam{
		"data": m,
	}, nil)
}

func editCategory(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.CategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := music.Category{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	m.UID = uid
	if err := svc.DBArchive.Model(&music.Category{}).Where("id=?", id).Update(m).Error; err != nil {
		log.Error("svc.editCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "index", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func indexCategory(c *bm.Context) {
	var (
		req    = c.Request.PostForm
		err    error
		cate   = atoi(req.Get("type"))
		column string
	)
	uid, uname := getUIDName(c)
	m := &music.IndexParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := music.Category{}
	if err = svc.DBArchive.Where("id=?", m.ID).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		log.Error("svc.indexCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	exist1 := music.Category{}
	if err = svc.DBArchive.Where("id=?", m.SwitchID).Where("state!=?", music.MusicDelete).First(&exist1).Error; err != nil {
		log.Error("svc.indexCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if cate == 1 {
		column = "camera_index"
	} else {
		column = "index"
	}
	m.UID = uid
	if err := svc.DBArchive.Model(&music.Category{}).Where("id=?", m.ID).Update(map[string]int64{column: m.SwitchIndex}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err := svc.DBArchive.Model(&music.Category{}).Where("id=?", m.SwitchID).Update(map[string]int64{column: m.Index}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}

	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "index", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func addCategory(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.CategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		log.Error("svc.addCategory bind error(%v)", err)
		return
	}
	m.UID = uid
	exist := &music.Category{}
	if err = svc.DBArchive.Where("state!=?", music.MusicDelete).Where("name=?", m.Name).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(nil, err)
		return
	}
	if exist.ID > 0 {
		c.JSON(map[string]int64{
			"id": exist.ID,
		}, nil)
		return
	}
	if err = svc.DBArchive.Create(m).Error; err != nil {
		log.Error("svc.addCategory  Create error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategory, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: m.Name})
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func delCategory(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	exist := music.Category{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	if err := svc.DBArchive.Model(music.Category{}).Where("id=?", id).Update(map[string]int{"state": music.MusicDelete}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.delCategory error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategory, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "del", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}
