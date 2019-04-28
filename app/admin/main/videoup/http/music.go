package http

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/music"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

const ()

func getUIDName(c *bm.Context) (uid int64, uname string) {
	unamei, ok := c.Get("username")
	if ok {
		uname = unamei.(string)
	}
	uidi, ok := c.Get("uid")
	if ok {
		uid = uidi.(int64)
	}
	return
}
func syncMusic(c *bm.Context) {
	var (
		err error
	)
	mp := &music.Param{}
	if err = c.BindWith(mp, binding.Form); err != nil {
		log.Error("vdaSvc.syncMusic  bind error(%+v)", err)
		log.Error("vdaSvc.syncMusic  bind error trace(%+v)", errors.Wrap(err, "sync bind  error"))
		return
	}
	m := &music.Music{Sid: mp.Sid, Name: mp.Name, Cover: mp.Cover, Stat: mp.Stat, Mid: mp.Mid, Musicians: mp.Musicians, Categorys: mp.Categorys, Playurl: mp.Playurl, PubTime: mp.PubTime, Duration: mp.Duration, Filesize: mp.Filesize, State: mp.State}
	if m.State != music.MusicDelete {
		m.State = music.MusicOpen
	}
	exist := music.Music{}
	if err = vdaSvc.DB.Where("sid=?", mp.Sid).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("vdaSvc.syncMusic  find error(%+v)", err)
		c.JSON(nil, errors.Wrap(err, "sync find sid error"))
		return
	}
	uid, uname := getUIDName(c)
	if exist.ID > 0 {
		m.ID = exist.ID
		if err = vdaSvc.DB.Model(&music.Music{}).Where("sid=?", mp.Sid).Update(m).Update(map[string]int8{"state": m.State}).Error; err != nil {
			log.Error("vdaSvc.syncMusic update error(%+v)", err)
			c.JSON(nil, errors.Wrap(err, "sync update error"))
			return
		}
		if m.State == music.MusicDelete {
			vdaSvc.SendMusicLog(c, archive.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "del", Name: m.Name})
		} else {
			vdaSvc.SendMusicLog(c, archive.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "update", Name: m.Name})
		}

	} else {
		if err = vdaSvc.DB.Create(m).Error; err != nil {
			log.Error("vdaSvc.syncMusic  Create error(%+v)", err)
			c.JSON(nil, errors.Wrap(err, "sync add error"))
			return
		}
		vdaSvc.SendMusicLog(c, archive.LogClientArchiveMusicTypeMusic, &music.LogParam{ID: m.ID, UID: uid, UName: uname, Action: "add", Name: m.Name})
	}
	c.JSON(map[string]int64{
		"id": m.ID,
	}, nil)
}
