package apollo

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/apolloconfig/agollo/v4/storage"
)

type watcher struct {
	event chan []*config.KeyValue
}

type customChangeListener struct {
	event chan []*config.KeyValue
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	kv := make([]*config.KeyValue, 0)
	for key, value := range changeEvent.Changes {
		kv = append(kv, &config.KeyValue{
			Key:   key,
			Value: []byte(value.NewValue.(string)),
		})
	}
	c.event <- kv
}

func (c *customChangeListener) OnNewestChange(changeEvent *storage.FullChangeEvent) {
	kv := make([]*config.KeyValue, 0)
	for key, value := range changeEvent.Changes {
		kv = append(kv, &config.KeyValue{
			Key:   key,
			Value: []byte(value.(string)),
		})
	}
	c.event <- kv
}

func NewWatcher(a *apollo) (config.Watcher, error) {
	e := make(chan []*config.KeyValue)
	a.client.AddChangeListener(&customChangeListener{event: e})
	return &watcher{event: e}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	return <-w.event, nil
}

func (w *watcher) Stop() error {
	close(w.event)
	return nil
}
