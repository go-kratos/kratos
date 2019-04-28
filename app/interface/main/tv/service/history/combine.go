package history

import (
	"context"

	hismdl "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/history"
	"go-common/library/log"
)

func (s *Service) pgcHisRes(ctx context.Context, res []*hismdl.Resource) (resMap map[int64]*history.HisRes, err error) {
	var (
		snMetas   map[int64]*model.SeasonCMS
		epMetas   map[int64]*model.EpCMS
		pickSids  []int64
		pickEpids []int64
	)
	resMap = make(map[int64]*history.HisRes)
	for _, v := range res {
		pickSids = append(pickSids, v.Sid)
		pickEpids = append(pickEpids, v.Epid)
	}
	if snMetas, err = s.cmsDao.LoadSnsCMSMap(ctx, pickSids); err != nil {
		log.Error("LoadSnsCMS Sids %v, Err %v", pickSids, err)
		return
	}
	if epMetas, err = s.cmsDao.LoadEpsCMS(ctx, pickEpids); err != nil {
		log.Warn("LoadEpsCMS Epids %v, Err %v", pickEpids, err)
		err = nil
	}
	for _, v := range res {
		his := hisTrans(v)
		his.Type = _typePGC
		his.Page = nil
		// season info
		snMeta, okS := snMetas[v.Sid]
		if !okS {
			log.Error("pgcHisRes Missing Info Sid %d", v.Sid)
			continue
		}
		his.Title = snMeta.Title
		his.Cover = snMeta.Cover
		if snMeta.NeedVip() { // add vip corner mark
			his.CornerMark = &(*s.conf.Cfg.SnVipCorner)
		}
		// ep info
		epMeta, okE := epMetas[v.Epid]
		if !okE {
			log.Warn("pgcHisRes Missing Info Epid %d", v.Epid)
		} else {
			his.EPMeta = &history.HisEP{
				EPID:      epMeta.EPID,
				Cover:     epMeta.Cover,
				Title:     epMeta.Subtitle,
				LongTitle: epMeta.Title,
			}
		}
		resMap[v.Sid] = his
	}
	return
}

func (s *Service) ugcHisRes(ctx context.Context, res []*hismdl.Resource) (resMap map[int64]*history.HisRes, err error) {
	var (
		arcMetas   map[int64]*model.ArcCMS
		videoMetas map[int64]*model.VideoCMS
		pickAids   []int64
		pickCids   []int64
	)
	resMap = make(map[int64]*history.HisRes)
	for _, v := range res {
		pickAids = append(pickAids, v.Oid)
		pickCids = append(pickCids, v.Cid)
	}
	if arcMetas, err = s.cmsDao.LoadArcsMediaMap(ctx, pickAids); err != nil {
		log.Error("LoadArcsMediaMap Sids %v, Err %v", pickAids, err)
		return
	}
	if videoMetas, err = s.cmsDao.LoadVideosMeta(ctx, pickCids); err != nil {
		log.Warn("LoadVideosMeta Epids %v, Err %v", pickCids, err)
		err = nil
	}
	for _, v := range res {
		his := hisTrans(v)
		his.Type = _typeUGC
		his.Page = nil
		// season info
		arcMeta, okS := arcMetas[v.Oid]
		if !okS {
			log.Error("ugcHisRes Missing Info Aid %d", v.Oid)
			continue
		}
		his.Title = arcMeta.Title
		his.Cover = arcMeta.Cover
		// ep info
		video, okE := videoMetas[v.Cid]
		if !okE {
			log.Warn("ugcHisRes Missing Info Cid %d", v.Cid)
		} else {
			his.Page = &history.HisPage{
				CID:  video.CID,
				Part: video.Title,
				Page: video.IndexOrder,
			}
		}
		resMap[v.Oid] = his
	}
	return
}

func hisTrans(res *hismdl.Resource) *history.HisRes {
	return &history.HisRes{
		Mid:      res.Mid,
		Oid:      res.Oid,
		Sid:      res.Sid,
		Epid:     res.Epid,
		Cid:      res.Cid,
		Business: res.Business,
		DT:       res.DT,
		Pro:      res.Pro,
		Unix:     res.Unix,
		Type:     _typePGC,
	}
}

func (s *Service) getDuration(ctx context.Context, res []*hismdl.Resource) (durs map[int64]int64) {
	var (
		aids []int64
	)
	durs = make(map[int64]int64)
	for _, v := range res {
		aids = append(aids, v.Oid)
	}
	resMeta := s.arcDao.LoadViews(ctx, aids)
	for _, v := range res {
		if view, ok := resMeta[v.Oid]; ok && len(view.Pages) > 0 {
			for _, vp := range view.Pages {
				if v.Cid == vp.Cid {
					durs[v.Oid] = vp.Duration
					break
				}
			}
		}
	}
	return
}
