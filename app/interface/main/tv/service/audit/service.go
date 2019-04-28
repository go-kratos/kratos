package audit

import (
	"go-common/app/interface/main/tv/conf"
	auditDao "go-common/app/interface/main/tv/dao/audit"
	"go-common/app/interface/main/tv/dao/cms"
)

// Service .
type Service struct {
	conf     *conf.Config
	auditDao *auditDao.Dao
	cmsDao   *cms.Dao
}

// New .
func New(c *conf.Config) *Service {
	srv := &Service{
		conf:     c,
		auditDao: auditDao.New(c),
		cmsDao:   cms.New(c),
	}
	return srv
}
