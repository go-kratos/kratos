package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_ps = 1000
)

var _emptyMsg = make(map[int64]*model.Message)

// AdminRuleByID .
func (s *Service) AdminRuleByID(c context.Context, id int64) (rs *model.FilterInfo, err error) {
	return s.dao.Filter(c, id)
}

// AdminSearch manager search filter.
func (s *Service) AdminSearch(c context.Context, msg, area, sourceStr, typeStr string, level int, state, deleted int, pn int64, ps int64) (rules []*model.FilterInfo, count int64, err error) {
	if sourceStr == "" {
		sourceStr = xstr.JoinInts(s.conf.Property.SourceMask)
	}
	if typeStr == "" {
		typeStr = xstr.JoinInts(s.conf.Property.FilterType)
	}
	var levelStr = ""
	if level == 0 {
		levelStr = xstr.JoinInts(s.conf.Property.Level)
	} else {
		levelStr = fmt.Sprintf("%d", level)
	}
	if count, err = s.dao.SearchCount(c, msg, area, sourceStr, typeStr, levelStr, state, deleted); err != nil {
		log.Error("s.dao.SearchCount err(%v)", err)
		return
	}
	start := (pn - 1) * ps
	offset := ps
	if start > count {
		return
	}
	rules, err = s.dao.Search(c, msg, area, sourceStr, typeStr, levelStr, state, deleted, start, offset)
	return
}

// AdminAdd .
func (s *Service) AdminAdd(c context.Context, areas, rules []string, level *model.AreaLevel, comment, name string, mode int8, tps []int64, adid, stime, etime int64, source, keyType int8) (err error) {
	var (
		r  *model.FilterInfo
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	for _, rule := range rules {
		if err = s.checkArea(c, areas); err != nil {
			return
		}
		if err = s.checkReg(mode, rule); err != nil {
			return
		}
		// 检查正常信息是否有大面积的误伤
		if err = s.checkWhiteSample(mode, rule); err != nil {
			return
		}
		if r, err = s.dao.FilterByContent(c, rule); err != nil {
			return
		}
		if r != nil && r.State == model.FilterStateNormal {
			return ecode.FilterDuplicateContent
		}
		var ruleID int64
		if ruleID, err = s.dao.UpsertRule(c, tx, rule, comment, level.Level, mode, source, keyType, time.Unix(stime, 0), time.Unix(etime, 0)); err != nil {
			tx.Rollback()
			return
		}
		if r != nil {
			ruleID = r.ID
		}
		if _, err = s.dao.DeleteAreaRules(c, tx, ruleID); err != nil {
			tx.Rollback()
			return
		}
		for _, area := range areas {
			var (
				areaLevel int8
				ok        bool
			)
			if areaLevel, ok = level.Area[area]; !ok {
				areaLevel = level.Level
			}
			for _, tp := range tps {
				if _, err = s.dao.UpsertAreaRule(c, tx, area, tp, ruleID, areaLevel); err != nil {
					tx.Rollback()
					return
				}
			}
		}
		if _, err = s.dao.InsertLog(c, tx, ruleID, adid, comment, name, model.LogStateAdd); err != nil {
			tx.Rollback()
			return
		}
	}
	s.mission(func() {
		s.notifySearch(context.TODO(), areas)
	})
	return
}

// AdminEdit manager edit filter.
func (s *Service) AdminEdit(c context.Context, areas []string, rule, comment, reason, name string, mode int8, level *model.AreaLevel, tps []int64, id, adid, stime, etime int64, source, keyType int8) (err error) {
	if err = s.checkArea(c, areas); err != nil {
		return
	}
	if err = s.checkReg(mode, rule); err != nil {
		return
	}
	if err = s.checkWhiteSample(mode, rule); err != nil {
		return
	}
	var afterRule *model.FilterInfo
	if afterRule, err = s.dao.FilterByContent(c, rule); err != nil {
		return
	}
	if afterRule != nil && afterRule.State == model.FilterStateNormal && afterRule.ID != id {
		return ecode.FilterDuplicateContent
	}

	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	// 先删除原敏感词
	if _, err = s.dao.TxUpdateRuleState(c, tx, id, model.FilterStateDeleted); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.DeleteAreaRules(c, tx, id); err != nil {
		tx.Rollback()
		return
	}
	var newID int64
	// insert 新敏感词 或 update 已有敏感词
	if newID, err = s.dao.UpsertRule(c, tx, rule, comment, level.Level, mode, source, keyType, time.Unix(stime, 0), time.Unix(etime, 0)); err != nil {
		tx.Rollback()
		return
	}
	if afterRule != nil {
		newID = afterRule.ID
	}
	// 如果新的敏感词和之前不一致，插入一条日志
	if newID != id {
		if _, err = s.dao.InsertLog(c, tx, id, adid, reason, name, model.LogStateEdit); err != nil {
			tx.Rollback()
			return
		}
	}
	// 插入area
	for _, area := range areas {
		var (
			areaLevel int8
			ok        bool
		)
		if areaLevel, ok = level.Area[area]; !ok {
			areaLevel = level.Level
		}
		for _, tp := range tps {
			if _, err = s.dao.UpsertAreaRule(c, tx, area, tp, newID, areaLevel); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	if _, err = s.dao.InsertLog(c, tx, newID, adid, reason, name, model.LogStateEdit); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	s.mission(func() {
		s.notifySearch(context.TODO(), areas)
	})
	return
}

// AdminDel .
func (s *Service) AdminDel(c context.Context, fid, adid int64, comment, name string) (err error) {
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.TxUpdateRuleState(c, tx, fid, model.FilterStateDeleted); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.DeleteAreaRules(c, tx, fid); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.InsertLog(c, tx, fid, adid, comment, name, model.LogStateDel); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()

	s.mission(func() {
		s.notifySearch(context.TODO(), []string{"common"})
	})
	return
}

// AdminLog .
func (s *Service) AdminLog(c context.Context, id int64) (ls []*model.Log, err error) {
	return s.dao.Logs(c, id)
}

// AdminOrigin .
func (s *Service) AdminOrigin(c context.Context, id int64, area string) (res *model.Message, err error) {
	var ct string
	if ct, err = s.dao.Content(c, id, area); err != nil {
		return
	}
	res = &model.Message{ID: id, Content: ct}
	return
}

// AdminOrigins .
func (s *Service) AdminOrigins(c context.Context, ids []int64, area string) (res map[int64]*model.Message, err error) {
	var (
		id int64
		ct string
	)
	res = make(map[int64]*model.Message)
	for _, id = range ids {
		if ct, err = s.dao.Content(c, id, area); err != nil {
			res = _emptyMsg
			return
		}
		res[id] = &model.Message{ID: id, Content: ct}
	}
	return
}

// expireFilter
func (s *Service) expireFilter() (err error) {
	maxID, err := s.dao.MaxFilterID(context.TODO())
	if err != nil {
		log.Error("s.dao.MaxFilterID, err(%v)", err)
		return
	}
	total := maxID / _ps
	for i := 0; int64(i) <= total; i++ {
		start := i * _ps
		end := (i + 1) * _ps
		IDs, err := s.dao.ExpiredRuleIDs(context.TODO(), start, end)
		if err != nil {
			log.Error("s.dao.GetOverdueFilterID", err)
			continue
		}
		if len(IDs) > 0 {
			log.Info("handle filter stage bewteen[%d, %d], ids [%+v]", start, end, IDs)
			if _, err := s.dao.UpdateRulesState(context.TODO(), IDs, model.FilterStateExpired); err != nil {
				log.Error("s.dao.UpdateState(%+v,%d,%+v)", IDs, model.FilterStateExpired, err)
			}
		}
	}
	return
}
