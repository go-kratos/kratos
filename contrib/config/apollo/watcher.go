package apollo

import (
	"context"
	"strings"

	"github.com/apolloconfig/agollo/v4/storage"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
)

type watcher struct {
	out <-chan []*config.KeyValue

	ctx      context.Context
	cancelFn func()
}

type customChangeListener struct {
	in     chan<- []*config.KeyValue
	apollo *apollo
}

func (c *customChangeListener) onChange(namespace string, changes map[string]*storage.ConfigChange) []*config.KeyValue {
	if strings.Contains(namespace, ".") && !strings.HasSuffix(namespace, "."+properties) && isOriginConfig(namespace) {
		value, err := c.apollo.client.GetConfigCache(namespace).Get("content")
		if err != nil {
			log.Warnw("apollo get config failed", "err", err)
			return nil
		}
		return []*config.KeyValue{
			{
				Key:    namespace,
				Value:  []byte(value.(string)),
				Format: format(namespace),
			},
		}
	}

	next := make(map[string]interface{})

	for key, change := range changes {
		resolve(genKey(namespace, key), change.NewValue, next)
	}

	f := format(namespace)
	val, err := encoding.GetCodec(f).Marshal(next)
	if err != nil {
		log.Warnf("apollo could not handle namespace %s: %v", namespace, err)
		return nil
	}
	return []*config.KeyValue{
		{
			Key:    namespace,
			Value:  val,
			Format: f,
		},
	}
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	change := c.onChange(changeEvent.Namespace, changeEvent.Changes)
	if len(change) == 0 {
		return
	}

	c.in <- change
}

func (c *customChangeListener) OnNewestChange(_ *storage.FullChangeEvent) {}

func newWatcher(a *apollo) (config.Watcher, error) {
	changeCh := make(chan []*config.KeyValue)
	listener := &customChangeListener{in: changeCh, apollo: a}
	a.client.AddChangeListener(listener)

	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{
		out: changeCh,

		ctx: ctx,
		cancelFn: func() {
			a.client.RemoveChangeListener(listener)
			cancel()
		},
	}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case kv := <-w.out:
		return kv, nil
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	if w.cancelFn != nil {
		w.cancelFn()
	}
	return nil
}
