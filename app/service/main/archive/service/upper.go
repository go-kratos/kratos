package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/time"

	"go-common/library/sync/errgroup"
)

// UpperCount upper count.
func (s *Service) UpperCount(c context.Context, mid int64) (count int, err error) {
	if count, err = s.arc.UpperCountCache(c, mid); err != nil {
		log.Error("s.arc.UpperCountCache(%d) error(%v)", mid, err)
	}
	if count < 0 {
		if count, err = s.arc.UpperCount(c, mid); err != nil {
			log.Error("s.arc.UpperCount(%d) error(%v)", mid, err)
			return
		}
		s.arc.AddUpperCountCache(c, mid, count)
	}
	return
}

// UppersCount uppers count
func (s *Service) UppersCount(c context.Context, mids []int64) (uc map[int64]int, err error) {
	var missed []int64
	uc, missed, _ = s.arc.UppersCountCache(c, mids)
	if len(missed) != 0 {
		var missedCnt map[int64]int
		if missedCnt, err = s.arc.UppersCount(c, missed); err != nil {
			return
		}
		for _, mid := range missed {
			uc[mid] = missedCnt[mid]
			tmid := mid
			tcnt := missedCnt[mid]
			s.addCache(func() {
				s.arc.AddUpperCountCache(context.TODO(), tmid, tcnt)
			})
		}
	}
	return
}

// UppersAidPubTime get aid and pime by mids
func (s *Service) UppersAidPubTime(c context.Context, mids []int64, pn, ps int) (mas map[int64][]*archive.AidPubTime, err error) {
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
		eg        errgroup.Group
		cacheAidM map[int64][]*archive.AidPubTime
		missAidM  = make(map[int64][]*archive.AidPubTime)
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
		if cacheAidM, err = s.arc.UppersPassedCacheWithScore(c, cachedUp, start, end); err != nil {
			log.Error("s.arc.UppersPassedCache(%v) error(%v)", cachedUp, err)
			return
		}
		return
	})
	eg.Go(func() (err error) {
		if len(missed) == 0 {
			return
		}
		var (
			ptimem     map[int64][]time.Time
			missM      map[int64][]int64
			copyrightM map[int64][]int8
		)
		if missM, ptimem, copyrightM, err = s.arc.UppersPassed(c, missed); err != nil {
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
				ccopyright = copyrightM[mid]
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
			for i := start; i < len(missM[mid]); i++ {
				missAidM[mid] = append(missAidM[mid], &archive.AidPubTime{Aid: aids[i], PubDate: time.Time(cptime[i]), Copyright: ccopyright[i]})
			}
		}
		return
	})
	if err = eg.Wait(); err != nil {
		log.Error("eg.Wait error(%v)", err)
		return
	}
	mas = make(map[int64][]*archive.AidPubTime, len(mids))
	for mid, v := range cacheAidM {
		mas[mid] = v
	}
	for mid, v := range missAidM {
		mas[mid] = v
	}
	return
}

// AddUpperPassedCache add up passed archive cache.
func (s *Service) AddUpperPassedCache(c context.Context, aid int64) (err error) {
	var a *api.Arc
	if a, err = s.arc.Archive3(c, aid); err != nil {
		return
	}
	if ok, _ := s.arc.ExpireUpperPassedCache(c, a.Author.Mid); !ok {
		alls, ptimes, copyrights, _ := s.arc.UpperPassed(c, a.Author.Mid)
		if len(alls) != 0 {
			s.arc.AddUpperPassedCache(context.Background(), a.Author.Mid, alls, ptimes, copyrights)
		}
	} else {
		if err = s.arc.AddUpperPassedCache(c, a.Author.Mid, []int64{aid}, []time.Time{a.PubDate}, []int8{int8(a.Copyright)}); err != nil {
			log.Error("s.arc.AddUpperPassedCache(%d,%d) error(%v)", a.Author.Mid, aid, err)
			return
		}
	}
	var count int
	if count, err = s.arc.UpperCount(c, a.Author.Mid); err != nil {
		log.Error("s.arc.UpperCount(%d) error(%v)", a.Author.Mid, err)
		return
	}
	s.arc.AddUpperCountCache(c, a.Author.Mid, count)
	return
}

// DelUpperPassedCache delete up passed archive cache.
func (s *Service) DelUpperPassedCache(c context.Context, aid, mid int64) (err error) {
	var a *api.Arc
	if a, err = s.arc.Archive3(c, aid); err != nil {
		return
	}
	var delMid = a.Author.Mid
	if mid != 0 {
		delMid = mid
	}
	if err = s.arc.DelUpperPassedCache(c, delMid, aid); err != nil {
		log.Error("s.arc.DelUpperPassedCache(%d, %d) error(%v)", delMid, aid)
		return
	}
	var count int
	if count, err = s.arc.UpperCount(c, delMid); err != nil {
		log.Error("s.arc.UpperCount(%d) error(%v)", delMid, err)
		return
	}
	s.arc.AddUpperCountCache(c, delMid, count)
	return
}

// UpperCache update upper cache.
func (s *Service) UpperCache(c context.Context, mid int64, action string) (err error) {
	var count int
	if count, err = s.UpperCount(c, mid); err != nil {
		log.Error("s.UpperCount(%d) error(%v)", mid, err)
		return
	}
	if count <= 0 {
		return
	}
	var as []*api.Arc
	if as, err = s.UpperPassed3(c, mid, 1, 1); err != nil {
		err = nil
		return
	}
	if len(as) == 0 {
		return
	}
	s.addCache(func() {
		s.arc.UpperCache(context.TODO(), mid, action, as[0].Author.Name, as[0].Author.Face)
	})
	return
}
