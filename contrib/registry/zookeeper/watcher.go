package zookeeper

import (
	"context"
	"errors"
	"path"
	"sync/atomic"

	"github.com/go-zookeeper/zk"

	"github.com/go-kratos/kratos/v2/registry"
)

var _ registry.Watcher = (*watcher)(nil)

var ErrWatcherStopped = errors.New("watcher stopped")

type watcher struct {
	ctx    context.Context
	event  chan zk.Event
	conn   *zk.Conn
	cancel context.CancelFunc

	first uint32
	// prefix for ZooKeeper paths or keys (used for filtering or identifying watched nodes)
	prefix string
	// the name of the service being watched in ZooKeeper
	serviceName string
}

func newWatcher(ctx context.Context, prefix, serviceName string, conn *zk.Conn) (*watcher, error) {
	w := &watcher{conn: conn, event: make(chan zk.Event, 1), prefix: prefix, serviceName: serviceName}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.watch(w.ctx)
	return w, nil
}

func (w *watcher) watch(ctx context.Context) {
	for {
		// since a single watch is only valid for one event, we need to loop to continue watching
		_, _, ch, err := w.conn.ChildrenW(w.prefix)
		if err != nil {
			// If the target service node has not been created
			if errors.Is(err, zk.ErrNoNode) {
				// Add watcher for the node exists
				_, _, ch, err = w.conn.ExistsW(w.prefix)
			}
			if err != nil {
				w.event <- zk.Event{Err: err}
				continue
			}
		}
		select {
		case <-ctx.Done():
			return
		case ev := <-ch:
			w.event <- ev
		}
	}
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	// TODO: multiple calls to Next may lead to inconsistent service instance information
	if atomic.CompareAndSwapUint32(&w.first, 0, 1) {
		return w.getServices()
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
		return w.getServices()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}

func (w *watcher) getServices() ([]*registry.ServiceInstance, error) {
	servicesID, _, err := w.conn.Children(w.prefix)
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(servicesID))
	for _, id := range servicesID {
		servicePath := path.Join(w.prefix, id)
		b, _, err := w.conn.Get(servicePath)
		if err != nil {
			return nil, err
		}
		item, err := unmarshal(b)
		if err != nil {
			return nil, err
		}

		// if the service name of the retrieved instance does not match the watcher's service name, skip it
		if item.Name != w.serviceName {
			continue
		}

		items = append(items, item)
	}
	return items, nil
}
