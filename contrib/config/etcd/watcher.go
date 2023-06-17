package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	source *source
	ch     clientv3.WatchChan

	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(s *source) *watcher {
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		source: s,
		ctx:    ctx,
		cancel: cancel,
	}

	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	w.ch = s.client.Watch(s.options.ctx, s.options.path, opts...)

	return w
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case resp := <-w.ch:
		if resp.Err() != nil {
			return nil, resp.Err()
		}
		return w.source.Load()
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}
