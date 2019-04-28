package service

import (
	"context"
	"regexp"

	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"

	"go-common/library/log"
)

var (
	ignoreRegexpsCondFn = func(r *model.Regexp) bool {
		return r.State != model.StateDeleted &&
			r.Operation == model.OperationIgnore
	}
	whiteRegexpsCondFn = func(r *model.Regexp) bool {
		return r.State != model.StateDeleted &&
			r.Operation == model.OperationPutToWhiteList
	}
	limitRegexpsCondFn = func(r *model.Regexp) bool {
		return r.State != model.StateDeleted &&
			(r.Operation == model.OperationLimit ||
				r.Operation == model.OperationRestrictLimit)
	}
)

var (
	regexps = []*model.Regexp{}
)

// RefreshRegexps .
func (s *SvcImpl) RefreshRegexps(ctx context.Context) {
	s.Lock()
	defer s.Unlock()

	dbRs, _, err := s.RegexpDao.GetByCond(ctx, ToDaoCond(&Condition{
		State: model.StateDefault,
	}))
	if err != nil || len(dbRs) == 0 {
		return
	}
	rs := ToModelRegexps(dbRs)
	for i, r := range rs {
		reg, err := regexp.Compile(r.Content)
		if err != nil {
			log.Warn("%v", err)
			rs[i] = nil
			continue
		}
		r.Reg = reg
	}

	regexps = regexps[:0]
	for _, r := range rs {
		if r != nil {
			regexps = append(regexps, r)
		}
	}
}

// GetRegexpByID .
func (s *SvcImpl) GetRegexpByID(ctx context.Context, id int64) (*model.Regexp, error) {
	r, err := s.RegexpDao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToModelRegexp(r), nil
}

// GetRegexpsByAreaAndCondFunc .
func (s *SvcImpl) GetRegexpsByAreaAndCondFunc(ctx context.Context, area string, condFunc func(*model.Regexp) bool) (rs []*model.Regexp) {
	if condFunc == nil {
		return
	}

	s.RLock()
	defer s.RUnlock()
	for _, r := range regexps {
		if condFunc(r) &&
			// mainsite_dm has its own regexps
			((area == model.AreaMainSiteDM && r.Area == area) ||
				// all other business areas (imessage, reply, live_dm, ...etc)
				// share the area reply's regexps
				(area != model.AreaMainSiteDM && r.Area == model.AreaReply)) {
			rs = append(rs, r)
		}
	}
	return
}

// GetRegexpByAreaAndContent query regexps by area and content
func (s *SvcImpl) GetRegexpByAreaAndContent(ctx context.Context, area, content string) (*model.Regexp, error) {
	r, err := s.RegexpDao.GetByAreaAndContent(ctx, ToDaoCond(&Condition{
		Area:     area,
		Contents: []string{content},
	}))
	if err != nil {
		return nil, err
	}
	return ToModelRegexp(r), nil
}

// GetRegexpsByCond query regexps by condition
func (s *SvcImpl) GetRegexpsByCond(ctx context.Context,
	cond *Condition) (rs []*model.Regexp, total int64, err error) {
	dbRs, total, err := s.RegexpDao.GetByCond(ctx, ToDaoCond(cond))
	if err == dao.ErrResourceNotExist {
		return []*model.Regexp{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return ToModelRegexps(dbRs), total, nil
}

// UpsertRegexp update or insert regexp
func (s *SvcImpl) UpsertRegexp(ctx context.Context, r *model.Regexp) (*model.Regexp, error) {
	var res *dao.Regexp
	var err error
	if r.ID > 0 {
		res, err = s.RegexpDao.Update(ctx, ToDaoRegexp(r))
	} else {
		res, err = s.RegexpDao.Insert(ctx, ToDaoRegexp(r))
	}
	if err != nil {
		return nil, err
	}
	return ToModelRegexp(res), nil
}

// DeleteRegexp delete regexp by id and adminID
// and delete cache
func (s *SvcImpl) DeleteRegexp(ctx context.Context, id, adminID int64) (*model.Regexp, error) {
	r, err := s.GetRegexpByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.State == model.StateDeleted {
		return r, nil
	}
	r.AdminID, r.State = adminID, model.StateDeleted
	if err = s.antiDao.DelRegexpCache(ctx); err != nil {
		log.Error("s.antiDao.DelRegexpCache() error(%v)", err)
		return nil, err
	}
	res, err := s.RegexpDao.Update(ctx, ToDaoRegexp(r))
	if err != nil {
		return nil, err
	}
	return ToModelRegexp(res), nil
}
