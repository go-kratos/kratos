package service

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// SnMsg returns the season auth msg
func (s *Service) SnMsg(sid int64) (ok bool, msg string, err error) {
	var season *model.SnAuth
	if season, err = s.cmsDao.SnAuth(ctx, sid); err != nil {
		log.Error("SnMsg LoadSeason Sid %d, Err %v", sid, err)
		return
	}
	ok, msg = s.cmsDao.SnErrMsg(season)
	return
}

// EpMsg returns the ep and its season auth msg
func (s *Service) EpMsg(epid, build int64) (ok bool, msg string, err error) {
	var (
		ep     *model.EpAuth
		season *model.SnAuth
		cfg    = s.conf.Cfg.VipMark.LoadepMsg
		epMeta *model.EpCMS
	)
	if ep, err = s.cmsDao.EpAuth(ctx, epid); err != nil {
		log.Error("EpMsg LoadEP epid %d, Err %v", epid, err)
		return
	}
	if ok, msg = s.cmsDao.EpErrMsg(ep); !ok { // if ep is already not ok, just return, no need to check its season
		return
	}
	if season, err = s.cmsDao.SnAuth(ctx, ep.SeasonID); err != nil { // ep ok, check season
		log.Error("SnMsg LoadSeason Sid %d, Err %v", ep.SeasonID, err)
		return
	}
	if ok, msg = s.cmsDao.SnErrMsg(season); !ok {
		return
	}
	if build < cfg.Build { // old version logic, remind upgrade
		if epMeta, err = s.cmsDao.LoadEpCMS(ctx, epid); err != nil {
			log.Error("EpMsg LoadEpCMS epid %d, Err %v", epid, err)
			return
		}
		if !epMeta.IsFree() { // if old version checks paid ep, remind the user by upgrade message
			return false, cfg.Msg, nil
		}
	}
	return
}
