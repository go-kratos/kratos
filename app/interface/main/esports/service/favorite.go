package service

import (
	"context"
	"time"

	"go-common/app/interface/main/esports/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_firstPs    = 5
	_firstAppPs = 50
	_favDay     = 15
)

var _empStime = make([]string, 0)

// AddFav add favorite contest.
func (s *Service) AddFav(c context.Context, mid, cid int64) (err error) {
	var (
		contest *model.Contest
		mapC    map[int64]*model.Contest
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	if mapC, err = s.dao.EpContests(c, []int64{cid}); err != nil {
		return
	}
	contest = mapC[cid]
	if contest == nil || contest.ID == 0 {
		err = ecode.EsportsContestNotExist
		return
	}
	if contest.LiveRoom <= 0 {
		err = ecode.EsportsContestFavNot
		return
	}
	nowTime := time.Now().Unix()
	if contest.Etime > 0 && nowTime >= contest.Etime {
		err = ecode.EsportsContestEnd
		return
	}
	if contest.Stime == 0 || nowTime >= contest.Stime {
		err = ecode.EsportsContestStart
		return
	}
	subDay := timeSub(contest.Stime)
	if subDay > _favDay {
		err = ecode.EsportsContestNotDay
		return
	}
	arg := &favmdl.ArgAdd{Type: favmdl.TypeEsports, Mid: mid, Oid: cid, Fid: 0, RealIP: ip}
	if err = s.fav.Add(c, arg); err != nil {
		log.Error("AddFav s.fav.Add(%+v) error(%v)", arg, err)
		return
	}
	if err = s.dao.DelFavCoCache(c, mid); err != nil {
		log.Error("AddFav s.dao.DelFavCoCache mid(%d) error(%v)", mid, err)
		return
	}
	return
}

func timeSub(stime int64) int {
	var (
		nowTime, endTime time.Time
	)
	nowTime = time.Now()
	endTime = time.Unix(stime, 0)
	nowTime = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.Local)
	return int(endTime.Sub(nowTime).Hours() / 24)
}

// DelFav delete favorite contest.
func (s *Service) DelFav(c context.Context, mid, cid int64) (err error) {
	var (
		contest *model.Contest
		mapC    map[int64]*model.Contest
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	if mapC, err = s.dao.EpContests(c, []int64{cid}); err != nil {
		return
	}
	contest = mapC[cid]
	if contest == nil || contest.ID == 0 {
		err = ecode.EsportsContestNotExist
		return
	}
	arg := &favmdl.ArgDel{Type: favmdl.TypeEsports, Mid: mid, Oid: cid, Fid: 0, RealIP: ip}
	if err = s.fav.Del(c, arg); err != nil {
		log.Error("DelFav s.fav.Del(%+v) error(%v)", arg, err)
		return
	}
	if err = s.dao.DelFavCoCache(c, mid); err != nil {
		log.Error("DelFav s.dao.DelFavCoCache mid(%d) error(%v)", mid, err)
		return
	}
	return
}

// ListFav list favorite contests.
func (s *Service) ListFav(c context.Context, mid, vmid int64, pn, ps int) (rs []*model.Contest, count int, err error) {
	var (
		isFirst                     bool
		uid                         int64
		favRes                      *favmdl.Favorites
		cids                        []int64
		ip                          = metadata.String(c, metadata.RemoteIP)
		teams, seasons              []*model.Filter
		cData                       map[int64]*model.Contest
		favContest                  []*model.Contest
		group                       *errgroup.Group
		contErr, teamErr, seasonErr error
	)

	if vmid > 0 {
		uid = vmid
	} else {
		uid = mid
	}
	isFirst = pn == 1 && ps == _firstPs
	if isFirst {
		if rs, count, err = s.dao.FavCoCache(c, uid); err != nil {
			err = nil
		}
		if len(rs) > 0 {
			s.fmtContest(c, rs, mid)
			return
		}
	}
	arg := &favmdl.ArgFavs{Type: favmdl.TypeEsports, Mid: mid, Vmid: vmid, Fid: 0, Pn: pn, Ps: ps, RealIP: ip}
	if favRes, err = s.fav.Favorites(c, arg); err != nil {
		log.Error("ListFav s.fav.Favorites(%+v) error(%v)", arg, err)
		return
	}
	count = favRes.Page.Count
	if favRes == nil || len(favRes.List) == 0 || count == 0 {
		rs = _emptContest
		return
	}
	for _, fav := range favRes.List {
		cids = append(cids, fav.Oid)
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if cData, contErr = s.dao.EpContests(c, cids); contErr != nil {
			log.Error("s.dao.Contest error(%v)", contErr)
		}
		return contErr
	})
	group.Go(func() error {
		if teams, teamErr = s.dao.Teams(errCtx); teamErr != nil {
			log.Error("s.dao.Teams error %v", teamErr)
		}
		return nil
	})
	group.Go(func() error {
		if seasons, seasonErr = s.dao.SeasonAll(errCtx); seasonErr != nil {
			log.Error("s.dao.SeasonAll error %v", seasonErr)
		}
		return nil
	})
	err = group.Wait()
	if err != nil {
		return
	}
	for _, fav := range favRes.List {
		if contest, ok := cData[fav.Oid]; ok {
			favContest = append(favContest, contest)
		}
	}
	rs = s.ContestInfo(c, cids, favContest, teams, seasons, mid)
	if isFirst {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetFavCoCache(c, uid, rs, count)
		})
	}
	return
}

// SeasonFav list favorite season.
func (s *Service) SeasonFav(c context.Context, mid int64, p *model.ParamSeason) (rs []*model.Season, count int, err error) {
	var (
		uid        int64
		elaContest []*model.ElaSub
		mapSeasons map[int64]*model.Season
		cids       []int64
		sids       []int64
		dbContests map[int64]*model.Contest
	)
	if p.VMID > 0 {
		uid = p.VMID
	} else {
		uid = mid
	}
	if elaContest, count, err = s.dao.SeasonFav(c, uid, p); err != nil {
		log.Error("s.dao.StimeFav error(%v)", err)
		return
	}
	for _, contest := range elaContest {
		cids = append(cids, contest.Oid)
		sids = append(sids, contest.Sid)
	}
	if len(cids) > 0 {
		if dbContests, err = s.dao.EpContests(c, cids); err != nil {
			log.Error("s.dao.EpContests error(%v)", err)
			return
		}
	} else {
		rs = _emptSeason
		return
	}
	if mapSeasons, err = s.dao.EpSeasons(c, sids); err != nil {
		log.Error("s.dao.EpSeasons error(%v)", err)
		return
	}
	ms := make(map[int64]struct{}, len(cids))
	for _, contest := range elaContest {
		if _, ok := ms[contest.Sid]; ok {
			continue
		}
		// del over contest stime.
		if contest, ok := dbContests[contest.Oid]; ok {
			if contest.Etime > 0 && time.Now().Unix() > contest.Etime {
				continue
			}
		}
		ms[contest.Sid] = struct{}{}
		if season, ok := mapSeasons[contest.Sid]; ok {
			rs = append(rs, season)
		}
	}
	if len(rs) == 0 {
		rs = _emptSeason
	}
	return
}

// StimeFav list favorite contests stime.
func (s *Service) StimeFav(c context.Context, mid int64, p *model.ParamSeason) (rs []string, count int, err error) {
	var (
		uid        int64
		elaContest []*model.ElaSub
		cids       []int64
		dbContests map[int64]*model.Contest
	)
	if p.VMID > 0 {
		uid = p.VMID
	} else {
		uid = mid
	}
	if elaContest, count, err = s.dao.StimeFav(c, uid, p); err != nil {
		log.Error("s.dao.StimeFav error(%v)", err)
	}
	for _, contest := range elaContest {
		cids = append(cids, contest.Oid)
	}
	if len(cids) > 0 {
		if dbContests, err = s.dao.EpContests(c, cids); err != nil {
			log.Error("s.dao.EpContests error(%v)", err)
			return
		}
	} else {
		rs = _empStime
		return
	}
	ms := make(map[string]struct{}, len(cids))
	for _, contest := range elaContest {
		tm := time.Unix(contest.Stime, 0)
		stime := tm.Format("2006-01-02")
		if _, ok := ms[stime]; ok {
			continue
		}
		ms[stime] = struct{}{}
		// del over contest stime.
		if contest, ok := dbContests[contest.Oid]; ok {
			if contest.Etime > 0 && time.Now().Unix() > contest.Etime {
				continue
			}
		}
		rs = append(rs, stime)
	}
	if len(rs) == 0 {
		rs = _empStime
	}
	return
}

// ListAppFav list favorite contests.
func (s *Service) ListAppFav(c context.Context, mid int64, p *model.ParamFav) (rs []*model.Contest, count int, err error) {
	var (
		uid                         int64
		cids                        []int64
		isFirst                     bool
		teams, seasons              []*model.Filter
		cData                       map[int64]*model.Contest
		favContest                  []*model.Contest
		group                       *errgroup.Group
		contErr, teamErr, seasonErr error
	)
	if p.VMID > 0 {
		uid = p.VMID
	} else {
		uid = mid
	}
	isFirst = p.Pn == 1 && p.Ps == _firstAppPs && p.Stime == "" && p.Etime == "" && len(p.Sids) == 0 && p.Sort == 0
	if isFirst {
		if rs, count, err = s.dao.FavCoAppCache(c, uid); err != nil {
			err = nil
		}
		if len(rs) > 0 {
			s.fmtContest(c, rs, uid)
			return
		}
	}
	if cids, count, err = s.dao.SearchFav(c, uid, p); err != nil {
		log.Error("s.dao.SearchFav error(%v)", err)
		return
	}
	if len(cids) == 0 || count == 0 {
		rs = _emptContest
		return
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if cData, contErr = s.dao.EpContests(c, cids); contErr != nil {
			log.Error("s.dao.Contest error(%v)", contErr)
		}
		return contErr
	})
	group.Go(func() error {
		if teams, teamErr = s.dao.Teams(errCtx); teamErr != nil {
			log.Error("s.dao.Teams error %v", teamErr)
		}
		return nil
	})
	group.Go(func() error {
		if seasons, seasonErr = s.dao.SeasonAll(errCtx); seasonErr != nil {
			log.Error("s.dao.SeasonAll error %v", seasonErr)
		}
		return nil
	})
	err = group.Wait()
	if err != nil {
		return
	}
	for _, cid := range cids {
		if contest, ok := cData[cid]; ok {
			favContest = append(favContest, contest)
		}
	}
	rs = s.ContestInfo(c, cids, favContest, teams, seasons, uid)
	if isFirst {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetAppFavCoCache(c, uid, rs, count)
		})
	}
	return
}

func (s *Service) isFavs(c context.Context, mid int64, cids []int64) (res map[int64]bool, err error) {
	if mid > 0 {
		ip := metadata.String(c, metadata.RemoteIP)
		if res, err = s.fav.IsFavs(c, &favmdl.ArgIsFavs{Type: favmdl.TypeEsports, Mid: mid, Oids: cids, RealIP: ip}); err != nil {
			log.Error("s.fav.IsFavs(%d,%+v) error(%d)", mid, cids, err)
			err = nil
		}
	}
	return
}
