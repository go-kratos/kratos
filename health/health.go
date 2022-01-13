package health

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Health struct {
	statusMap map[string]Status
	mutex     sync.RWMutex
	opts      options
	watchers  map[string]map[string]chan Status
}

type Option func(*options)

type options struct{}

func New(opts ...Option) *Health {
	h := &Health{
		statusMap: make(map[string]Status),
		mutex:     sync.RWMutex{},
		watchers:  make(map[string]map[string]chan Status),
	}
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	_ = h.opts
	return h
}

func (h *Health) SetStatus(service string, status Status) {
	h.mutex.Lock()
	h.statusMap[service] = status
	h.mutex.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	h.mutex.RLock()
	for _, w := range h.watchers[service] {
		ch := w
		eg.Go(func() error {
			select {
			case ch <- status:
			case <-ctx.Done():
			}
			return nil
		})
	}
	h.mutex.RUnlock()
	_ = eg.Wait()
}

func (h *Health) GetStatus(service string) (status Status, ok bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	status, ok = h.statusMap[service]
	return status, ok
}

func (h *Health) Watch(service string, id string) (ch chan Status) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if _, ok := h.watchers[service]; !ok {
		h.watchers[service] = make(map[string]chan Status)
	}
	if _, ok := h.watchers[service][id]; !ok {
		h.watchers[service][id] = make(chan Status, 1)
	}
	return h.watchers[service][id]
}
