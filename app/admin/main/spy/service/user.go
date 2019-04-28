package service

import (
	"context"

	"go-common/app/admin/main/spy/model"
	"go-common/library/log"
)

// UserInfo get UserInfo by mid , from cache or db or generate.
func (s *Service) UserInfo(c context.Context, mid int64) (u *model.UserInfoDto, err error) {
	var (
		udb *model.UserInfo
	)
	if udb, err = s.spyDao.Info(c, mid); err != nil {
		log.Error("s.spyDao.Info(%d) error(%v)", mid, err)
		return
	}

	// init user score by rpc call
	if udb == nil {
		if _, err = s.spyDao.UserScore(c, mid); err != nil {
			log.Error("s.spyDao.UserScore(%d) error(%v)", mid, err)
			return
		}
		if udb, err = s.spyDao.Info(c, mid); err != nil {
			log.Error("s.spyDao.Info(%d) error(%v)", mid, err)
			return
		}
	}

	if udb == nil {
		log.Error("UserInfo init failed , still nil")
		return
	}
	u = &model.UserInfoDto{
		ID:          udb.ID,
		Mid:         udb.Mid,
		Score:       udb.Score,
		BaseScore:   udb.BaseScore,
		EventScore:  udb.EventScore,
		State:       udb.State,
		ReliveTimes: udb.ReliveTimes,
		Mtime:       udb.Mtime.Unix(),
	}
	if m, err1 := s.spyDao.AccInfo(c, mid); err1 == nil && m != nil {
		u.Name = m.Name
	}
	return
}

// HisoryPage history page.
func (s *Service) HisoryPage(c context.Context, h *model.HisParamReq) (page *model.HistoryPage, err error) {
	totalCount, err := s.spyDao.HistoryPageTotalC(c, h)
	if err != nil {
		log.Error("userDao HistoryPageTotalC(%v) error(%v)", h, err)
		return
	}
	page = &model.HistoryPage{}
	items, err := s.spyDao.HistoryPage(c, h)
	if err != nil {
		log.Error("spyDao.HistoryPage(%v) error(%v)", h, err)
		return
	}
	page.TotalCount = totalCount
	page.Items = items
	page.Pn = h.Pn
	page.Ps = h.Ps
	return
}

// ResetBase reset user base score.
func (s *Service) ResetBase(c context.Context, mid int64, operator string) (err error) {
	if err = s.spyDao.ResetBase(c, mid, operator); err != nil {
		log.Error("s.spyDao.ResetBase(%d,%s) error(%v)", mid, operator, err)
		return
	}
	return
}

// Refresh reset user base score.
func (s *Service) RefreshBase(c context.Context, mid int64, operator string) (err error) {
	if err = s.spyDao.RefreshBase(c, mid, operator); err != nil {
		log.Error("s.spyDao.RefreshBase(%d,%s) error(%v)", mid, operator, err)
		return
	}
	return
}

// ResetEvent reset user event score.
func (s *Service) ResetEvent(c context.Context, mid int64, operator string) (err error) {
	if err = s.spyDao.ResetEvent(c, mid, operator); err != nil {
		log.Error("s.spyDao.ResetEvent(%d,%s) error(%v)", mid, operator, err)
		return
	}
	return
}

// ClearCount clear count.
func (s *Service) ClearCount(c context.Context, mid int64, operator string) (err error) {
	if err = s.spyDao.ClearCount(c, mid, operator); err != nil {
		log.Error("s.spyDao.ClearCount(%d, %s) error(%v)", mid, operator, err)
		return
	}
	return
}

// ReportList report list.
func (s *Service) ReportList(c context.Context, ps, pn int) (page *model.ReportPage, err error) {
	count, err := s.spyDao.ReportCount(c)
	if err != nil {
		log.Error("s.spyDao.ReportCount error(%v)", err)
		return
	}
	page = &model.ReportPage{}
	items, err := s.spyDao.ReportList(c, ps, pn)
	if err != nil {
		log.Error("s.spyDao.ReportPage(%d,%d) error(%v)", ps, pn, err)
		return
	}
	page.TotalCount = count
	page.Items = items
	page.Pn = pn
	page.Ps = ps
	return
}
