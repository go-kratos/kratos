package service

import (
	"context"

	plamdl "go-common/app/interface/main/playlist/model"
	"go-common/app/job/main/playlist/conf"
	"go-common/app/job/main/playlist/dao"
	"go-common/app/job/main/playlist/model"
	"go-common/library/log"
)

func (s *Service) viewproc(i int64) {
	defer s.waiter.Done()
	var (
		c  = context.TODO()
		ch = s.statCh[i]
	)
	for {
		sm, ok := <-ch
		if !ok {
			log.Warn("statproc(%d) quit", i)
			return
		}
		// filter view count
		if conf.Conf.Job.InterceptOn && sm.Count != nil && *sm.Count > 0 {
			if s.intercept(sm) {
				log.Info("intercept view count (pid:%d, aid:%d, ip:%s)", sm.ID, sm.Aid, sm.IP)
				dao.PromInfo("stat:访问计数拦截")
				continue
			}
		}
		s.upStat(c, sm, model.ViewCountType)
	}
}

func (s *Service) upStat(c context.Context, sm *model.StatM, tp string) {
	// update cache
	s.updateCache(c, sm, tp)
	// update db
	s.updateDB(c, sm, tp)
}

// updateDB update stat in db.
func (s *Service) updateDB(c context.Context, stat *model.StatM, tp string) (err error) {
	if _, err = s.dao.Update(context.TODO(), stat, tp); err != nil {
		return
	}
	log.Info("update db success "+tp+" pid(%d) count(%d) ", stat.ID, *stat.Count)
	dao.PromInfo("stat:更新计数DB")
	return
}

// updateCache update stat in cache
func (s *Service) updateCache(c context.Context, sm *model.StatM, tp string) (err error) {
	var (
		mid  int64
		fid  int64
		st   *plamdl.PlStat
		stat *plamdl.ArgStats
	)
	if st, err = s.dao.Stat(c, sm.ID); err != nil {
		log.Error("s.dao.Stat(%d) error(%v)", sm.ID, err)
		return
	}
	switch tp {
	case model.ViewCountType:
		stat = &plamdl.ArgStats{PlStat: &plamdl.PlStat{
			ID:    sm.ID,
			Mid:   st.Mid,
			Fid:   st.Fid,
			View:  *sm.Count,
			Reply: st.Reply,
			Fav:   st.Fav,
			Share: st.Share,
			MTime: st.MTime,
		}}
	case model.FavCountType:
		stat = &plamdl.ArgStats{PlStat: &plamdl.PlStat{
			ID:    sm.ID,
			Mid:   st.Mid,
			Fid:   st.Fid,
			View:  st.View,
			Reply: st.Reply,
			Fav:   *sm.Count,
			Share: st.Share,
			MTime: st.MTime,
		}}
	case model.ReplyCountType:
		stat = &plamdl.ArgStats{PlStat: &plamdl.PlStat{
			ID:    sm.ID,
			Mid:   st.Mid,
			Fid:   st.Fid,
			View:  st.View,
			Reply: *sm.Count,
			Fav:   st.Fav,
			Share: st.Share,
			MTime: st.MTime,
		}}
	case model.ShareCountType:
		stat = &plamdl.ArgStats{PlStat: &plamdl.PlStat{
			ID:    sm.ID,
			Mid:   st.Mid,
			Fid:   st.Fid,
			View:  st.View,
			Reply: st.Reply,
			Fav:   st.Fav,
			Share: *sm.Count,
			MTime: st.MTime,
		}}
	}
	if err = s.playlistRPC.SetStat(c, stat); err != nil {
		dao.PromError("stat:更新计数缓存",
			"s.playlistRPC.SetStat "+tp+" pid(%d) fid(%d) mid(%d) view(%d) favorite(%d) reply(%d) share(%d) error(%v)",
			sm.ID, fid, mid, *sm.Count, st.Fav, st.Reply, st.Share, err)
	}
	log.Info("update cache success "+tp+"  pid(%d)  aid(%d)  fid(%d) mid(%d) view(%d) favorite(%d) reply(%d) share(%d)",
		sm.ID, sm.Aid, fid, mid, *sm.Count, st.Fav, st.Reply, st.Share)
	dao.PromInfo("stat:更新计数缓存")
	return
}

// intercept intercepts illegal views.
func (s *Service) intercept(stat *model.StatM) bool {
	return s.dao.Intercept(context.TODO(), stat.ID, stat.Aid, stat.IP)
}
