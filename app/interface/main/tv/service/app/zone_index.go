package service

import (
	"context"
	"math"

	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/search"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/log"
)

// LoadZoneIdx loads zone index page data
func (s *Service) LoadZoneIdx(page int, category int) (idxSns []*model.IdxSeason, pager *model.IdxPager, err error) {
	pagesize := s.conf.Cfg.ZonePs
	var (
		sids    []int64
		count   int
		start   = (page - 1) * pagesize
		end     = page*pagesize - 1
		seasons []*model.SeasonCMS
		arcs    []*model.ArcCMS
	)
	idxSns = make([]*model.IdxSeason, 0)
	// pick up the page sids
	if sids, count, err = s.dao.ZrevrangeList(ctx, category, start, end); err != nil {
		log.Error("LoadZoneIdx - ZrevrangeList - Category %d Start %d End %d, Error %v", category, start, end, err)
		return
	}
	pager = &model.IdxPager{
		CurrentPage: page,
		TotalItems:  count,
		TotalPages:  int(math.Ceil(float64(count) / float64(pagesize))),
		PageSize:    pagesize,
	}
	if len(sids) == 0 {
		return
	}
	if s.catIsUGC(category) {
		if arcs, err = s.cmsDao.LoadArcsMedia(ctx, sids); err != nil {
			log.Error("PickDBeiPage - UGC - LoadArcsMedia - Sids %v, Error %v", sids, err)
			return
		}
		for _, v := range arcs {
			idxSns = append(idxSns, v.ToIdxSn())
		}
	} else {
		if seasons, _, err = s.cmsDao.LoadSnsCMS(ctx, sids); err != nil {
			log.Error("LoadZoneIdx - LoadSnCMS - Category %d Start %d End %d, Error %v", category, start, end, err)
			return
		}
		for _, v := range seasons {
			idxSn := v.IdxSn()
			if idxShow, ok := s.PGCIndexShow[v.SeasonID]; ok { // use pgc index_show to replace the DB upinfo
				idxSn.Upinfo = idxShow
			}
			idxSns = append(idxSns, idxSn)
		}
	}
	log.Info("Combine Info for %d Sids, Page %d, Category %d", len(idxSns), page, category)
	return
}

func (s *Service) catIsUGC(category int) bool {
	for _, v := range s.conf.Cfg.ZonesInfo.UGCZonesID {
		if category == int(v) {
			return true
		}
	}
	return false
}

// esInterv treats the es index intervention
func (s *Service) esInterv(ctx context.Context, req *search.ReqIdxInterv) (resIDs []int64, err error) {
	var intervIDs []int64
	if s.IdxIntervs != nil {
		if req.IsPGC {
			if sids, ok := s.IdxIntervs.Pgc[req.Category]; ok {
				intervIDs = sids
			}
		} else {
			if aids, ok := s.IdxIntervs.Ugc[req.Category]; ok {
				intervIDs = aids
			}
		}
	}
	resIDs = applyInterv(req.EsIDs, intervIDs, req.Pn)
	return
}

// EsPgcIdx returns the elastic search index page result
func (s *Service) EsPgcIdx(ctx context.Context, req *search.ReqPgcIdx) (res *search.EsPager, err error) {
	var (
		data             *search.EsPgcResult
		esSids, authSids []int64
		target           []*model.Card
		idxCards         []*search.EsCard
		authMap          map[int64]*model.SnAuth
	)
	if data, err = s.searchDao.PgcIdx(ctx, req); err != nil {
		log.Error("EsPgcIdx Req %v, Err %v", req, err)
		return
	}
	for _, v := range data.Result {
		esSids = append(esSids, v.SeasonID)
	}
	if authMap, err = s.cmsDao.LoadSnsAuthMap(ctx, esSids); err != nil {
		log.Error("EsPgcIdx Sids %v, LoadSnsAuthMap Err %v", esSids, err)
		return
	}
	for _, v := range esSids {
		if snAuth, ok := authMap[v]; ok {
			if snAuth.CanPlay() {
				authSids = append(authSids, v)
			}
		}
	}
	log.Info("EsPgcIdx EsSids Len %d, Detail %v. AuthSids Len %d,  %v", len(esSids), esSids, len(authSids), authSids)
	if req.IsDefault() { // if it's default condition, we apply the interventions
		reqIdx := &search.ReqIdxInterv{}
		reqIdx.FromPGC(authSids, req)
		if authSids, err = s.esInterv(ctx, reqIdx); err != nil {
			return
		}
	}
	target, _ = s.transformCards(authSids)
	for _, v := range target {
		card := &search.EsCard{}
		card.FromPgc(v)
		idxCards = append(idxCards, card)
	}
	res = &search.EsPager{
		Page:   data.Page,
		Result: idxCards,
		Title:  req.Title(),
	}
	return
}

// EsUgcIdx def.
func (s *Service) EsUgcIdx(ctx context.Context, req *search.ReqUgcIdx) (res *search.EsPager, err error) {
	var (
		srvReq = &search.SrvUgcIdx{
			ReqEsPn: req.ReqEsPn,
		}
		childTps    []*arcwar.Tp
		esRes       *search.EsUgcResult
		aids        []int64
		ugcCardsMap map[int64]*model.ArcCMS
		tp          *arcwar.Tp
	)
	if req.PubTime != "" && req.PubTime != search.AllLabel { // pubtime treatment
		if srvReq.PubTime, err = req.TimeStr(); err != nil {
			log.Error("EsUgcIdx TimeJson Err %v", err)
			return
		}
	}
	if req.SecondTID > 0 { // typeid treatment
		srvReq.TIDs = []int32{req.SecondTID}
	} else {
		if childTps, err = s.arcDao.TypeChildren(req.ParentTID); err != nil {
			log.Error("EsUgcIdx TypeChildren parentTid %d, Err %v", req.ParentTID, err)
			return
		}
		for _, v := range childTps {
			srvReq.TIDs = append(srvReq.TIDs, v.ID)
		}
	}
	if esRes, err = s.searchDao.UgcIdx(ctx, srvReq); err != nil { // pick es result
		log.Error("EsUgcIdx Req %v, Err %v", req, err)
		return
	}
	if tp, err = s.arcDao.TypeInfo(req.ParentTID); err != nil {
		log.Error("EsUgcIdx ParentTID %d, Err %v", req.ParentTID, err)
		return
	}
	res = &search.EsPager{
		Page:  esRes.Page,
		Title: tp.Name,
	}
	if len(esRes.Result) == 0 {
		return
	}
	for _, v := range esRes.Result {
		aids = append(aids, v.AID)
	}
	if req.IsDefault() { // if it's default condition, we apply the interventions
		reqIdx := &search.ReqIdxInterv{}
		reqIdx.FromUGC(aids, req)
		if aids, err = s.esInterv(ctx, reqIdx); err != nil {
			return
		}
	}
	if ugcCardsMap, err = s.cmsDao.LoadArcsMediaMap(ctx, aids); err != nil { // transform UGC
		log.Error("[EsUgcIdx] Can't Pick MediaCache Data, Aids: %v, Err: %v", aids, err)
		return
	}
	for _, v := range aids { // if canPlay, add it into the final result
		if arcCMS, ok := ugcCardsMap[v]; ok {
			if arcCMS.CanPlay() {
				card := &search.EsCard{}
				card.FromUgc(arcCMS.ToCard())
				res.Result = append(res.Result, card)
			}
		}
	}
	return
}
