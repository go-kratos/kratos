package thirdp

import (
	"math"

	"go-common/app/interface/main/tv/dao/thirdp"
	"go-common/app/interface/main/tv/model"
	tpMdl "go-common/app/interface/main/tv/model/thirdp"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func (s *Service) buildPager(req *tpMdl.ReqDBeiPages) (pager *model.IdxPager, err error) {
	var (
		addCache bool
		count    int
	)
	// pick up the count from redis, otherwise pick it from DB
	if count, err = s.dao.GetThirdpCnt(ctx, req.TypeC); err != nil {
		if count, err = s.dao.ThirdpCnt(ctx, req.TypeC); err != nil {
			return // if db count still error, fatal error
		}
		log.Error("PickDBeiPage - Can't Get Count, Pass by DB, Page %d", req.Page)
		addCache = true
	}
	pager = &model.IdxPager{
		CurrentPage: int(req.Page),
		TotalItems:  count,
		TotalPages:  int(math.Ceil(float64(count) / float64(req.Ps))),
		PageSize:    int(req.Ps),
	}
	if req.Page > int64(pager.TotalPages) {
		err = ecode.TvDangbeiPageNotExist
		return
	}
	// async Reset the DB data: Count & CurrentPage ID in MC for next time
	if addCache {
		cache.Save(func() {
			s.dao.SetThirdpCnt(ctx, count, req.TypeC)
			log.Info("PickDBeiPage Set Count %d Into Cache", count)
		})
	}
	return
}

// PickDBeiPage picks the dangbei's page
func (s *Service) PickDBeiPage(page int64, typeC string) (data *tpMdl.DBeiPage, err error) {
	var (
		cPageID int64
		sids    []int64 // this page's season ids
		sns     []*model.SeasonCMS
		arcs    []*model.ArcCMS
		dbeiSns []*tpMdl.DBeiSeason
		pager   *model.IdxPager
	)
	req := &tpMdl.ReqDBeiPages{
		Ps:    s.conf.Cfg.Dangbei.Pagesize,
		TypeC: typeC,
		Page:  page,
	}
	if pager, err = s.buildPager(req); err != nil {
		return
	}
	if req.LastID, err = s.dao.LoadPageID(ctx, req); err != nil {
		log.Error("MangoPage getPageID LastPage %d Miss, Pass by offset", page-1)
		return
	}
	if sids, cPageID, err = s.dao.DBeiPages(ctx, req); err != nil {
		return
	}
	if len(sids) == 0 {
		err = errors.Wrapf(ecode.NothingFound, "Type_C [%s], Page [%d], Offset Result Empty", typeC, page)
		return
	}
	// load data from cache and transform data to Dangbei structure
	if typeC == thirdp.DBeiPGC { // pgc - seasonCMS
		if sns, _, err = s.cmsDao.LoadSnsCMS(ctx, sids); err != nil {
			log.Error("PickDBeiPage - PGC - LoadSnsCMS - Sids %v, Error %v", sids, err)
			return
		}
		for _, v := range sns { // transform the object to DbeiSeason
			dbeiSns = append(dbeiSns, tpMdl.DBeiSn(v))
		}
	} else if typeC == thirdp.DBeiUGC { // ugc - arcCMS
		if arcs, err = s.cmsDao.LoadArcsMedia(ctx, sids); err != nil {
			log.Error("PickDBeiPage - UGC - LoadArcsMedia - Sids %v, Error %v", sids, err)
			return
		}
		for _, v := range arcs {
			first, second := s.arcDao.GetPTypeName(int32(v.TypeID))
			dbeiSns = append(dbeiSns, tpMdl.DbeiArc(v, first, second))
		}
	}
	// async Reset the DB data: CurrentPage ID in MC for next time
	cache.Save(func() {
		s.dao.SetPageID(ctx, &tpMdl.ReqPageID{Page: page, ID: cPageID, TypeC: typeC})
	})
	data = &tpMdl.DBeiPage{
		List:  dbeiSns,
		Pager: pager,
	}
	return
}
