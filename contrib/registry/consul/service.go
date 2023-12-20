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

func (s *serviceSet) broadcast(ss map[string][]*registry.ServiceInstance) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if ss == nil {
		return
	}

	if s.services == nil {
		s.services = ss
	} else {
		for cluster, service := range ss {
			s.services[cluster] = service
		}
	}

	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}

// getInstances get service instances
func (s *serviceSet) getInstances() []*registry.ServiceInstance {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var services []*registry.ServiceInstance
	for _, instances := range s.services {
		services = append(services, instances...)
	}
	return services
}

// getInstances get service instances
func (s *serviceSet) getInstancesMap(cluster string) map[string][]*registry.ServiceInstance {
	s.lock.RLock()
	defer s.lock.RUnlock()

	tmp := make(map[string][]*registry.ServiceInstance)
	for k, instances := range s.services {
		if k == cluster {
			tmp[k] = instances
		}
	}

	return tmp
}

func (s *serviceSet) addWatcher(w *watcher) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.watcher[w] = struct{}{}
}

func (s *serviceSet) delWatcher(w *watcher) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.watcher, w)
}
