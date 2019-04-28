package v1

import (
	"context"
	"sync"

	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/dao/v1"
	"go-common/app/infra/config/model"
	xtime "go-common/library/time"
)

// Service service.
type Service struct {
	dao         *v1.Dao
	vLock       sync.RWMutex
	versions    map[string]int64 // serviceName_buildVersion > configVersion
	eLock       sync.RWMutex
	events      map[string]chan *model.Version
	PollTimeout xtime.Duration
	token       map[string]string
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = new(Service)
	s.dao = v1.New(c)
	s.versions = make(map[string]int64)
	s.events = make(map[string]chan *model.Version)
	s.PollTimeout = c.PollTimeout
	s.token = make(map[string]string)
	return
}

// Ping check is ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close resources.
func (s *Service) Close() {
	s.dao.Close()
}
