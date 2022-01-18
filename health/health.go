package health

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Health struct {
	mutex sync.RWMutex
	// key: service name, type is string
	// value: service Status type is Status
	statusMap sync.Map
	updates   map[string]map[string]chan Status
	ticker    *time.Ticker
}

func New(opts ...Option) *Health {
	option := options{
		watchTime: time.Second * 5,
	}
	for _, o := range opts {
		o(&option)
	}
	h := &Health{
		ticker:  time.NewTicker(option.watchTime),
		updates: map[string]map[string]chan Status{},
	}
	return h
}

func (h *Health) GetStatus(service string) (Status, bool) {
	status, ok := h.statusMap.Load(service)
	if !ok {
		return Status_UNKNOWN, false
	}
	return status.(Status), ok
}

func (h *Health) SetStatus(service string, status Status) error {
	h.statusMap.Store(service, status)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	h.mutex.RLock()
	for _, w := range h.updates[service] {
		ch := w
		eg.Go(func() error {
			select {
			case ch <- status:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}
	h.mutex.RUnlock()

	return eg.Wait()
}

func (h *Health) Update(service string, id string) chan Status {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if _, ok := h.updates[service]; !ok {
		h.updates[service] = make(map[string]chan Status)
	}
	if _, ok := h.updates[service][id]; !ok {
		h.updates[service][id] = make(chan Status, 1)
	}
	return h.updates[service][id]
}

func (h *Health) DelUpdate(service string, id string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if _, ok := h.updates[service]; ok {
		delete(h.updates[service], id)
	}
}

func (h *Health) Ticker() *time.Ticker {
	return h.ticker
}
