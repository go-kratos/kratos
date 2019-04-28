package service

import (
	"context"
	"time"

	"go-common/app/job/main/mcn/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// UpMcnSignStateCron .
func (s *Service) UpMcnSignStateCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err     error
		c       = context.TODO()
		now     = time.Now()
		nowDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
		mss     []*model.MCNSignInfo
	)
	if mss, err = s.dao.McnSigns(c); err != nil {
		log.Error("s.dao.McnSigns error(%+v)", err)
		return
	}
	if len(mss) == 0 {
		log.Warn("mcn sign data is empty!")
		return
	}
	for _, v := range mss {
		var state int8
		switch {
		case v.State.NotDealState():
			continue
		case nowDate > v.EndDate.Time().Unix() && nowDate-v.EndDate.Time().Unix() <= model.ThirtyDayUnixTime && v.State != model.MCNSignStateOnCooling:
			state = int8(model.MCNSignStateOnCooling)
		case nowDate > v.EndDate.Time().Unix() && nowDate-v.EndDate.Time().Unix() > model.ThirtyDayUnixTime:
			state = int8(model.MCNSignStateOnExpire)
		case nowDate < v.BeginDate.Time().Unix() && v.State != model.MCNSignStateOnPreOpen:
			state = int8(model.MCNSignStateOnPreOpen)
		case v.BeginDate.Time().Unix() <= nowDate && nowDate <= v.EndDate.Time().Unix() && v.State != model.MCNSignStateOnSign && v.State == model.MCNSignStateOnPreOpen:
			state = int8(model.MCNSignStateOnSign)
		default:
			continue
		}
		if _, err = s.dao.UpMcnSignStateOP(c, v.SignID, state); err != nil {
			log.Error("s.dao.UpMcnSignStateOP(%d,%d) error(%+v)", v.SignID, state, err)
			continue
		}
		if err = s.dao.DelMcnSignCache(c, v.McnMid); err != nil {
			log.Error("s.dao.DelMcnSignCache(%d) error(%+v)", v.McnMid, err)
			continue
		}
		log.Info("signID(%d) change old state(%d) to new state(%d)", v.SignID, v.State, state)
	}
}

// UpExpirePayCron .
func (s *Service) UpExpirePayCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err error
		c   = context.TODO()
		sps []*model.SignPayInfo
	)
	if sps, err = s.dao.McnSignPayWarns(c); err != nil {
		log.Error("s.dao.McnSignPayWarns error(%+v)", err)
		return
	}
	if len(sps) == 0 {
		log.Warn("mcn sign pay date is empty!")
		return
	}
	ms := make(map[int64]struct{})
	for _, v := range sps {
		ms[v.SignID] = struct{}{}
	}
	for signID := range ms {
		if _, err = s.dao.UpMcnSignPayExpOP(c, signID); err != nil {
			log.Error("s.dao.UpMcnSignPayExpOP(%d) error(%+v)", signID, err)
			continue
		}
		log.Info("sign_id(%d) change pay data warn state to 2", signID)
	}
}
