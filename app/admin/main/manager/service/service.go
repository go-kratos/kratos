package service

import (
	"context"
	"time"

	"go-common/app/admin/main/manager/conf"
	"go-common/app/admin/main/manager/dao"
	"go-common/app/admin/main/manager/model"
)

// Service biz service def.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rbac may not change frequent, can update every few seconds. only assignment must get from db.
	points        map[int64]*model.AuthItem
	pointList     []*model.AuthItem
	groupAuth     map[int64][]int64
	orgAuth       map[int64]*model.AuthOrg // group + role info
	roleAuth      map[int64][]int64
	admins        map[int64]bool
	userNames     map[int64]string // users' name
	userNicknames map[int64]string // user's nickname
	userDeps      map[int64]string // users' department
	userIds       map[string]int64 // users' ids
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.syncRbac()
	s.loadUnames()
	go s.syncRbacproc()
	go s.loadUnamesproc()
	return s
}

func (s *Service) syncRbacproc() {
	for {
		time.Sleep(time.Second * 10)
		s.syncRbac()
	}
}

func (s *Service) syncRbac() {
	if points, mpoints, err := s.ptrs(); err != nil {
		return
	} else if len(mpoints) > 0 {
		s.pointList = points
		s.points = mpoints
	}
	if admins, err := s.adms(); err != nil {
		return
	} else if len(admins) > 0 {
		s.admins = admins
	}
	if ra, err := s.roleAuths(); err != nil {
		return
	} else if len(ra) > 0 {
		s.roleAuth = ra
	}
	if ga, err := s.groupAuths(); err != nil {
		return
	} else if len(ga) > 0 {
		s.groupAuth = ga
	}
	if oa, err := s.orgAuths(); err != nil {
		return
	} else if len(oa) > 0 {
		s.orgAuth = oa
	}
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
