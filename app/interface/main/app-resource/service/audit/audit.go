package audit

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	auditdao "go-common/app/interface/main/app-resource/dao/audit"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Service audit service.
type Service struct {
	dao *auditdao.Dao
	// tick
	tick time.Duration
	// cache
	auditCache map[string]map[int]struct{}
}

// New new a audit service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: auditdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// cache
		auditCache: map[string]map[int]struct{}{},
	}
	s.loadAuditCache()
	go s.cacheproc()
	return
}

// Audit
func (s *Service) Audit(c context.Context, mobiApp string, build int) (err error) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			return ecode.OK
		}
	}
	return ecode.NotModified
}

// cacheproc load all cache.
func (s *Service) cacheproc() {
	for {
		time.Sleep(s.tick)
		s.loadAuditCache()
	}
}

func (s *Service) loadAuditCache() {
	as, err := s.dao.Audits(context.TODO())
	if err != nil {
		log.Error("s.dao.Audits error(%v)", err)
		return
	}
	s.auditCache = as
}
