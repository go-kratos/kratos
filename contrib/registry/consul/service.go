package consul

import (
	"sync"

	"github.com/go-kratos/kratos/v2/registry"
)

type serviceSet struct {
	serviceName string
	watcher     map[*watcher]struct{}
	services    map[string][]*registry.ServiceInstance
	lock        sync.RWMutex
}

// Please lock outside
func (s *serviceSet) flatServices() []*registry.ServiceInstance {
	var ss []*registry.ServiceInstance
	for _, v := range s.services {
		ss = append(ss, v...)
	}
	return ss
}

func (s *serviceSet) broadcast(ss map[string][]*registry.ServiceInstance) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.services = ss
	for k := range s.watcher {
		// non-blocking send
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
