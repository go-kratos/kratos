package service

import (
	"context"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SetExp set user exp.
// NOTE: only for admin manager.
func (s *Service) SetExp(c context.Context, arg *model.ArgAddExp) (err error) {
	var base *model.BaseInfo
	if base, err = s.BaseInfo(c, arg.Mid); err != nil {
		log.Error("s.BaseInfo(%d) error(%v)", arg.Mid, err)
		return
	}
	if base.Rank < 10000 {
		err = ecode.UserNoMember
		return
	}
	var exp int64
	if exp, err = s.mbDao.Exp(c, arg.Mid); err != nil {
		log.Error("s.mbDao.Exp(%d) error(%v)", arg.Mid, err)
		return
	}
	if _, err = s.mbDao.SetExp(c, arg.Mid, int64(arg.Count*model.ExpMulti)); err != nil {
		log.Error("s.mbDao.SetExp(%d) error(%v)", arg.Mid, err)
		return
	}
	if err = s.mbDao.AddExplog(c, arg.Mid, exp/model.ExpMulti, int64(arg.Count), arg.Operate, arg.Reason, arg.IP); err != nil {
		log.Error("s.mbDao.AddExplog(%d) fromExp(%d) toExp(%d) oper(%s) reason(%s) ip(%s) error(%v)", arg.Mid, exp, int64(arg.Count), arg.Operate, arg.Reason, arg.IP, err)
	} else {
		log.Info("s.mbDao.AddExplog(%d) fromExp(%d) toExp(%d) oper(%s) reason(%s) ip(%s)", arg.Mid, exp, int64(arg.Count), arg.Operate, arg.Reason, arg.IP)
	}
	return
}

// UpdateExp update user exp.
func (s *Service) UpdateExp(c context.Context, arg *model.ArgAddExp) (err error) {
	var base *model.BaseInfo
	if base, err = s.BaseInfo(c, arg.Mid); err != nil {
		log.Error("s.BaseInfo(%d) error(%v)", arg.Mid, err)
		return
	}
	if base.Rank < 10000 {
		err = ecode.UserNoMember
		return
	}
	if arg.Count == 0 {
		log.Info("s.UpdateExp(%d) arg(%+v) count eq(0) continue", arg.Mid, arg)
		return
	}
	var exp int64
	if exp, err = s.mbDao.Exp(c, arg.Mid); err != nil {
		log.Error("s.mbDao.Exp(%d) error(%v)", arg.Mid, err)
		return
	}
	if exp == 0 {
		if _, err = s.mbDao.SetExp(c, arg.Mid, int64(arg.Count*model.ExpMulti)); err != nil {
			log.Error("s.mbDao.SetExp(%d) error(%v)", arg.Mid, err)
			return
		}
	} else {
		if _, err = s.mbDao.UpdateExp(c, arg.Mid, int64(arg.Count*model.ExpMulti)); err != nil {
			log.Error("s.mbDao.UpdateExp(%d) error(%v)", arg.Mid, err)
			return
		}
	}
	if err = s.mbDao.AddExplog(c, arg.Mid, exp/model.ExpMulti, (int64(arg.Count*model.ExpMulti)+exp)/model.ExpMulti, arg.Operate, arg.Reason, arg.IP); err != nil {
		log.Error("s.mbDao.AddExplog(%d) fromExp(%d) toExp(%d) oper(%s) reason(%s) ip(%s) error(%v)", arg.Mid, exp/model.ExpMulti, (int64(arg.Count*model.ExpMulti)+exp)/model.ExpMulti, arg.Operate, arg.Reason, arg.IP, err)
	} else {
		log.Info("s.mbDao.AddExplog(%d) fromExp(%d) toExp(%d) oper(%s) reason(%s) ip(%s)", arg.Mid, exp/model.ExpMulti, (int64(arg.Count*model.ExpMulti)+exp)/model.ExpMulti, arg.Operate, arg.Reason, arg.IP)
	}
	return
}

// Exp get user exp.
func (s *Service) Exp(c context.Context, mid int64) (exp *model.LevelInfo, err error) {
	var count int64
	if count, err = s.mbDao.Exp(c, mid); err != nil {
		log.Error("s.mbDao.Exp(%d) error(%v)", mid, err)
		return
	}
	exp = new(model.LevelInfo)
	exp.BuildLevel(count, true)
	return
}

// Level get user level info.
func (s *Service) Level(c context.Context, mid int64) (lv *model.LevelInfo, err error) {
	var count int64
	if count, err = s.mbDao.Exp(c, mid); err != nil {
		log.Error("s.mbDao.Exp(%d) error(%v)", mid, err)
		return
	}
	lv = new(model.LevelInfo)
	lv.BuildLevel(count, false)
	return
}

// Official is.
func (s *Service) Official(c context.Context, mid int64) (of *model.OfficialInfo, err error) {
	if of, err = s.mbDao.Official(c, mid); err != nil {
		return
	}
	if of == nil {
		err = ecode.NothingFound
	}
	return
}

// Exps get exps by mids.
func (s *Service) Exps(c context.Context, mids []int64) (lvs map[int64]*model.LevelInfo, err error) {
	var exps map[int64]int64
	if exps, err = s.mbDao.Exps(c, mids); err != nil {
		log.Error("s.mbDao.Exps(%v) error(%v)", mids, err)
		return
	}
	lvs = make(map[int64]*model.LevelInfo, len(exps))
	for mid, exp := range exps {
		lv := new(model.LevelInfo)
		lv.BuildLevel(exp, true)
		lvs[mid] = lv
	}
	return
}

// ExpLog get user exp log.
func (s *Service) ExpLog(c context.Context, mid int64, ip string) (lg []*model.UserLog, err error) {
	return s.mbDao.ExpLog(c, mid, ip)
}

// Stat get user exp stat.
func (s *Service) Stat(c context.Context, mid int64) (st *model.ExpStat, err error) {
	return s.mbDao.StatCache(c, mid, int64(time.Now().Day()))
}
