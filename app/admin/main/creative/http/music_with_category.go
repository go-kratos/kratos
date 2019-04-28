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

func searchCategoryRelation(c *bm.Context) {
	var (
		req   = c.Request.Form
		err   error
		items []*music.CategoryList
		count int64
		mid   = atoi(req.Get("mid"))
		tid   = atoi(req.Get("tid"))
		sid   = atoi(req.Get("sid"))
		name  = req.Get("name")
		page  = atoi(req.Get("page"))
		size  = 20
	)
	if page == 0 {
		page = 1
	}
	//init sql
	db := svc.DBArchive.Model(&music.CategoryList{}).
		Joins("left join music on music.sid=music_with_category.sid").
		Joins("left join music_category on music_category.id=music_with_category.tid").
		Where("music_with_category.state !=?", music.MusicDelete).
		Where("music.state !=?", music.MusicDelete).
		Select("music_with_category.*,music_category.name as category_name,music.name,music.cover,music.playurl,music.duration,music.mid,music.musicians,music.state as music_state,music.tags,music.timeline,music.frontname,music.cooperate")
	if mid != 0 {
		db = db.Where("mid=?", mid)
	}
	if sid != 0 {
		db = db.Where("music_with_category.sid=?", sid)
	}
	if name != "" {
		db = db.Where("music_category.name=?", name)
	}
	if tid != 0 {
		db = db.Where("music_with_category.tid=?", tid)
	}
	//count total
	db.Model(&music.CategoryList{}).Count(&count)
	if err = db.Order("music_with_category.index").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &music.CategoryListPager{
		Pager: &music.Pager{Num: page, Size: size, Total: count},
		Items: items,
	}
	c.JSON(pager, nil)
}

func musicCategoryRelationInfo(c *bm.Context) {
	var (
		req = c.Request.Form
		id  = parseInt(req.Get("id"))
		err error
	)
	m := &music.WithCategoryParam{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&m).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(map[string]*music.WithCategoryParam{
		"data": m,
	}, nil)
}

func editCategoryRelation(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.WithCategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		httpCode(c, "参数错误", ecode.RequestErr)
		return
	}
	exist := music.WithCategory{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DBArchive.Model(&music.WithCategory{}).Where("id=?", id).Update(m).Error; err != nil {
		log.Error("svc.editCategoryRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "update", Name: string(m.ID)})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func indexCategoryRelation(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.IndexParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		return
	}
	exist := music.WithCategory{}
	if err = svc.DBArchive.Where("id=?", m.ID).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	exist1 := music.WithCategory{}
	if err = svc.DBArchive.Where("id=?", m.SwitchID).Where("state!=?", music.MusicDelete).First(&exist1).Error; err != nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DBArchive.Model(&music.WithCategory{}).Where("id=?", m.ID).Update(map[string]int64{"index": m.SwitchIndex}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexCategoryRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = svc.DBArchive.Model(&music.WithCategory{}).Where("id=?", m.SwitchID).Update(map[string]int64{"index": m.Index}).Update(map[string]int64{"uid": uid}).Error; err != nil {
		log.Error("svc.indexCategoryRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "index", Name: string(m.ID)})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func addCategoryRelation(c *bm.Context) {
	var (
		err error
	)
	uid, uname := getUIDName(c)
	m := &music.WithCategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		log.Error("svc.addCategoryRelation bind error(%v)", err)
		return
	}
	exist := &music.WithCategory{}
	if err = svc.DBArchive.Where("state!=?", music.MusicDelete).Where("sid=?", m.Sid).Where("tid=?", m.Tid).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
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
		log.Error("svc.addCategoryRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}

	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: string(m.ID)})
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}

func batchAddCategoryRelation(c *bm.Context) {
	var (
		err           error
		sids, sendIds []int64
		sidsNotify    = make(map[int64]*music.SidNotify)
	)
	uid, uname := getUIDName(c)
	m := &music.BatchMusicWithCategoryParam{}
	if err = c.BindWith(m, binding.Form); err != nil {
		log.Error("svc.batchAddCategoryRelation bind error(%v)", err)
		return
	}
	//todo tid 必须是二级分类
	if sids, err = xstr.SplitInts(m.SidList); err != nil {
		log.Error("svc.batchAddCategoryRelation SplitInts error(%v)", err)
		c.JSON(nil, err)
		return
	}
	max := music.WithCategory{}
	if err = svc.DBArchive.Model(&music.WithCategory{}).Where("tid=?", m.Tid).Where("state!=?", music.MusicDelete).Order("music_with_category.index desc").First(&max).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("svc.batchAddMaterialRelation max error(%v)", err)
		c.JSON(nil, err)
		return
	}
	thisIndex := max.Index
	i := int64(1)
	for _, sid := range sids {
		sidnotify := music.SidNotify{Sid: sid}
		//find mid
		musicExist := music.Music{}
		if err = svc.DBArchive.Model(&music.Music{}).Where("sid=?", sid).First(&musicExist).Error; err != nil {
			log.Error("svc.batchAddMaterialRelation check exist sid (%d)  error(%v)", sid, err)
			continue
		}
		//check mid first bind
		midExists := music.WithCategory{}
		if err = svc.DBArchive.Model(&music.WithCategory{}).Joins("left join music on music.sid=music_with_category.sid").Where("music.mid=?", musicExist.Mid).First(&midExists).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("svc.batchAddMaterialRelation check mid exist sid mid (%d,%d)  error(%v)", sid, musicExist.Mid, err)
			c.JSON(nil, err)
			return
		}
		if err == gorm.ErrRecordNotFound {
			sidnotify.MidFirst = true
		}
		//check sid bind
		sidExists := music.WithCategory{}
		if err = svc.DBArchive.Model(&music.WithCategory{}).Where("sid=?", sid).First(&sidExists).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("svc.batchAddMaterialRelation check sid exist sid mid (%d,%d)  error(%v)", sid, musicExist.Mid, err)
			c.JSON(nil, err)
			return
		}
		if err == gorm.ErrRecordNotFound {
			sidnotify.SidFirst = true
		}
		//check sid-tid bind exists
		exists := music.WithCategory{}
		if err = svc.DBArchive.Model(&music.WithCategory{}).Where("sid=?", sid).Where("tid=?", m.Tid).Where("state!=?", music.MusicDelete).First(&exists).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("svc.batchAddMaterialRelation check exist sid tid (%d,%d)  error(%v)", sid, m.Tid, err)
			c.JSON(nil, err)
			return
		}
		if exists.ID > 0 {
			//pass
			continue
		}
		//should create bind
		mp := &music.WithCategoryParam{Sid: sid, Tid: m.Tid, Index: thisIndex + i}
		if err = svc.DBArchive.Create(mp).Error; err != nil {
			log.Error("svc.batchAddMaterialRelation Create error(%v)", err)
			continue
		}
		sidsNotify[sid] = &sidnotify
		i++
		svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: mp.Sid, UID: uid, UName: uname, Action: "add", Name: string(mp.Sid)})
	}
	log.Info("svc.SendNotify param SendList(%+v) sidsNotify(%+v)", m.SendList, sidsNotify)

	if m.SendList != "" {
		if sendIds, err = xstr.SplitInts(m.SendList); err != nil {
			log.Error("svc.batchAddCategoryRelation SplitInts  SendList error(%v)", err)
			c.JSON(nil, err)
			return
		}
		log.Info("svc.SendNotify param sendIds(%+v) sidsNotify(%+v)", sendIds, sidsNotify)
		svc.SendNotify(c, sendIds, sidsNotify)
	}

	c.JSON(map[string]int64{
		"code": 0,
	}, nil)
}

func batchDeleteCategoryRelation(c *bm.Context) {
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
		if err = svc.DBArchive.Model(music.WithCategory{}).Where("id=?", id).Update(map[string]int{"state": music.MusicDelete}).Error; err != nil {
			log.Error("svc.batchDeleteCategoryRelation error(%v)", err)
			err = nil
			continue
		}
		svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "delete", Name: string(id)})
	}
	c.JSON(map[string]int64{
		"code": 0,
	}, nil)
}

func delCategoryRelation(c *bm.Context) {
	var (
		req = c.Request.PostForm
		id  = parseInt(req.Get("id"))
		err error
	)
	uid, uname := getUIDName(c)
	exist := music.WithCategory{}
	if err = svc.DBArchive.Where("id=?", id).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DBArchive.Model(music.WithCategory{}).Where("id=?", id).Update(map[string]int{"state": music.MusicDelete}).Error; err != nil {
		log.Error("svc.delCategoryRelation error(%v)", err)
		c.JSON(nil, err)
		return
	}
	svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeCategoryRelation, &music.LogParam{ID: id, UID: uid, UName: uname, Action: "del", Name: string(exist.ID)})
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}
