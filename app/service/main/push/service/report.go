package service

import (
	"context"

	"go-common/app/service/main/push/model"
	"go-common/library/log"
)

// AddReport add report.
func (s *Service) AddReport(ctx context.Context, r *model.Report) (err error) {
	var old *model.Report
	if old, err = s.dao.Report(ctx, r.DeviceToken); err != nil {
		return
	}
	if old != nil {
		r.ID = old.ID
		if err = s.dao.UpdateReport(ctx, r); err != nil {
			return
		}
	} else {
		if r.ID, err = s.dao.AddReport(ctx, r); err != nil {
			return
		}
	}
	if r.NotifySwitch == model.SwitchOn {
		s.reportCache.Save(func() { s.dao.AddTokenCache(context.Background(), r.DeviceToken, r) })
		if r.Mid > 0 {
			s.reportCache.Save(func() { s.AddReportCache(context.Background(), r) })
		}
	}
	if old != nil && old.Dtime == 0 && old.Mid > 0 && (old.Mid != r.Mid || r.NotifySwitch == model.SwitchOff) {
		s.reportCache.Save(func() { s.dao.DelReportCache(context.TODO(), old.Mid, old.APPID, old.DeviceToken) })
	}
	return
}

// AddReportCache add report cache.
func (s *Service) AddReportCache(c context.Context, r *model.Report) (err error) {
	res, err := s.dao.ReportsCacheByMid(c, r.Mid)
	if err != nil {
		return
	}
	if len(res) == 0 {
		if res, err = s.dao.ReportsByMid(c, r.Mid); err != nil {
			return
		}
	}
	if len(res) == 0 {
		return
	}
	m := make(map[string]*model.Report)
	for _, v := range res {
		m[v.DeviceToken] = v
	}
	m[r.DeviceToken] = r
	var rs []*model.Report
	for _, v := range m {
		rs = append(rs, v)
	}
	mrs := map[int64][]*model.Report{r.Mid: rs}
	return s.dao.AddReportsCacheByMids(c, mrs)
}

// AddUserReportCache add user report cache.
func (s *Service) AddUserReportCache(c context.Context, mid int64, rs []*model.Report) (err error) {
	mrs := map[int64][]*model.Report{mid: rs}
	return s.dao.AddReportsCacheByMids(c, mrs)
}

// AddTokenCache add token cache.
func (s *Service) AddTokenCache(ctx context.Context, r *model.Report) (err error) {
	err = s.dao.AddTokenCache(ctx, r.DeviceToken, r)
	return
}

// AddTokensCache add token cache.
func (s *Service) AddTokensCache(ctx context.Context, rs map[string]*model.Report) (err error) {
	for k := range rs {
		log.Info("AddTokensCache token(%s)", k)
	}
	err = s.dao.AddTokensCache(ctx, rs)
	return
}

// DelReport delelte report & its cache.
func (s *Service) DelReport(c context.Context, appID, mid int64, token string) (err error) {
	if _, err = s.dao.DelReport(c, token); err != nil {
		log.Error("delete token(%s) error(%v)", token, err)
		return
	}
	log.Info("delete report app(%d) mid(%d) token(%s)", appID, mid, token)
	s.reportCache.Save(func() {
		if err = s.dao.DelTokenCache(context.Background(), token); err != nil {
			log.Error("s.dao.DelTokeCache(%s) error(%v)", token, err)
		}
		if mid > 0 {
			if err = s.dao.DelReportCache(context.Background(), mid, appID, token); err != nil {
				log.Error("s.dao.DelReportCache(%d,%d,%s) error(%v)", mid, appID, token, err)
			}
		}
	})
	return
}

// DelInvalidReports deletes invalid reports.
func (s *Service) DelInvalidReports(c context.Context, tp int) (err error) {
	switch tp {
	case model.DelMiFeedback:
		s.reportCache.Save(func() { s.dao.DelMiInvalid(context.TODO()) })
	case model.DelMiUninstalled:
		s.cache.Save(func() { s.dao.DelMiUninstalled(context.TODO()) })
	default:
		log.Error("delete invalid reports type error(%d)", tp)
	}
	return
}
