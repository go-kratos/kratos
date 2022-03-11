package etcd

import (
	"github.com/SeeMusic/kratos/v2/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	source    *source
	ch        clientv3.WatchChan
	closeChan chan struct{}
}

func newWatcher(s *source) *watcher {
	w := &watcher{
		source:    s,
		closeChan: make(chan struct{}),
	}

	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	w.ch = s.client.Watch(s.options.ctx, s.options.path, opts...)

	return w
}

func (s *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case _, ok := <-s.ch:
		if !ok {
			return nil, nil
		}
		return s.source.Load()
	case <-s.closeChan:
		return nil, nil
	}
}

func (s *watcher) Stop() error {
	close(s.closeChan)
	return nil
}
