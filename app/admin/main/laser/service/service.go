package service

import (
	"context"

	"go-common/app/admin/main/laser/conf"
	"go-common/app/admin/main/laser/dao"
)

// Service struct
type Service struct {
	conf *conf.Config
	dao  *dao.Dao
}

// New is new instance
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  dao.New(c),
	}
	return
}

// Ping is check dao connected
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close is close dao connection
func (s *Service) Close() (err error) {
	return s.dao.Close(context.TODO())
}
