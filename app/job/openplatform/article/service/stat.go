package service

import (
	"context"
	"strconv"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/log"
)

func (s *Service) statproc(i int64) {
	defer s.waiter.Done()
	var (
		err  error
		ls   *lastTimeStat
		c    = context.TODO()
		ch   = s.statCh[i]
		last = s.statLastTime[i]
	)
	for {
		stat, ok := <-ch
		if !ok {
			log.Warn("statproc(%d) quit", i)
			s.multiUpdateDB(i, last)
			return
		}
		// filter view count
		if stat.View != nil && *stat.View > 0 {
			var ban bool
			var reason, valid string
			if ban = s.intercept(stat); ban {
				log.Info("intercept view count (aid:%d, ip:%s, mid:%d)", stat.Aid, stat.IP, stat.Mid)
				dao.PromInfo("stat:访问计数拦截")
				reason = "访问计数拦截"
			} else if ban = s.dao.DupViewIntercept(c, stat.Aid, stat.Mid); ban {
				log.Info("dupintercept view count (aid:%d, ip:%s, mid:%d)", stat.Aid, stat.IP, stat.Mid)
				dao.PromInfo("stat:重复访问计数拦截")
				reason = "重复访问计数拦截"
			} else if stat.CheatInfo != nil {
				viewLv, _ := strconv.Atoi(stat.CheatInfo.Lv)
				if limitLv, ok1 := s.cheatArts[stat.Aid]; ok1 && (viewLv <= limitLv) {
					ban = true
					log.Info("lvintercept view count (aid:%d, ip:%s, mid:%d)", stat.Aid, stat.IP, stat.Mid)
					dao.PromInfo("stat:等级访问计数拦截")
					reason = "等级访问计数拦截"
				}
			}
			if ban {
				valid = "0"
			} else {
				valid = "1"
			}
			if stat.CheatInfo != nil {
				stat.CheatInfo.Valid = valid
				stat.CheatInfo.Reason = reason
			}
			s.cheatInfo(stat.CheatInfo)
			if ban {
				continue
			}
		}
		// get stat
		if ls, ok = last[stat.Aid]; !ok {
			var st *artmdl.StatMsg
			if st, err = s.dao.Stat(c, stat.Aid); err != nil {
				log.Error("s.dao.Stat(%d) error(%+v)", stat.Aid, err)
				continue
			}
			ls = &lastTimeStat{}
			if st == nil {
				ls.stat = &artmdl.StatMsg{Aid: stat.Aid, View: new(int64), Like: new(int64), Dislike: new(int64), Favorite: new(int64), Reply: new(int64), Share: new(int64), Coin: new(int64)}
				ls.time = 0 // NOTE: make sure update db in first.
			} else {
				ls.stat = st
				ls.time = time.Now().Unix()
			}
			last[stat.Aid] = ls
		}
		changed := model.Merge(ls.stat, stat)
		// update cache
		s.updateCache(c, ls.stat, 0)
		s.updateSortCache(c, ls.stat.Aid, changed)
		// update db after 120s
		if time.Now().Unix()-ls.time > s.updateDbInterval {
			s.updateDB(c, ls.stat, 0)
			s.updateRecheckDB(c, ls.stat)
			s.updateSearchStats(c, ls.stat)
			delete(last, stat.Aid) // NOTE: delete ensures that memory should be normal in 120s after channel has been closed.
		}
	}
}

// updateCache purge stat info in cache
func (s *Service) updateCache(c context.Context, sm *artmdl.StatMsg, count int) (err error) {
	stat := &artmdl.ArgStats{
		Aid: sm.Aid,
		Stats: &artmdl.Stats{
			View:     *sm.View,
			Like:     *sm.Like,
			Dislike:  *sm.Dislike,
			Favorite: *sm.Favorite,
			Reply:    *sm.Reply,
			Share:    *sm.Share,
			Coin:     *sm.Coin,
		},
	}
	if err = s.articleRPC.SetStat(context.TODO(), stat); err != nil {
		log.Error("s.articleRPC.SetStat aid(%d) view(%d) likes(%d) dislike(%d) favorite(%d) reply(%d) share(%d) coin(%d) error(%+v)",
			sm.Aid, *sm.View, *sm.Like, *sm.Dislike, *sm.Favorite, *sm.Reply, *sm.Share, *sm.Coin, err)
		dao.PromError("stat:更新计数缓存")
		s.dao.PushStat(c, &dao.StatRetry{
			Action: dao.RetryUpdateStatCache,
			Count:  count,
			Data:   sm,
		})
		return
	}
	log.Info("update cache success aid(%d) view(%d) likes(%d) dislike(%d) favorite(%d) reply(%d) share(%d) coin(%d)",
		sm.Aid, *sm.View, *sm.Like, *sm.Dislike, *sm.Favorite, *sm.Reply, *sm.Share, *sm.Coin)
	dao.PromInfo("stat:更新计数缓存")
	return
}

// updateSortCache update sort cache
func (s *Service) updateSortCache(c context.Context, aid int64, changed [][2]int64) (err error) {
	if len(changed) == 0 {
		return
	}
	arg := &artmdl.ArgSort{
		Aid:     aid,
		Changed: changed,
	}
	if err = s.articleRPC.UpdateSortCache(context.TODO(), arg); err != nil {
		log.Error("s.articleRPC.UpdateSortCache(aid:%v arg: %+v)", aid, arg)
		dao.PromError("stat:更新排序缓存")
		return
	}
	log.Info("success s.articleRPC.UpdateSortCache(aid:%v arg: %+v)", aid, arg)
	dao.PromInfo("stat:更新排序缓存")
	return
}

// updateDB update stat in db.
func (s *Service) updateDB(c context.Context, stat *artmdl.StatMsg, count int) (err error) {
	if _, err = s.dao.Update(context.TODO(), stat); err != nil {
		dao.PromError("stat:更新计数DB")
		s.dao.PushStat(c, &dao.StatRetry{
			Action: dao.RetryUpdateStatDB,
			Count:  count,
			Data:   stat,
		})
		return
	}
	log.Info("update db success aid(%d) view(%d) likes(%d) dislike(%d) favorite(%d) reply(%d) share(%d) coin(%d)",
		stat.Aid, *stat.View, *stat.Like, *stat.Dislike, *stat.Favorite, *stat.Reply, *stat.Share, *stat.Coin)
	return
}

// multiUpdateDB update some stat in db.
func (s *Service) multiUpdateDB(i int64, last map[int64]*lastTimeStat) (err error) {
	log.Info("multiUpdateDB close(%d) start", i)
	c := context.TODO()
	for aid, ls := range last {
		s.updateDB(c, ls.stat, 0)
		log.Info("multiUpdateDB close(%d) update stats aid: %d, value: %+v", i, aid, ls.stat)
	}
	log.Info("multiUpdateDB close(%d) end", i)
	return
}

// intercept intercepts illegal views.
func (s *Service) intercept(stat *artmdl.StatMsg) bool {
	return s.dao.Intercept(context.TODO(), stat.Aid, stat.Mid, stat.IP)
}

func (s *Service) cheatInfo(cheat *artmdl.CheatInfo) {
	if cheat == nil {
		return
	}
	log.Info("cheatInfo %+v", cheat)
	if err := s.cheatInfoc.Info(cheat.Valid, cheat.Client, cheat.Cvid, cheat.Mid, cheat.Lv, cheat.Ts, cheat.IP, cheat.UA, cheat.Refer, cheat.Sid, cheat.Buvid, cheat.DeviceID, cheat.Build, cheat.Reason); err != nil {
		log.Error("cheatInfo error(%+v)", err)
	}
}

func (s *Service) updateRecheckDB(c context.Context, stat *artmdl.StatMsg) (err error) {
	var (
		publishTime int64
		checkState  int
	)
	if s.setting.Recheck.View == 0 || s.setting.Recheck.Day == 0 {
		return
	}
	if isRecheck, _ := s.dao.GetRecheckCache(c, stat.Aid); isRecheck {
		return
	}
	if publishTime, checkState, err = s.dao.GetRecheckInfo(c, stat.Aid); err != nil || checkState != 0 {
		return
	}

	if s.dao.IsAct(c, stat.Aid) {
		return
	}
	if *(stat.View) > s.setting.Recheck.View {
		if publishTime+60*60*24*s.setting.Recheck.Day+s.updateDbInterval > time.Now().Unix() {
			if err = s.dao.UpdateRecheck(c, stat.Aid); err == nil {
				log.Info("update recheck success aid(%d)", stat.Aid)
				s.dao.SetRecheckCache(c, stat.Aid)
			}
		}
	}
	return
}
