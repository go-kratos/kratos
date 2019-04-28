package service

import (
	"context"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) initStat(c context.Context, mid int64, fid int64) (err error) {
	var (
		s1, s2 *model.Stat
		es     = &model.Stat{}
		now    = time.Now()
	)
	if s1, err = s.Stat(c, mid); err != nil {
		return
	} else if s1.Empty() {
		if _, err = s.dao.AddStat(c, mid, es, now); err != nil {
			return
		}
	}
	if s2, err = s.Stat(c, fid); err != nil {
		return
	} else if s2.Empty() {
		if _, err = s.dao.AddStat(c, fid, es, now); err != nil {
			return
		}
	}
	return
}

func (s *Service) txStat(c context.Context, tx *sql.Tx, mid, fid int64) (mst *model.Stat, sst *model.Stat, err error) {
	// NOTE avoid db deadlock
	if mid < fid {
		if mst, err = s.dao.TxStat(c, tx, mid); err != nil {
			return
		}
		if sst, err = s.dao.TxStat(c, tx, fid); err != nil {
			return
		}
	} else {
		if sst, err = s.dao.TxStat(c, tx, fid); err != nil {
			return
		}
		if mst, err = s.dao.TxStat(c, tx, mid); err != nil {
			return
		}
	}
	return
}

// Stat get stat.
func (s *Service) Stat(c context.Context, mid int64) (stat *model.Stat, err error) {
	if mid <= 0 {
		return
	}
	var mc = true
	if stat, err = s.dao.StatCache(c, mid); err != nil {
		err = nil // ignore error
		mc = false
	} else if stat != nil {
		return
	}
	if stat, err = s.dao.Stat(c, mid); err != nil {
		return
	}
	if stat == nil {
		stat = &model.Stat{
			Mid: mid,
		}
	}
	if mc {
		s.addCache(func() {
			s.dao.SetStatCache(context.TODO(), mid, stat)
		})
	}
	return
}

// Stats get stats.
func (s *Service) Stats(c context.Context, mids []int64) (sts map[int64]*model.Stat, err error) {
	for _, v := range mids {
		if v <= 0 {
			return
		}
	}
	var (
		cache = true
		miss  []int64
	)
	if sts, miss, err = s.dao.StatsCache(c, mids); err != nil {
		err = nil // ignore error
		cache = false
	} else if len(miss) == 0 {
		return
	}
	for _, i := range miss {
		mid := i
		var stat *model.Stat
		if stat, err = s.dao.Stat(c, mid); err != nil {
			return
		}
		if stat == nil {
			stat = &model.Stat{
				Mid: mid,
			}
		}
		sts[mid] = stat
		if cache {
			s.addCache(func() {
				s.dao.SetStatCache(context.TODO(), mid, stat)
			})
		}
	}
	return
}

// SetStat set stat.
func (s *Service) SetStat(c context.Context, mid int64, st *model.Stat) (err error) {
	if mid <= 0 {
		return
	}
	var (
		tx  *sql.Tx
		nst *model.Stat
		now = time.Now()
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if nst, err = s.dao.TxStat(c, tx, mid); err != nil {
		return
	}
	if nst == nil {
		nst = new(model.Stat)
	}
	nst.Fill(st)
	_, err = s.dao.TxSetStat(c, tx, mid, nst, now)
	return
}

// DelStatCache delete stat cache.
func (s *Service) DelStatCache(c context.Context, mid int64) (err error) {
	err = s.dao.DelStatCache(c, mid)
	return
}
