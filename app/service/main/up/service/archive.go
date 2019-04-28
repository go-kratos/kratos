package service

import (
	"context"

	arcgrpc "go-common/app/service/main/archive/api"
	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"

	"go-common/library/sync/errgroup"
)

// UpArcs upper passed.
func (s *Service) UpArcs(c context.Context, arg *upgrpc.UpArcsReq) (res *upgrpc.UpArcsReply, err error) {
	var (
		upcnt *upgrpc.UpCountReply
		mid   = arg.Mid
		ps    = int(arg.Ps)
		pn    = int(arg.Pn)
	)
	res = new(upgrpc.UpArcsReply)
	if upcnt, err = s.UpCount(c, &upgrpc.UpCountReq{Mid: mid}); err != nil {
		return
	}
	if upcnt.Count == 0 {
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
			return
		}
		if err = s.buildUpperStaff(c, mid, &alls, &ptimes, &copyrights); err != nil {
			return
		}
		length := len(alls)
		if length == 0 || length < start {
			err = ecode.NothingFound
			return
		}
		s.cacheWorker.Do(c, func(c context.Context) {
			s.arc.AddUpperPassedCache(c, mid, alls, ptimes, copyrights)
		})
		if length > end+1 {
			aids = alls[start : end+1]
		} else {
			aids = alls[start:]
		}
		s.pCacheMiss.Add("up_arcs_mid", 1)
	} else {
		if aids, err = s.arc.UpperPassedCache(c, mid, start, end); err != nil {
			return
		}
		s.pCacheHit.Add("up_arcs_mid", 1)
	}
	if len(aids) == 0 {
		return
	}
	var am *arcgrpc.ArcsReply
	if am, err = global.GetArcClient().Arcs(c, &arcgrpc.ArcsRequest{Aids: aids}); err != nil {
		return
	}
	for _, aid := range aids {
		if a, ok := am.Arcs[aid]; ok {
			res.Archives = append(res.Archives, a)
		}
	}
	return
}

// UpsArcs get archives by mids.
func (s *Service) UpsArcs(c context.Context, arg *upgrpc.UpsArcsReq) (res *upgrpc.UpsArcsReply, err error) {
	var (
		cachedUp []int64
		missed   []int64
		mids     = arg.Mids
		ps       = int(arg.Ps)
		pn       = int(arg.Pn)
	)
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 20
	}
	res = new(upgrpc.UpsArcsReply)
	if cachedUp, missed, err = s.arc.ExpireUppersCountCache(c, mids); err != nil {
		return
	}
	s.pCacheMiss.Add("up_arcs_mid", int64(len(missed)))
	s.pCacheHit.Add("up_arcs_mid", int64(len(cachedUp)))
	var (
		start     = (pn - 1) * ps
		end       = start + ps - 1
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
			log.Error("s.arc.UppersPassedCache(%+v) error(%+v)", cachedUp, err)
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
			ptimeM     map[int64][]time.Time // <mid,ptime-list>
			copyrightM map[int64][]int8      // <mid,copyright-list>
		)
		if missM, ptimeM, copyrightM, err = s.arc.UppersPassed(c, missed); err != nil {
			log.Error("s.arc.UppersPassed(%+v) error(%+v)", missed, err)
			return
		}
		if err = s.buildUppersStaffs(c, missed, missM, ptimeM, copyrightM); err != nil {
			log.Error("s.arc.UppersPassed(%+v,%+v,%+v,%+v) error(%+v)", missed, missM, ptimeM, copyrightM, err)
			return
		}
		for mid, aids := range missM {
			length := len(aids)
			if length == 0 || length < start {
				continue
			}
			var (
				cmid       = mid
				clength    = int64(length)
				cptime     = ptimeM[mid]
				ccopyright = copyrightM[mid]
				caids      = aids
			)
			s.cacheWorker.Do(c, func(c context.Context) {
				s.arc.AddUpperCountCache(c, cmid, clength)
				s.arc.AddUpperPassedCache(c, cmid, caids, cptime, ccopyright)
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
				s.cacheWorker.Do(c, func(c context.Context) {
					s.arc.AddUpperCountCache(c, cmid, 0)
				})
			}
		}
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}
	var (
		allAids = append(cacheAids, missAids...)
		am      *arcgrpc.ArcsReply
	)
	if len(allAids) == 0 {
		return
	}
	if am, err = global.GetArcClient().Arcs(c, &arcgrpc.ArcsRequest{Aids: allAids}); err != nil {
		return
	}
	res.Archives = make(map[int64]*upgrpc.UpArcsReply, len(am.Arcs))
	s.buildUpArcsReply(cacheM, am, res)
	s.buildUpArcsReply(missM, am, res)
	return
}

func (s *Service) buildUpArcsReply(ms map[int64][]int64, am *arcgrpc.ArcsReply, ups *upgrpc.UpsArcsReply) {
	for mid, aids := range ms {
		arcs := &upgrpc.UpArcsReply{
			Archives: make([]*arcgrpc.Arc, 0, len(ms)),
		}
		for _, aid := range aids {
			if a, ok := am.Arcs[aid]; ok {
				arcs.Archives = append(arcs.Archives, a)
			}
		}
		ups.Archives[mid] = arcs
	}
}

// UpCount upper count.
func (s *Service) UpCount(c context.Context, arg *upgrpc.UpCountReq) (res *upgrpc.UpCountReply, err error) {
	var mid = arg.Mid
	res = new(upgrpc.UpCountReply)
	if res.Count, err = s.arc.UpperCountCache(c, mid); err != nil {
		return
	}
	if res.Count > 0 {
		return
	}
	if res.Count, err = s.buildMidCount(c, mid); err != nil {
		return
	}
	s.arc.AddUpperCountCache(c, mid, res.Count)
	return
}

// UpsCount uppers count
func (s *Service) UpsCount(c context.Context, arg *upgrpc.UpsCountReq) (res *upgrpc.UpsCountReply, err error) {
	var (
		missed []int64
		mids   = arg.Mids
	)
	res = new(upgrpc.UpsCountReply)
	res.Count = make(map[int64]int64, len(mids))
	res.Count, missed, _ = s.arc.UppersCountCache(c, mids)
	if len(missed) == 0 {
		return
	}
	var missedCnt map[int64]int64 // <mid,count>
	if missedCnt, err = s.buildMidCounts(c, missed); err != nil {
		return
	}
	for _, mid := range missed {
		res.Count[mid] = missedCnt[mid]
		tmid := mid
		tcnt := missedCnt[mid]
		s.cacheWorker.Do(c, func(c context.Context) {
			s.arc.AddUpperCountCache(c, tmid, tcnt)
		})
	}
	return
}

// UpsAidPubTime get aid and pime by mids
func (s *Service) UpsAidPubTime(c context.Context, arg *upgrpc.UpsArcsReq) (res *upgrpc.UpsAidPubTimeReply, err error) {
	var (
		cachedUp []int64
		missed   []int64
		mids     = arg.Mids
		ps       = int(arg.Ps)
		pn       = int(arg.Pn)
	)
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 20
	}
	res = new(upgrpc.UpsAidPubTimeReply)
	if cachedUp, missed, err = s.arc.ExpireUppersCountCache(c, mids); err != nil {
		return
	}
	s.pCacheMiss.Add("up_arcs_mid", int64(len(missed)))
	s.pCacheHit.Add("up_arcs_mid", int64(len(cachedUp)))
	var (
		eg        errgroup.Group
		cacheAidM map[int64][]*model.AidPubTime
		missAidM  = make(map[int64][]*model.AidPubTime)
		start     = (pn - 1) * ps
		end       = start + ps - 1
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
			ptimeM     map[int64][]time.Time // <mid,pubdate-list>
			missM      map[int64][]int64     // <mid,aid-list>
			copyrightM map[int64][]int8      // <mid,copyright-list>
		)
		if missM, ptimeM, copyrightM, err = s.arc.UppersPassed(c, missed); err != nil {
			log.Error("s.arc.UppersPassed(%v) error(%v)", missed, err)
			return
		}
		if err = s.buildUppersStaffs(c, missed, missM, ptimeM, copyrightM); err != nil {
			log.Error("s.arc.UppersPassed(%+v,%+v,%+v,%+v) error(%+v)", missed, missM, ptimeM, copyrightM, err)
			return
		}
		for mid, aids := range missM {
			length := len(aids)
			if length == 0 || length < start {
				continue
			}
			var (
				cmid       = mid
				clength    = int64(length)
				cptime     = ptimeM[mid]
				ccopyright = copyrightM[mid]
				caids      = aids
			)
			s.cacheWorker.Do(c, func(c context.Context) {
				s.arc.AddUpperCountCache(c, cmid, clength)
				s.arc.AddUpperPassedCache(c, cmid, caids, cptime, ccopyright)
			})
			if len(aids) > end+1 {
				missM[mid] = aids[start : end+1]
			} else {
				missM[mid] = aids[start:]
			}
			for i := start; i < len(missM[mid]); i++ {
				missAidM[mid] = append(missAidM[mid], &model.AidPubTime{Aid: aids[i], PubDate: time.Time(cptime[i]), Copyright: ccopyright[i]})
			}
		}
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}
	res.Archives = make(map[int64]*upgrpc.UpAidPubTimeReply, len(mids))
	s.buildUpAidPubTimeReply(cacheAidM, res)
	s.buildUpAidPubTimeReply(missAidM, res)
	return
}

func (s *Service) buildUpAidPubTimeReply(ms map[int64][]*model.AidPubTime, ups *upgrpc.UpsAidPubTimeReply) {
	for mid, aps := range ms {
		arcs := &upgrpc.UpAidPubTimeReply{
			Archives: make([]*upgrpc.AidPubTime, 0, len(ms)),
		}
		for _, v := range aps {
			arcs.Archives = append(arcs.Archives, &upgrpc.AidPubTime{Aid: v.Aid, PubDate: v.PubDate, Copyright: int32(v.Copyright)})
		}
		ups.Archives[mid] = arcs
	}
}

// AddUpPassedCacheByStaff add up passed archive cache by staff.
func (s *Service) AddUpPassedCacheByStaff(c context.Context, arg *upgrpc.UpCacheReq) (res *upgrpc.NoReply, err error) {
	var (
		count int64
		a     *arcgrpc.ArcReply
	)
	res = new(upgrpc.NoReply)
	if a, err = global.GetArcClient().Arc(c, &arcgrpc.ArcRequest{Aid: arg.Aid}); err != nil {
		return
	}
	if ok, _ := s.arc.ExpireUpperPassedCache(c, arg.Mid); !ok {
		var (
			alls       []int64
			ptimes     []time.Time
			copyrights []int8
		)
		if alls, ptimes, copyrights, err = s.arc.UpperPassed(c, arg.Mid); err != nil {
			return
		}
		if err = s.buildUpperStaff(c, arg.Mid, &alls, &ptimes, &copyrights); err != nil {
			return
		}
		if len(alls) != 0 {
			s.arc.AddUpperPassedCache(context.Background(), arg.Mid, alls, ptimes, copyrights)
		}
		count = int64(len(alls))
	} else {
		if err = s.arc.AddUpperPassedCache(c, arg.Mid, []int64{a.Arc.Aid}, []time.Time{a.Arc.PubDate}, []int8{int8(a.Arc.Copyright)}); err != nil {
			return
		}
		if count, err = s.buildMidCount(c, arg.Mid); err != nil {
			return
		}
	}
	s.arc.AddUpperCountCache(c, arg.Mid, count)
	return
}

// AddUpPassedCache add up passed archive cache.
func (s *Service) AddUpPassedCache(c context.Context, arg *upgrpc.UpCacheReq) (res *upgrpc.NoReply, err error) {
	var (
		midsCnt                   map[int64]int64 // <mid,count>
		expireOKs                 map[int64]bool  // <mid,bool>
		a                         *arcgrpc.ArcReply
		mids, cacheMids, missMids []int64
	)
	res = new(upgrpc.NoReply)
	if a, err = global.GetArcClient().Arc(c, &arcgrpc.ArcRequest{Aid: arg.Aid}); err != nil {
		return
	}
	if mids, err = s.arc.StaffAid(c, arg.Aid); err != nil {
		return
	}
	mids = append(mids, arg.Mid)
	expireOKs, _ = s.arc.ExpireUppersPassedCache(c, mids)
	for mid, ok := range expireOKs {
		if ok {
			cacheMids = append(cacheMids, mid)
		} else {
			missMids = append(missMids, mid)
		}
	}
	midsCnt = make(map[int64]int64, len(mids))
	if len(missMids) != 0 {
		var (
			ptimeM     map[int64][]time.Time // <mid,pubdate-list>
			missM      map[int64][]int64     // <mid,aid-list>
			copyrightM map[int64][]int8      // <mid,copyright-list>
		)
		if missM, ptimeM, copyrightM, err = s.arc.UppersPassed(c, missMids); err != nil {
			return
		}
		if err = s.buildUppersStaffs(c, missMids, missM, ptimeM, copyrightM); err != nil {
			return
		}
		for mid, aids := range missM {
			var (
				cmid       = mid
				clength    = int64(len(aids))
				caids      = aids
				cptime     = ptimeM[mid]
				ccopyright = copyrightM[mid]
			)
			midsCnt[mid] = clength
			s.cacheWorker.Do(c, func(c context.Context) {
				s.arc.AddUpperPassedCache(c, cmid, caids, cptime, ccopyright)
			})
		}
	}
	if len(cacheMids) != 0 {
		var cntM map[int64]int64
		for _, mid := range cacheMids {
			var cmid = mid
			s.cacheWorker.Do(c, func(c context.Context) {
				if err = s.arc.AddUpperPassedCache(c, cmid, []int64{arg.Aid}, []time.Time{a.Arc.PubDate}, []int8{int8(a.Arc.Copyright)}); err != nil {
					return
				}
			})
		}
		if cntM, err = s.buildMidCounts(c, mids); err != nil {
			return
		}
		for mid, count := range cntM {
			midsCnt[mid] = count
		}
	}
	for mid, count := range midsCnt {
		var (
			cmid   = mid
			ccount = count
		)
		s.cacheWorker.Do(c, func(c context.Context) {
			s.arc.AddUpperCountCache(c, cmid, ccount)
		})
	}
	return
}

// DelUpPassedCacheByStaff delete up passed archive cache by staff.
func (s *Service) DelUpPassedCacheByStaff(c context.Context, arg *upgrpc.UpCacheReq) (res *upgrpc.NoReply, err error) {
	var (
		mid = arg.Mid
		aid = arg.Aid
	)
	res = new(upgrpc.NoReply)
	if err = s.arc.DelUpperPassedCache(c, mid, aid); err != nil {
		return
	}
	var count int64
	if count, err = s.buildMidCount(c, mid); err != nil {
		return
	}
	s.arc.AddUpperCountCache(c, mid, count)
	return
}

// DelUpPassedCache delete up passed archive cache.
func (s *Service) DelUpPassedCache(c context.Context, arg *upgrpc.UpCacheReq) (res *upgrpc.NoReply, err error) {
	var (
		mids    []int64
		midsCnt map[int64]int64
	)
	res = new(upgrpc.NoReply)
	if mids, err = s.arc.StaffAid(c, arg.Aid); err != nil {
		return
	}
	mids = append(mids, arg.Mid)
	for _, mid := range mids {
		var cmid = mid
		s.cacheWorker.Do(c, func(c context.Context) {
			if err = s.arc.DelUpperPassedCache(c, cmid, arg.Aid); err != nil {
				return
			}
		})
	}
	if midsCnt, err = s.buildMidCounts(c, mids); err != nil {
		return
	}
	for mid, count := range midsCnt {
		var (
			cmid   = mid
			ccount = count
		)
		s.cacheWorker.Do(c, func(c context.Context) {
			s.arc.AddUpperCountCache(c, cmid, ccount)
		})
	}
	return
}

func (s *Service) buildUpperStaff(c context.Context, mid int64, aids *[]int64, ptimes *[]time.Time, copyrights *[]int8) (err error) {
	var (
		staffAlls, staffAids []int64
		staffPtimes          []time.Time
		staffCopyrights      []int8
	)
	if staffAlls, err = s.arc.Staff(c, mid); err != nil {
		return
	}
	if len(staffAlls) == 0 {
		return
	}
	if staffAids, staffPtimes, staffCopyrights, _, err = s.arc.ArcsAids(c, staffAlls); err != nil {
		return
	}
	if len(staffAids) == 0 || len(staffPtimes) == 0 || len(staffCopyrights) == 0 {
		return
	}
	*aids = append(*aids, staffAids...)
	*ptimes = append(*ptimes, staffPtimes...)
	*copyrights = append(*copyrights, staffCopyrights...)
	return
}

func (s *Service) buildUppersStaffs(c context.Context, mids []int64, midM map[int64][]int64, ptimeM map[int64][]time.Time, copyrightM map[int64][]int8) (err error) {
	var (
		staffAids      []int64
		staffAidPtimeM map[int64]*model.AidPubTime // <aid,model.AidPubTime>
		staffAidM      map[int64][]int64           // <mid,aid-list>
	)
	if staffAidM, err = s.arc.Staffs(c, mids); err != nil {
		return
	}
	if len(staffAidM) == 0 || staffAidM == nil {
		return
	}
	for _, aids := range staffAidM {
		staffAids = append(staffAids, aids...)
	}
	if len(staffAids) == 0 {
		return
	}
	if _, _, _, staffAidPtimeM, err = s.arc.ArcsAids(c, staffAids); err != nil {
		return
	}
	for mid, aids := range staffAidM {
		for _, aid := range aids {
			if apt, ok := staffAidPtimeM[aid]; ok {
				ptimeM[mid] = append(ptimeM[mid], apt.PubDate)
				copyrightM[mid] = append(copyrightM[mid], apt.Copyright)
				midM[mid] = append(midM[mid], aid)
			}
		}
	}
	return
}

func (s *Service) buildMidCount(c context.Context, mid int64) (count int64, err error) {
	if count, err = s.arc.UpperCount(c, mid); err != nil {
		return
	}
	var aids, staffAids []int64
	if aids, err = s.arc.Staff(c, mid); err != nil {
		return
	}
	if len(aids) == 0 {
		return
	}
	if staffAids, _, _, _, err = s.arc.ArcsAids(c, aids); err != nil {
		return
	}
	count += int64(len(staffAids))
	return
}

func (s *Service) buildMidCounts(c context.Context, mids []int64) (midsCnt map[int64]int64, err error) {
	if len(mids) == 0 {
		return
	}
	var (
		staffAlls, staffAids []int64
		aidM                 map[int64]struct{} // <mid,struct>
		staffAidM            map[int64][]int64  // <mid,aid-list>
	)
	midsCnt = make(map[int64]int64, len(mids))
	defer func() {
		if err == nil {
			for _, mid := range mids {
				if _, ok := midsCnt[mid]; !ok {
					midsCnt[mid] = 0
				}
			}
		}
	}()
	if midsCnt, err = s.arc.UppersCount(c, mids); err != nil {
		return
	}
	if staffAidM, err = s.arc.Staffs(c, mids); err != nil {
		return
	}
	for _, aids := range staffAidM {
		staffAlls = append(staffAlls, aids...)
	}
	if len(staffAlls) == 0 {
		return
	}
	if staffAids, _, _, _, err = s.arc.ArcsAids(c, staffAlls); err != nil {
		return
	}
	aidM = make(map[int64]struct{}, len(staffAids))
	for _, aid := range staffAids {
		aidM[aid] = struct{}{}
	}
	for mid, aids := range staffAidM {
		for _, aid := range aids {
			if _, ok := aidM[aid]; ok {
				midsCnt[mid]++
			}
		}
	}

	return
}
