package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/member/model"
	"go-common/library/log"
)

const (
	expMulti = 100
	level1   = 1
	level2   = 200
	level3   = 1500
	level4   = 4500
	level5   = 10800
	level6   = 28800
)

func (s *Service) initExp(c context.Context, mid int64) (err error) {
	var opers []*model.ExpOper
	if opers, err = s.CheckExpInit(c, mid); err != nil {
		log.Error("s.CheckExpInit(%d) error(%v)", mid, err)
		return
	}
	if len(opers) == 0 {
		log.Info("s.CheckExpInit(%d) opers eq(0) continue", mid)
		return
	}
	var exp *model.NewExp
	if exp, err = s.dao.SelExp(c, mid); err != nil {
		log.Error("s.dao.SelExp(%d) error(%v)", mid, err)
		return
	}
	if exp.Mid == 0 {
		if err = s.dao.InitExp(c, mid); err != nil {
			log.Error("s.dao.InitExp(%d) init user exp completed error(%v)", mid, err)
			return
		}
	}
	var (
		rows int64
		now  = time.Now().Unix()
	)
	for _, oper := range opers {
		if rows, err = s.dao.UpdateExpAped(c, mid, oper.Count*100, oper.Flag); err != nil {
			log.Error("s.dao.UpdateExpAped(%d) error(%v)", mid, err)
			return
		}
		if rows == 0 {
			log.Info("s.dao.UpdateExpAped(%d) exp(%d) flag(%d) rows affected eq(0) continue", mid, oper.Count*100, oper.Flag)
			continue
		}
		if err = s.dao.DatabusAddLog(c, mid, exp.Exp/100, (exp.Exp+oper.Count*100)/100, now, oper.Oper, oper.Reason, ""); err != nil {
			log.Error("s.dao.DatabusAddLog(%d) fromExp(%d) toExp(%d) ts(%d) oper(%s) reason(%s) error(%v)", mid, exp.Exp/100, (exp.Exp+oper.Count*100)/100, now, oper.Oper, oper.Reason, err)
			err = nil
			continue
		} else {
			log.Info("s.dao.DatabusAddLog(%d) fromExp(%d) toExp(%d) ts(%d) oper(%s) reason(%s) msg published", mid, exp.Exp/100, (exp.Exp+oper.Count*100)/100, now, oper.Oper, oper.Reason)
			exp.Exp = exp.Exp + oper.Count*100
			now++
		}
	}
	return
}

func (s *Service) delayUpdateExp() {
	s.limiter.UpdateExp.Wait(context.Background())
}

func (s *Service) addExp(c context.Context, e *model.AddExp) (err error) {
	if e.Mid <= 0 {
		return
	}
	now := time.Unix(e.Ts, 0)
	exp, eo, added, ok, err := s.checkExpAdd(c, e.Mid, e.Event, now)
	if err != nil || added || !ok {
		log.Info("s.checkExpAdd(%d) event(%s) result added(%v) ok(%v) err(%v)", e.Mid, e.Event, added, ok, err)
		return
	}

	// 写数据库限速，防止写入过大导致主从延迟
	s.delayUpdateExp()

	var rows int64
	if rows, err = s.dao.UpdateExpAped(c, e.Mid, eo.Count*100, eo.Flag); err != nil {
		log.Error("s.dao.UpdateExpAped(%d) exp(%d) flag(%d) error(%v) ", e.Mid, eo.Count*100, eo.Flag, err)
		return
	}
	if rows == 0 {
		log.Info("s.dao.UpdateExpAped(%d) exp(%d) flag(%d) rows affected eq(0) continue!", e.Mid, eo.Count*100, eo.Flag)
		return
	}
	if _, err = s.dao.SetExpAdded(context.Background(), e.Mid, now.Day(), eo.Oper); err != nil {
		log.Error("s.dao.SetExpAdded(%d) oper(%s)", e.Mid, eo.Oper)
		err = nil
	}
	if err = s.dao.DatabusAddLog(context.Background(), e.Mid, (exp.Exp)/100, (exp.Exp+eo.Count*100)/100, e.Ts, eo.Oper, eo.Reason, e.IP); err != nil {
		log.Error("s.dao.DatabusAddLog(%d) oper(%s) reason(%s) error(%v)", e.Mid, eo.Oper, eo.Reason, err)
		err = nil
	} else {
		log.Info("s.dao.DatabusAddLog(%d) oper(%s) reason(%s) msg published!", e.Mid, eo.Oper, eo.Reason)
	}
	return
}

func (s *Service) awardDo(ms []interface{}) {
	for _, m := range ms {
		l, ok := m.(*model.LoginLogIPString)
		if !ok {
			continue
		}
		s.addExp(context.TODO(), &model.AddExp{
			Mid:   l.Mid,
			IP:    l.Loginip,
			Ts:    l.Timestamp,
			Event: "login",
		})
		s.recoverMoral(context.TODO(), l.Mid)
		log.Info("consumer mid:%d,ts: %d", l.Mid, l.Timestamp)
	}
}

func isExpAndLevelChange(mu *model.Message) (bool, bool) {
	if mu.Action == "insert" {
		return true, true
	}
	if len(mu.Old) <= 0 || len(mu.New) <= 0 {
		return false, false
	}
	old := &model.ExpMessage{}
	new := &model.ExpMessage{}
	if err := json.Unmarshal(mu.New, new); err != nil {
		return false, false
	}
	if err := json.Unmarshal(mu.Old, old); err != nil {
		return false, false
	}
	expChange := false
	levelChange := false
	if old.Exp != new.Exp {
		expChange = true
	}
	if level(old.Exp) != level(new.Exp) {
		levelChange = true
	}
	return expChange, levelChange
}

func level(exp int64) int8 {
	exp = exp / expMulti
	switch {
	case exp < level1:
		return 0
	case exp < level2:
		return 1
	case exp < level3:
		return 2
	case exp < level4:
		return 3
	case exp < level5:
		return 4
	case exp < level6:
		return 5
	default:
		return 6
	}
}
