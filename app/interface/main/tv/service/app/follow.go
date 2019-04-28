package service

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// FollowData gets the follow data from pgc api
func (s *Service) FollowData(ctx context.Context, accessKey string) (res []*model.Follow) {
	var (
		err    error
		result []*model.Follow
	)
	result, err = s.dao.FollowData(ctx, s.TVAppInfo, accessKey)
	if err != nil {
		log.Error("[LoadHP] Can't Pick PGC Follow Data, Err: %v", err)
	}
	res = s.followInterv(result)
	if res == nil {
		log.Error("Follow Data is Nil!")
	}
	return
}

// Intervention
func (s *Service) followInterv(result []*model.Follow) (newRes []*model.Follow) {
	for _, v := range result {
		var (
			sid     = int64(atoi(v.SeasonID))
			epid    int64
			err     error
			snCache *model.SeasonCMS
			epCache *model.EpCMS
		)
		// filter not passed data
		season, err := s.cmsDao.SnAuth(ctx, sid)
		if err != nil {
			log.Error("followInterv LoadSeason[%d] ERROR [%v]", sid, err)
			continue
		}
		if season == nil {
			log.Info("followInterv LoadSeason[%d] Can't Found", sid)
			continue
		}
		if !(season.IsDeleted == 0 && season.Check == 1 && season.Valid == 1) {
			log.Info("SEASON[%d] is not authorized to play, DETAILS: %v", sid, season)
			continue
		}
		// season intervention
		snCache, err = s.cmsDao.GetSnCMSCache(ctx, sid)
		if err != nil {
			log.Error("[cardInterv] ErrorCache GetSnCMSCache SeasonID(%d) (%v)", sid, err)
		} else if snCache != nil {
			if snCache.Title != "" {
				v.Title = snCache.Title
			}
			if snCache.Cover != "" {
				v.Cover = snCache.Cover
			}
			if snCache.NeedVip() {
				v.CornerMark = &(*s.conf.Cfg.SnVipCorner)
			}
		}
		// ep intervention
		if v.NewEP == nil {
			continue
		}
		epid = int64(atoi(v.NewEP.EpisodeID))
		epCache, err = s.cmsDao.GetEpCMSCache(ctx, epid)
		if err != nil {
			log.Error("[cardInterv] ErrorCache GetEpCMSCache EPID(%d) (%v)", epid, err)
		} else if epCache != nil {
			if epCache.Title != "" {
				v.NewEP.IndexTitle = epCache.Title
			}
			if epCache.Cover != "" {
				v.NewEP.Cover = epCache.Cover
			}
		}
		newRes = append(newRes, v)
	}
	return
}
