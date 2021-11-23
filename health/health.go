package health

import (
	"sync"
)

type Status int8

const (
	StatusUnknown Status = iota
	StatusServing
	StatusNotServing
	StatusServiceUnknown
)

type Health struct {
	status Status
	mutex  sync.RWMutex
	opts   options
	watchers []func(status Status)
}
type Option func(*options)

type options struct {
}

func New(opts ...Option) *Health {
	h := &Health{
		status: StatusUnknown,
		mutex:  sync.RWMutex{},
	}
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	return h
}

func (h *Health) SetStatus(status Status) {
	h.mutex.Lock()
	h.status = status
	h.mutex.Unlock()
	go func() {
		for _, w := range h.watchers {
			w(status)
		}
	}()
}

func (h *Health) GetStatus() Status {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.status
}

func (h *Health) Watch(f func(status Status)) {
	h.watchers = append(h.watchers, f)
}