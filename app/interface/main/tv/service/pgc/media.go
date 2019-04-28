package pgc

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// snAuth picks the season's auth status and validate it
func (s *Service) snAuth(sid int64) (msg string, err error) {
	var season *model.SnAuth
	if season, err = s.cmsDao.SnAuth(ctx, sid); err != nil {
		log.Error("snVerify Sid %d, Err %v", sid, err)
		err = ecode.NothingFound
		return
	}
	if !season.CanPlay() {
		err = ecode.CopyrightLimit
		_, msg = s.cmsDao.SnErrMsg(season) // season auth failure msg
	}
	return
}

func (s *Service) snDecor(core *model.SnDetailCore) {
	if snCMS, err := s.cmsDao.LoadSnCMS(ctx, core.SeasonID); err == nil {
		core.CmsInterv(snCMS)
	} else {
		log.Warn("snDecor LoadSnCMS %d, err %v", core.SeasonID, err)
	}
	core.StyleLabel = s.styleLabel[core.SeasonID]
}

// epAuthDecor checks epids status, hide auditing episodes
func (s *Service) epAuthDecor(epids []int64) (res map[int64]*model.EpDecor, msg string, err error) {
	var (
		epsMetaMap map[int64]*model.EpCMS
		epsAuthMap map[int64]*model.EpAuth
	)
	res = make(map[int64]*model.EpDecor, len(epids))
	if epsAuthMap, err = s.cmsDao.LoadEpsAuthMap(ctx, epids); err != nil {
		log.Warn("FullIntervs LoadEpsAuthMap epids %v, err %v", epids, err)
		return
	}
	if epsMetaMap, err = s.cmsDao.LoadEpsCMS(ctx, epids); err != nil {
		log.Warn("FullIntervs LoadEpsCMS epids %v, err %v", epids, err)
		return
	}
	for _, v := range epids {
		if epAuth, ok := epsAuthMap[v]; ok {
			if epAuth.CanPlay() { // we hide not passed episodes
				decor := &model.EpDecor{
					Watermark: epAuth.Whitelist(),
				}
				if epMeta, okMeta := epsMetaMap[v]; okMeta { // ep intervention
					decor.EpCMS = epMeta
				}
				res[v] = decor
			}
		}
	}
	if len(res) == 0 { // return err
		err = ecode.CopyrightLimit
		msg = s.conf.Cfg.AuthMsg.PGCOffline
		log.Warn("Epids %v, After filter empty", epids)
	}
	return
}

// SnDetail validates the season is authorized to play and involve the intervention
func (s *Service) SnDetail(c context.Context, param *model.MediaParam) (detail *model.SeasonDetail, msg string, err error) {
	var (
		sid    = param.SeasonID
		epids  []int64
		decors map[int64]*model.EpDecor
		cfg    = s.conf.Cfg.VipMark
	)
	if msg, err = s.snAuth(sid); err != nil {
		return
	}
	if detail, err = s.dao.Media(ctx, param); err != nil { // pgc media api
		log.Error("DAO MediaDetail Sid %d, Error (%v)", sid, err)
		return
	}
	// filter auditing eps, and do ep intervention and watermark logic
	for _, v := range detail.Episodes {
		if cfg.V1HideChargeable { // before vip version goes online, we still hide the chargeable episodes
			if v.EpisodeStatus != cfg.EpFree {
				continue
			}
		}
		epids = append(epids, v.EPID)
	}
	if decors, msg, err = s.epAuthDecor(epids); err != nil {
		return
	}
	offAuditing := make([]*model.Episode, 0, len(decors))
	for _, v := range detail.Episodes {
		if decor, ok := decors[v.EPID]; ok {
			if decor.EpCMS != nil {
				v.CmsInterv(decor.EpCMS)
			}
			v.WaterMark = decor.Watermark
			offAuditing = append(offAuditing, v)
		}
	}
	detail.Episodes = offAuditing
	s.snDecor(&detail.SnDetailCore)
	return
}

// SnDetailV2 validates the season is authorized to play and involve the intervention
func (s *Service) SnDetailV2(c context.Context, param *model.MediaParam) (detail *model.SnDetailV2, msg string, err error) {
	var (
		sid      = param.SeasonID
		epids    []int64
		decors   map[int64]*model.EpDecor
		cmarkCfg = s.conf.Cfg.VipMark
	)
	if msg, err = s.snAuth(sid); err != nil {
		return
	}
	if detail, err = s.dao.MediaV2(ctx, param); err != nil || detail == nil { // pgc media api v2
		log.Error("DAO MediaDetail Sid %d, Error (%v)", sid, err)
		return
	}
	detail.TypeTrans()
	if len(detail.Section) > 0 { // pgc media api v2 logic, prevues are in the sections, we need to pick them up and re-insert into the episodes list
		for _, v := range detail.Section {
			detail.Episodes = append(detail.Episodes, v.Episodes...)
		}
	}
	for _, v := range detail.Episodes {
		epids = append(epids, v.ID)
	}
	if decors, msg, err = s.epAuthDecor(epids); err != nil {
		return
	}
	offAuditing := make([]*model.EpisodeV2, 0, len(decors))
	for _, v := range detail.Episodes {
		if decor, ok := decors[v.ID]; ok {
			if decor.EpCMS != nil {
				v.CmsInterv(decor.EpCMS)
			}
			v.WaterMark = decor.Watermark
			if v.Status != cmarkCfg.EpFree { // if ep is not free, put the corner mark
				v.CornerMark = &(*cmarkCfg.EP)
			}
			offAuditing = append(offAuditing, v)
		}
	}
	detail.Episodes = offAuditing
	s.snDecor(&detail.SnDetailCore)
	return
}

// EpControl validates the ep is authorized to play and involve the intervention
func (s *Service) EpControl(c context.Context, epid int64) (sid int64, msg string, err error) {
	var ep *model.EpAuth
	if ep, err = s.cmsDao.EpAuth(c, epid); err != nil {
		log.Error("LoadEP Epid %d Error(%v)", epid, err)
		err = ecode.NothingFound
		return
	}
	if !ep.CanPlay() {
		err = ecode.CopyrightLimit
		_, msg = s.cmsDao.EpErrMsg(ep) // ep auth failure msg
		return
	}
	sid = ep.SeasonID
	return
}

func (s *Service) upStyleCache() {
	for {
		res, err := s.dao.GetLabelCache(context.Background())
		if err != nil {
			log.Error("s.dao.GetLabelCache upStyleCache error(%s)", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(res) > 0 {
			s.styleLabel = res
		}
		time.Sleep(time.Duration(s.conf.Style.LabelSpan))
	}
}

func (s *Service) styleCache() {
	res, err := s.dao.GetLabelCache(context.Background())
	if err != nil {
		log.Error("s.dao.GetLabelCache error(%v)", err)
		panic(err)
	}
	s.styleLabel = res
}
