package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// DM dm list.
func (s *Service) DM(c context.Context, tp int32, aid, oid int64) (res *model.DMSeg, err error) {
	var (
		total, size int64 = 1, model.DefaultVideoEnd
		mu          sync.Mutex
	)
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	res = &model.DMSeg{Elems: make([]*model.Elem, 0, 2*sub.Maxlimit)}
	duration, err := s.videoDuration(c, aid, oid)
	if err != nil {
		return
	}
	if duration != 0 {
		total = int64(math.Ceil(float64(duration) / float64(model.DefaultPageSize)))
		size = model.DefaultPageSize
	}
	g, ctx := errgroup.WithContext(c)
	for i := int64(1); i <= total; i++ {
		num := i
		g.Go(func() (err error) {
			var dmseg *model.DMSeg
			if dmseg, err = s.dao.DMSegCache(ctx, tp, oid, total, num); err != nil {
				return
			}
			if dmseg == nil {
				ps := (num - 1) * size
				pe := num * size
				fmt.Println(ps, pe, total, num)
				if dmseg, err = s.dmSegV2(ctx, sub, total, num, ps, pe); err != nil {
					return
				}
			}
			if dmseg != nil {
				mu.Lock()
				res.Elems = append(res.Elems, dmseg.Elems...)
				mu.Unlock()
			}
			return
		})
	}
	err = g.Wait()
	return
}

// DMSegV2 dm segment new.
func (s *Service) DMSegV2(c context.Context, tp int32, mid, aid, oid, pn int64, plat int32) (res *model.DMSegResp, err error) {
	page, err := s.pageinfo(c, tp, aid, oid, pn)
	if err != nil {
		return
	}
	ps := (page.Num - 1) * page.Size
	pe := page.Num * page.Size
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	res = &model.DMSegResp{
		Flag: model.DefaultFlag,
	}
	if sub.State == model.SubStateClosed {
		return
	}
	flag, err := s.dao.RecFlag(c, mid, aid, oid, 2*sub.Maxlimit, ps, pe, plat)
	if err == nil {
		res.Flag = flag
	}
	dmseg, err := s.dao.DMSegCache(c, tp, oid, page.Total, page.Num)
	if err != nil {
		return
	}
	if dmseg != nil {
		res.Dms = dmseg.Elems
		return
	}
	if dmseg, err = s.dmSegV2(c, sub, page.Total, page.Num, ps, pe); err != nil {
		return
	}
	res.Dms = dmseg.Elems
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.SetDMSegCache(ctx, tp, oid, page.Total, page.Num, dmseg) // add mc cache
	})
	return
}

func (s *Service) dmSegV2(c context.Context, sub *model.Subject, total, num, ps, pe int64) (res *model.DMSeg, err error) {
	var (
		cache = true
		limit = 2 * sub.Maxlimit
		dmids = make([]int64, 0, limit)
	)
	res = &model.DMSeg{Elems: make([]*model.Elem, 0, limit)}
	normalIds, err := s.dmNormalIds(c, sub.Type, sub.Oid, total, num, ps, pe, limit)
	if err != nil {
		return
	}
	dmids = append(dmids, normalIds...)
	if sub.Childpool > 0 {
		var subtitleIds []int64
		if subtitleIds, err = s.dmSegSubtitlesIds(c, sub.Type, sub.Oid, ps, pe, limit); err != nil {
			return
		}
		dmids = append(dmids, subtitleIds...)
	}
	if len(dmids) <= 0 {
		return
	}
	elemsCache, missed, err := s.dao.IdxContentCacheV2(c, sub.Type, sub.Oid, dmids)
	if err != nil {
		missed = dmids
		cache = false
	} else {
		res.Elems = append(res.Elems, elemsCache...)
	}
	if len(missed) == 0 {
		return
	}
	dms, err := s.dmsSeg(c, sub.Type, sub.Oid, missed)
	if err != nil {
		return
	}
	for _, dm := range dms {
		if e := dm.ToElem(); e != nil {
			res.Elems = append(res.Elems, e)
		}
	}
	if cache && len(dms) > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddIdxContentCache(ctx, sub.Type, sub.Oid, dms, false) // add memcache,realname=false
		})
	}
	return
}

// pageinfo get page info of oid.
func (s *Service) pageinfo(c context.Context, tp int32, aid, oid, pn int64) (p *model.Page, err error) {
	var duration int64
	data, ok := s.localCache[keyDuration(tp, oid)]
	if ok {
		duration, err = strconv.ParseInt(string(data), 10, 64)
	} else {
		duration, err = s.videoDuration(c, aid, oid)
	}
	if err != nil {
		return
	}
	if duration == 0 {
		p = &model.Page{
			Num:   pn,
			Size:  model.DefaultVideoEnd,
			Total: 1,
		}
	} else {
		p = &model.Page{
			Num:   pn,
			Size:  model.DefaultPageSize,
			Total: int64(math.Ceil(float64(duration) / float64(model.DefaultPageSize))),
		}
	}
	if pn > p.Total {
		log.Warn("oid:%d pn:%d larger than total page:%d", oid, pn, p.Total)
		err = ecode.NotModified
		return
	}
	return
}
