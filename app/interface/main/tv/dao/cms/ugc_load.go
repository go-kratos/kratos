package cms

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// LoadArcsMediaMap loads a batch of arc meta
func (d *Dao) LoadArcsMediaMap(ctx context.Context, aids []int64) (resMetas map[int64]*model.ArcCMS, err error) {
	var (
		cachedMetas map[int64]*model.ArcCMS // cache hit seasons
		missedMetas map[int64]*model.ArcCMS // cache miss seasons, pick from DB
		missed      []int64                 // cache miss seasons
		addCache    = true                  // whether we need to fill DB data in MC
	)
	resMetas = make(map[int64]*model.ArcCMS) // merge info from MC and from DB
	// pick up the information for these season ids
	if cachedMetas, missed, err = d.ArcsMetaCache(ctx, aids); err != nil {
		log.Error("LoadArcsMedia ArcsMetaCache Aids:%v, Error:%v", aids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.ArcMetas(ctx, missed); err != nil {
			log.Error("LoadArcsMedia ArcMetas Sids:%v, Error:%v", missed, err)
			return
		}
	}
	// merge info from DB and the info from MC
	for sid, v := range cachedMetas {
		resMetas[sid] = v
	}
	for sid, v := range missedMetas {
		resMetas[sid] = v
	}
	// async Reset the DB data in MC for next time
	log.Info("Set Sids [%d], MissedMetas [%d] Data in MC", len(aids), len(missedMetas))
	if addCache && len(missedMetas) > 0 {
		for _, art := range missedMetas {
			d.AddArcMetaCache(art)
		}
	}
	return
}

// LoadVideosMeta picks the videos meta info
func (d *Dao) LoadVideosMeta(ctx context.Context, cids []int64) (resMetas map[int64]*model.VideoCMS, err error) {
	var (
		cachedMetas map[int64]*model.VideoCMS // cache hit seasons
		missedMetas map[int64]*model.VideoCMS // cache miss seasons, pick from DB
		missed      []int64                   // cache miss seasons
		addCache    = true                    // whether we need to fill DB data in MC
	)
	resMetas = make(map[int64]*model.VideoCMS) // merge info from MC and from DB
	// pick up the information for these season ids
	if cachedMetas, missed, err = d.VideosMetaCache(ctx, cids); err != nil {
		log.Error("LoadVideosMeta VideosMetaCache Aids:%v, Error:%v", cids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.VideoMetas(ctx, missed); err != nil {
			log.Error("LoadVideosMeta VideoMetas Sids:%v, Error:%v", missed, err)
			return
		}
	}
	// merge info from DB and the info from MC
	for sid, v := range cachedMetas {
		resMetas[sid] = v
	}
	for sid, v := range missedMetas {
		resMetas[sid] = v
	}
	// async Reset the DB data in MC for next time
	log.Info("Set Sids [%d], MissedMetas [%d] Data in MC", len(cids), len(missedMetas))
	if addCache && len(missedMetas) > 0 {
		for _, art := range missedMetas {
			d.AddVideoMetaCache(art)
		}
	}
	return
}

// LoadArcsMedia loads the arc meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadArcsMedia(ctx context.Context, aids []int64) (arcs []*model.ArcCMS, err error) {
	var (
		resMetas map[int64]*model.ArcCMS // merge info from MC and from DB
	)
	if resMetas, err = d.LoadArcsMediaMap(ctx, aids); err != nil {
		log.Error("LoadArcsMedia LoadArcsMediaMap Aids: %v, Err: %v", aids, err)
		return
	}
	// re-arrange the info, according to the order got from Redis
	for _, v := range aids {
		if arcCMS, ok := resMetas[v]; !ok {
			log.Error("PickDBeiPage LoadArcsMedia Miss Info for Sid: %d", v)
			continue
		} else {
			arcs = append(arcs, arcCMS)
		}
	}
	return
}

// LoadArcMeta loads the arc meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadArcMeta(ctx context.Context, aid int64) (arcMeta *model.ArcCMS, err error) {
	if arcMeta, err = d.ArcMetaCache(ctx, aid); err != nil { // mc error
		log.Error("LoadArcMedia Get Aid [%d] from CMS Error (%v)", aid, err)
		return
	}
	if arcMeta != nil { // mc found
		return
	}
	if arcMeta, err = d.ArcMetaDB(ctx, aid); err != nil { // db error
		log.Error("LoadArcMedia ArcMetaDB Aid ERROR (%d) (%v)", aid, err)
		return
	}
	if arcMeta == nil { // db not found
		err = ecode.NothingFound
		return
	}
	d.AddArcMetaCache(arcMeta) // db found, re-fill the cache
	return
}

// LoadVideoMeta loads the video meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadVideoMeta(ctx context.Context, cid int64) (videoMeta *model.VideoCMS, err error) {
	if videoMeta, err = d.VideoMetaCache(ctx, cid); err != nil { // mc error
		log.Error("LoadVideoMeta Get Cid [%d] from CMS Error (%v)", cid, err)
		return
	}
	if videoMeta != nil { // mc found
		return
	}
	if videoMeta, err = d.VideoMetaDB(ctx, cid); err != nil { // db error
		log.Error("LoadArcMedia ArcMetaDB Aid ERROR (%d) (%v)", cid, err)
		return
	}
	if videoMeta == nil { // db not found
		err = ecode.NothingFound
		return
	}
	d.AddVideoMetaCache(videoMeta) // db found, re-fill the cache
	return
}
