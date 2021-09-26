package apollo

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"

	"github.com/apolloconfig/agollo/v4/storage"
)

type watcher struct {
	event chan []*config.KeyValue
}

type customChangeListener struct {
	event chan []*config.KeyValue
}

func (c *customChangeListener) onChange(
	namespace string, changes map[string]*storage.ConfigChange) []*config.KeyValue {
	kv := make([]*config.KeyValue, 0, 2)
	next := make(map[string]interface{})

	for key, change := range changes {
		convertProperties(genKey(namespace, key), change.NewValue, next)
	}

	f := format(namespace)
	codec := encoding.GetCodec(f)
	val, err := codec.Marshal(next)
	if err != nil {
		fmt.Printf("Warn: apollo could not handle namespace %s: %v\n", namespace, err)
		return nil
	}
	kv = append(kv, &config.KeyValue{
		Key:    namespace,
		Value:  val,
		Format: f,
	})

	return kv
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	c.event <- c.onChange(changeEvent.Namespace, changeEvent.Changes)
}

func (c *customChangeListener) OnNewestChange(changeEvent *storage.FullChangeEvent) {
	// TODO(@yeqown): finish this callback method. but it's not necessarily now.
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
