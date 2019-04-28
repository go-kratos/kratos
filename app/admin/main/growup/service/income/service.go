package income

import (
	"context"

	"go-common/app/admin/main/growup/conf"
	upD "go-common/app/admin/main/growup/dao"
	incomeD "go-common/app/admin/main/growup/dao/income"
	"go-common/app/admin/main/growup/dao/message"
)

// Service struct
type Service struct {
	conf  *conf.Config
	dao   *incomeD.Dao
	msg   *message.Dao
	upDao *upD.Dao
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:  c,
		dao:   incomeD.New(c),
		msg:   message.New(c),
		upDao: upD.New(c),
	}
	return s
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
