package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/member/model"
	"go-common/library/log"
)

func (s *Service) checkExpAdd(c context.Context, mid int64, event string, now time.Time) (exp *model.NewExp, eo *model.ExpOper, added, ok bool, err error) {
	if eo, ok = model.ExpFlagOper[event]; !ok {
		log.Info("s.checkExpAdd(%d) oper(%s) not found", mid, event)
		return
	}
	var base *model.BaseInfo
	if base, err = s.dao.BaseInfo(c, mid); err != nil {
		log.Error("s.dao.BaseInfo(%d) error(%v)", mid, err)
		return
	}
	if base == nil {
		err = fmt.Errorf("No base info with mid(%v)", mid)
		log.Error("Failed to checkExpAdd with mid(%d) error: %+v", mid, err)
		return
	}
	if ok = !(base.Rank < 10000); !ok {
		log.Info("s.checkExpAdd(%d) mid.Rank<10000", mid)
		return
	}
	if added, err = s.dao.ExpAdded(c, mid, now.Day(), eo.Oper); err != nil || added {
		log.Info("s.dao.ExpAdded(%d) error(%v) added(%v)", mid, err, added)
		return
	}
	if exp, err = s.dao.SelExp(c, mid); err != nil {
		log.Error("s.dao.SelExp(%d) error(%v)", mid, err)
		return
	}
	if now.Unix()-int64(exp.Addtime) < 24*60*60 {
		added = exp.Flag&eo.Flag == eo.Flag
		return
	}
	if err = s.dao.InitExp(c, mid); err != nil {
		log.Error("s.dao.InitExp(%d) error(%v)", mid, err)
		return
	}
	exp.FlagDailyReset(now)
	if err = s.dao.UpdateExpFlag(c, mid, exp.Flag, exp.Addtime); err != nil {
		log.Error("s.dao.UpdateExpFlag(%d) flag(%d) addtime(%v)", mid, exp.Flag, exp.Addtime)
		return
	}
	added = exp.Flag&eo.Flag == eo.Flag
	return
}

// CheckExpInit check init user exp if exp not exist.
func (s *Service) CheckExpInit(c context.Context, mid int64) (opers []*model.ExpOper, err error) {
	var aso *model.MemberAso
	if aso, err = s.dao.AsoStatus(c, mid); err != nil {
		log.Error("s.dao.AsoStatus(%d) error(%v)", mid, err)
		return
	}
	if aso.Spacesta >= 0 && len(aso.Email) != 0 {
		opers = append(opers, model.ExpFlagOper["email"])
	}
	if len(aso.Telphone) != 0 {
		opers = append(opers, model.ExpFlagOper["phone"])
	}
	if aso.SafeQs != 0 {
		opers = append(opers, model.ExpFlagOper["safe"])
	}
	var ri *model.RealnameInfo
	if ri, err = s.dao.RealnameInfo(c, mid); err != nil {
		log.Error("s.dao.RealnameInfo(%d) error(%+v)", mid, err)
		return
	}
	if ri != nil && ri.Status.IsPass() {
		opers = append(opers, model.ExpFlagOper["identify"])
	}
	log.Info("exp init opers with mid: %d: %+v", mid, opers)
	return
}
func sameAccInfo(base *model.BaseInfo, res *model.AccountInfo) (same bool) {
	return sameName(base, res)
}

func sameName(base *model.BaseInfo, res *model.AccountInfo) bool {
	return base.Name == res.Name
}
