package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/registry"
)

var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	key         string
	ctx         context.Context
	cancel      context.CancelFunc
	client      *clientv3.Client
	watchChan   clientv3.WatchChan
	watcher     clientv3.Watcher
	kv          clientv3.KV
	first       bool
	serviceName string
}

func newWatcher(ctx context.Context, key, name string, client *clientv3.Client) (*watcher, error) {
	w := &watcher{
		key:         key,
		client:      client,
		watcher:     clientv3.NewWatcher(client),
		kv:          clientv3.NewKV(client),
		first:       true,
		serviceName: name,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0), clientv3.WithKeysOnly())
	err := w.watcher.RequestProgress(w.ctx)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if w.first {
		item, err := w.getInstance()
		w.first = false
		return item, err
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case watchResp, ok := <-w.watchChan:
		if !ok || watchResp.Err() != nil {
			time.Sleep(time.Second)
			err := w.reWatch()
			if err != nil {
				return nil, err
			}
		}
		return w.getInstance()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *watcher) getInstance() ([]*registry.ServiceInstance, error) {
	resp, err := w.kv.Get(w.ctx, w.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != w.serviceName {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (w *watcher) reWatch() error {
	w.watcher.Close()
	w.watcher = clientv3.NewWatcher(w.client)
	w.watchChan = w.watcher.Watch(w.ctx, w.key, clientv3.WithPrefix(), clientv3.WithRev(0), clientv3.WithKeysOnly())
	return w.watcher.RequestProgress(w.ctx)
}
