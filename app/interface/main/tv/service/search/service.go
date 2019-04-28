package search

import (
	"go-common/app/interface/main/tv/conf"
	arcdao "go-common/app/interface/main/tv/dao/archive"
	cmsDao "go-common/app/interface/main/tv/dao/cms"
	"go-common/app/interface/main/tv/dao/search"
)

// Service .
type Service struct {
	conf   *conf.Config
	dao    *search.Dao
	arcDao *arcdao.Dao
	cmsDao *cmsDao.Dao
}

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:   c,
		dao:    search.New(c),
		arcDao: arcdao.New(c),
		cmsDao: cmsDao.New(c),
	}
	return srv
}
