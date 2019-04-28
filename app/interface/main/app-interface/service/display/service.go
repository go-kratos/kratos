package display

import (
	"go-common/app/interface/main/app-interface/conf"
	locdao "go-common/app/interface/main/app-interface/dao/location"
)

// Service is zone service.
type Service struct {
	// ip
	loc *locdao.Dao
}

// New initial display service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		loc: locdao.New(c),
	}
	return
}
