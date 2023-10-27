package consul

import (
	"sync"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    *atomic.Value
	lock        sync.RWMutex
}

func (s *serviceSet) broadcast(ss map[string][]*registry.ServiceInstance) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if ss == nil {
		return
	}
	
	if s.services.Load() == nil {
		s.services.Store(ss)
	} else {
		for cluster, service := range ss {
			ms := s.services.Load().(map[string][]*registry.ServiceInstance)
			ms[cluster] = service
			s.services.Store(ms)
		}
	}

	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
