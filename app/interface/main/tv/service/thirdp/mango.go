package thirdp

import (
	"context"
	"math"

	"go-common/app/interface/main/tv/dao/thirdp"
	"go-common/app/interface/main/tv/model"
	tpMdl "go-common/app/interface/main/tv/model/thirdp"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// mangoPage picks the Mango's page
func (s *Service) mangoPage(page int64, typeC string) (pager *model.IdxPager, dataSet []*tpMdl.RespSid, err error) {
	var (
		cPageID int64
		req     = &tpMdl.ReqDBeiPages{
			Page:  page,
			TypeC: typeC,
			Ps:    int64(s.conf.Cfg.Dangbei.MangoPS),
		}
	)
	if pager, err = s.buildPager(req); err != nil {
		return
	}
	if req.LastID, err = s.dao.LoadPageID(ctx, req); err != nil {
		log.Error("MangoPage getPageID Page %d Miss, Pass by offset", page-1)
		return
	}
	if dataSet, cPageID, err = s.dao.MangoPages(ctx, req); err != nil {
		return
	}
	if len(dataSet) == 0 {
		err = errors.Wrapf(ecode.NothingFound, "Type_C [%s], Page [%d], Offset Result Empty", typeC, page)
		return
	}
	cache.Save(func() {
		s.dao.SetPageID(ctx, &tpMdl.ReqPageID{Page: page, ID: cPageID, TypeC: typeC})
	})
	return
}

// MangoSns picks mango season pages
func (s *Service) MangoSns(ctx context.Context, page int64) (data *tpMdl.MangoSnPage, err error) {
	var (
		pager       *model.IdxPager
		dataSet     []*tpMdl.RespSid
		snMetas     = map[int64]*model.SeasonCMS{}
		snAuths     = map[int64]*model.SnAuth{}
		newestEpids []int64
		epMetas     = map[int64]*model.EpCMS{}
	)
	if pager, dataSet, err = s.mangoPage(page, thirdp.MangoPGC); err != nil {
		return
	}
	data = &tpMdl.MangoSnPage{
		Pager: pager,
	}
	sids := tpMdl.PickSids(dataSet)
	if snMetas, err = s.cmsDao.LoadSnsCMSMap(ctx, sids); err != nil {
		log.Error("MangoSns - PGC - LoadSnsCMS - Sids %v, Error %v", sids, err)
		return
	}
	if snAuths, err = s.cmsDao.LoadSnsAuthMap(ctx, sids); err != nil {
		log.Error("MangoSns - PGC - LoadSnsAuthMap - Sids %v, Error %v", sids, err)
		return
	}
	for _, v := range snMetas { // pick newestEpids
		if v.NewestEPID != 0 {
			newestEpids = append(newestEpids, v.NewestEPID)
		}
	}
	if len(newestEpids) > 0 { // pick eps cms meta info
		if epMetas, err = s.cmsDao.LoadEpsCMS(ctx, newestEpids); err != nil {
			log.Error("MangoSns - PGC - LoadEpsCMS - Epids %v, Error %v", newestEpids, err)
			return
		}
	}
	for _, v := range dataSet { // transform the object to DbeiSeason
		var (
			snMeta         *model.SeasonCMS
			snAuth         *model.SnAuth
			okMeta, okAuth bool
		)
		if snMeta, okMeta = snMetas[v.Sid]; okMeta {
			if snAuth, okAuth = snAuths[v.Sid]; okAuth {
				mangoSn := tpMdl.ToMangoSn(snMeta, v.Mtime, snAuth.CanPlay())
				if newestEp := snMeta.NewestEPID; newestEp != 0 {
					if epMeta, ok := epMetas[snMeta.NewestEPID]; ok {
						mangoSn.EpCover = epMeta.Cover
					}
				}
				data.List = append(data.List, mangoSn)
				continue
			}
		}
		log.Warn("MangoSns Sid %d Missing Info, Meta %v, Auth %v", v.Sid, okMeta, okAuth)
	}
	return
}

// MangoArcs picks mango archive pages
func (s *Service) MangoArcs(ctx context.Context, page int64) (data *tpMdl.MangoArcPage, err error) {
	var (
		pager    *model.IdxPager
		dataSet  []*tpMdl.RespSid
		arcMetas map[int64]*model.ArcCMS
	)
	if pager, dataSet, err = s.mangoPage(page, thirdp.MangoUGC); err != nil {
		return
	}
	data = &tpMdl.MangoArcPage{
		Pager: pager,
	}
	sids := tpMdl.PickSids(dataSet)
	if arcMetas, err = s.cmsDao.LoadArcsMediaMap(ctx, sids); err != nil {
		log.Error("MangoArcs - UGC - LoadArcsMediaMap - Sids %v, Error %v", sids, err)
		return
	}
	for _, v := range dataSet { // transform the object to DbeiSeason
		if arcMeta, ok := arcMetas[v.Sid]; ok {
			cat1, cat2 := s.arcDao.GetPTypeName(int32(arcMeta.TypeID))
			data.List = append(data.List, tpMdl.ToMangoArc(arcMeta, v.Mtime, cat1, cat2))
			continue
		}
		log.Warn("MangoSns Aid %d Missing Info", v.Sid)
	}
	return
}

// MangoEps returns mango eps data
func (s *Service) MangoEps(ctx context.Context, sid int64, page int) (data *tpMdl.MangoEpPage, err error) {
	var (
		count    int
		pagesize = s.conf.Cfg.Dangbei.MangoPS
		resp     []*tpMdl.RespSid
		epMetas  map[int64]*model.EpCMS
		epAuths  map[int64]*model.EpAuth
	)
	if count, err = s.dao.LoadSnCnt(ctx, true, sid); err != nil {
		log.Error("MangoEps LoadSnCnt Sid %d, Err %v", sid, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pagesize)))
	if page > totalPages {
		return nil, ecode.TvDangbeiPageNotExist
	}
	data = &tpMdl.MangoEpPage{
		SeasonID: sid,
		Pager: &model.IdxPager{
			CurrentPage: page,
			TotalItems:  count,
			TotalPages:  int(math.Ceil(float64(count) / float64(pagesize))),
			PageSize:    int(pagesize),
		},
	}
	if resp, err = s.dao.MangoSnOffset(ctx, true, sid, page, pagesize); err != nil {
		log.Error("MangoEps MangoSnOffset Sid %d, Err %v", sid, err)
		return
	}
	epids := tpMdl.PickSids(resp)
	if epMetas, err = s.cmsDao.LoadEpsCMS(ctx, epids); err != nil {
		log.Error("MangoEps LoadEpsCMS Sid %d, Err %v", sid, err)
		return
	}
	if epAuths, err = s.cmsDao.LoadEpsAuthMap(ctx, epids); err != nil {
		log.Error("MangoEps LoadEpsAuthMap Sid %d, Err %v", sid, err)
		return
	}
	for _, v := range resp {
		var (
			epMeta         *model.EpCMS
			epAuth         *model.EpAuth
			okMeta, okAuth bool
		)
		if epMeta, okMeta = epMetas[v.Sid]; okMeta {
			if epAuth, okAuth = epAuths[v.Sid]; okAuth {
				data.List = append(data.List, &tpMdl.MangoEP{
					EpCMS:     *epMeta,
					SeasonID:  sid,
					Mtime:     v.Mtime,
					Autorised: epAuth.CanPlay(),
				})
				continue
			}
		}
		log.Warn("MangoEps Sid %d, Epids %d Missing Info, Meta %v, Auth %v", sid, v.Sid, okMeta, okAuth)
	}
	return
}

// MangoVideos returns mango videos data
func (s *Service) MangoVideos(ctx context.Context, sid int64, page int) (data *tpMdl.MangoVideoPage, err error) {
	var (
		count      int
		pagesize   = s.conf.Cfg.Dangbei.MangoPS
		resp       []*tpMdl.RespSid
		videoMetas map[int64]*model.VideoCMS
		vp         *arcwar.ViewReply
		vPages     = make(map[int64]*arcwar.Page)
	)
	if count, err = s.dao.LoadSnCnt(ctx, false, sid); err != nil {
		log.Error("MangoVideos LoadSnCnt Sid %d, Err %v", sid, err)
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(pagesize)))
	if page > totalPages {
		return nil, ecode.TvDangbeiPageNotExist
	}
	data = &tpMdl.MangoVideoPage{
		AVID: sid,
		Pager: &model.IdxPager{
			CurrentPage: page,
			TotalItems:  count,
			TotalPages:  int(math.Ceil(float64(count) / float64(pagesize))),
			PageSize:    int(pagesize),
		},
	}
	if resp, err = s.dao.MangoSnOffset(ctx, false, sid, page, pagesize); err != nil {
		log.Error("MangoVideos MangoSnOffset Sid %d, Err %v", sid, err)
		return
	}
	epids := tpMdl.PickSids(resp)
	if videoMetas, err = s.cmsDao.LoadVideosMeta(ctx, epids); err != nil {
		log.Error("MangoVideos LoadEpsCMS Sid %d, Err %v", sid, err)
		return
	}
	if vp, err = s.arcDao.GetView(ctx, sid); err != nil {
		log.Error("MangoVideos ViewPage getView Aid:%d, Err:%v", sid, err)
		return
	}
	for _, v := range vp.Pages {
		vPages[v.Cid] = v
	}
	for _, v := range resp {
		var (
			vMeta          *model.VideoCMS
			view           *arcwar.Page
			okMeta, okView bool
		)
		if vMeta, okMeta = videoMetas[v.Sid]; okMeta {
			if view, okView = vPages[v.Sid]; okView {
				data.List = append(data.List, &tpMdl.MangoVideo{
					CID:       v.Sid,
					Page:      vMeta.IndexOrder,
					Desc:      view.Desc,
					Title:     vMeta.Title,
					Duration:  view.Duration,
					Autorised: vMeta.CanPlay(),
					Mtime:     v.Mtime,
				})
				continue
			}
		}
		log.Warn("MangoViews Sid %d, Epids %d Missing Info, Meta %v, View %v", sid, v.Sid, okMeta, okView)
	}
	return
}
