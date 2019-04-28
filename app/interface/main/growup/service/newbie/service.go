package newbie

import (
	"context"

	"go-common/app/interface/main/growup/conf"
	"go-common/app/interface/main/growup/dao/newbiedao"
)

// Service is growup service
type Service struct {
	conf *conf.Config
	dao  *newbiedao.Dao
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		dao:  newbiedao.New(c),
	}
	return s
}

// Ping fn
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
