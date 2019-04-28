package notice

import (
	"context"
	"time"

	"go-common/app/admin/main/app/conf"
	noticedao "go-common/app/admin/main/app/dao/notice"
	"go-common/app/admin/main/app/model/notice"
	"go-common/library/log"
)

// Service notice service.
type Service struct {
	dao *noticedao.Dao
}

// New new a notice service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: noticedao.New(c),
	}
	return
}

// Notices select All
func (s *Service) Notices(c context.Context) (res []*notice.Notice, err error) {
	if res, err = s.dao.Notices(c); err != nil {
		log.Error("s.dao.Notices error(%v)", err)
		return
	}
	return
}

// NoticeByID by id
func (s *Service) NoticeByID(c context.Context, id int64) (res *notice.Notice, err error) {
	if res, err = s.dao.NoticeByID(c, id); err != nil {
		log.Error("s.dao.NoticeByID error(%v)", err)
		return
	}
	return
}

// UpdateNotice update notice
func (s *Service) UpdateNotice(c context.Context, a *notice.Param, now time.Time) (err error) {
	if err = s.dao.Update(c, a, now); err != nil {
		log.Error("s.dao.Update(%v)", err)
		return
	}
	return
}

// UpdateBuild update notice
func (s *Service) UpdateBuild(c context.Context, a *notice.Param, now time.Time) (err error) {
	if err = s.dao.UpdateBuild(c, a, now); err != nil {
		log.Error("s.dao.UpdateBuild(%v)", err)
		return
	}
	return
}

// UpdateState update state
func (s *Service) UpdateState(c context.Context, a *notice.Param, now time.Time) (err error) {
	if err = s.dao.UpdateState(c, a, now); err != nil {
		log.Error("s.dao.UpdateState(%v)", err)
		return
	}
	return
}

// Insert insert
func (s *Service) Insert(c context.Context, a *notice.Param, now time.Time) (err error) {
	if err = s.dao.Insert(c, a, now); err != nil {
		log.Error("s.dao.Insert error(%v)", err)
		return
	}
	return
}
