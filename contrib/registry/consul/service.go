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

func (s *serviceSet) broadcast(cluster string, instances []*registry.ServiceInstance) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.services == nil {
		s.services = make(map[string][]*registry.ServiceInstance)
	}

	hasInstances := false
	if len(instances) == 0 {
		for c, ins := range s.services {
			if cluster == c {
				continue
			}
			if len(ins) != 0 {
				hasInstances = true
				delete(s.services, cluster)
			}
		}
	} else {
		hasInstances = true
		s.services[cluster] = instances
	}

	// If there is no instance, no need to notify the watcher
	if !hasInstances {
		return
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
