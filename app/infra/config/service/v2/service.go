package v2

import (
	"context"
	"sync"

	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/dao/v2"
	"go-common/app/infra/config/model"
	xtime "go-common/library/time"
)

// Service service.
type Service struct {
	dao         *v2.Dao
	PollTimeout xtime.Duration
	//config2 app
	aLock    sync.RWMutex
	bLock    sync.RWMutex
	apps     map[string]*model.App
	services map[string]*model.App
	//config2 tags
	tLock sync.RWMutex
	tags  map[string]*curTag

	eLock  sync.RWMutex
	events map[string]chan *model.Diff
	//force version
	fLock  sync.RWMutex
	forces map[string]int64

	//config2 tag force
	tfLock    sync.RWMutex
	forceType map[string]int8
	// config2 tagID
	tagIDLock sync.RWMutex
	tagID     map[string]int64

	//config2 last force version
	lfvLock   sync.RWMutex
	lfvforces map[string]int64
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = new(Service)
	s.dao = v2.New(c)
	s.PollTimeout = c.PollTimeout
	s.tags = make(map[string]*curTag)
	s.apps = make(map[string]*model.App)
	s.services = make(map[string]*model.App)
	s.events = make(map[string]chan *model.Diff)
	s.forces = make(map[string]int64)
	s.forceType = make(map[string]int8)
	s.tagID = make(map[string]int64)
	s.lfvforces = make(map[string]int64)
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
