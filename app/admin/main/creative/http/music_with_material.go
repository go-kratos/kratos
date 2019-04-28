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

func musicMaterialRelationInfo(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &music.WithMaterialParam{}
	if err = svc.DBArchive.Where("id=?", id).First(&m).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]*music.WithMaterialParam{
		"data": m,
	}, nil)
}

func editMaterialRelation(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.WithMaterialParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := music.WithMaterial{}
	if err = svc.DBArchive.Where("sid=?", m.Sid).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	m.ID = exist.ID
	m.UID = uid
	m.Index = exist.Index
	if err = svc.DBArchive.Model(&music.WithMaterial{}).Where("id=?", exist.ID).Update(m).Error; err != nil {
		log.Error("svc.editMaterialRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterialRelation, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "update", Name: string(m.ID)})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func addMaterialRelation(c *bm.Context) {
	var (
		err error
	)
	m := &music.WithMaterialParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := &music.WithMaterial{}
	//一期  sid tid 是 一一对应的关系
	//check sid bind
	if err = svc.DBArchive.Where("state!=?", music.MusicDelete).Where("sid=?", m.Sid).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
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
		log.Error("svc.addMaterialRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterialRelation, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: string(m.ID)})
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func batchAddMaterialRelation(c *bm.Context) {
	var (
		err  error
		sids []int64
	)
	uid, uname := getUIDName(c)
	m := &music.BatchMusicWithMaterialParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		log.Error("svc.batchAddMaterialRelation bind error(%v)", err)
		return
	}
	//todo tid 必须是二级分类
	if sids, err = xstr.SplitInts(m.SidList); err != nil {
		log.Error("svc.batchAddMaterialRelation SplitInts error(%v)", err)
		c.JSON(nil, err)
		return
	}
	max := music.WithMaterial{}
	if err = svc.DBArchive.Model(&music.WithMaterial{}).Where("tid=?", m.Tid).Where("state!=?", music.MusicDelete).Order("music_with_material.index desc").First(&max).Error; err != nil && err != gorm.ErrRecordNotFound {
		//sql err
		log.Error("svc.batchAddMaterialRelation max index error(%v)", err)
		return
	}
	thisIndex := max.Index
	i := int64(1)
	for _, sid := range sids {
		//check exists
		exists := music.WithMaterial{}
		if err = svc.DBArchive.Model(&music.WithMaterial{}).Where("sid=?", sid).Where("state!=?", music.MusicDelete).First(&exists).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("svc.batchAddMaterialRelation check exist sid tid (%d,%d)  error(%v)", sid, m.Tid, err)
			c.JSON(nil, err)
			return
		}
		if exists.ID > 0 {
			//覆盖
			if err = svc.DBArchive.Model(&music.WithMaterial{}).Where("sid=?", sid).Where("state!=?", music.MusicDelete).Update("tid", m.Tid).Error; err != nil && err != gorm.ErrRecordNotFound {
				log.Error("svc.batchAddMaterialRelation Update error(%v)", err)
				continue
			}
			svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterialRelation, &music.LogParam{ID: exists.Sid, UID: uid, UName: uname, Action: "update", Name: string(exists.Sid)})
		} else {
			mp := &music.WithMaterialParam{Sid: sid, Tid: m.Tid, Index: thisIndex + i}
			if err = svc.DBArchive.Create(mp).Error; err != nil {
				log.Error("svc.batchAddMaterialRelation Create error(%v)", err)
				continue
			}
			i++
			svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMaterialRelation, &music.LogParam{ID: mp.Sid, UID: uid, UName: uname, Action: "add", Name: string(mp.Sid)})
		}
	}
	c.JSON(map[string]int64{
		"code": 0,
	}, nil)
}
