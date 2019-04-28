package cms

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// LoadSnsAuthMap loads a batch of arc meta
func (d *Dao) LoadSnsAuthMap(ctx context.Context, sids []int64) (resMetas map[int64]*model.SnAuth, err error) {
	var (
		cachedMetas map[int64]*model.SnAuth // cache hit seasons
		missedMetas map[int64]*model.SnAuth // cache miss seasons, pick from DB
		missed      []int64                 // cache miss seasons
		addCache    = true                  // whether we need to fill DB data in MC
	)
	resMetas = make(map[int64]*model.SnAuth) // merge info from MC and from DB
	if cachedMetas, missed, err = d.snAuthCache(ctx, sids); err != nil {
		log.Error("LoadSnsAuthMap snAuthCache Sids:%v, Error:%v", sids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.SnsAuthDB(ctx, missed); err != nil {
			log.Error("LoadSnsAuthMap SnsAuthDB Sids:%v, Error:%v", missed, err)
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
	if addCache && len(missedMetas) > 0 {
		for _, snAuth := range missedMetas {
			d.addCache(func() {
				d.AddSnAuthCache(ctx, snAuth)
			})
		}
	}
	return
}

// LoadEpsAuthMap loads a batch of arc meta
func (d *Dao) LoadEpsAuthMap(ctx context.Context, epids []int64) (resMetas map[int64]*model.EpAuth, err error) {
	var (
		cachedMetas map[int64]*model.EpAuth // cache hit seasons
		missedMetas map[int64]*model.EpAuth // cache miss seasons, pick from DB
		missed      []int64                 // cache miss seasons
		addCache    = true                  // whether we need to fill DB data in MC
	)
	resMetas = make(map[int64]*model.EpAuth) // merge info from MC and from DB
	if cachedMetas, missed, err = d.epAuthCache(ctx, epids); err != nil {
		log.Error("LoadEpsAuthMap epAuthCache epids:%v, Error:%v", epids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.EpsAuthDB(ctx, missed); err != nil {
			log.Error("LoadEpsAuthMap EpsAuthDB epids:%v, Error:%v", missed, err)
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
	if addCache && len(missedMetas) > 0 {
		for _, epAuth := range missedMetas {
			d.addCache(func() {
				d.AddEpAuthCache(ctx, epAuth)
			})
		}
	}
	return
}

// LoadSnsCMSMap loads season cms meta data from cache and db
func (d *Dao) LoadSnsCMSMap(ctx context.Context, sids []int64) (resMetas map[int64]*model.SeasonCMS, err error) {
	var (
		cachedMetas, missedMetas map[int64]*model.SeasonCMS // cache hit seasons
		missed                   []int64                    // cache miss seasons
		addCache                 = true                     // whether we need to fill DB data in MC
	)
	resMetas = make(map[int64]*model.SeasonCMS)
	// pick up the information for these season ids
	if cachedMetas, missed, err = d.SeasonsMetaCache(ctx, sids); err != nil {
		log.Error("LoadSnsCMS SeasonMetaCache Sids:%v, Error:%v", sids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.SeasonMetas(ctx, missed); err != nil {
			log.Error("LoadSnsCMS SeasonMetas Sids:%v, Error:%v", sids, err)
			return
		}
	}
	log.Info("Set Sids [%d], HitMetas [%d], MissedMetas [%d][%d] Data in MC", len(sids), len(cachedMetas), len(missed), len(missedMetas))
	// merge info from DB and the info from MC
	for sid, v := range cachedMetas {
		resMetas[sid] = v
	}
	for sid, v := range missedMetas {
		resMetas[sid] = v
	}
	// async Reset the DB data in MC for next time
	if addCache && len(missedMetas) > 0 {
		for _, art := range missedMetas {
			d.addCache(func() {
				d.AddSeasonMetaCache(ctx, art)
			})
		}
	}
	return
}

// LoadSnsCMS loads the seasons meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadSnsCMS(ctx context.Context, sids []int64) (seasons []*model.SeasonCMS, newestEpids []int64, err error) {
	var (
		resMetas map[int64]*model.SeasonCMS // merge info from MC and from DB
	)
	if resMetas, err = d.LoadSnsCMSMap(ctx, sids); err != nil {
		log.Error("LoadSnsCMS Sids %v, Err %v", sids, err)
		return
	}
	// re-arrange the info, according to the order got from Redis
	for _, v := range sids {
		if SnCMS, ok := resMetas[v]; !ok {
			log.Error("LoadSnsCMS Miss Info for Sid: %d", v)
			continue
		} else {
			seasons = append(seasons, SnCMS)
			newestEpids = append(newestEpids, SnCMS.NewestEPID)
		}
	}
	return
}

// LoadSnCMS loads the sn meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadSnCMS(ctx context.Context, sid int64) (sn *model.SeasonCMS, err error) {
	if sn, err = d.GetSnCMSCache(ctx, sid); err != nil {
		log.Error("LoadSnsCMS Get Season[%d] from CMS Error (%v)", sid, err) // cache set/get error
		return
	}
	if sn != nil { // if cache hit, return
		return
	}
	if sn, err = d.SeasonCMS(ctx, sid); err != nil {
		log.Error("[LoadSnCMS] SeasonCMS SeasonID ERROR (%d) (%v)", sid, err)
		return
	} else if sn == nil {
		err = ecode.NothingFound
		return
	}
	d.addCache(func() {
		d.AddSeasonMetaCache(ctx, sn)
	})
	return
}

// LoadEpCMS loads the sn meta cms data from cache, for missed ones, pick them from the DB
func (d *Dao) LoadEpCMS(ctx context.Context, epid int64) (ep *model.EpCMS, err error) {
	if ep, err = d.GetEpCMSCache(ctx, epid); err != nil {
		log.Error("LoadEpCMS Get EP[%d] from CMS Error (%v)", epid, err) // cache set/get error
		return
	} else if ep == nil {
		if ep, err = d.EpCMS(ctx, epid); err != nil {
			log.Error("[LoadEpCMS] EpCMS Epid ERROR (%d) (%v)", epid, err)
			return
		} else if ep == nil {
			err = ecode.NothingFound
			return
		}
	}
	d.addCache(func() {
		d.SetEpCMSCache(ctx, ep)
	})
	return
}

// LoadEpsCMS picks ep meta information from Cache & DB
func (d *Dao) LoadEpsCMS(ctx context.Context, epids []int64) (resMetas map[int64]*model.EpCMS, err error) {
	var (
		cachedMetas, missedMetas map[int64]*model.EpCMS
		missed                   []int64
		addCache                 = true
	)
	resMetas = make(map[int64]*model.EpCMS)
	// pick up the information for these season ids
	if cachedMetas, missed, err = d.EpMetaCache(ctx, epids); err != nil {
		log.Error("loadEpCMS EpMetaCache Sids:%v, Error:%v", epids, err)
		err = nil
		addCache = false // mc error, we don't add
	}
	if len(missed) > 0 {
		if missedMetas, err = d.EpMetas(ctx, missed); err != nil {
			log.Error("loadEpCMS EpMetas Sids:%v, Error:%v", epids, err)
			return
		}
	}
	// merge info from DB and the info from MC
	resMetas = make(map[int64]*model.EpCMS, len(epids))
	for sid, v := range cachedMetas {
		resMetas[sid] = v
	}
	for sid, v := range missedMetas {
		resMetas[sid] = v
	}
	log.Info("Combine Info for %d Epids, Origin Length %d", len(epids), len(resMetas))
	if addCache && len(missedMetas) > 0 { // async Reset the DB data in MC for next time
		log.Info("Set MissedMetas %d Data in MC", missedMetas)
		for _, art := range missedMetas {
			d.addCache(func() {
				d.AddEpMetaCache(ctx, art)
			})
		}
	}
	return
}
