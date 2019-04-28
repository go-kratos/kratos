package http

import (
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

func searchMaterial(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*music.MaterialMixParent
		count int64
		pid   = atoi(req.Get("pid"))
		name  = req.Get("name")
	)

	db := svc.DBArchive.Where("music_material.state!=?", music.MusicDelete)
	if pid != 0 {
		db = db.Where("music_material.pid=?", pid)
	}
	if name != "" {
		db = db.Where("music_material.name=?", name)
	}
	db.Model(&music.Material{}).Count(&count)
	if err = db.Model(&music.Material{}).
		Joins("left join music_material p on p.id=music_material.pid").
		Select("music_material.*,p.name as p_name").
		Order("music_material.index").Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &music.MaterialMixParentPager{
		Pager: &music.Pager{Num: 1, Size: int(count), Total: count},
		Items: items,
	}
	c.JSON(pager, nil)
}

func materialInfo(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &music.MaterialParam{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&m).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]*music.MaterialParam{
		"data": m,
	}, nil)
}

func editMaterial(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &music.MaterialParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := music.Material{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	uid, uname := getUIDName(c)
	m.UID = uid
	if err = svc.DBArchive.Model(&music.Material{}).Where("id=?", id).Update(m).Error; err != nil {
		log.Error("svc.editMaterial error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterial, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "update", Name: m.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func addMaterial(c *bm.Context) {
	var (
		err error
	)
	m := &music.MaterialParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	uid, uname := getUIDName(c)
	m.UID = uid
	exist := &music.Material{}
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
		log.Error("svc.addMaterial error(%v)", err)
		c.JSON(nil, err)
		return
	}

	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterial, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: m.Name})
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func delMaterial(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	exist := music.Material{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	uid, uname := getUIDName(c)
	if err = svc.DBArchive.Model(music.Material{}).Where("id=?", id).Update(map[string]int{"state": music.MusicDelete}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.delMaterial error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterial, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "del", Name: exist.Name})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func batchDeleteMaterial(c *bm.Context) {
	var (
		err    error
		req    = c.Request.PostForm
		ids    = req.Get("ids")
		arrids []int64
	)
	uid, uname := getUIDName(c)
	if arrids, err = xstr.SplitInts(ids); err != nil {
		log.Error("svc.batchDeleteCategoryRelation SplitInts error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, id := range arrids {
		if err = svc.DBArchive.Model(music.Material{}).Where("id=?", id).Update(map[string]int{"state": music.MusicDelete}).Error; err != nil {
			log.Error("svc.batchDeleteMaterial error(%v)", err)
			err = nil
			continue
		}
		//for log
		exist := music.Material{}
		if err = svc.DBArchive.Where("id=?", id).First(&exist).Error; err != nil {
			err = nil
			continue
		}
		svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterial, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "del", Name: exist.Name})
	}
	c.JSON(map[string]int64{
		"code": 0,
	}, nil)
}
