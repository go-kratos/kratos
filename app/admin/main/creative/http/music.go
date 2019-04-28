package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/creative/model/logcli"
	"go-common/app/admin/main/creative/model/music"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func searchMusic(c *bm.Context) {
	var (
		req      = c.Request.Form
		err      error
		items    []*music.Music
		count    int64
		mid      = atoi(req.Get("mid"))
		sid      = atoi(req.Get("sid"))
		tid      = req.Get("tid")
		musicans = req.Get("musicans")
		name     = req.Get("name")
		page     = atoi(req.Get("page"))
		size     = 20
	)
	if page == 0 {
		page = 1
	}
	db := svc.DBArchive.Model(&music.Music{}).Where("music.state!=?", music.MusicDelete).
		Joins("left join music_with_material on music.sid=music_with_material.sid and music_with_material.state !=? ", music.MusicDelete).
		Joins("left join music_material on music_with_material.tid=music_material.id").
		Select("music.*,music_material.name as material_name,music_material.pid as pid,music_material.id as tid,music_with_material.id as rid")
	if mid != 0 {
		db = db.Where("music.mid=?", mid)
	}
	if sid != 0 {
		db = db.Where("music.sid=?", sid)
	}
	//未添加分类的BGM filter tid=0 filter
	if tid != "" {
		if atoi(tid) == 0 {
			db = db.Where("music_with_material.tid is null")
		} else {
			db = db.Where("music_with_material.tid=?", atoi(tid))
		}
	}
	if name != "" {
		db = db.Where("music.name=?", name)
	}
	if musicans != "" {
		db = db.Where("music.musicians=?", musicans)
	}
	db.Count(&count)
	//投稿时间ctime 倒序排列
	if err = db.Order("music.ctime DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		log.Error("%v\n", err)
		c.JSON(nil, err)
		return
	}
	pager := &music.ResultPager{
		Items: items,
		Pager: &music.Pager{Num: page, Size: size, Total: count},
	}
	c.JSON(pager, nil)
}

func editMusicTimeline(c *bm.Context) {
	var (
		req      = c.Request.PostForm
		sid      = parseInt(req.Get("sid"))
		timeline = req.Get("timeline")
		err      error
	)
	if len(timeline) > 0 {
		type timelineItem struct {
			//转换成毫秒存储
			Point     int64  `json:"point"`
			Comment   string `json:"comment"`
			Recommend int8   `json:"recommend"`
		}
		var timelineExp []*timelineItem
		if err = json.Unmarshal([]byte(timeline), &timelineExp); err != nil {
			httpCode(c, "timeline json is wrong", ecode.RequestErr)
			return
		}
	}
	exist := music.Music{}
	if err = svc.DBArchive.Where("sid=?", sid).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		httpCode(c, "sid不存在", ecode.RequestErr)
		return
	}
	if err := svc.DBArchive.Model(&music.Music{}).Where("sid=?", sid).Update(map[string]string{"timeline": timeline}).Error; err != nil {
		log.Error("vdaSvc.editMusicTimeline error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

//开放来修改一些不复杂的数据
func editMusic(c *bm.Context) {
	var (
		req       = c.Request.PostForm
		sid       = parseInt(req.Get("sid"))
		cooperate = req.Get("cooperate")
		err       error
	)
	exist := music.Music{}
	if err = svc.DBArchive.Where("sid=?", sid).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		httpCode(c, "sid不存在", ecode.RequestErr)
		return
	}
	if err := svc.DBArchive.Model(&music.Music{}).Where("sid=?", sid).Update(map[string]int{"cooperate": atoi(cooperate)}).Error; err != nil {
		log.Error("vdaSvc.editMusic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func editMusicTags(c *bm.Context) {
	var (
		req  = c.Request.PostForm
		sid  = parseInt(req.Get("sid"))
		tags = req.Get("tags")
		err  error
	)
	exist := music.Music{}
	if err = svc.DBArchive.Where("sid=?", sid).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		httpCode(c, "sid不存在", ecode.RequestErr)
		return
	}
	if err := svc.DBArchive.Model(&music.Music{}).Where("sid=?", sid).Update(map[string]string{"tags": tags}).Error; err != nil {
		log.Error("vdaSvc.editMusicTags error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func batchEditMusicTags(c *bm.Context) {
	var (
		req    = c.Request.PostForm
		sids   = req.Get("sids")
		tags   = req.Get("tags")
		arrids []int64
		err    error
	)
	if arrids, err = xstr.SplitInts(sids); err != nil {
		log.Error("svc.batchEditMusicTags SplitInts error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if len(arrids) < 2 {
		c.JSON(nil, fmt.Errorf("sids参数错误"))
	}
	if err = svc.DBArchive.Model(&music.Music{}).Where("sid IN (?)", arrids).Update(map[string]string{"tags": tags}).Error; err != nil {
		log.Error("vdaSvc.batchEditMusicTags error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

//用户侧BGM名称 front_name
func editMusicFrontName(c *bm.Context) {
	var (
		req       = c.Request.PostForm
		sid       = parseInt(req.Get("sid"))
		frontname = req.Get("frontname")
		err       error
	)
	exist := music.Music{}
	if err = svc.DBArchive.Where("sid=?", sid).Where("state!=?", music.MusicDelete).First(&exist).Error; err != nil {
		httpCode(c, "sid不存在", ecode.RequestErr)
		return
	}
	if err := svc.DBArchive.Model(&music.Music{}).Where("sid=?", sid).Update(map[string]string{"frontname": frontname}).Error; err != nil {
		log.Error("vdaSvc.editMusicName error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int{
		"code": 0,
	}, nil)
}

func syncMusic(c *bm.Context) {
	var (
		err error
	)
	mp := &music.Param{}
	if err = c.BindWith(mp, binding.Form); err != nil {
		log.Error("svc.syncMusic  bind error(%+v)", err)
		log.Error("svc.syncMusic  bind error trace(%+v)", errors.Wrap(err, "sync bind  error"))
		return
	}
	m := &music.Music{Sid: mp.Sid, Name: mp.Name, Cover: mp.Cover, Stat: mp.Stat, Mid: mp.Mid, Musicians: mp.Musicians, Categorys: mp.Categorys, Playurl: mp.Playurl, PubTime: mp.PubTime, Duration: mp.Duration, Filesize: mp.Filesize, State: mp.State}
	if m.State != music.MusicDelete {
		m.State = music.MusicOpen
	}
	exist := music.Music{}
	if err = svc.DBArchive.Where("sid=?", mp.Sid).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("svc.syncMusic  find error(%+v)", err)
		c.JSON(nil, errors.Wrap(err, "sync find sid error"))
		return
	}
	uid, uname := getUIDName(c)
	if exist.ID > 0 {
		m.ID = exist.ID
		if err = svc.DBArchive.Model(&music.Music{}).Where("sid=?", mp.Sid).Update(m).Update(map[string]int8{"state": m.State}).Error; err != nil {
			log.Error("svc.syncMusic update error(%+v)", err)
			c.JSON(nil, errors.Wrap(err, "sync update error"))
			return
		}
		if m.State == music.MusicDelete {
			svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "del", Name: m.Name})
		} else {
			svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "update", Name: m.Name})
		}

	} else {
		if err = svc.DBArchive.Create(m).Error; err != nil {
			log.Error("svc.syncMusic  Create error(%+v)", err)
			c.JSON(nil, errors.Wrap(err, "sync add error"))
			return
		}
		svc.SendMusicLog(c, logcli.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: m.Name})
	}
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}
