package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Status string

const (
	Up   Status = "UP"
	Down Status = "DOWN"
)

type Checker interface {
	Check(ctx context.Context) error
}

type CheckerFunc func(ctx context.Context) error

func (f CheckerFunc) Check(ctx context.Context) error {
	return f(ctx)
}

type Health struct {
	ctx    context.Context
	done   context.CancelFunc
	ticker *time.Ticker
	lock   sync.RWMutex

	status   Status
	checkers map[string]Checker
	errors   map[string]error

	// can configurable
	timeout      time.Duration
	intervalTime time.Duration

	startOnce sync.Once
	stopOnce  sync.Once
}

func New(opts ...Option) *Health {
	h := &Health{
		lock:         sync.RWMutex{},
		status:       Down,
		checkers:     make(map[string]Checker),
		errors:       make(map[string]error),
		timeout:      time.Second * 2,
		intervalTime: time.Second * 10,
		startOnce:    sync.Once{},
		stopOnce:     sync.Once{},
	}

	for _, opt := range opts {
		opt(h)
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.ctx = ctx
	h.done = cancel
	h.ticker = time.NewTicker(h.intervalTime)
	return h
}

func (h *Health) CheckAll(_ context.Context) Result {
	r := Result{Status: h.status}
	h.lock.RLock()
	defer h.lock.RUnlock()
	for c, e := range h.errors {
		d := Detail{}
		if e == nil {
			d.Status = Up
		} else {
			d.Status = Down
			d.Error = e.Error()
		}
		r.Details[c] = d
	}
	return r
}

func (h *Health) Check(ctx context.Context, service string) error {
	h.lock.RLock()
	defer h.lock.RUnlock()
	err, ok := h.errors[service]
	if !ok {
		return ErrServiceNotFind
	}
	return err
}

func (h *Health) Start(ctx context.Context) error {
	h.startOnce.Do(func() {
		h.status = Up
		h.check()
		for {
			select {
			case <-h.ctx.Done():
				return
			case <-h.ticker.C:
				h.check()
			}
		}
	})
	return nil
}

func (h *Health) Stop(ctx context.Context) error {
	h.stopOnce.Do(func() {
		h.done()
		h.ticker.Stop()
		h.status = Down
	})
	return nil
}

func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	if service == "" {
		res := h.CheckAll(r.Context())
		if res.Status == Down {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		_ = json.NewEncoder(w).Encode(res)
	} else {
		d := &Detail{Status: Up}
		code := http.StatusOK
		err := h.Check(r.Context(), service)
		if err != nil {
			code = http.StatusInternalServerError
			d.Status = Down
			d.Error = err.Error()
		}
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(d)
	}
}

func (h *Health) check() {
	ctx, cancel := context.WithTimeout(h.ctx, h.timeout)
	defer cancel()

	var wg sync.WaitGroup

	for component, checker := range h.checkers {
		wg.Add(1)
		go func(ctx context.Context, component string, checker Checker) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				err := checker.Check(ctx)
				h.lock.Lock()
				h.errors[component] = err
				h.lock.Unlock()
			}
		}(ctx, component, checker)
	}
	wg.Wait()
}

type Result struct {
	Status  Status            `json:"status"`
	Details map[string]Detail `json:"details"`
}

type Detail struct {
	Status Status `json:"status"`
	Error  string `json:"error"`
}
