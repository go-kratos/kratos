package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"

	"go-common/library/sync/errgroup"
)

// UpperReommend from archive_recommend by aid
func (s *Service) UpperReommend(c context.Context, aid int64) (as []*api.Arc, err error) {
	var reAids []int64
	if reAids, err = s.arc.UpperReommend(c, aid); err != nil {
		return
	}
	if len(reAids) == 0 {
		err = ecode.NothingFound
		return
	}
	var am map[int64]*api.Arc
	if am, err = s.arc.Archives3(c, reAids); err != nil {
		return
	}
	for _, rid := range reAids {
		if a, ok := am[rid]; ok {
			as = append(as, a)
		}
	}
	return
}

// UpperPassed3 upper passed.
func (s *Service) UpperPassed3(c context.Context, mid int64, pn, ps int) (as []*api.Arc, err error) {
	var cnt int
	if cnt, err = s.UpperCount(c, mid); err != nil {
		log.Error("s.UpperCount(%d) error(%v)", mid, err)
		return
	}
	if cnt == 0 {
		err = ecode.NothingFound
		return
	}
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 20
	}
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1
		aids  []int64
	)
	if ok, _ := s.arc.ExpireUpperPassedCache(c, mid); !ok {
		var (
			alls       []int64
			ptimes     []time.Time
			copyrights []int8
		)
		if alls, ptimes, copyrights, err = s.arc.UpperPassed(c, mid); err != nil {
			log.Error("s.arc.UpperPassed(%d) error(%v)", mid, err)
			return
		}
		length := len(alls)
		if length == 0 || length < start {
			err = ecode.NothingFound
			return
		}
		s.addCache(func() {
			s.arc.AddUpperPassedCache(context.TODO(), mid, alls, ptimes, copyrights)
		})
		if length > end+1 {
			aids = alls[start : end+1]
		} else {
			aids = alls[start:]
		}
		s.missProm.Add("mid", 1)
	} else {
		if aids, err = s.arc.UpperPassedCache(c, mid, start, end); err != nil {
			log.Error("s.arc.UpperPassedCache(%d) error(%v)", mid, err)
			return
		}
		s.hitProm.Add("mid", 1)
	}
	if len(aids) == 0 {
		return
	}
	var am map[int64]*api.Arc
	if am, err = s.arc.Archives3(c, aids); err != nil {
		log.Error("s.arc.Archives(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(am))
	for _, aid := range aids {
		if a, ok := am[aid]; ok {
			as = append(as, a)
		}
	}
	return
}

// UppersPassed3 get archives by mids.
func (s *Service) UppersPassed3(c context.Context, mids []int64, pn, ps int) (mas map[int64][]*api.Arc, err error) {
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 20
	}
	var (
		cachedUp []int64
		missed   []int64
		start    = (pn - 1) * ps
		end      = start + ps - 1
	)
	if cachedUp, missed, err = s.arc.ExpireUppersCountCache(c, mids); err != nil {
		return
	}
	s.missProm.Add("mid", int64(len(missed)))
	s.hitProm.Add("mid", int64(len(cachedUp)))
	var (
		cacheAids []int64
		cacheM    map[int64][]int64
		missAids  []int64
		missM     map[int64][]int64
		eg        errgroup.Group
	)
	eg.Go(func() (err error) {
		if len(cachedUp) == 0 {
			return
		}
		_, err = s.arc.ExpireUppersPassedCache(c, cachedUp)
		return
	})
	eg.Go(func() (err error) {
		if len(cachedUp) == 0 {
			return
		}
		if cacheM, err = s.arc.UppersPassedCache(c, cachedUp, start, end); err != nil {
			log.Error("s.arc.UppersPassedCache(%v) error(%v)", cachedUp, err)
			return
		}
		for _, aids := range cacheM {
			if len(aids) > 0 {
				cacheAids = append(cacheAids, aids...)
			}
		}
		return
	})
	eg.Go(func() (err error) {
		if len(missed) == 0 {
			return
		}
		var (
			ptimem     map[int64][]time.Time
			copyrightm map[int64][]int8
		)
		if missM, ptimem, copyrightm, err = s.arc.UppersPassed(c, missed); err != nil {
			log.Error("s.arc.UppersPassed(%v) error(%v)", missed, err)
			return
		}
		for mid, aids := range missM {
			length := len(aids)
			if length == 0 || length < start {
				continue
			}
			var (
				cmid       = mid
				clength    = length
				cptime     = ptimem[mid]
				ccopyright = copyrightm[mid]
				caids      = aids
			)
			s.addCache(func() {
				s.arc.AddUpperCountCache(context.TODO(), cmid, clength)
				s.arc.AddUpperPassedCache(context.TODO(), cmid, caids, cptime, ccopyright)
			})
			if len(aids) > end+1 {
				missM[mid] = aids[start : end+1]
			} else {
				missM[mid] = aids[start:]
			}
			missAids = append(missAids, missM[mid]...)
		}
		for _, mid := range missed {
			if _, ok := missM[mid]; !ok {
				var cmid = mid
				s.addCache(func() {
					s.arc.AddUpperCountCache(context.TODO(), cmid, 0)
				})
			}
		}
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
		return
	}
	var (
		allAids = append(cacheAids, missAids...)
		am      map[int64]*api.Arc
	)
	if len(allAids) == 0 {
		return
	}
	if am, err = s.arc.Archives3(c, allAids); err != nil {
		log.Error("s.arc.Archives3(%v) error(%v)", allAids, err)
		return
	}
	mas = make(map[int64][]*api.Arc, len(am))
	for mid, aids := range cacheM {
		for _, aid := range aids {
			if a, ok := am[aid]; ok {
				mas[mid] = append(mas[mid], a)
			}
		}
	}
	for mid, aids := range missM {
		for _, aid := range aids {
			if a, ok := am[aid]; ok {
				mas[mid] = append(mas[mid], a)
			}
		}
	}
	return
}
