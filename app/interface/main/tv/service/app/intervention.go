package service

import (
	"time"

	"context"
	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/search"
	"go-common/library/log"
)

// cardIntervSn, makes season intervention effective for cards
func (s *Service) cardIntervSn(cards []*model.Card) (err error) {
	var (
		snMetas     map[int64]*model.SeasonCMS
		epMetas     map[int64]*model.EpCMS
		sids        []int64
		newestEPIDs []int64
	)
	for _, card := range cards {
		if card.NewEP != nil {
			sids = append(sids, int64(card.SeasonID))
			newestEPIDs = append(newestEPIDs, card.NewEP.ID)
		}
	}
	if snMetas, err = s.cmsDao.LoadSnsCMSMap(ctx, sids); err != nil {
		log.Error("[cardIntervSn] loadSnCMSMap Sids: %v, Err %v", sids, err)
		return
	}
	if epMetas, err = s.cmsDao.LoadEpsCMS(ctx, newestEPIDs); err != nil {
		log.Error("[cardIntervSn] loadEpCMS epids: %v, Err %v", newestEPIDs, err)
		return
	}
	for _, card := range cards {
		// season intervention
		sid := int64(card.SeasonID)
		if snCache, ok := snMetas[sid]; !ok {
			log.Error("LoadSnsCMS Miss Info for Sid: %d", sid)
			continue
		} else { // intervention
			if snCache.Title != "" {
				card.Title = snCache.Title
			}
			if snCache.Cover != "" {
				card.Cover = snCache.Cover
			}
			if snCache.NeedVip() { // card add vip corner mark
				card.CornerMark = &(*s.conf.Cfg.SnVipCorner)
			}
		}
		// ep intervention
		epid := card.NewEP.ID
		if epCache, ok := epMetas[epid]; !ok {
			log.Error("LoadSnsCMS Miss Info for Sid: %d", epid)
			continue
		} else { // intervention
			if epCache.Cover != "" {
				card.NewEP.Cover = epCache.Cover
			}
		}
	}
	return
}

// load index show proc
func (s *Service) indexShowproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.IndexShowReload))
		s.indexShow()
	}
}

func (s *Service) indexShow() (err error) {
	defer elapsed("indexShow")()
	var (
		indexShows map[int64]string
		sids       []int64
	)
	if sids, err = s.dao.PassedSns(ctx); err != nil {
		log.Error("[loadIndexShow] AllIntervs Error %v", err)
		return
	}
	if indexShows, err = s.PgcCards(sids); err != nil {
		log.Error("[loadIndexShow] AllIntervs Error %v", err)
		return
	}
	log.Info("Reload Types For Index_Show, Origin:%d, Length: %d", len(sids), len(indexShows))
	if len(indexShows) > 0 {
		s.PGCIndexShow = indexShows
	}
	return
}

func (s *Service) filterIntervs(ctx context.Context) (err error) {
	defer elapsed("filterIntervs")()
	var (
		sids, aids, rmSids, rmAids []int64
		pgcAuth                    map[int64]*model.SnAuth
		ugcAuth                    map[int64]*model.ArcCMS
	)
	if sids, aids, err = s.dao.AllIntervs(ctx); err != nil {
		log.Error("[filterIntervs] AllIntervs Error %v", err)
		return
	}
	if pgcAuth, err = s.cmsDao.LoadSnsAuthMap(ctx, sids); err != nil {
		log.Error("[filterIntervs] LoadSnsAuthMap Error %v, Sids %v", err, sids)
		return
	}
	if ugcAuth, err = s.cmsDao.LoadArcsMediaMap(ctx, aids); err != nil {
		log.Error("[filterIntervs] LoadArcsMediaMap Error %v, Aids %v", err, aids)
		return
	}
	for sid, pAuth := range pgcAuth {
		if !pAuth.CanPlay() {
			rmSids = append(rmSids, sid)
		}
	}
	for aid, arcAuth := range ugcAuth {
		if !arcAuth.CanPlay() {
			rmAids = append(rmAids, aid)
		}
	}
	if len(rmSids) > 0 || len(rmAids) > 0 {
		if err = s.dao.RmInterv(ctx, rmAids, rmSids); err != nil {
			log.Error("RmInterv Aids %v, Sids %v, Err %v", aids, sids, err)
			return
		}
		log.Warn("[filterIntervs] Ori Sids %d, Aids %d, To Remove Sids %v, Aids %v", len(sids), len(aids), rmSids, rmAids)
	}
	return
}

func (s *Service) ppIdxIntev(ctx context.Context) (err error) {
	var (
		newIdx *search.IdxIntervSave
	)
	if newIdx, err = s.dao.IdxIntervs(ctx); err != nil {
		log.Error("prepareIdxIntervs Err %v")
		return
	}
	s.IdxIntervs = newIdx
	log.Info("ppIdxInterv refresh, PGC %d, UGC %d", len(newIdx.Pgc), len(newIdx.Ugc))
	return
}

func (s *Service) ppIdxIntervproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.EsIntervReload))
		s.ppIdxIntev(context.Background())
	}
}
