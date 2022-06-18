package zookeeper

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
)

var ErrWatcherStopped = errors.New("watcher stopped")

type watcher struct {
	ctx    context.Context
	source *source
	event  chan zk.Event
	cancel context.CancelFunc
	first  uint32
}

func newWatcher(s *source) *watcher {
	ctx, cancel := context.WithCancel(s.options.ctx)
	w := &watcher{
		ctx:    ctx,
		source: s,
		event:  make(chan zk.Event, 1),
		cancel: cancel,
	}

	go w.watch(w.ctx)
	return w
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	// todo 如果多处调用 next 可能会导致多实例信息不同步
	if atomic.CompareAndSwapUint32(&w.first, 0, 1) {
		return w.getConfig()
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case e := <-w.event:
		if e.State == zk.StateDisconnected {
			return nil, ErrWatcherStopped
		}
		if e.Err != nil {
			return nil, e.Err
		}
		return w.getConfig()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}

func (w *watcher) watch(ctx context.Context) {
	fullPath := path.Join(w.source.options.namespace, w.source.options.key)
	for {
		// 每次 watch 只有一次有效期 所以循环 watch
		_, _, ch, err := w.source.conn.ChildrenW(fullPath)
		if err != nil {
			w.event <- zk.Event{Err: err}
		}
		select {
		case <-ctx.Done():
			return
		default:
			w.event <- <-ch
		}
	}
}

func (w *watcher) getConfig() ([]*config.KeyValue, error) {
	fullPath := path.Join(w.source.options.namespace, w.source.options.key)
	res, _, err := w.source.conn.Get(fullPath)
	if err != nil {
		return nil, err
	}

	return []*config.KeyValue{{
		Key:    w.source.options.key,
		Value:  res,
		Format: strings.TrimPrefix(filepath.Ext(w.source.options.key), "."),
	}}, nil
}
