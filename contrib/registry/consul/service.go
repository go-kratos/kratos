package consul

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/registry"
)

type serviceSet struct {
	registry    *Registry
	serviceName string
	watcher     map[*watcher]struct{}
	ref         atomic.Int32
	services    *atomic.Value
	lock        sync.RWMutex

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *serviceSet) broadcast(ss []*registry.ServiceInstance) {
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}

func (s *serviceSet) delete(w *watcher) {
	s.lock.Lock()
	delete(s.watcher, w)
	s.lock.Unlock()
	s.registry.tryDelete(s)
}
