package service

import (
	"context"
	"time"
	"unicode/utf8"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// TipList tip list.
func (s *Service) TipList(c context.Context, platform int8, state int8, position int8) (ts []*model.Tips, err error) {
	var now = time.Now().Unix()
	if ts, err = s.dao.TipList(c, platform, state, now, position); err != nil {
		err = errors.WithStack(err)
	}
	for _, v := range ts {
		v.TipState(v.StartTime, v.EndTime, now)
	}
	return
}

// TipByID tip by id.
func (s *Service) TipByID(c context.Context, id int64) (t *model.Tips, err error) {
	if t, err = s.dao.TipByID(c, id); err != nil {
		err = errors.WithStack(err)
	}
	if t == nil {
		err = ecode.VipTipNotFoundErr
		return
	}
	t.TipState(t.StartTime, t.EndTime, time.Now().Unix())
	return
}

// AddTip add tip.
func (s *Service) AddTip(c context.Context, t *model.Tips) (err error) {
	if t.StartTime >= t.EndTime {
		err = ecode.RequestErr
		return
	}
	if utf8.RuneCountInString(t.Tip) > _maxTipLen {
		err = ecode.VipTipTooLoogErr
		return
	}
	t.Ctime = xtime.Time(time.Now().Unix())
	if _, err = s.dao.AddTip(c, t); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// TipUpdate update tip.
func (s *Service) TipUpdate(c context.Context, t *model.Tips) (err error) {
	var (
		old *model.Tips
		now = time.Now().Unix()
	)
	if t.ID == 0 {
		err = ecode.RequestErr
		return
	}
	if old, err = s.TipByID(c, t.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if t.StartTime >= t.EndTime {
		err = ecode.VipTipTimeErr
		return
	}
	if utf8.RuneCountInString(t.Tip) > _maxTipLen {
		err = ecode.VipTipTooLoogErr
		return
	}
	if old.StartTime != t.StartTime && old.StartTime < now {
		err = ecode.VipTipStartTimeCatNotModifyErr
		return
	}
	if old.EndTime != t.EndTime && old.EndTime < now {
		err = ecode.VipTipEndTimeCatNotModifyErr
		return
	}
	if _, err = s.dao.TipUpdate(c, t); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DeleteTip delete tip.
func (s *Service) DeleteTip(c context.Context, id int64, operator string) (err error) {
	var (
		old *model.Tips
		now = time.Now().Unix()
	)
	if old, err = s.TipByID(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if old.StartTime <= now {
		err = ecode.VipTipCatNotDeleteErr
		return
	}
	if _, err = s.dao.DeleteTip(c, id, model.Delete, operator); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// ExpireTip expire tip.
func (s *Service) ExpireTip(c context.Context, id int64, operator string) (err error) {
	var (
		old *model.Tips
		now = time.Now().Unix()
	)
	if old, err = s.TipByID(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if old.StartTime > now || old.EndTime < now {
		err = ecode.VipTipCatNotExpireErr
		return
	}
	if _, err = s.dao.ExpireTip(c, id, operator, now); err != nil {
		err = errors.WithStack(err)
	}
	return
}
