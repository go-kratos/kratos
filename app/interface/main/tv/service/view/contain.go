package view

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/app/interface/main/tv/model/view"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_rate = map[int]int64{15: 464, 16: 464, 32: 1028, 48: 1328, 64: 2192, 74: 3192, 80: 3192, 112: 6192, 116: 6192, 66: 1820}
)

func initPage(v *arcwar.Page, isBangumi bool) (page *view.Page) {
	page = &view.Page{}
	metas := make([]*view.Meta, 0, 4)
	for q, r := range _rate {
		meta := &view.Meta{
			Quality: q,
			Size:    int64(float64(r*v.Duration) * 1.1 / 8.0),
		}
		metas = append(metas, meta)
	}
	if isBangumi {
		v.From = "bangumi"
	}
	page.Page = v
	page.Metas = metas
	return
}

func (s *Service) initPages(c context.Context, vs *view.Static, ap []*arcwar.Page) (err error) {
	var (
		cids      []int64
		pages     = make([]*view.Page, 0, len(ap))
		vsAuth    map[int64]*model.VideoCMS
		isBangumi = vs.AttrVal(archive.AttrBitIsBangumi) == archive.AttrYes
		emptyArc  = true
	)
	for _, v := range ap {
		cids = append(cids, v.Cid)
	}
	if vsAuth, err = s.cmsDao.LoadVideosMeta(c, cids); err != nil {
		log.Error("initPages LoadVideosMeta Cid %v, Err %v", cids, err)
		return
	}
	for _, v := range ap {
		if auth, ok := vsAuth[v.Cid]; ok { // auditing data can't show
			if !auth.Auditing() {
				pages = append(pages, initPage(v, isBangumi))
			}
			if auth.CanPlay() {
				emptyArc = false
			}
		}
	}
	if emptyArc { // if the arc doesn't have any video that can play, we put its valid field to 0 in an asynchronous manner
		log.Info("emptyArc add Aid %d, Cids %v", vs.Aid, cids)
		s.emptyArcCh <- vs.Aid
	}
	if len(pages) == 0 {
		err = ecode.TvAllDataAuditing
		return
	}
	vs.Pages = pages
	return
}

// initRelates init Relates
func (s *Service) initRelates(c context.Context, v *view.View, ip string, now time.Time) {
	var (
		rls []*view.Relate
		err error
	)
	if rls, err = s.dealRcmdRelate(ctx, v.Aid, ip); err != nil {
		log.Error("initRelates For Aid %d, Error %v", v.Aid, err)
		return
	}
	if len(rls) == 0 {
		s.prom.Incr("zero_relates")
		return
	}
	v.Relates = rls
}

func (s *Service) dealRcmdRelate(c context.Context, aid int64, ip string) (rls []*view.Relate, err error) {
	if rls, err = s.arcDao.RelatesCache(c, aid); err != nil { // mc error
		return
	}
	if len(rls) != 0 {
		s.pHit.Incr("relate_cache")
		return
	}
	var (
		aids     []int64
		as       map[int64]*arcwar.Arc
		arcMetas map[int64]*model.ArcCMS
	)
	s.pMiss.Incr("relate_cache")
	if aids, err = s.arcDao.RelateAids(c, aid, ip); err != nil { // backsource
		return
	}
	if len(aids) == 0 {
		return
	}
	g, errCtx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		as, err = s.arcDao.Archives(errCtx, aids)
		return
	})
	g.Go(func() (err error) {
		arcMetas, err = s.cmsDao.LoadArcsMediaMap(errCtx, aids)
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("dealRcmdRelate For Aid %d, Err %v", aid, err)
		return
	}
	for _, aid := range aids {
		if a, ok := as[aid]; ok {
			// auth, filter can't play ones
			if arcCMS, okCMS := arcMetas[aid]; !okCMS {
				log.Error("LoadArcsMediaMap Missing Aid %d Info", aid)
				continue
			} else if canplay, _ := s.cmsDao.UgcErrMsg(arcCMS.Deleted, arcCMS.Result, arcCMS.Valid); !canplay {
				log.Warn("LoadArcsMediaMap Aid %d Can't play, Struct %v", aid, arcCMS)
				continue
			}
			// can play, init them
			r := &view.Relate{}
			r.FromAv(a, "")
			rls = append(rls, r)
		}
	}
	if len(rls) != 0 {
		s.arcDao.AddRelatesCache(aid, rls)
	}
	return
}
