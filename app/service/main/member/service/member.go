package service

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

var (
	_emptyBaseExps = make(map[int64]*model.BaseExp)
	// _emptyDetails  = make(map[int64]*model.OldMember)
)

func (s *Service) baseExp(c context.Context, mid int64) (m *model.BaseExp, err error) {
	var (
		b *model.BaseInfo
		l *model.LevelInfo
	)
	eg, _ := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		if b, err = s.BaseInfo(c, mid); err != nil {
			log.Error("s.BaseInfo mid(%d), err(%v)", mid, err)
		}
		return
	})
	eg.Go(func() (err error) {
		if l, err = s.Exp(c, mid); err != nil {
			log.Error("s.Exp mid(%d), err(%v)", mid, err)
		}
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}
	if b == nil {
		return
	}
	return &model.BaseExp{BaseInfo: b, LevelInfo: l}, err
}

func (s *Service) baseExps(c context.Context, mids []int64) (res map[int64]*model.BaseExp, err error) {
	if len(mids) > 100 {
		err = ecode.MemberOverLimit
		return
	}
	g, _ := errgroup.WithContext(c)
	var bs map[int64]*model.BaseInfo
	g.Go(func() (err error) {
		if bs, err = s.BatchBaseInfo(c, mids); err != nil {
			log.Error("s.BatchBaseInfo mids(%v), err(%v)", mids, err)
		}
		return
	})
	var lvs map[int64]*model.LevelInfo
	g.Go(func() (err error) {
		if lvs, err = s.Exps(c, mids); err != nil {
			log.Error("s.Exps mids(%v), err(%v)", mids, err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	if len(bs) == 0 {
		res = _emptyBaseExps
		return
	}
	res = make(map[int64]*model.BaseExp, len(bs))
	for _, b := range bs {
		mem := new(model.BaseExp)
		mem.BaseInfo = b
		if l, ok := lvs[b.Mid]; ok {
			mem.LevelInfo = l
		}
		res[b.Mid] = mem
	}
	return
}

// Member get the full information within member-service.
func (s *Service) Member(c context.Context, mid int64) (m *model.Member, err error) {
	var (
		b *model.BaseExp
		o *model.OfficialInfo
	)
	if b, err = s.baseExp(c, mid); err != nil {
		log.Error("s.BaseExp mid(%d), err(%v)", mid, err)
		return
	}
	if s.officials != nil {
		o = s.officials[mid]
	}
	return &model.Member{BaseInfo: b.BaseInfo, LevelInfo: b.LevelInfo, OfficialInfo: o}, err
}

// Members get the full information within member-service from a set of mids.
func (s *Service) Members(c context.Context, mids []int64) (res map[int64]*model.Member, err error) {
	if len(mids) > 100 {
		err = ecode.MemberOverLimit
		return
	}
	var (
		bs map[int64]*model.BaseExp
	)
	if bs, err = s.baseExps(c, mids); err != nil {
		log.Error("s.BaseExps mids(%v), err(%v)", mids, err)
		return
	}
	res = map[int64]*model.Member{}
	for _, b := range bs {
		bo := &model.Member{BaseInfo: b.BaseInfo, LevelInfo: b.LevelInfo}
		if s.officials != nil {
			bo.OfficialInfo = s.officials[b.Mid]
		}
		res[b.Mid] = bo
	}
	return
}

// SetOfficialDoc is.
func (s *Service) SetOfficialDoc(ctx context.Context, aod *model.ArgOfficialDoc) error {
	od := &model.OfficialDoc{
		Mid:          aod.Mid,
		Name:         aod.Name,
		Role:         aod.Role,
		Title:        aod.Title,
		Desc:         aod.Desc,
		SubmitSource: aod.SubmitSource,
		SubmitTime:   xtime.Time(time.Now().Unix()),

		OfficialExtra: model.OfficialExtra{
			Realname:          aod.Realname,
			Operator:          aod.Operator,
			Telephone:         aod.Telephone,
			Email:             aod.Email,
			Address:           aod.Address,
			Company:           aod.Company,
			CreditCode:        aod.CreditCode,
			Organization:      aod.Organization,
			OrganizationType:  aod.OrganizationType,
			BusinessLicense:   aod.BusinessLicense,
			BusinessScale:     aod.BusinessScale,
			BusinessLevel:     aod.BusinessLevel,
			BusinessAuth:      aod.BusinessAuth,
			Supplement:        aod.Supplement,
			Professional:      aod.Professional,
			Identification:    aod.Identification,
			OfficialSite:      aod.OfficialSite,
			RegisteredCapital: aod.RegisteredCapital,
		},
	}

	if !od.Validate() {
		log.Error("Failed to validate official doc: %+v", od)
		return ecode.RequestErr
	}
	if err := s.mbDao.SetOfficialDoc(ctx, od); err != nil {
		log.Error("Failed to set official doc: %+v", err)
		return ecode.SubmitOfficialDocFailed
	}

	if od.CreditCode != "" {
		if err := s.mbDao.SetOfficialDocAddit(ctx, od.Mid, "credit_code", od.CreditCode); err != nil {
			log.Error("Failed to set official doc addit, mid(%d), error(%+v)", od.Mid, err)
		}
	}
	return nil
}

// OfficialDoc is.
func (s *Service) OfficialDoc(c context.Context, mid int64) (*model.OfficialDoc, error) {
	od, err := s.mbDao.OfficialDoc(c, mid)
	if errors.Cause(err) == sql.ErrNoRows {
		return nil, ecode.NoOfficialDoc
	}
	return od, nil
}

// DelCache del card and info cache.
func (s *Service) DelCache(c context.Context, mid int64, action string, ak string, sd string) (err error) {
	return
}
