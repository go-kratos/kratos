package apollo

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/apolloconfig/agollo/v4/storage"
)

type watcher struct {
	out      <-chan []*config.KeyValue
	cancelFn func()
}

type customChangeListener struct {
	in     chan<- []*config.KeyValue
	apollo *apollo
}

func (c *customChangeListener) onChange(namespace string, changes map[string]*storage.ConfigChange) []*config.KeyValue {
	kv := make([]*config.KeyValue, 0, 2)
	value, err := c.apollo.client.GetConfigCache(namespace).Get("content")
	if err != nil {
		log.Warnw("apollo get config failed", "err", err)
	}
	kv = append(kv, &config.KeyValue{
		Key:    namespace,
		Value:  []byte(value.(string)),
		Format: format(namespace),
	})

	return kv
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	change := c.onChange(changeEvent.Namespace, changeEvent.Changes)
	if len(change) == 0 {
		return
	}

	c.in <- change
}

func (c *customChangeListener) OnNewestChange(changeEvent *storage.FullChangeEvent) {}

func newWatcher(a *apollo) (config.Watcher, error) {
	changeCh := make(chan []*config.KeyValue)
	listener := &customChangeListener{in: changeCh, apollo: a}
	a.client.AddChangeListener(listener)
	return &watcher{
		out: changeCh,
		cancelFn: func() {
			a.client.RemoveChangeListener(listener)
			close(changeCh)
		},
	}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	kv, ok := <-w.out
	if !ok {
		return nil, context.Canceled
	}
	return kv, nil
}

func (w *watcher) Stop() error {
	if w.cancelFn != nil {
		w.cancelFn()
	}

	return nil
}
