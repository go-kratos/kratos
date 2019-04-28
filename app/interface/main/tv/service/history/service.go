package history

import (
	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/dao/archive"
	"go-common/app/interface/main/tv/dao/cms"
	"go-common/app/interface/main/tv/dao/history"
)

// Service .
type Service struct {
	conf   *conf.Config
	dao    *history.Dao
	cmsDao *cms.Dao
	arcDao *archive.Dao
}

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:   c,
		dao:    history.New(c),
		cmsDao: cms.New(c),
		arcDao: archive.New(c),
	}
	return srv
}
