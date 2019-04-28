package service

import (
	"context"

	"go-common/app/admin/main/filter/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

var _emptyRules = []*model.KeyInfo{}

// AddKey .
func (s *Service) AddKey(c context.Context, areas []string, key, rule, comment, name string, mode, level int8, adid, stime, etime int64) (err error) {
	var (
		filterID int64
		fRule    *model.KeyInfo
	)
	if err = s.checkArea(c, areas); err != nil {
		return
	}
	if err = s.checkReg(mode, rule); err != nil {
		return
	}
	if err = s.checkWhiteSample(mode, rule); err != nil {
		return
	}
	// 存在判断
	if fRule, err = s.dao.ConKey(c, rule, key); err != nil {
		return
	}
	if fRule == nil {
		if filterID, err = s.dao.InsertConkey(c, rule, key, comment, mode, level, stime, etime); err != nil {
			return
		}
	} else {
		filterID = fRule.ID
		if _, err = s.dao.DelKeyFid(c, key, filterID); err != nil {
			return
		}
		if _, err = s.dao.UpdateConkey(c, rule, key, comment, mode, level, stime, etime); err != nil {
			return
		}
	}
	s.cacheCh.Save(func() {
		s.dao.InsertFkLog(context.TODO(), key, name, comment, adid, model.LogStateAdd)
	})
	// 插入操作
	for _, area := range areas {
		if _, err = s.dao.InsertKey(c, area, key, filterID); err != nil {
			return
		}
	}
	// 清理cache
	s.cacheCh.Save(func() {
		for _, area := range areas {
			s.dao.DelKeyAreaCache(context.TODO(), key, area)
		}
	})
	return
}

// DelKeyFid .
func (s *Service) DelKeyFid(c context.Context, key string, fid, adid int64, comment, name, reason string) (err error) {
	var (
		areas []string
	)
	if areas, err = s.dao.KeyArea(c, key, fid); err != nil {
		return
	}
	if err = s.delkey(c, fid, key); err != nil {
		return
	}
	// 清理cache
	s.cacheCh.Save(func() {
		ctx := context.TODO()
		// log
		s.dao.InsertFkLog(c, key, name, reason, adid, model.LogStateDel)
		for _, area := range areas {
			s.dao.DelKeyAreaCache(ctx, key, area)
		}
	})
	return
}

// EditInfo .
func (s *Service) EditInfo(c context.Context, key string, id int64) (fil *model.KeyInfo, err error) {
	if fil, err = s.dao.ConKeyByID(c, id, key); err != nil {
		return
	}
	var areas []string
	if areas, err = s.dao.KeyArea(c, key, fil.ID); err != nil {
		return
	}
	if len(areas) == 0 {
		areas = []string{}
		fil.State = 1
	}
	fil.TpIDs = []int64{}
	fil.Areas = areas
	return
}

// EditKey .
func (s *Service) EditKey(c context.Context, key string, oldFid int64, areas []string, mode int8, rule string, level int8,
	stime, etime int64, adid int64, name, comment, reason string) (err error) {
	var (
		filterID    int64
		beforeAreas []string
		fRule       *model.KeyInfo
	)
	if err = s.checkReg(mode, rule); err != nil {
		return
	}
	if err = s.checkReg(mode, rule); err != nil {
		return
	}
	if err = s.checkWhiteSample(mode, rule); err != nil {
		return
	}
	if beforeAreas, err = s.dao.KeyArea(c, key, oldFid); err != nil {
		return
	}
	// 这里需要把老的删掉
	if err = s.delkey(c, oldFid, key); err != nil {
		return
	}
	// 存在判断
	if fRule, err = s.dao.ConKey(c, rule, key); err != nil {
		return
	}
	if fRule == nil {
		if filterID, err = s.dao.InsertConkey(c, rule, key, comment, mode, level, stime, etime); err != nil {
			return
		}
	} else {
		filterID = fRule.ID
		if _, err = s.dao.DelKeyFid(c, key, filterID); err != nil {
			return
		}
		if _, err = s.dao.UpdateConkey(c, rule, key, comment, mode, level, stime, etime); err != nil {
			return
		}
	}
	s.cacheCh.Save(func() {
		s.dao.InsertFkLog(context.TODO(), key, name, reason, adid, model.LogStateEdit)
	})
	// 更新新的关系
	for _, area := range areas {
		if _, err = s.dao.InsertKey(c, area, key, filterID); err != nil {
			return
		}
	}
	// 清理旧有cache must here
	s.cacheCh.Save(func() {
		ctx := context.TODO()
		for _, area := range beforeAreas {
			s.dao.DelKeyAreaCache(ctx, key, area)
		}
	})
	return
}

// SearchKey .
func (s *Service) SearchKey(c context.Context, key, comment string, pn, ps int64, state int8) (total int64, rs []*model.KeyInfo, err error) {
	start := (pn - 1) * ps
	if total, err = s.dao.CountKey(c, key, comment, state); err != nil {
		return
	}
	if rs, err = s.dao.SearchKey(c, key, comment, start, ps, state); err != nil {
		return
	}
	if len(rs) == 0 {
		rs = _emptyRules
		return
	}
	// 拼装area
	var r *model.KeyInfo
	for _, r = range rs {
		var areas []string
		if areas, err = s.dao.KeyArea(c, r.Key, r.ID); err != nil {
			return
		}
		r.TpIDs = []int64{}
		if len(areas) == 0 {
			areas = []string{}
			r.State = 1
		}
		r.Areas = areas
		r.Shelve = s.shelve(r.Stime.Time(), r.Etime.Time())
		if int64(r.Etime-r.Stime) == 3600*24*365*10 { // 默认时间
			r.Etime = 0
		}
	}
	return
}

// FkLog .
func (s *Service) FkLog(c context.Context, key string) (ls []*model.Log, err error) {
	return s.dao.FkLogs(c, key)
}

func (s *Service) delkey(c context.Context, fid int64, key string) (err error) {
	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.TxDelCon(tx, fid); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxDelKeyFid(tx, key, fid); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
	}
	return
}

func (s *Service) shelve(stime, etime time.Time) bool {
	if time.Now().Unix() >= stime.Unix() && time.Now().Unix() <= etime.Unix() {
		return true
	}
	return false
}

func (s *Service) notifySearch(c context.Context, areas []string) {
	if s.conf.HTTPClient.Off {
		return
	}
	if err := s.searchDao.Notify(c, areas); err != nil {
		log.Error("notifySearch err(%v)", err)
	}
}
