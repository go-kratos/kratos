package service

import (
	"context"
	"time"

	"go-common/app/service/main/vip/model"

	"github.com/pkg/errors"
)

const (
	_startMonth = -12
	_endMonth   = 36
	_hours      = 24
)

// BcoinGive bcoin give.
func (s *Service) BcoinGive(c context.Context, mid int64) (res *model.BcoinSalaryResp, err error) {
	var (
		sl         []*model.VipBcoinSalary
		nextDays   int32
		start, end time.Time
	)
	start = time.Now().AddDate(0, _startMonth, 0)
	end = time.Now().AddDate(0, _endMonth, 0)
	if sl, err = s.dao.BcoinSalaryList(c, mid, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}
	if nextDays, err = s.NextGiveBpDay(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	res = &model.BcoinSalaryResp{BcoinList: sl, DaysNextGive: nextDays}
	return
}

// NextGiveBpDay next give bp day.
func (s *Service) NextGiveBpDay(c context.Context, mid int64) (nextdays int32, err error) {
	var (
		giveBpDay = int(s.c.Property.GiveBpDay)
		now       = time.Now()
		salarys   []*model.VipBcoinSalary
	)
	_, _, d := now.Date()
	if d == giveBpDay {
		nextdays = s.daysBetween(c, now.AddDate(0, 1, 0), now)
		return
	}
	if d > giveBpDay {
		nextdays = s.daysBetween(c, now.AddDate(0, 1, int(giveBpDay)-d), now)
		return
	}
	if salarys, err = s.currentMonthSalary(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(salarys) == 1 && salarys[0].Status == model.Grant {
		nextdays = s.daysBetween(c, now.AddDate(0, 1, int(giveBpDay)-d), now)
	} else {
		nextdays = int32(giveBpDay - d)
	}
	return
}

func (s *Service) daysBetween(c context.Context, ta, tb time.Time) int32 {
	ta = ta.Truncate(_hours * time.Hour)
	tb = tb.Truncate(_hours * time.Hour)
	return int32(ta.Sub(tb).Hours() / _hours)
}

func (s *Service) currentMonthSalary(c context.Context, mid int64) (res []*model.VipBcoinSalary, err error) {
	year, month, _ := time.Now().Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)
	if res, err = s.dao.BcoinSalaryList(c, mid, start, end); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
