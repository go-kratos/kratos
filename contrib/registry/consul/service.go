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

func (s *serviceSet) broadcast(cluster string, ss []*registry.ServiceInstance) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.services.Load() == nil {
		s.services.Store(map[string][]*registry.ServiceInstance{cluster: ss})
	} else {
		ms := s.services.Load().(map[string][]*registry.ServiceInstance)
		ms[cluster] = ss
		s.services.Store(ms)
	}

	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
