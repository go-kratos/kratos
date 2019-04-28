package service

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_retry = 3
)

// RecomFilter gets the recommend data from PGC API and filter the not passed seasons
func (s *Service) RecomFilter(sid string, stype string) (res []*model.Recom, err error) {
	var (
		sids   []int64
		result []*model.Recom
		cmsRes map[int64]*model.SeasonCMS
	)
	log.Info("[RecomFilter] Sid: %s, Stype: %s", sid, stype)
	for i := 0; i < _retry; i++ {
		if result, err = s.dao.RecomData(ctx, s.TVAppInfo, sid, stype); err == nil {
			break
		}
	}
	if err != nil {
		log.Error("[RecomFilter] Can't Pick PGC Recom Data, Err: %v", err)
		return
	}
	if len(result) == 0 {
		log.Error("[RecomFilter] No need to filter for Sid: %s. Length = 0", sid)
		return
	}
	for _, v := range result {
		season, err2 := s.cmsDao.SnAuth(ctx, v.SeasonID)
		if err != nil {
			log.Error("[RecomFilter] LoadSeason[%d] ERROR [%v]", sid, err2)
			continue
		}
		if season == nil {
			log.Info("[RecomFilter] LoadSeason[%d] Can't Found", v.SeasonID)
			continue
		}
		if !season.CanPlay() {
			log.Info("[RecomFilter] SEASON[%d] is not authorized to play, DETAILS: %v", season.ID, season)
			continue
		}
		res = append(res, v)
		sids = append(sids, v.SeasonID)
	}
	// add vip corner mark
	if cmsRes, err = s.cmsDao.LoadSnsCMSMap(ctx, sids); err != nil {
		log.Error("[recomm.RecomFilter] sids(%s) error(%v)", xstr.JoinInts(sids), err)
		return
	}
	for idx, v := range res {
		if r, ok := cmsRes[v.SeasonID]; ok && r.NeedVip() {
			res[idx].CornerMark = &(*s.conf.Cfg.SnVipCorner)
		}
	}
	log.Info("[RecomFilter] PGC Data: %d, Filtered Data: %d", len(result), len(res))
	return
}
