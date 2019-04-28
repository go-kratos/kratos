package relation

import (
	"go-common/app/interface/main/app-interface/conf"
	accdao "go-common/app/interface/main/app-interface/dao/account"
	livedao "go-common/app/interface/main/app-interface/dao/live"
	reldao "go-common/app/interface/main/app-interface/dao/relation"
)

// Service is favorite.
type Service struct {
	c *conf.Config
	// dao
	accDao  *accdao.Dao
	relDao  *reldao.Dao
	liveDao *livedao.Dao
}

// New new favorite。
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		accDao:  accdao.New(c),
		relDao:  reldao.New(c),
		liveDao: livedao.New(c),
	}
	return s
}
