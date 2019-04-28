package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/stat/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixStatPB  = "stp_"
	_prefixClickPB = "clkp_"
)

func statPBKey(aid int64) string {
	return _prefixStatPB + strconv.FormatInt(aid, 10)
}

func clickPBKey(aid int64) string {
	return _prefixClickPB + strconv.FormatInt(aid, 10)
}

func (s *Service) statDealproc(i int64) {
	defer s.waiter.Done()
	var (
		ch  = s.subStatCh[i]
		sm  = s.statSM[i]
		c   = context.TODO()
		ls  *lastTmStat
		err error
	)
	for {
		now := time.Now().Unix()
		ms, ok := <-ch
		if !ok {
			s.multiUpdateDB(i, sm)
			log.Warn("statDealproc(%d) quit", i)
			return
		}
		if s.maxAid > 0 && s.maxAid+300 < ms.Aid {
			log.Warn("aid(%d) too big maxAid(%d)", ms.Aid, s.maxAid)
			continue
		}
		// get stat
		if ls, ok = sm[ms.Aid]; !ok {
			var stat *api.Stat
			if stat, err = s.dao.Stat(c, ms.Aid); err != nil {
				log.Error("s.dao.Stat(%d) error(%v)", ms.Aid, err)
				continue
			}
			ls = &lastTmStat{}
			if stat == nil {
				ls.stat = &api.Stat{Aid: ms.Aid}
				ls.last = 0 // NOTE: make sure update db in first.
			} else {
				ls.stat = stat
				ls.last = time.Now().Unix()
			}
			sm[ms.Aid] = ls
		}
		model.Merge(ms, ls.stat)
		if now-ms.Ts < 60 {
			// update cache
			s.updateCache(ls.stat)
		}
		// update db when after 60s
		if time.Now().Unix()-ls.last > 120 {
			s.updateDB(ls.stat)
			delete(sm, ms.Aid) // NOTE: delete make sure the normal scope of memory and can be save all in 120s when close chan.
		}
	}
}

// updateDB update stat in db.
func (s *Service) updateDB(stat *api.Stat) (err error) {
	if _, err := s.dao.Update(context.TODO(), stat); err != nil {
		log.Error("s.dao.Update(%v) error(%v)", stat, err)
	}
	log.Info("update db aid(%d) stat(%+v) success", stat.Aid, stat)
	return
}

// multiUpdateDB update some stat in db.
func (s *Service) multiUpdateDB(yu int64, sm map[int64]*lastTmStat) (err error) {
	log.Info("start close(%d) multi update stat start", yu)
	var (
		c     = context.TODO()
		alloc = [1000]*api.Stat{}
		stats = alloc[:0]
		i     int
	)
	for aid, ls := range sm {
		stats = append(stats, ls.stat)
		if i > 0 && i%1000 == 0 {
			s.dao.MultiUpdate(c, yu, stats...)
		} else if i+1 == len(sm) {
			s.dao.MultiUpdate(c, yu, stats...)
		} else {
			log.Info("start close(%d) aid(%d) append", i, aid)
			continue
		}
		log.Info("start close(%d) multi update stat %d", i, aid)
		stats = alloc[:0]
	}
	log.Info("start close(%d) multi update stat endm", yu)
	return
}

// updateCache purge stat info in cache
func (s *Service) updateCache(st *api.Stat) (err error) {
	var (
		stat3 = &api.Stat{
			Aid:     st.Aid,
			View:    int32(st.View),
			Danmaku: int32(st.Danmaku),
			Reply:   int32(st.Reply),
			Fav:     int32(st.Fav),
			Coin:    int32(st.Coin),
			Share:   int32(st.Share),
			NowRank: int32(st.NowRank),
			HisRank: int32(st.HisRank),
			Like:    int32(st.Like),
			DisLike: 0,
		}
		click   *archive.Click3
		upclick = true
	)
	if click, err = s.dao.Click(context.TODO(), st.Aid); err != nil {
		upclick = false
	}
	if click == nil {
		click = &archive.Click3{}
	}
	for _, mc := range s.memcaches {
		var c = context.TODO()
		conn := mc.Get(c)
		if err = conn.Set(&memcache.Item{Key: statPBKey(stat3.Aid), Object: stat3, Flags: memcache.FlagProtobuf, Expiration: 0}); err != nil {
			log.Error("conn1.Set(%s, %+v) error(%v)", statPBKey(stat3.Aid), stat3, err)
		}
		if upclick {
			if err = conn.Set(&memcache.Item{Key: clickPBKey(stat3.Aid), Object: click, Flags: memcache.FlagProtobuf, Expiration: 0}); err != nil {
				log.Error("conn1.Set(%s, %+v) error(%v)", clickPBKey(stat3.Aid), click, err)
			}
		}
		if err == nil {
			log.Info("update cache aid(%d) stat(%+v) success", st.Aid, stat3)
			log.Info("update cache aid(%d) click(%+v) success", st.Aid, click)
		}
		conn.Close()
	}
	return
}

// Purge purge arc's stat cache
func (s *Service) Purge(c context.Context, aids []int64) (err error) {
	for _, aid := range aids {
		var stat *api.Stat
		if stat, err = s.dao.Stat(c, aid); err != nil {
			return
		}
		s.updateCache(stat)
	}
	return
}
