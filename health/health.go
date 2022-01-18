package health

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Health struct {
	// key: service name, type is string
	// value: service Status type is Status
	statusMap sync.Map
	// key: service name, type is string
	// value: update, type is sync.Map
	// update: key is id, type is string, value: chan Status
	updates sync.Map
	ticker  *time.Ticker
}

func New(opts ...Option) *Health {
	option := options{
		watchTime: time.Second * 5,
	}
	for _, o := range opts {
		o(&option)
	}
	h := &Health{
		ticker: time.NewTicker(option.watchTime),
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
	u, _ := h.updates.LoadOrStore(service, sync.Map{})
	update := u.(sync.Map)

	update.Range(func(key, value interface{}) bool {
		ch := value.(chan Status)
		eg.Go(func() error {
			select {
			case ch <- status:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
		return true
	})

	return eg.Wait()
}

func (h *Health) Update(service string, id string) chan Status {
	u, _ := h.updates.LoadOrStore(service, sync.Map{})
	update := u.(sync.Map)
	ch, _ := update.LoadOrStore(id, make(chan Status, 1))
	return ch.(chan Status)
}

func (h *Health) DelUpdate(service string, id string) {
	if u, ok := h.updates.Load(service); ok {
		update := u.(sync.Map)
		update.Delete(id)
	}
}

func (h *Health) Ticker() *time.Ticker {
	return h.ticker
}
