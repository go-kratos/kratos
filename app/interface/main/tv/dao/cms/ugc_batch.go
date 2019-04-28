package cms

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

const (
	_arcCMSDeleted = 1
	_mcArcCMSKey   = "arc_cms_%d"
	_mcVideoCMSKey = "video_cms_%d"
)

// ArcCMSCacheKey .
func (d *Dao) ArcCMSCacheKey(aid int64) string {
	return fmt.Sprintf(_mcArcCMSKey, aid)
}

// VideoCMSCacheKey .
func (d *Dao) VideoCMSCacheKey(cid int64) string {
	return fmt.Sprintf(_mcVideoCMSKey, cid)
}

// ArcsMetaCache pick archive cms meta cache
func (d *Dao) ArcsMetaCache(c context.Context, ids []int64) (cached map[int64]*model.ArcCMS, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.ArcCMS, len(ids))
	idmap, allKeys := keysTreat(ids, d.ArcCMSCacheKey)
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取Archive信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		err = nil
		return
	}
	for key, item := range replys {
		art := &model.ArcCMS{}
		if err = conn.Scan(item, art); err != nil {
			PromError("mc:获取Archive信息缓存json解析")
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		cached[idmap[key]] = art
		delete(idmap, key)
	}
	missed = missedTreat(idmap, len(cached))
	return
}

// VideosMetaCache pick video cms meta cache
func (d *Dao) VideosMetaCache(c context.Context, ids []int64) (cached map[int64]*model.VideoCMS, missed []int64, err error) {
	if len(ids) == 0 {
		return
	}
	cached = make(map[int64]*model.VideoCMS, len(ids))
	idmap, allKeys := keysTreat(ids, d.VideoCMSCacheKey)
	conn := d.mc.Get(c)
	defer conn.Close()
	replys, err := conn.GetMulti(allKeys)
	if err != nil {
		PromError("mc:获取Video信息缓存")
		log.Error("conn.Gets(%v) error(%v)", allKeys, err)
		err = nil
		return
	}
	for key, item := range replys {
		art := &model.VideoCMS{}
		if err = conn.Scan(item, art); err != nil {
			PromError("mc:获取Video信息缓存缓存json解析")
			log.Error("item.Scan(%s) error(%v)", item.Value, err)
			err = nil
			continue
		}
		if art.Deleted == _arcCMSDeleted { // if it's deleted, we ignore it
			log.Info("ArcCMS deleted, %v, %v", item, art)
			continue
		}
		cached[idmap[key]] = art
		delete(idmap, key)
	}
	missed = missedTreat(idmap, len(cached))
	return
}
