package service

import (
	"context"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"

	"go-common/library/log"
)

// GetRuleByArea .
func (s *SvcImpl) GetRuleByArea(ctx context.Context, area string) ([]*model.Rule, error) {
	rs, err := s.RuleDao.GetByArea(ctx, ToDaoCond(&Condition{
		Area:  area,
		State: model.StateDefault,
	}))
	if err == dao.ErrResourceNotExist {
		return []*model.Rule{}, nil
	}
	if err != nil {
		return nil, err
	}
	return ToModelRules(rs), nil
}

// GetRuleByAreaAndLimitTypeAndScope .
func (s *SvcImpl) GetRuleByAreaAndLimitTypeAndScope(ctx context.Context, area, limitType, limitScope string) (*model.Rule, error) {
	cond := Condition{
		Area:       area,
		LimitType:  limitType,
		LimitScope: limitScope,
	}
	r, err := s.RuleDao.GetByAreaAndTypeAndScope(ctx, ToDaoCond(&cond))
	if err != nil {
		return nil, err
	}
	return ToModelRule(r), nil
}

var (
	rules = []*model.Rule{}
)

// RefreshRules .
func (s *SvcImpl) RefreshRules(ctx context.Context) {
	s.Lock()
	defer s.Unlock()

	rs, _, err := s.RuleDao.GetByCond(ctx, ToDaoCond(&Condition{State: model.StateDefault}))
	if err != nil {
		return
	}
	rules = ToModelRules(rs)
}

// GetAggregateRuleByAreaAndLimitType .
func (s *SvcImpl) GetAggregateRuleByAreaAndLimitType(ctx context.Context, area, limitType string) (*model.AggregateRule, error) {
	s.RLock()
	defer s.RUnlock()

	res := &model.AggregateRule{}
	for _, r := range rules {
		if r.Area == area && r.LimitType == limitType {
			if r.LimitScope == model.LimitScopeGlobal {
				res.GlobalDurationSec = r.DurationSec
				res.GlobalAllowedCounts = r.AllowedCounts
			}
			if r.LimitScope == model.LimitScopeLocal {
				res.LocalDurationSec = r.DurationSec
				res.LocalAllowedCounts = r.AllowedCounts
			}
		}
	}

	if res.GlobalAllowedCounts > conf.Conf.MaxAllowedCounts {
		res.GlobalAllowedCounts = conf.Conf.MaxAllowedCounts
	}
	if res.LocalAllowedCounts > conf.Conf.MaxAllowedCounts {
		res.LocalAllowedCounts = conf.Conf.MaxAllowedCounts
	}
	if res.GlobalDurationSec > conf.Conf.MaxDurationSec {
		res.GlobalDurationSec = conf.Conf.MaxDurationSec
	}
	if res.LocalDurationSec > conf.Conf.MaxDurationSec {
		res.LocalDurationSec = conf.Conf.MaxDurationSec
	}

	return res, nil
}

// UpsertRule .
func (s *SvcImpl) UpsertRule(ctx context.Context, r *model.Rule) (*model.Rule, error) {
	_, err := s.RuleDao.GetByAreaAndTypeAndScope(ctx, ToDaoCond(&Condition{
		Area:       r.Area,
		LimitType:  r.LimitType,
		LimitScope: r.LimitScope,
	}))
	var res *dao.Rule
	if err == nil {
		res, err = s.RuleDao.Update(ctx, ToDaoRule(r))
	} else {
		res, err = s.RuleDao.Insert(ctx, ToDaoRule(r))
	}
	if err != nil {
		return nil, err
	}
	if err := s.antiDao.DelRulesCache(ctx, r.Area, r.LimitType); err != nil {
		log.Error("s.antiDao.DelRulesCache(%s,%s) error(%v)", r.Area, r.LimitType, err)
		return nil, err
	}
	return ToModelRule(res), nil
}
