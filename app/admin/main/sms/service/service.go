package service

import (
	"context"
	"sync"

	"go-common/app/admin/main/sms/conf"
	"go-common/app/admin/main/sms/dao"

	"github.com/jinzhu/gorm"
)

// Service is service.
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	wg     sync.WaitGroup
	db     *gorm.DB
	closed bool
}

// New is workflow-admin service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.db = s.dao.DB
	return s
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// Close consumer close.
func (s *Service) Close() {
	s.closed = true
	s.dao.Close()
	s.wg.Wait()
}
