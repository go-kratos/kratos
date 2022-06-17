package consul

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type watcher struct {
	source *source
	ch     chan interface{}
	wp     *watch.Plan

	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) handle(idx uint64, data interface{}) {
	if data == nil {
		return
	}

	_, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	w.ch <- struct{}{}
}

func newWatcher(s *source) (*watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		source: s,
		ch:     make(chan interface{}),

		ctx:    ctx,
		cancel: cancel,
	}

	wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": s.options.path})
	if err != nil {
		return nil, err
	}

	wp.Handler = w.handle
	w.wp = wp

	// wp.Run is a blocking call and will prevent newWatcher from returning
	go func() {
		err := wp.RunWithClientAndHclog(s.client, nil)
		if err != nil {
			panic(err)
		}
	}()

	return w, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ch:
		return w.source.Load()
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.wp.Stop()
	w.cancel()
	return nil
}
