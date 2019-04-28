package service

import (
	"context"
	"time"

	"go-common/app/job/main/mcn/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// UpMcnUpStateCron .
func (s *Service) UpMcnUpStateCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err     error
		page    = 1
		limit   = 100
		c       = context.TODO()
		now     = time.Now()
		nowDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
		mus     []*model.MCNUPInfo
	)
	for {
		offset := int64((page - 1) * limit)
		if mus, err = s.dao.McnUps(c, offset, int64(limit)); err != nil {
			log.Error("s.dao.McnUps(%d,%d) error(%+v)", offset, limit, err)
			return
		}
		if len(mus) == 0 {
			log.Warn("mcn up data is empty!")
			return
		}
		for _, v := range mus {
			var state int8
			switch {
			case v.State.NotDealState():
				continue
			case v.BeginDate.Time().Unix() <= nowDate && nowDate <= v.EndDate.Time().Unix() && v.State != model.MCNUPStateOnSign && v.State == model.MCNUPStateOnPreOpen:
				state = int8(model.MCNUPStateOnSign)
			case nowDate > v.EndDate.Time().Unix() && v.State != model.MCNUPStateOnExpire:
				state = int8(model.MCNUPStateOnExpire)
			default:
				continue
			}
			if _, err = s.dao.UpMcnUpStateOP(c, v.SignUpID, state); err != nil {
				log.Error("s.dao.UpMcnUpStateOP(%d,%d) error(%+v)", v.SignUpID, state, err)
				continue
			}
			log.Info("signUpID(%d) change old state(%d) to new state(%d)", v.SignUpID, v.State, state)
		}
		page++
	}
}
